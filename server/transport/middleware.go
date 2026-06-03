package transport

/*
	Leonidas C2 Framework
	Copyright (C) 2026  Leonidas C2 Project

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sync"
	"time"

	"github.com/leonidas-c2/leonidas/protobuf/clientpb"
	"github.com/leonidas-c2/leonidas/server/configs"
	"github.com/leonidas-c2/leonidas/server/core"
	"github.com/leonidas-c2/leonidas/server/db"
	"github.com/leonidas-c2/leonidas/server/db/models"
	"github.com/leonidas-c2/leonidas/server/log"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_tags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	serverConfig  = configs.GetServerConfig()
	middlewareLog = log.NamedLogger("transport", "middleware")
)

type contextKey int

const (
	Transport contextKey = iota
	Operator
)

// initMiddleware - Initialize middleware
func initMiddleware(enableAuth bool) []grpc.ServerOption {
	logrusEntry := log.NamedLogger("transport", "grpc")
	logrusOpts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(codeToLevel),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)
	if enableAuth {
		return []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				grpc_auth.UnaryServerInterceptor(tokenAuthFunc),
				permissionsUnaryServerInterceptor(),
				auditLogUnaryServerInterceptor(),
				grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
				grpc_logrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
				grpc_logrus.PayloadUnaryServerInterceptor(logrusEntry, deciderUnary),
			),
			grpc.ChainStreamInterceptor(
				grpc_auth.StreamServerInterceptor(tokenAuthFunc),
				permissionsStreamServerInterceptor(),
				trackOperatorStreamInterceptor(),
				grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
				grpc_logrus.StreamServerInterceptor(logrusEntry, logrusOpts...),
				grpc_logrus.PayloadStreamServerInterceptor(logrusEntry, deciderStream),
			),
		}
	} else {
		return []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				grpc_auth.UnaryServerInterceptor(serverAuthFunc),
				auditLogUnaryServerInterceptor(),
				grpc_tags.UnaryServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
				grpc_logrus.UnaryServerInterceptor(logrusEntry, logrusOpts...),
				grpc_logrus.PayloadUnaryServerInterceptor(logrusEntry, deciderUnary),
			),
			grpc.ChainStreamInterceptor(
				grpc_auth.StreamServerInterceptor(serverAuthFunc),
				trackOperatorStreamInterceptor(),
				grpc_tags.StreamServerInterceptor(grpc_tags.WithFieldExtractor(grpc_tags.CodeGenRequestFieldExtractor)),
				grpc_logrus.StreamServerInterceptor(logrusEntry, logrusOpts...),
				grpc_logrus.PayloadStreamServerInterceptor(logrusEntry, deciderStream),
			),
		}
	}

}

var (
	tokenCache = sync.Map{}
)

// ClearTokenCache - Clear the auth token cache
func ClearTokenCache() {
	tokenCache = sync.Map{}
}

func serverAuthFunc(ctx context.Context) (context.Context, error) {
	newCtx := context.WithValue(ctx, Transport, "local")
	newCtx = context.WithValue(newCtx, Operator, "server")
	return newCtx, nil
}

func tokenAuthFunc(ctx context.Context) (context.Context, error) {
	mtlsLog.Debugf("Auth interceptor checking operator token ...")
	rawToken, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		mtlsLog.Errorf("Authentication failure: %s", err)
		return nil, status.Error(codes.Unauthenticated, "Authentication failure")
	}

	// Check auth cache
	digest := sha256.Sum256([]byte(rawToken))
	token := hex.EncodeToString(digest[:])
	newCtx := context.WithValue(ctx, Transport, "mtls")
	if op, ok := tokenCache.Load(token); ok {
		mtlsLog.Debugf("Token in cache!")
		newCtx = context.WithValue(newCtx, Operator, op.(*models.Operator))
		return newCtx, nil
	}
	operator, err := db.OperatorByToken(token)
	if err != nil || operator == nil {
		mtlsLog.Errorf("Authentication failure: %v", err)
		return nil, status.Error(codes.Unauthenticated, "Authentication failure")
	}
	mtlsLog.Debugf("Valid token for %s", operator.Name)
	tokenCache.Store(token, operator)

	newCtx = context.WithValue(newCtx, Operator, operator)
	return newCtx, nil
}

var (
	// Builder - Allowed methods
	builderMethods = map[string]bool{
		"/rpcpb.LeonidasRPC/GetVersion": true,

		"/rpcpb.LeonidasRPC/GenerateExternalGetBuildConfig": true,
		"/rpcpb.LeonidasRPC/GenerateExternalSaveBuild":      true,
		"/rpcpb.LeonidasRPC/BuilderRegister":                true,
		"/rpcpb.LeonidasRPC/BuilderTrigger":                 true,
		"/rpcpb.LeonidasRPC/Builders":                       true,
	}
	// Crackstation - Allowed methods
	crackstationMethods = map[string]bool{
		"/rpcpb.LeonidasRPC/GetVersion": true,

		"/rpcpb.LeonidasRPC/CrackstationRegister":   true,
		"/rpcpb.LeonidasRPC/CrackstationTrigger":    true,
		"/rpcpb.LeonidasRPC/CrackstationBenchmark":  true,
		"/rpcpb.LeonidasRPC/Crackstations":          true,
		"/rpcpb.LeonidasRPC/CrackTaskByID":          true,
		"/rpcpb.LeonidasRPC/CrackTaskUpdate":        true,
		"/rpcpb.LeonidasRPC/CrackFilesList":         true,
		"/rpcpb.LeonidasRPC/CrackFileCreate":        true,
		"/rpcpb.LeonidasRPC/CrackFileChunkUpload":   true,
		"/rpcpb.LeonidasRPC/CrackFileChunkDownload": true,
		"/rpcpb.LeonidasRPC/CrackFileComplete":      true,
		"/rpcpb.LeonidasRPC/CrackFileDelete":        true,
	}
)

func permissionsUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		operator := ctx.Value(Operator).(*models.Operator)
		if operator == nil {
			return nil, status.Error(codes.Unauthenticated, "Authentication failure")
		}
		if operator.PermissionAll {
			return handler(ctx, req)
		}
		if operator.PermissionBuilder {
			if ok, _ := builderMethods[info.FullMethod]; ok {
				return handler(ctx, req)
			}
		}
		if operator.PermissionCrackstation {
			if ok, _ := crackstationMethods[info.FullMethod]; ok {
				return handler(ctx, req)
			}
		}
		mtlsLog.Warnf("Permission denied for %s attempting to access %s", operator.Name, info.FullMethod)
		return nil, status.Error(codes.PermissionDenied, "Token has insufficient permissions")
	}
}

func permissionsStreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		operator := ss.Context().Value(Operator).(*models.Operator)
		if operator == nil {
			return status.Error(codes.Unauthenticated, "Authentication failure")
		}
		if operator.PermissionAll {
			return handler(srv, ss)
		}
		if operator.PermissionBuilder {
			if ok, _ := builderMethods[info.FullMethod]; ok {
				return handler(srv, ss)
			}
		}
		if operator.PermissionCrackstation {
			if ok, _ := crackstationMethods[info.FullMethod]; ok {
				return handler(srv, ss)
			}
		}
		mtlsLog.Warnf("Permission denied for %s attempting to access %s", operator.Name, info.FullMethod)
		return status.Error(codes.PermissionDenied, "Token has insufficient permissions")
	}
}

func deciderUnary(_ context.Context, _ string, _ interface{}) bool {
	return serverConfig.Logs.GRPCUnaryPayloads
}

func deciderStream(_ context.Context, _ string, _ interface{}) bool {
	return serverConfig.Logs.GRPCStreamPayloads
}

// Maps a grpc response code to a logging level
func codeToLevel(code codes.Code) logrus.Level {
	switch code {
	case codes.OK:
		return logrus.InfoLevel
	case codes.Canceled:
		return logrus.InfoLevel
	case codes.Unknown:
		return logrus.ErrorLevel
	case codes.InvalidArgument:
		return logrus.InfoLevel
	case codes.DeadlineExceeded:
		return logrus.WarnLevel
	case codes.NotFound:
		return logrus.InfoLevel
	case codes.AlreadyExists:
		return logrus.InfoLevel
	case codes.PermissionDenied:
		return logrus.WarnLevel
	case codes.Unauthenticated:
		return logrus.InfoLevel
	case codes.ResourceExhausted:
		return logrus.WarnLevel
	case codes.FailedPrecondition:
		return logrus.WarnLevel
	case codes.Aborted:
		return logrus.WarnLevel
	case codes.OutOfRange:
		return logrus.WarnLevel
	case codes.Unimplemented:
		return logrus.ErrorLevel
	case codes.Internal:
		return logrus.ErrorLevel
	case codes.Unavailable:
		return logrus.WarnLevel
	case codes.DataLoss:
		return logrus.ErrorLevel
	default:
		return logrus.ErrorLevel
	}
}

// auditUnaryLogMsg - Structured audit log entry written to audit.json for every
// operator RPC call. Used for post-incident review and compliance.
type auditUnaryLogMsg struct {
	Timestamp  time.Time `json:"timestamp"`
	Method     string    `json:"method"`
	OperatorID string    `json:"operator_id,omitempty"`
	Operator   string    `json:"operator,omitempty"`
	RemoteIP   string    `json:"remote_ip"`
	Session    string    `json:"session,omitempty"`
	Beacon     string    `json:"beacon,omitempty"`
	Request    string    `json:"request"`
	StatusCode string    `json:"status_code"`
	Err        string    `json:"error,omitempty"`
}

func auditLogUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		start := time.Now().UTC()

		rawRequest, err := json.Marshal(req)
		if err != nil {
			middlewareLog.Errorf("Failed to serialize %s", err)
			return
		}
		middlewareLog.Debugf("Raw request: %s", string(rawRequest))
		session, beacon, targetErr := getActiveTarget(rawRequest)
		if targetErr != nil {
			middlewareLog.Errorf("Middleware failed to insert details: %s", targetErr)
		}

		p, _ := peer.FromContext(ctx)

		// Extract operator information set by the auth interceptor.
		operatorName := getUser(p)
		operatorID := ""
		if op, ok := ctx.Value(Operator).(*models.Operator); ok && op != nil {
			operatorName = op.Name
			operatorID = op.ID.String()
		}

		// Execute the real RPC handler.
		resp, handlerErr := handler(ctx, req)

		// Derive a gRPC status code string from the handler result.
		statusCode := codes.OK.String()
		errMsg := ""
		if handlerErr != nil {
			if st, ok := status.FromError(handlerErr); ok {
				statusCode = st.Code().String()
				errMsg = st.Message()
			} else {
				statusCode = codes.Unknown.String()
				errMsg = handlerErr.Error()
			}
		}

		// Construct and write audit log entry.
		msg := &auditUnaryLogMsg{
			Timestamp:  start,
			Method:     info.FullMethod,
			OperatorID: operatorID,
			Operator:   operatorName,
			RemoteIP:   p.Addr.String(),
			Request:    string(rawRequest),
			StatusCode: statusCode,
			Err:        errMsg,
		}
		if session != nil {
			sessionJSON, _ := json.Marshal(session)
			msg.Session = string(sessionJSON)
		}
		if beacon != nil {
			beaconJSON, _ := json.Marshal(beacon)
			msg.Beacon = string(beaconJSON)
		}

		msgData, _ := json.Marshal(msg)
		log.AuditLogger.Info(string(msgData))

		return resp, handlerErr
	}
}

func getUser(client *peer.Peer) string {
	tlsAuth, ok := client.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return ""
	}
	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return ""
	}
	if tlsAuth.State.VerifiedChains[0][0].Subject.CommonName != "" {
		return tlsAuth.State.VerifiedChains[0][0].Subject.CommonName
	}
	return ""
}

func getActiveTarget(rawRequest []byte) (*clientpb.Session, *clientpb.Beacon, error) {

	var activeBeacon *clientpb.Beacon
	var activeSession *clientpb.Session

	var request map[string]interface{}
	err := json.Unmarshal(rawRequest, &request)
	if err != nil {
		return nil, nil, err
	}

	// RPC is not a session/beacon request
	if _, ok := request["Request"]; !ok {
		return nil, nil, nil
	}

	rpcRequest := request["Request"].(map[string]interface{})

	middlewareLog.Debugf("RPC Request: %v", rpcRequest)

	if rawBeaconID, ok := rpcRequest["BeaconID"]; ok {
		beaconID := rawBeaconID.(string)
		middlewareLog.Debugf("Found Beacon ID: %s", beaconID)
		beacon, err := db.BeaconByID(beaconID)
		middlewareLog.Infof("query complete")
		if err != nil {
			middlewareLog.Errorf("Failed to get beacon %s: %s", beaconID, err)
		} else if beacon != nil {
			activeBeacon = beacon.ToProtobuf()
		}
	}

	if rawSessionID, ok := rpcRequest["SessionID"]; ok {
		sessionID := rawSessionID.(string)
		middlewareLog.Debugf("Found Session ID: %s", sessionID)
		session := core.Sessions.Get(sessionID)
		if session != nil {
			activeSession = session.ToProtobuf()
		}
	}

	return activeSession, activeBeacon, nil
}

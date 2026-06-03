package rpc

import (
	"context"

	"github.com/leonidas-c2/leonidas/protobuf/clientpb"
	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"github.com/leonidas-c2/leonidas/server/certs"
	"github.com/leonidas-c2/leonidas/server/generate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GenerateWGClientConfig - Generate a client config for a WG interface
func (rpc *Server) GenerateWGClientConfig(ctx context.Context, _ *commonpb.Empty) (*clientpb.WGClientConfig, error) {
	clientIP, privkey, pubkey, err := generate.GenerateUniqueWGPeerKeys()
	if err != nil {
		rpcLog.Errorf("Could not generate WG keys: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	_, serverPubKey, err := certs.GetWGServerKeys()
	if err != nil {
		rpcLog.Errorf("Could not get WG server keys: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &clientpb.WGClientConfig{
		ClientPrivateKey: privkey,
		ClientIP:         clientIP,
		ClientPubKey:     pubkey,
		ServerPubKey:     serverPubKey,
	}

	return resp, nil
}

// WGStartPortForward - Start a port forward
func (rpc *Server) WGStartPortForward(ctx context.Context, req *leonidaspb.WGPortForwardStartReq) (*leonidaspb.WGPortForward, error) {
	resp := &leonidaspb.WGPortForward{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WGStopPortForward - Stop a port forward
func (rpc *Server) WGStopPortForward(ctx context.Context, req *leonidaspb.WGPortForwardStopReq) (*leonidaspb.WGPortForward, error) {
	resp := &leonidaspb.WGPortForward{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WGAddForwarder - Add a TCP forwarder
func (rpc *Server) WGStartSocks(ctx context.Context, req *leonidaspb.WGSocksStartReq) (*leonidaspb.WGSocks, error) {
	resp := &leonidaspb.WGSocks{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WGStopForwarder - Stop a TCP forwarder
func (rpc *Server) WGStopSocks(ctx context.Context, req *leonidaspb.WGSocksStopReq) (*leonidaspb.WGSocks, error) {
	resp := &leonidaspb.WGSocks{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rpc *Server) WGListSocksServers(ctx context.Context, req *leonidaspb.WGSocksServersReq) (*leonidaspb.WGSocksServers, error) {
	resp := &leonidaspb.WGSocksServers{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// WGAddForwarder - List wireguard forwarders
func (rpc *Server) WGListForwarders(ctx context.Context, req *leonidaspb.WGTCPForwardersReq) (*leonidaspb.WGTCPForwarders, error) {
	resp := &leonidaspb.WGTCPForwarders{}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

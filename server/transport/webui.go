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

// WebUI gRPC-Web proxy
//
// Starts an HTTPS server that:
//   1. Serves the embedded React Web UI static assets on GET requests.
//   2. Proxies gRPC-Web (Content-Type: application/grpc-web+proto) requests to
//      the local gRPC server, translating between HTTP/1.1 and HTTP/2 frames.
//
// The operator gRPC API still requires a Bearer token (same as leonidas-client).
// HTTPS uses optional client certificates so local dev tools (e.g. Vite proxy)
// can connect without browser TLS client certs when bound to loopback.
//
// Usage:
//   _, err := transport.StartWebUI(webBindHost, webPort, grpcBackendAddr)
//
// grpcBackendAddr must match the operator gRPC listener (e.g. "127.0.0.1:31337").

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/leonidas-c2/leonidas/server/certs"
	"github.com/leonidas-c2/leonidas/server/log"
	"github.com/leonidas-c2/leonidas/server/webui"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	webuiLog = log.NamedLogger("transport", "webui")
)

const (
	// DefaultWebUIPort - Default port for the Leonidas Web UI
	DefaultWebUIPort = 8443

	grpcWebContentType     = "application/grpc-web+proto"
	grpcWebTextContentType = "application/grpc-web-text+proto"
	grpcContentType        = "application/grpc"
)

// StartWebUI - Starts the Leonidas Web UI HTTPS server.
// grpcBackendAddr is the TCP address of the operator gRPC API (same port as multiplayer daemon).
// Returns the listening address or an error.
func StartWebUI(host string, port uint16, grpcBackendAddr string) (string, error) {
	if port == 0 {
		port = DefaultWebUIPort
	}
	bind := fmt.Sprintf("%s:%d", host, port)
	if grpcBackendAddr == "" {
		grpcBackendAddr = "127.0.0.1:31337"
	}

	// Operator PKI — same leaf pattern as multiplayer gRPC (see getOperatorServerTLSConfig).
	caCertPtr, _, err := certs.GetCertificateAuthority(certs.OperatorCA)
	if err != nil {
		return "", fmt.Errorf("[webui] load operator CA: %w", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCertPtr)

	webTLSHost := "leonidas-webui"
	_, _, err = certs.OperatorServerGetCertificate(webTLSHost)
	if err == certs.ErrCertDoesNotExist {
		_, _, err = certs.OperatorServerGenerateCertificate(webTLSHost)
		if err != nil {
			return "", fmt.Errorf("[webui] generate operator server certificate: %w", err)
		}
	}
	certPEM, keyPEM, err := certs.OperatorServerGetCertificate(webTLSHost)
	if err != nil {
		return "", fmt.Errorf("[webui] operator server certificate: %w", err)
	}
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return "", fmt.Errorf("[webui] parse server cert: %w", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		ClientCAs:    caCertPool,
		// Allow browsers / Vite dev proxy without a client certificate while still
		// validating certs when operators present them. Bind Web UI to loopback only.
		ClientAuth: tls.VerifyClientCertIfGiven,
		MinVersion: tls.VersionTLS12,
	}

	// Dial the local gRPC server (loopback). Incoming HTTP requests forward Bearer tokens as metadata.
	grpcConn, err := grpc.NewClient(grpcBackendAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return "", fmt.Errorf("[webui] dial gRPC server: %w", err)
	}

	mux := http.NewServeMux()

	// Static assets — served from the embedded webui FS
	mux.Handle("/", http.FileServer(http.FS(webui.FS)))

	// gRPC-Web proxy — all /rpcpb.LeonidasRPC/* paths
	mux.HandleFunc("/rpcpb.LeonidasRPC/", func(w http.ResponseWriter, r *http.Request) {
		ct := r.Header.Get("Content-Type")
		if !strings.HasPrefix(ct, grpcWebContentType) && !strings.HasPrefix(ct, grpcWebTextContentType) {
			http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
			return
		}
		proxyGRPCWeb(w, r, grpcConn)
	})

	srv := &http.Server{
		Addr:      bind,
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	ln, err := net.Listen("tcp", bind)
	if err != nil {
		return "", fmt.Errorf("[webui] listen %s: %w", bind, err)
	}
	tlsLn := tls.NewListener(ln, tlsConfig)

	go func() {
		webuiLog.Infof("Leonidas Web UI listening on https://%s", bind)
		if err := srv.Serve(tlsLn); err != nil && err != http.ErrServerClosed {
			webuiLog.Errorf("Web UI server error: %v", err)
		}
	}()

	return bind, nil
}

// proxyGRPCWeb translates a gRPC-Web HTTP/1.1 request into a gRPC HTTP/2
// call forwarded to the internal gRPC server via the pre-established conn.
//
// gRPC-Web wire format (binary):
//   - Request body: standard gRPC length-prefixed messages (5-byte header + data).
//   - Response:     same format; trailers appended as a trailing message with
//                   the high bit of the compression flag set to 1 (0x80).
func proxyGRPCWeb(w http.ResponseWriter, r *http.Request, grpcConn *grpc.ClientConn) {
	ctx := r.Context()

	// Build gRPC metadata from incoming HTTP headers.
	md := make(map[string]string)
	for k, vs := range r.Header {
		k = strings.ToLower(k)
		if k == "content-type" || strings.HasPrefix(k, ":") {
			continue
		}
		md[k] = strings.Join(vs, ",")
	}

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(md))

	// Read all request frames.
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the gRPC message frames (5-byte header: 1 byte flags + 4 bytes length).
	var payload []byte
	offset := 0
	for offset+5 <= len(body) {
		// compression flag (ignored) + 4-byte big-endian length
		msgLen := int(binary.BigEndian.Uint32(body[offset+1 : offset+5]))
		offset += 5
		if offset+msgLen > len(body) {
			break
		}
		payload = append(payload, body[offset:offset+msgLen]...)
		offset += msgLen
	}

	// Invoke via gRPC ClientConn using the raw method path.
	methodPath := r.URL.Path // e.g. /rpcpb.LeonidasRPC/GetSessions
	var respBuf bytes.Buffer
	codec := rawBytesCodec{}

	err = grpcConn.Invoke(
		ctx,
		methodPath,
		payload,
		&respBuf,
		grpc.ForceCodec(codec),
		grpc.Header(nil),
	)

	w.Header().Set("Content-Type", grpcWebContentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "grpc-status,grpc-message")

	if err != nil {
		w.Header().Set("grpc-status", "2")
		w.Header().Set("grpc-message", err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	// Wrap response in a gRPC-Web data frame.
	respBytes := respBuf.Bytes()
	frameHeader := make([]byte, 5)
	frameHeader[0] = 0 // no compression
	binary.BigEndian.PutUint32(frameHeader[1:], uint32(len(respBytes)))
	w.Write(frameHeader)
	w.Write(respBytes)

	// Write gRPC status trailer frame (flag byte 0x80 = trailer).
	trailer := "grpc-status: 0\r\n"
	trailerHeader := make([]byte, 5)
	trailerHeader[0] = 0x80
	binary.BigEndian.PutUint32(trailerHeader[1:], uint32(len(trailer)))
	w.Write(trailerHeader)
	w.Write([]byte(trailer))
}

// rawBytesCodec is a grpc.Codec that passes pre-serialised bytes through
// without any further marshalling — needed for the raw proxy path.
type rawBytesCodec struct{}

func (rawBytesCodec) Name() string { return "proto" }

func (rawBytesCodec) Marshal(v interface{}) ([]byte, error) {
	if b, ok := v.([]byte); ok {
		return b, nil
	}
	return nil, fmt.Errorf("rawBytesCodec: expected []byte, got %T", v)
}

func (rawBytesCodec) Unmarshal(data []byte, v interface{}) error {
	if buf, ok := v.(*bytes.Buffer); ok {
		buf.Write(data)
		return nil
	}
	return fmt.Errorf("rawBytesCodec: expected *bytes.Buffer, got %T", v)
}


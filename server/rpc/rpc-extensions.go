package rpc

import (
	"context"

	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
)

// RegisterExtension registers a new extension in the implant
func (rpc *Server) RegisterExtension(ctx context.Context, req *leonidaspb.RegisterExtensionReq) (*leonidaspb.RegisterExtension, error) {
	resp := &leonidaspb.RegisterExtension{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ListExtensions lists the registered extensions
func (rpc *Server) ListExtensions(ctx context.Context, req *leonidaspb.ListExtensionsReq) (*leonidaspb.ListExtensions, error) {
	resp := &leonidaspb.ListExtensions{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CallExtension calls a specific export of the loaded extension
func (rpc *Server) CallExtension(ctx context.Context, req *leonidaspb.CallExtensionReq) (*leonidaspb.CallExtension, error) {
	resp := &leonidaspb.CallExtension{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

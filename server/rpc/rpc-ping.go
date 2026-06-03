package rpc

import (
	"context"

	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
)

// Ping - Try to send a round trip message to the implant
func (rpc *Server) Ping(ctx context.Context, req *leonidaspb.Ping) (*leonidaspb.Ping, error) {
	resp := &leonidaspb.Ping{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

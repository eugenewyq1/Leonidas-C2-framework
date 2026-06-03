package rpc

import (
	"context"

	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
)

// Services - List and control services
func (rpc *Server) Services(ctx context.Context, req *leonidaspb.ServicesReq) (*leonidaspb.Services, error) {
	resp := &leonidaspb.Services{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rpc *Server) ServiceDetail(ctx context.Context, req *leonidaspb.ServiceDetailReq) (*leonidaspb.ServiceDetail, error) {
	resp := &leonidaspb.ServiceDetail{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// StartService creates and starts a Windows service on a remote host
func (rpc *Server) StartService(ctx context.Context, req *leonidaspb.StartServiceReq) (*leonidaspb.ServiceInfo, error) {
	resp := &leonidaspb.ServiceInfo{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (rpc *Server) StartServiceByName(ctx context.Context, req *leonidaspb.StartServiceByNameReq) (*leonidaspb.ServiceInfo, error) {
	resp := &leonidaspb.ServiceInfo{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// StopService stops a remote service
func (rpc *Server) StopService(ctx context.Context, req *leonidaspb.StopServiceReq) (*leonidaspb.ServiceInfo, error) {
	resp := &leonidaspb.ServiceInfo{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RemoveService deletes a service from the remote system
func (rpc *Server) RemoveService(ctx context.Context, req *leonidaspb.RemoveServiceReq) (*leonidaspb.ServiceInfo, error) {
	resp := &leonidaspb.ServiceInfo{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

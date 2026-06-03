package rpc

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

	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
)

// Ps - List the processes on the remote machine
func (rpc *Server) Ps(ctx context.Context, req *leonidaspb.PsReq) (*leonidaspb.Ps, error) {
	resp := &leonidaspb.Ps{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ProcessDump - Dump the memory of a remote process
func (rpc *Server) ProcessDump(ctx context.Context, req *leonidaspb.ProcessDumpReq) (*leonidaspb.ProcessDump, error) {
	resp := &leonidaspb.ProcessDump{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Terminate - Terminate a remote process
func (rpc *Server) Terminate(ctx context.Context, req *leonidaspb.TerminateReq) (*leonidaspb.Terminate, error) {
	resp := &leonidaspb.Terminate{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

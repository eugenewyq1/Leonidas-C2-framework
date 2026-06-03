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
	"crypto/sha256"
	"fmt"

	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"github.com/leonidas-c2/leonidas/server/core"
	"github.com/leonidas-c2/leonidas/server/db"
	"github.com/leonidas-c2/leonidas/server/db/models"
	"github.com/leonidas-c2/leonidas/server/log"
	"github.com/leonidas-c2/leonidas/util/encoders"
)

var (
	fsLog = log.NamedLogger("rcp", "fs")
)

// Ls - List a directory
func (rpc *Server) Ls(ctx context.Context, req *leonidaspb.LsReq) (*leonidaspb.Ls, error) {
	resp := &leonidaspb.Ls{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Mv - Move or rename a file
func (rpc *Server) Mv(ctx context.Context, req *leonidaspb.MvReq) (*leonidaspb.Mv, error) {
	resp := &leonidaspb.Mv{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Cp - Copy a file to another location
func (rpc *Server) Cp(ctx context.Context, req *leonidaspb.CpReq) (*leonidaspb.Cp, error) {
	resp := &leonidaspb.Cp{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Rm - Remove file or directory
func (rpc *Server) Rm(ctx context.Context, req *leonidaspb.RmReq) (*leonidaspb.Rm, error) {
	resp := &leonidaspb.Rm{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Mkdir - Make a directory
func (rpc *Server) Mkdir(ctx context.Context, req *leonidaspb.MkdirReq) (*leonidaspb.Mkdir, error) {
	resp := &leonidaspb.Mkdir{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Cd - Change directory
func (rpc *Server) Cd(ctx context.Context, req *leonidaspb.CdReq) (*leonidaspb.Pwd, error) {
	resp := &leonidaspb.Pwd{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Pwd - Print working directory
func (rpc *Server) Pwd(ctx context.Context, req *leonidaspb.PwdReq) (*leonidaspb.Pwd, error) {
	resp := &leonidaspb.Pwd{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Download - Download a file from the remote file system
func (rpc *Server) Download(ctx context.Context, req *leonidaspb.DownloadReq) (*leonidaspb.Download, error) {
	resp := &leonidaspb.Download{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Upload - Upload a file from the remote file system
func (rpc *Server) Upload(ctx context.Context, req *leonidaspb.UploadReq) (*leonidaspb.Upload, error) {
	resp := &leonidaspb.Upload{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	if req.IsIOC {
		go trackIOC(req, resp)
	}
	return resp, nil
}

// Chmod - Change permission on a file or directory
func (rpc *Server) Chmod(ctx context.Context, req *leonidaspb.ChmodReq) (*leonidaspb.Chmod, error) {
	resp := &leonidaspb.Chmod{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Chown - Change owner on a file or directory
func (rpc *Server) Chown(ctx context.Context, req *leonidaspb.ChownReq) (*leonidaspb.Chown, error) {
	resp := &leonidaspb.Chown{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Chtimes - Change file access and modification times on a file or directory
func (rpc *Server) Chtimes(ctx context.Context, req *leonidaspb.ChtimesReq) (*leonidaspb.Chtimes, error) {
	resp := &leonidaspb.Chtimes{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MemfilesList - List memfiles
func (rpc *Server) MemfilesList(ctx context.Context, req *leonidaspb.MemfilesListReq) (*leonidaspb.Ls, error) {
	resp := &leonidaspb.Ls{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MemfilesAdd - Add memfile
func (rpc *Server) MemfilesAdd(ctx context.Context, req *leonidaspb.MemfilesAddReq) (*leonidaspb.MemfilesAdd, error) {
	resp := &leonidaspb.MemfilesAdd{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// MemfilesRm - Close memfile
func (rpc *Server) MemfilesRm(ctx context.Context, req *leonidaspb.MemfilesRmReq) (*leonidaspb.MemfilesRm, error) {
	resp := &leonidaspb.MemfilesRm{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func hashUploadData(encoder string, data []byte) [32]byte {
	if encoder == "gzip" {
		decodedData, err := new(encoders.Gzip).Decode(data)
		if err != nil {
			return sha256.Sum256(nil)
		}
		return sha256.Sum256(decodedData)
	} else {
		return sha256.Sum256(data)
	}
}

func trackIOC(req *leonidaspb.UploadReq, resp *leonidaspb.Upload) {
	fsLog.Debugf("Adding IOC to database ...")
	request := req.GetRequest()
	if request == nil {
		fsLog.Error("No request for upload")
		return
	}
	session := core.Sessions.Get(request.SessionID)
	if session == nil {
		fsLog.Error("No session for upload request")
		return
	}
	host, err := db.HostByHostUUID(session.UUID)
	if err != nil {
		fsLog.Errorf("No host for session uuid %v", session.UUID)
		return
	}

	sum := hashUploadData(req.Encoder, req.Data)
	ioc := &models.IOC{
		HostID:   host.HostUUID,
		Path:     resp.Path,
		FileHash: fmt.Sprintf("%x", sum),
	}
	if db.Session().Create(ioc).Error != nil {
		fsLog.Error("Failed to create IOC")
	}
}

// Grep - Search a file or directory for text matching a regex
func (rpc *Server) Grep(ctx context.Context, req *leonidaspb.GrepReq) (*leonidaspb.Grep, error) {
	resp := &leonidaspb.Grep{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Mount - Get information on mounted filesystems
func (rpc *Server) Mount(ctx context.Context, req *leonidaspb.MountReq) (*leonidaspb.Mount, error) {
	resp := &leonidaspb.Mount{Response: &commonpb.Response{}}
	err := rpc.GenericHandler(req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

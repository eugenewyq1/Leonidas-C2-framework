package filesystem

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

	"github.com/leonidas-c2/leonidas/client/console"
	"github.com/leonidas-c2/leonidas/protobuf/clientpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

// PwdCmd - Print the remote working directory.
func PwdCmd(cmd *cobra.Command, con *console.SliverClient, args []string) {
	session, beacon := con.ActiveTarget.GetInteractive()
	if session == nil && beacon == nil {
		return
	}
	pwd, err := con.Rpc.Pwd(context.Background(), &leonidaspb.PwdReq{
		Request: con.ActiveTarget.Request(cmd),
	})
	if err != nil {
		con.PrintErrorf("%s\n", err)
		return
	}
	if pwd.Response != nil && pwd.Response.Async {
		con.AddBeaconCallback(pwd.Response.TaskID, func(task *clientpb.BeaconTask) {
			err = proto.Unmarshal(task.Response, pwd)
			if err != nil {
				con.PrintErrorf("Failed to decode response %s\n", err)
				return
			}
			PrintPwd(pwd, con)
		})
		con.PrintAsyncResponse(pwd.Response)
	} else {
		PrintPwd(pwd, con)
	}
}

// PrintPwd - Print the remote working directory.
func PrintPwd(pwd *leonidaspb.Pwd, con *console.SliverClient) {
	if pwd.Response != nil && pwd.Response.Err != "" {
		con.PrintErrorf("%s\n", pwd.Response.Err)
		return
	}
	con.PrintInfof("%s\n", pwd.Path)
}

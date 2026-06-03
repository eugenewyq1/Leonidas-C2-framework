package privilege

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

	"google.golang.org/protobuf/proto"

	"github.com/spf13/cobra"

	"github.com/leonidas-c2/leonidas/client/console"
	"github.com/leonidas-c2/leonidas/protobuf/clientpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
)

// RevToSelfCmd - Drop any impersonated tokens
func RevToSelfCmd(cmd *cobra.Command, con *console.SliverClient, args []string) {
	session, beacon := con.ActiveTarget.GetInteractive()
	if session == nil && beacon == nil {
		return
	}

	revert, err := con.Rpc.RevToSelf(context.Background(), &leonidaspb.RevToSelfReq{
		Request: con.ActiveTarget.Request(cmd),
	})
	if err != nil {
		con.PrintErrorf("%s\n", err)
		return
	}

	if revert.Response != nil && revert.Response.Async {
		con.AddBeaconCallback(revert.Response.TaskID, func(task *clientpb.BeaconTask) {
			err = proto.Unmarshal(task.Response, revert)
			if err != nil {
				con.PrintErrorf("Failed to decode response %s\n", err)
				return
			}
			PrintRev2Self(revert, con)
		})
		con.PrintAsyncResponse(revert.Response)
	} else {
		PrintRev2Self(revert, con)
	}
}

// PrintRev2Self - Print the result of revert to self
func PrintRev2Self(revert *leonidaspb.RevToSelf, con *console.SliverClient) {
	if revert.Response != nil && revert.Response.GetErr() != "" {
		con.PrintErrorf("%s\n", revert.Response.GetErr())
		return
	}
	con.PrintInfof("Back to self...")
}

package registry

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

// RegListSubKeysCmd - List sub registry keys
func RegListSubKeysCmd(cmd *cobra.Command, con *console.SliverClient, args []string) {
	session, beacon := con.ActiveTarget.GetInteractive()
	if session == nil && beacon == nil {
		return
	}
	targetOS := getOS(session, beacon)
	if targetOS != "windows" {
		con.PrintErrorf("Registry operations can only target Windows\n")
		return
	}

	regPath := args[0]
	hive, _ := cmd.Flags().GetString("hive")
	hostname, _ := cmd.Flags().GetString("hostname")

	regList, err := con.Rpc.RegistryListSubKeys(context.Background(), &leonidaspb.RegistrySubKeyListReq{
		Hive:     hive,
		Hostname: hostname,
		Path:     regPath,
		Request:  con.ActiveTarget.Request(cmd),
	})
	if err != nil {
		con.PrintErrorf("%s\n", err)
		return
	}

	if regList.Response != nil && regList.Response.Async {
		con.AddBeaconCallback(regList.Response.TaskID, func(task *clientpb.BeaconTask) {
			err = proto.Unmarshal(task.Response, regList)
			if err != nil {
				con.PrintErrorf("Failed to decode response %s\n", err)
				return
			}
			PrintListSubKeys(regList, hive, regPath, con)
		})
		con.PrintAsyncResponse(regList.Response)
	} else {
		PrintListSubKeys(regList, hive, regPath, con)
	}
}

// PrintListSubKeys - Print the list sub keys command result
func PrintListSubKeys(regList *leonidaspb.RegistrySubKeyList, hive string, regPath string, con *console.SliverClient) {
	if regList.Response != nil && regList.Response.Err != "" {
		con.PrintErrorf("%s\n", regList.Response.Err)
		return
	}
	if 0 < len(regList.Subkeys) {
		con.PrintInfof("Sub keys under %s:\\%s:\n", hive, regPath)
	}
	for _, subKey := range regList.Subkeys {
		con.Println(subKey)
	}
}

// RegListValuesCmd - List registry values
func RegListValuesCmd(cmd *cobra.Command, con *console.SliverClient, args []string) {
	session, beacon := con.ActiveTarget.GetInteractive()
	if session == nil && beacon == nil {
		return
	}

	regPath := args[0]
	hive, _ := cmd.Flags().GetString("hive")
	hostname, _ := cmd.Flags().GetString("hostname")

	regList, err := con.Rpc.RegistryListValues(context.Background(), &leonidaspb.RegistryListValuesReq{
		Hive:     hive,
		Hostname: hostname,
		Path:     regPath,
		Request:  con.ActiveTarget.Request(cmd),
	})
	if err != nil {
		con.PrintErrorf("%s\n", err)
		return
	}

	if regList.Response != nil && regList.Response.Async {
		con.AddBeaconCallback(regList.Response.TaskID, func(task *clientpb.BeaconTask) {
			err = proto.Unmarshal(task.Response, regList)
			if err != nil {
				con.PrintErrorf("Failed to decode response %s\n", err)
				return
			}
			PrintListValues(regList, hive, regPath, con)
		})
		con.PrintAsyncResponse(regList.Response)
	} else {
		PrintListValues(regList, hive, regPath, con)
	}
}

// PrintListValues - Print the registry list values
func PrintListValues(regList *leonidaspb.RegistryValuesList, hive string, regPath string, con *console.SliverClient) {
	if regList.Response != nil && regList.Response.Err != "" {
		con.PrintErrorf("%s\n", regList.Response.Err)
		return
	}
	if 0 < len(regList.ValueNames) {
		con.PrintInfof("Values under %s:\\%s:\n", hive, regPath)
	}
	for _, val := range regList.ValueNames {
		con.Println(val)
	}
}

package network

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
	"fmt"
	"net"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"

	"github.com/leonidas-c2/leonidas/client/command/settings"
	"github.com/leonidas-c2/leonidas/client/console"
	"github.com/leonidas-c2/leonidas/protobuf/clientpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
)

// NetstatCmd - Display active network connections on the remote system
func NetstatCmd(cmd *cobra.Command, con *console.SliverClient, args []string) {
	session, beacon := con.ActiveTarget.GetInteractive()
	if session == nil && beacon == nil {
		return
	}

	listening, _ := cmd.Flags().GetBool("listen")
	ip4, _ := cmd.Flags().GetBool("ip4")
	ip6, _ := cmd.Flags().GetBool("ip6")
	tcp, _ := cmd.Flags().GetBool("tcp")
	udp, _ := cmd.Flags().GetBool("udp")
	numeric, _ := cmd.Flags().GetBool("numeric")

	implantPID := getPID(session, beacon)
	activeC2 := getActiveC2(session, beacon)

	netstat, err := con.Rpc.Netstat(context.Background(), &leonidaspb.NetstatReq{
		Request:   con.ActiveTarget.Request(cmd),
		TCP:       tcp,
		UDP:       udp,
		Listening: listening,
		IP4:       ip4,
		IP6:       ip6,
	})
	if err != nil {
		con.PrintErrorf("%s\n", err)
		return
	}
	if netstat.Response != nil && netstat.Response.Async {
		con.AddBeaconCallback(netstat.Response.TaskID, func(task *clientpb.BeaconTask) {
			err = proto.Unmarshal(task.Response, netstat)
			if err != nil {
				con.PrintErrorf("Failed to decode response %s\n", err)
				return
			}
			PrintNetstat(netstat, implantPID, activeC2, numeric, con)
		})
		con.PrintAsyncResponse(netstat.Response)
	} else {
		PrintNetstat(netstat, implantPID, activeC2, numeric, con)
	}
}

func PrintNetstat(netstat *leonidaspb.Netstat, implantPID int32, activeC2 string, numeric bool, con *console.SliverClient) {
	lookup := func(skaddr *leonidaspb.SockTabEntry_SockAddr) string {
		addr := skaddr.Ip
		names, err := net.LookupAddr(addr)
		if err == nil && len(names) > 0 {
			addr = names[0]
		}
		return fmt.Sprintf("%s:%d", addr, skaddr.Port)
	}

	tw := table.NewWriter()
	tw.SetStyle(settings.GetTableStyle(con))
	tw.AppendHeader(table.Row{"Protocol", "Local Address", "Foreign Address", "State", "PID/Program name"})

	for _, entry := range netstat.Entries {
		pid := ""
		if entry.Process != nil {
			pid = fmt.Sprintf("%d/%s", entry.Process.Pid, entry.Process.Executable)
		}
		srcAddr := fmt.Sprintf("%s:%d", entry.LocalAddr.Ip, entry.LocalAddr.Port)
		dstAddr := fmt.Sprintf("%s:%d", entry.RemoteAddr.Ip, entry.RemoteAddr.Port)
		if !numeric {
			srcAddr = lookup(entry.LocalAddr)
			dstAddr = lookup(entry.RemoteAddr)
		}
		if entry.Process != nil && entry.Process.Pid == implantPID {
			tw.AppendRow(table.Row{
				console.StyleGreen.Render(entry.Protocol),
				console.StyleGreen.Render(srcAddr),
				console.StyleGreen.Render(dstAddr),
				console.StyleGreen.Render(entry.SkState),
				console.StyleGreen.Render(pid),
			})
		} else {
			tw.AppendRow(table.Row{entry.Protocol, srcAddr, dstAddr, entry.SkState, pid})
		}
	}
	con.Printf("%s\n", tw.Render())
}

func getActiveC2(session *clientpb.Session, beacon *clientpb.Beacon) string {
	if session != nil {
		return session.ActiveC2
	}
	if beacon != nil {
		return beacon.ActiveC2
	}
	return ""
}

func getPID(session *clientpb.Session, beacon *clientpb.Beacon) int32 {
	if session != nil {
		return session.PID
	}
	if beacon != nil {
		return beacon.PID
	}
	return -1
}

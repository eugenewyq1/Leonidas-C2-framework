package info

import (
	"context"

	"github.com/leonidas-c2/leonidas/client/console"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"github.com/leonidas-c2/leonidas/util"
	"github.com/spf13/cobra"
)

// PingCmd - Send a round trip C2 message to an implant (does not use ICMP).
func PingCmd(cmd *cobra.Command, con *console.SliverClient, args []string) {
	session := con.ActiveTarget.GetSessionInteractive()
	if session == nil {
		return
	}

	nonce := util.Intn(999999)
	con.PrintInfof("Ping %d\n", nonce)
	pong, err := con.Rpc.Ping(context.Background(), &leonidaspb.Ping{
		Nonce:   int32(nonce),
		Request: con.ActiveTarget.Request(cmd),
	})
	if err != nil {
		con.PrintErrorf("%s\n", err)
	} else {
		con.PrintInfof("Pong %d\n", pong.Nonce)
	}
}

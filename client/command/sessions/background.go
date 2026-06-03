package sessions

import (
	"github.com/leonidas-c2/leonidas/client/console"
	"github.com/spf13/cobra"
)

// BackgroundCmd - Background the active session.
func BackgroundCmd(cmd *cobra.Command, con *console.SliverClient, args []string) {
	con.ActiveTarget.Background()
	con.PrintInfof("Background ...\n")
}

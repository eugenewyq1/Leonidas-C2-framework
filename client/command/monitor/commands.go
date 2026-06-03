package monitor

import (
	"github.com/spf13/cobra"

	"github.com/leonidas-c2/leonidas/client/console"
	consts "github.com/leonidas-c2/leonidas/client/constants"
)

// Commands returns the “ command and its subcommands.
func Commands(con *console.SliverClient) []*cobra.Command {
	monitorCmd := &cobra.Command{
		Use:     consts.MonitorStr,
		Short:   "Monitor threat intel platforms for Sliver implants",
		GroupID: consts.SliverHelpGroup,
	}
	monitorCmd.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "Start the monitoring loops",
		Run: func(cmd *cobra.Command, args []string) {
			MonitorStartCmd(cmd, con, args)
		},
	})
	monitorCmd.AddCommand(&cobra.Command{
		Use:   "stop",
		Short: "Stop the monitoring loops",
		Run: func(cmd *cobra.Command, args []string) {
			MonitorStopCmd(cmd, con, args)
		},
	})

	return []*cobra.Command{monitorCmd}
}

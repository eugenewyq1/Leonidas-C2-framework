package socks

import (
	"github.com/leonidas-c2/leonidas/client/command/flags"
	"github.com/leonidas-c2/leonidas/client/command/help"
	"github.com/leonidas-c2/leonidas/client/console"
	consts "github.com/leonidas-c2/leonidas/client/constants"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Commands returns the “ command and its subcommands.
func RootCommands(con *console.SliverClient) []*cobra.Command {
	socksCmd := &cobra.Command{
		Use:         consts.Socks5Str,
		Short:       "In-band SOCKS5 Proxy",
		Long:        help.GetHelpFor([]string{consts.Socks5Str}),
		GroupID:     consts.NetworkHelpGroup,
		Annotations: flags.RestrictTargets(consts.SessionCmdsFilter),
		Run: func(cmd *cobra.Command, args []string) {
			SocksCmd(cmd, con, args)
		},
	}
	flags.Bind("", true, socksCmd, func(f *pflag.FlagSet) {
		f.Int64P("timeout", "t", flags.DefaultTimeout, "grpc timeout in seconds")
	})

	socksStopCmd := &cobra.Command{
		Use:   consts.StopStr,
		Short: "Stop a SOCKS5 proxy",
		Long:  help.GetHelpFor([]string{consts.Socks5Str}),
		Run: func(cmd *cobra.Command, args []string) {
			SocksStopCmd(cmd, con, args)
		},
	}
	socksCmd.AddCommand(socksStopCmd)
	flags.Bind("", false, socksStopCmd, func(f *pflag.FlagSet) {
		f.Uint64P("id", "i", 0, "id of portfwd to remove")
	})
	flags.BindFlagCompletions(socksStopCmd, func(comp *carapace.ActionMap) {
		(*comp)["id"] = SocksIDCompleter(con)
	})

	return []*cobra.Command{socksCmd}
}

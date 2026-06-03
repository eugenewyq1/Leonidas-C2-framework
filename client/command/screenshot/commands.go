package screenshot

import (
	"github.com/leonidas-c2/leonidas/client/command/completers"
	"github.com/leonidas-c2/leonidas/client/command/flags"
	"github.com/leonidas-c2/leonidas/client/command/help"
	"github.com/leonidas-c2/leonidas/client/console"
	consts "github.com/leonidas-c2/leonidas/client/constants"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Commands returns the “ command and its subcommands.
func Commands(con *console.SliverClient) []*cobra.Command {
	screenshotCmd := &cobra.Command{
		Use:   consts.ScreenshotStr,
		Short: "Take a screenshot",
		Long:  help.GetHelpFor([]string{consts.ScreenshotStr}),
		Run: func(cmd *cobra.Command, args []string) {
			ScreenshotCmd(cmd, con, args)
		},
		GroupID: consts.InfoHelpGroup,
	}
	flags.Bind("", false, screenshotCmd, func(f *pflag.FlagSet) {
		f.StringP("save", "s", "", "save to file (will overwrite if exists)")
		f.BoolP("loot", "X", false, "save output as loot")
		f.StringP("name", "n", "", "name to assign loot (optional)")

		f.Int64P("timeout", "t", flags.DefaultTimeout, "grpc timeout in seconds")
	})
	flags.BindFlagCompletions(screenshotCmd, func(comp *carapace.ActionMap) {
		(*comp)["save"] = carapace.ActionFiles()
	})
	completers.RegisterLocalFilePathFlagCompletion(screenshotCmd, "save")

	return []*cobra.Command{screenshotCmd}
}

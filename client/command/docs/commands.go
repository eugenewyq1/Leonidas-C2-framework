package docs

import (
	"github.com/leonidas-c2/leonidas/client/command/help"
	"github.com/leonidas-c2/leonidas/client/console"
	consts "github.com/leonidas-c2/leonidas/client/constants"
	"github.com/spf13/cobra"
)

// Commands returns the docs command.
func Commands(con *console.SliverClient) []*cobra.Command {
	return []*cobra.Command{newDocsCommand(consts.SliverCoreHelpGroup, con)}
}

// ServerCommands returns the docs command for the top-level client REPL.
func ServerCommands(con *console.SliverClient) []*cobra.Command {
	return []*cobra.Command{newDocsCommand(consts.GenericHelpGroup, con)}
}

func newDocsCommand(groupID string, con *console.SliverClient) *cobra.Command {
	return &cobra.Command{
		Use:     consts.DocsStr,
		Short:   "Browse the embedded Sliver docs in a TUI",
		Long:    help.GetHelpFor([]string{consts.DocsStr}),
		Args:    cobra.NoArgs,
		GroupID: groupID,
		Run: func(cmd *cobra.Command, args []string) {
			DocsCmd(cmd, con, args)
		},
	}
}

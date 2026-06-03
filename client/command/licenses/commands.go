package licenses

import (
	"github.com/spf13/cobra"

	"github.com/leonidas-c2/leonidas/client/command/help"
	"github.com/leonidas-c2/leonidas/client/console"
	consts "github.com/leonidas-c2/leonidas/client/constants"
	"github.com/leonidas-c2/leonidas/client/licenses"
)

// Commands returns the `licences` command.
func Commands(con *console.SliverClient) []*cobra.Command {
	licensesCmd := &cobra.Command{
		Use:   consts.LicensesStr,
		Short: "Open source licenses",
		Long:  help.GetHelpFor([]string{consts.LicensesStr}),
		Run: func(cmd *cobra.Command, args []string) {
			con.Println(licenses.All)
		},
		GroupID: consts.GenericHelpGroup,
	}

	return []*cobra.Command{licensesCmd}
}

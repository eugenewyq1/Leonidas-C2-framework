//go:build !client

package serverctx

import (
	"github.com/leonidas-c2/leonidas/client/console"
	"github.com/spf13/cobra"
)

// Commands is a no-op when building without the `client` build tag (e.g. leonidas-server).
func Commands(_ *console.SliverClient) []*cobra.Command {
	return nil
}

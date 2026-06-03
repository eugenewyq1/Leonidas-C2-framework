package cli

import (
	"github.com/leonidas-c2/leonidas/client/transport"
	"github.com/spf13/cobra"
)

const (
	enableWGFlag = "enable-wg"
)

func applyMultiplayerConnectMode(cmd *cobra.Command) error {
	if cmd == nil {
		transport.SetMultiplayerConnectMode(transport.MultiplayerConnectDirect)
		return nil
	}

	enableWG, err := cmd.Flags().GetBool(enableWGFlag)
	if err != nil {
		return err
	}

	mode := transport.MultiplayerConnectDirect
	if enableWG {
		mode = transport.MultiplayerConnectEnableWG
	}
	transport.SetMultiplayerConnectMode(mode)
	return nil
}

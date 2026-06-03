package command

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
	"fmt"

	"github.com/leonidas-c2/leonidas/client/assets"
	"github.com/leonidas-c2/leonidas/client/command/ai"
	"github.com/leonidas-c2/leonidas/client/command/aka"
	"github.com/leonidas-c2/leonidas/client/command/alias"
	"github.com/leonidas-c2/leonidas/client/command/backdoor"
	"github.com/leonidas-c2/leonidas/client/command/cursed"
	"github.com/leonidas-c2/leonidas/client/command/dllhijack"
	docscmd "github.com/leonidas-c2/leonidas/client/command/docs"
	"github.com/leonidas-c2/leonidas/client/command/edit"
	"github.com/leonidas-c2/leonidas/client/command/environment"
	"github.com/leonidas-c2/leonidas/client/command/exec"
	"github.com/leonidas-c2/leonidas/client/command/extensions"
	"github.com/leonidas-c2/leonidas/client/command/filesystem"
	"github.com/leonidas-c2/leonidas/client/command/hexedit"
	"github.com/leonidas-c2/leonidas/client/command/info"
	"github.com/leonidas-c2/leonidas/client/command/kill"
	"github.com/leonidas-c2/leonidas/client/command/network"
	"github.com/leonidas-c2/leonidas/client/command/pivots"
	"github.com/leonidas-c2/leonidas/client/command/portfwd"
	"github.com/leonidas-c2/leonidas/client/command/privilege"
	"github.com/leonidas-c2/leonidas/client/command/processes"
	"github.com/leonidas-c2/leonidas/client/command/reconfig"
	"github.com/leonidas-c2/leonidas/client/command/registry"
	"github.com/leonidas-c2/leonidas/client/command/rportfwd"
	"github.com/leonidas-c2/leonidas/client/command/screenshot"
	"github.com/leonidas-c2/leonidas/client/command/sessions"
	"github.com/leonidas-c2/leonidas/client/command/shell"
	"github.com/leonidas-c2/leonidas/client/command/socks"
	"github.com/leonidas-c2/leonidas/client/command/tasks"
	"github.com/leonidas-c2/leonidas/client/command/wasm"
	"github.com/leonidas-c2/leonidas/client/command/wireguard"
	client "github.com/leonidas-c2/leonidas/client/console"
	consts "github.com/leonidas-c2/leonidas/client/constants"
	"github.com/reeflective/console"
	"github.com/spf13/cobra"
)

// SliverCommands returns all commands bound to the implant menu.
func SliverCommands(con *client.SliverClient) console.Commands {
	sliverCommands := func() *cobra.Command {
		sliver := &cobra.Command{
			Short: "Implant commands",
			CompletionOptions: cobra.CompletionOptions{
				HiddenDefaultCmd: true,
			},
		}
		if !con.IsCLI {
			sliver.SilenceErrors = true
			sliver.SilenceUsage = true
		}

		// Utility function to be used for binding new commands to
		// the sliver menu: call the function with the name of the
		// group under which this/these commands should be added,
		// and the group will be automatically created if needed.
		bind := makeBind(sliver, con)

		// [ Core ]
		bind(consts.SliverCoreHelpGroup,
			ai.Commands,
			docscmd.Commands,
			reconfig.Commands,
			// sessions.Commands,
			sessions.SliverCommands,
			kill.Commands,
			// use.Commands,
			tasks.Commands,
			pivots.Commands,
			aka.ImplantCommands,
		)

		// [ Info ]
		bind(consts.InfoHelpGroup,
			// info.Commands,
			info.SliverCommands,
			screenshot.Commands,
			environment.Commands,
			registry.Commands,
			extensions.SliverCommands,
		)

		// [ Filesystem ]
		bind(consts.FilesystemHelpGroup,
			edit.Commands,
			hexedit.Commands,
			filesystem.Commands,
		)

		// [ Network tools ]
		bind(consts.NetworkHelpGroup,
			network.Commands,
			rportfwd.Commands,
			portfwd.Commands,
			socks.Commands,
			wireguard.SliverCommands,
		)

		// [ Execution ]
		bind(consts.ExecutionHelpGroup,
			shell.Commands,
			exec.Commands,
			backdoor.Commands,
			dllhijack.Commands,
			cursed.Commands,
			wasm.Commands,
		)

		// [ Privileges ]
		bind(consts.PrivilegesHelpGroup,
			privilege.Commands,
		)

		// [ Processes ]
		bind(consts.ProcessHelpGroup,
			processes.Commands,
		)

		// [ Aliases ]
		bind(consts.AliasHelpGroup)

		// [ Extensions ]
		bind(consts.ExtensionHelpGroup)

		// [ Post-command declaration setup ]----------------------------------------

		// Load Aliases
		aliasManifests := assets.GetInstalledAliasManifests()
		for _, manifest := range aliasManifests {
			_, err := alias.LoadAlias(manifest, sliver, con)
			if err != nil {
				con.PrintErrorf("Failed to load alias: %s", err)
				continue
			}
		}

		// Load Extensions
		extensionManifests := extensions.GetAllExtensionManifests()
		for _, manifest := range extensionManifests {
			mext, err := extensions.LoadExtensionManifest(manifest)
			// Absorb error in case there's no extensions manifest
			if err != nil {
				//con doesn't appear to be initialised here?
				//con.PrintErrorf("Failed to load extension: %s", err)
				fmt.Printf("Failed to load extension: %s\n", err)
				continue
			}

			for _, ext := range mext.ExtCommand {
				extensions.ExtensionRegisterCommand(ext, sliver, con)
			}
		}

		// [ Post-command declaration setup ]----------------------------------------

		// Everything below this line should preferably not be any command binding
		// (although you can do so without fear). If there are any final modifications
		// to make to the server menu command tree, it time to do them here.

		sliver.InitDefaultHelpCmd()
		sliver.SetHelpCommandGroupID(consts.SliverCoreHelpGroup)

		// Compute which commands should be available based on the current session/beacon.
		con.ExposeCommands()

		return sliver
	}

	return sliverCommands
}

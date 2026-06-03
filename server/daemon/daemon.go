package daemon

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
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/leonidas-c2/leonidas/server/configs"
	"github.com/leonidas-c2/leonidas/server/log"
	"github.com/leonidas-c2/leonidas/server/transport"
)

var (
	serverConfig = configs.GetServerConfig()
	daemonLog    = log.NamedLogger("daemon", "main")

	// BlankHost is a blank hostname
	BlankHost = "-"
	// BlankPort is a blank port number
	BlankPort = uint16(0)
)

// Start - Start as daemon process
func Start(host string, port uint16, tailscale bool, enableWG bool) {
	var (
		ln  net.Listener
		err error
	)
	// cli args take president over config
	if host == BlankHost {
		daemonLog.Info("No cli lhost, using config file or default value")
		host = serverConfig.DaemonConfig.Host
	}
	if port == BlankPort {
		daemonLog.Info("No cli lport, using config file or default value")
		port = uint16(serverConfig.DaemonConfig.Port)
	}

	daemonLog.Infof("Starting Leonidas daemon %s:%d ...", host, port)
	if tailscale {
		_, ln, err = transport.StartTsNetClientListener(host, port)
	} else if enableWG {
		_, ln, err = transport.StartWGWrappedMtlsClientListener(host, port)
	} else {
		_, ln, err = transport.StartMtlsClientListener(host, port)
	}
	if err != nil {
		fmt.Printf("[!] Failed to start daemon %s\n", err)
		fmt.Printf("[*] If you previously run the multiplayer command, that automatically starts a listener which might conflict with the daemon execution (default port 31337)\n")
		fmt.Printf("[*] If you want to use the daemon mode kill the multiplayer job and try to start the daemon again.\n")
		daemonLog.Errorf("Error starting client listener %s", err)
		os.Exit(1)
	}

	grpcBackend := fmt.Sprintf("127.0.0.1:%d", port)
	if _, err := transport.StartWebUI("127.0.0.1", transport.DefaultWebUIPort, grpcBackend); err != nil {
		daemonLog.Warnf("Web UI (HTTPS) failed to start — Vite/dashboard may show connection errors: %v", err)
	} else {
		daemonLog.Infof("Web UI listening on https://127.0.0.1:%d (use operator Bearer token for API)", transport.DefaultWebUIPort)
	}

	done := make(chan bool)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	go func() {
		<-signals
		daemonLog.Infof("Received SIGTERM, exiting ...")
		ln.Close()
		done <- true
	}()
	<-done
}

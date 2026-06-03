package c2

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

// ICMP C2 Transport
//
// Protocol design:
//   - Implants send ICMP Echo Request packets to the server.
//   - Each packet carries a Session ID encoded in the ICMP Identifier field
//     (16-bit), and a sequence number in the Sequence field.
//   - The C2 payload (a serialised leonidaspb.Envelope) is split across the
//     ICMP data field if needed; each fragment is preceded by a 4-byte
//     big-endian total-length header so the receiver knows when to reassemble.
//   - The server replies with an ICMP Echo Reply carrying outbound task data
//     destined for that Session ID.
//
// Limitations:
//   - Requires CAP_NET_RAW (root / administrator) on the server host.
//   - One raw socket is bound per listener; all sessions share it.
//   - Designed for low-bandwidth covert comms; not suitable for large payloads.

import (
	"encoding/binary"
	"fmt"
	"net"
	"sync"

	consts "github.com/leonidas-c2/leonidas/client/constants"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"github.com/leonidas-c2/leonidas/server/core"
	serverHandlers "github.com/leonidas-c2/leonidas/server/handlers"
	"github.com/leonidas-c2/leonidas/server/log"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"google.golang.org/protobuf/proto"
)

const (
	icmpProto        = 1   // ICMP protocol number
	icmpMaxDataBytes = 508 // conservative MTU-safe data payload per packet
	icmpHeaderBytes  = 8   // ICMP header size (type+code+checksum+id+seq)
)

var (
	icmpLog = log.NamedLogger("c2", consts.ICMPStr)
)

// icmpSession tracks an in-progress receive buffer for a single implant session.
type icmpSession struct {
	conn       *core.ImplantConnection
	recvBuf    []byte
	totalLen   uint32
	rxMu       sync.Mutex
}

// StartICMPListener opens a raw ICMP socket on host and starts accepting
// implant connections encoded in ICMP Echo Request data payloads.
// Returns a stop channel — send true to terminate the listener.
func StartICMPListener(host string) (chan bool, error) {
	listenAddr := host
	if listenAddr == "" {
		listenAddr = "0.0.0.0"
	}

	conn, err := icmp.ListenPacket("ip4:icmp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("[icmp] listen %s: %w", listenAddr, err)
	}
	icmpLog.Infof("ICMP listener started on %s", listenAddr)

	sessions := &sync.Map{} // map[uint16]*icmpSession keyed by ICMP identifier

	stopCh := make(chan bool, 1)

	go func() {
		<-stopCh
		icmpLog.Infof("ICMP listener stopping")
		conn.Close()
	}()

	go icmpReadLoop(conn, sessions)

	return stopCh, nil
}

// icmpReadLoop reads packets forever and dispatches them by session ID.
func icmpReadLoop(conn *icmp.PacketConn, sessions *sync.Map) {
	buf := make([]byte, 1500)
	for {
		n, peer, err := conn.ReadFrom(buf)
		if err != nil {
			// socket closed — normal shutdown
			return
		}
		msg, err := icmp.ParseMessage(icmpProto, buf[:n])
		if err != nil {
			continue
		}
		if msg.Type != ipv4.ICMPTypeEcho {
			continue
		}
		echo, ok := msg.Body.(*icmp.Echo)
		if !ok || len(echo.Data) < 4 {
			continue
		}

		sessionID := uint16(echo.ID)
		data := echo.Data

		// First 4 bytes of data payload encode the total envelope length.
		totalLen := binary.BigEndian.Uint32(data[:4])
		fragment := data[4:]

		sessRaw, _ := sessions.LoadOrStore(sessionID, &icmpSession{
			conn: core.NewImplantConnection(consts.ICMPStr, peer.String()),
		})
		sess := sessRaw.(*icmpSession)

		sess.rxMu.Lock()
		if totalLen > 0 && sess.totalLen == 0 {
			sess.totalLen = totalLen
			sess.recvBuf = make([]byte, 0, totalLen)
		}
		sess.recvBuf = append(sess.recvBuf, fragment...)

		if uint32(len(sess.recvBuf)) >= sess.totalLen && sess.totalLen > 0 {
			raw := sess.recvBuf[:sess.totalLen]
			sess.recvBuf = nil
			sess.totalLen = 0
			sess.rxMu.Unlock()

			envelope := &leonidaspb.Envelope{}
			if err := proto.Unmarshal(raw, envelope); err != nil {
				icmpLog.Warnf("[icmp] failed to unmarshal envelope from session %d: %v", sessionID, err)
				continue
			}

			// Dispatch to the same handler map as other transports (payload = envelope.Data).
			handlers := serverHandlers.GetHandlers()
			if handler, ok := handlers[envelope.Type]; ok {
				go func(implantConn *core.ImplantConnection, env *leonidaspb.Envelope) {
					defer recoverAndLogPanic(icmpLog.Errorf, "icmp message handler")
					respEnvelope := handler(implantConn, env.Data)
					if respEnvelope != nil {
						implantConn.Send <- respEnvelope
					}
				}(sess.conn, envelope)
			}

			// Drain any pending outbound envelopes and send replies.
			go icmpDrainOutbound(conn, peer, echo.ID, sess.conn)
		} else {
			sess.rxMu.Unlock()
		}
	}
}

// icmpDrainOutbound sends all queued outbound envelopes back to the implant
// as ICMP Echo Reply packets.
func icmpDrainOutbound(conn *icmp.PacketConn, peer net.Addr, id int, implantConn *core.ImplantConnection) {
	for {
		select {
		case envelope, ok := <-implantConn.Send:
			if !ok {
				return
			}
			raw, err := proto.Marshal(envelope)
			if err != nil {
				icmpLog.Warnf("[icmp] marshal error: %v", err)
				continue
			}
			if err := icmpSendFragmented(conn, peer, id, raw); err != nil {
				icmpLog.Warnf("[icmp] send error: %v", err)
				return
			}
		default:
			return
		}
	}
}

// icmpSendFragmented sends payload as one or more ICMP Echo Reply packets.
// The first fragment includes a 4-byte big-endian total length prefix.
func icmpSendFragmented(conn *icmp.PacketConn, peer net.Addr, id int, payload []byte) error {
	totalLen := uint32(len(payload))
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, totalLen)

	seq := 0
	offset := 0
	firstPkt := true
	for offset < len(payload) {
		var chunk []byte
		if firstPkt {
			end := offset + icmpMaxDataBytes - 4
			if end > len(payload) {
				end = len(payload)
			}
			chunk = append(header, payload[offset:end]...)
			offset = end
			firstPkt = false
		} else {
			end := offset + icmpMaxDataBytes
			if end > len(payload) {
				end = len(payload)
			}
			chunk = payload[offset:end]
			offset = end
		}

		msg := &icmp.Message{
			Type: ipv4.ICMPTypeEchoReply,
			Code: 0,
			Body: &icmp.Echo{
				ID:   id,
				Seq:  seq,
				Data: chunk,
			},
		}
		raw, err := msg.Marshal(nil)
		if err != nil {
			return err
		}
		if _, err := conn.WriteTo(raw, peer); err != nil {
			return err
		}
		seq++
	}
	return nil
}

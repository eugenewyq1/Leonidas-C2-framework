//go:build !darwin && !linux && !freebsd && !openbsd && !dragonfly

package tunnel_handlers

import (
	"github.com/leonidas-c2/leonidas/implant/leonidas/transports"
	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	"google.golang.org/protobuf/proto"
)

func ShellResizeReqHandler(envelope *leonidaspb.Envelope, connection *transports.Connection) {
	resp, _ := proto.Marshal(&commonpb.Empty{})
	connection.Send <- &leonidaspb.Envelope{
		ID:   envelope.ID,
		Data: resp,
	}
}

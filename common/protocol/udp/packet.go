package udp

import (
	"github.com/qazz-shyper/website/common/buf"
	"github.com/qazz-shyper/website/common/net"
)

// Packet is a UDP packet together with its source and destination address.
type Packet struct {
	Payload *buf.Buffer
	Source  net.Destination
	Target  net.Destination
}

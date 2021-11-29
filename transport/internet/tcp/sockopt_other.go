//go:build !linux && !freebsd
// +build !linux,!freebsd

package tcp

import (
	"github.com/qazz-shyper/website/common/net"
	"github.com/qazz-shyper/website/transport/internet/stat"
)

func GetOriginalDestination(conn stat.Connection) (net.Destination, error) {
	return net.Destination{}, nil
}

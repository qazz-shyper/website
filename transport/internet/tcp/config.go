package tcp

import (
	"github.com/qazz-shyper/website/common"
	"github.com/qazz-shyper/website/transport/internet"
)

const protocolName = "tcp"

func init() {
	common.Must(internet.RegisterProtocolConfigCreator(protocolName, func() interface{} {
		return new(Config)
	}))
}

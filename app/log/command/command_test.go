package command_test

import (
	"context"
	"testing"

	"github.com/qazz-shyper/website/app/dispatcher"
	"github.com/qazz-shyper/website/app/log"
	. "github.com/qazz-shyper/website/app/log/command"
	"github.com/qazz-shyper/website/app/proxyman"
	_ "github.com/qazz-shyper/website/app/proxyman/inbound"
	_ "github.com/qazz-shyper/website/app/proxyman/outbound"
	"github.com/qazz-shyper/website/common"
	"github.com/qazz-shyper/website/common/serial"
	"github.com/qazz-shyper/website/core"
)

func TestLoggerRestart(t *testing.T) {
	v, err := core.New(&core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&log.Config{}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.InboundConfig{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	})
	common.Must(err)
	common.Must(v.Start())

	server := &LoggerServer{
		V: v,
	}
	common.Must2(server.RestartLogger(context.Background(), &RestartLoggerRequest{}))
}

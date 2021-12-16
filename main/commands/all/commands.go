package all

import (
	"github.com/qazz-shyper/website/main/commands/all/api"
	"github.com/qazz-shyper/website/main/commands/all/tls"
	"github.com/qazz-shyper/website/main/commands/base"
)

// go:generate go run github.com/qazz-shyper/website/common/errors/errorgen

func init() {
	base.RootCommand.Commands = append(
		base.RootCommand.Commands,
		api.CmdAPI,
		// cmdConvert,
		tls.CmdTLS,
		cmdUUID,
	)
}

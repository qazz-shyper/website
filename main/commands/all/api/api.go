package api

import (
	"github.com/qazz-shyper/website/main/commands/base"
)

// CmdAPI calls an API in an Website process
var CmdAPI = &base.Command{
	UsageLine: "{{.Exec}} api",
	Short:     "Call an API in an Website process",
	Long: `{{.Exec}} {{.LongName}} provides tools to manipulate Website via its API.
`,
	Commands: []*base.Command{
		cmdRestartLogger,
		cmdGetStats,
		cmdQueryStats,
		cmdSysStats,
		cmdAddInbounds,
		cmdAddOutbounds,
		cmdRemoveInbounds,
		cmdRemoveOutbounds,
	},
}

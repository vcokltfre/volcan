package meta

import (
	"github.com/vcokltfre/volcan/src/bot/core"
)

var statusCommand = &core.Command{
	Name: "status",
	Aliases: []string{
		"ping",
	},
	Description: "Get the status of the bot.",
	Usage:       "status",
	Handler: func(ctx *core.Context) error {
		return nil
	},
}

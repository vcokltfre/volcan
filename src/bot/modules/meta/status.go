package meta

import (
	"fmt"

	"github.com/vcokltfre/volcan/src/api/modules/meta"
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
		status := &meta.Status{}
		err := ctx.Request("GET", "/status", nil, status)
		if err != nil {
			ctx.Error(fmt.Errorf("an error occurred while fetching the bot's status: %v", err))
			return nil
		}

		_, err = ctx.Reply("Status from internal API: %s", status.Status)
		if err != nil {
			return err
		}

		return nil
	},
}

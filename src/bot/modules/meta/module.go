package meta

import "github.com/vcokltfre/volcan/src/bot/core"

var MetaModule = &core.Module{
	Name:        "meta",
	Description: "Meta commands for the bot.",
	Commands: []*core.Command{
		statusCommand,
	},
}

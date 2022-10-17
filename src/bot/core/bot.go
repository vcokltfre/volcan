package core

import (
	"context"
	"fmt"
	"os"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/vcokltfre/volcan/src/utils"
)

type Bot struct {
	Modules []*Module
	Client  bot.Client

	modules map[string]*Module
}

func (b *Bot) Build() error {
	commandNames := []string{}

	for _, module := range b.modules {
		if _, ok := b.modules[module.Name]; ok {
			return fmt.Errorf("module %s already exists", module.Name)
		}

		cmds, err := module.Build()
		if err != nil {
			return err
		}

		for _, cmd := range cmds {
			if utils.Contains(commandNames, cmd) {
				return fmt.Errorf("command %s already exists", cmd)
			}

			commandNames = append(commandNames, cmd)
		}

		b.modules[module.Name] = module
	}

	return nil
}

func (b *Bot) Start() error {
	client, err := disgo.New(os.Getenv("BOT_TOKEN"), bot.WithGatewayConfigOpts(
		gateway.WithIntents(gateway.IntentsAll),
	))
	if err != nil {
		return err
	}

	b.Client = client

	return client.OpenGateway(context.Background())
}

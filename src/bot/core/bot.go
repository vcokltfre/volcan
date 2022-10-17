package core

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/sirupsen/logrus"
	"github.com/vcokltfre/glex"
	"github.com/vcokltfre/volcan/src/config"
	"github.com/vcokltfre/volcan/src/utils"
)

var BotInstance *Bot

type Bot struct {
	Modules      []*Module
	Client       bot.Client
	CommandCount int

	modules  map[string]*Module
	commands map[string]*Command
}

func (b *Bot) Build() error {
	commandNames := []string{}
	commandCount := 0

	b.modules = map[string]*Module{}
	b.commands = map[string]*Command{}

	for _, module := range b.Modules {
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
		commandCount += len(module.Commands)

		for _, command := range module.Commands {
			b.commands[command.Name] = command

			for _, alias := range command.Aliases {
				b.commands[alias] = command
			}
		}

		logrus.Info("Registered module ", module.Name)
	}

	b.CommandCount = commandCount

	BotInstance = b

	return nil
}

func (b *Bot) Start() error {
	client, err := disgo.New(
		os.Getenv("BOT_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentsAll),
		),
		bot.WithEventListenerFunc(b.handleMessage),
	)
	if err != nil {
		return err
	}

	b.Client = client

	return client.OpenGateway(context.Background())
}

func (b *Bot) handleMessage(msg *events.MessageCreate) {
	prefix := getPrefix(msg)
	if len(prefix) == 0 {
		return
	}

	if !strings.HasPrefix(msg.Message.Content, prefix) {
		return
	}

	commandString := strings.TrimPrefix(msg.Message.Content, prefix)
	commandParts := strings.Split(commandString, " ")
	command, num := b.findCommand(commandParts)

	if command == nil {
		return
	}

	ctx := &Context{
		Bot:     b,
		Message: &msg.Message,
		Event:   msg,
		args:    map[string]string{},
		flags:   map[string]string{},
		bools:   map[string]bool{},
		varArgs: []string{},
	}

	args, err := glex.SplitCommand(msg.Message.Content)
	if err != nil {
		ctx.Error(fmt.Errorf("parsing failed: %v", err))
		return
	}

	err = command.Run(ctx, args[num:])
	if err != nil {
		ctx.Error(err)
		return
	}
}

func (b *Bot) findCommand(parts []string) (*Command, int) {
	if len(parts) == 0 {
		return nil, 0
	}

	command, ok := b.commands[parts[0]]
	if !ok {
		return nil, 0
	}

	return command.Find(parts[1:], 1)
}

func getPrefix(msg *events.MessageCreate) string {
	guild, ok := msg.Guild()
	if !ok {
		return config.Config.Prefixes.Default
	}

	prefix, ok := config.Config.Prefixes.Servers[guild.ID.String()]
	if !ok {
		return config.Config.Prefixes.Default
	}

	return prefix
}

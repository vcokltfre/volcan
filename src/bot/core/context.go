package core

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/vcokltfre/volcan/src/config"
)

type Context struct {
	Bot     *Bot
	Event   *events.MessageCreate
	Message *discord.Message

	args    map[string]string
	flags   map[string]string
	bools   map[string]bool
	varArgs []string
}

func (c *Context) GetLevel() int {
	topLevel := 0

	userLevel := config.Config.Levels[c.Message.Author.ID.String()]
	if userLevel > topLevel {
		topLevel = userLevel
	}

	guild, ok := c.Event.Guild()
	if !ok {
		return topLevel
	}

	guildLevel := config.Config.Levels[guild.ID.String()]
	if guildLevel > topLevel {
		topLevel = guildLevel
	}

	if c.Event.Message.Member == nil {
		return topLevel
	}

	for _, role := range c.Event.Message.Member.RoleIDs {
		roleLevel := config.Config.Levels[role.String()]
		if roleLevel > topLevel {
			topLevel = roleLevel
		}
	}

	return topLevel
}

func (c *Context) Error(err error) {
	c.Bot.Client.Rest().CreateMessage(
		c.Event.ChannelID,
		discord.NewMessageCreateBuilder().SetContent("An error occurred: "+err.Error()).SetAllowedMentions(&discord.AllowedMentions{
			Parse: []discord.AllowedMentionType{},
		}).Build(),
	)
}

func (c *Context) Reply(message string, format ...interface{}) (*discord.Message, error) {
	return c.Bot.Client.Rest().CreateMessage(
		c.Event.ChannelID,
		discord.NewMessageCreateBuilder().SetContent(fmt.Sprintf(message, format...)).SetAllowedMentions(&discord.AllowedMentions{
			Parse: []discord.AllowedMentionType{},
		}).Build(),
	)
}

func (c *Context) Arg(name string) string {
	return c.args[name]
}

func (c *Context) Flag(name string) string {
	return c.flags[name]
}

func (c *Context) Bool(name string) bool {
	return c.bools[name]
}

func (c *Context) Request(method string, path string, body any, out any, args ...any) error {
	return c.Bot.Request(method, path, body, out, args...)
}

package core

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/vcokltfre/volcan/src/config"
)

type Context struct {
	Bot     *Bot
	Event   *events.MessageCreate
	Message *discord.Message
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

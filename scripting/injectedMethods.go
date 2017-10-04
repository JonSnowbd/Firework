package scripting

import (
	"github.com/bwmarrin/discordgo"
	lua "github.com/yuin/gopher-lua"
)

func lua_say(Client *discordgo.Session, Message *discordgo.MessageCreate) func(*lua.LState) int {
	return func(state *lua.LState) int {
		msg := state.ToString(1)
		Client.ChannelMessageSend(Message.ChannelID, msg)
		return 0
	}
}

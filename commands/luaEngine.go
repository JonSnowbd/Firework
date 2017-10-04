package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yuin/gopher-lua"
)

func lua_say(Client *discordgo.Session, Message *discordgo.MessageCreate) func(*lua.LState) int {
	return func(state *lua.LState) int {
		msg := state.ToString(1)
		Client.ChannelMessageSend(Message.ChannelID, msg)
		return 0
	}
}

func MakeEngine(Client *discordgo.Session, Message *discordgo.MessageCreate) *lua.LState {
	L := lua.NewState()

	L.SetGlobal("say", L.NewFunction(lua_say(Client, Message)))
	return L
}

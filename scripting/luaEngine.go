package scripting

import (
	"github.com/bwmarrin/discordgo"
	"github.com/yuin/gopher-lua"
)

// MakeEngine returns a ready to use engine hooked up with all the required globals.
func MakeEngine(Client *discordgo.Session, Message *discordgo.MessageCreate) *lua.LState {
	L := lua.NewState()

	L.SetGlobal("say", L.NewFunction(lua_say(Client, Message)))
	return L
}

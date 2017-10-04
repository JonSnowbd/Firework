package bot

import (
	"github.com/bwmarrin/discordgo"
)

// Command is an interface representing everything the bot needs to know about running a command.
type Command interface {
	// Initializes the command.
	Init()
	// Determines if the provided token should trigger this
	Match(token string, isUser bool) bool
	// Performs the action
	Run(Client *discordgo.Session, Message *discordgo.MessageCreate)
}

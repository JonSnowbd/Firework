package x // Replace.

import (
	"github.com/bwmarrin/discordgo"
	"github.com/JonSnowbd/Firework/util"
)

// SimpleCommand is a barebones, do nothing command for people to copy
//     in order to build their own commands without needing some magical
//     terminal based tool.
type SimpleCommand struct {
}

// Called when this command is placed into a bot. Make sure to do anything you need
// here, before the bot uses this.
func (command SimpleCommand) Init() {

}

// Match returns true if provided token matches this command's identifier.
func (command SimpleCommand) Match(token string, isUser bool) bool {
	return util.SimplePublicCommand("ping", token, isUser)
}

// Run performs the command's logic.
func (command SimpleCommand) Run(Client *discordgo.Session, Message *discordgo.MessageCreate) {
}

package commands // Replace.

import (
	"database/sql"
	"strings"

	"github.com/JonSnowbd/Firework/bot"
	"github.com/JonSnowbd/Firework/util"
	"github.com/bwmarrin/discordgo"
)

// ScriptCommand is a barebones, do nothing command for people to copy
//     in order to build their own commands without needing some magical
//     terminal based tool.
type ScriptCommand struct {
	pref     string
	Database *sql.DB
}

// Called when this command is placed into a bot. Make sure to do anything you need
// here, before the bot uses this.
func (command ScriptCommand) Init(bot bot.State) {
	command.pref = bot.Prefix
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS scripts (
		id INTEGER NOT NULL PRIMARY KEY,
		author TEXT NOT NULL,
		date DATETIME NOT NULL,
		name TEXT NOT NULL,
		content TEXT NOT NULL,
		uses INTEGER NOT NULL
	)
	`
	_, err := command.Database.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
}

// Match returns true if provided token matches this command's identifier.
func (command ScriptCommand) Match(token string, isUser bool) bool {
	return util.SimplePrivateCommand("script", token, isUser)
}

// Run performs the command's logic.
func (command ScriptCommand) Run(Client *discordgo.Session, Message *discordgo.MessageCreate) error {

	log := bot.GetLogger()

	defer func() {
		if r := recover(); r != nil {
			log.Error("Something went horribly wrong with a Script command:", r)
			Client.ChannelMessageSend(Message.ChannelID, "Error parsing command, are you sure you didn't forget a subcommand and target?")
		}
	}()

	clean := strings.TrimPrefix(Message.Content, command.pref)
	args := strings.Split(clean, " ")
	if len(args) < 2 {
		Client.ChannelMessageSend(Message.ChannelID, "Youre going to need to give me more than that.")
	}
	argSubcommand := args[1]

	// When creating a tag
	if argSubcommand == "add" || argSubcommand == "+" || argSubcommand == "create" {
		err := command.handleAdding(Client, Message, args)
		if err != nil {
			return err
		}
		log.Info(Message.Author.Username, " created command: ", args[2])
		return nil
	}

	// When deleting a tag
	if argSubcommand == "-" || argSubcommand == "delete" || argSubcommand == "delet" || argSubcommand == "remove" {
		err := command.handleRemoving(Client, Message, args)
		if err != nil {
			return err
		}
		log.Info(Message.Author.Username, " deleted command: ", args[2])
		return nil
	}

	// When asking to view the raw data of any tag
	if argSubcommand == "view" || argSubcommand == "raw" || argSubcommand == " " {
		err := command.handleViewing(Client, Message, args)
		if err != nil {
			return err
		}
		return nil
	}

	err := command.handleExecution(Client, Message, args)
	if err != nil {
		return err
	}
	return nil
}

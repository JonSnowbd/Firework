package commands // Replace.

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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
	return util.SimplePublicCommand("script", token, isUser)
}

// Run performs the command's logic.
func (command ScriptCommand) Run(Client *discordgo.Session, Message *discordgo.MessageCreate) {

	defer func() {
		if r := recover(); r != nil {
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
		argName := args[2]
		argScript := strings.Join(args[3:], " ")
		argActualContent := strings.TrimPrefix(argScript, "```lua")
		argActualContent = strings.TrimSuffix(argActualContent, "```")
		argActualContent = strings.Trim(argActualContent, "\n")

		err := command.AddTag(argName, argActualContent, Message.Author.Username)
		if err != nil {
			Client.ChannelMessageSend(Message.ChannelID, "There was an error creating that script. **See your latest log file.**")
			fmt.Println(err)
			return
		}

		Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("Succesfully created ", "**"+argName+"**"))
		return
	}

	if argSubcommand == "-" || argSubcommand == "delete" || argSubcommand == "delet" || argSubcommand == "remove" {
		target := args[2]
		err := command.DeleteTag(target)
		if err != nil {
			fmt.Println(err)
			Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("Failed to delete ", "**"+target+"**"))
			return
		}
		Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("Successfully deleted ", "**"+target+"**"))
		return
	}

	if argSubcommand == "view" || argSubcommand == "raw" {
		target := args[2]
		s, err := command.GetRawScript(target)
		if err != nil {
			fmt.Println(err)
			Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("**"+target+"**", " does not exist."))
			return
		}
		Client.ChannelMessageSend(Message.ChannelID, "```lua\n"+s+"\n```")
		return
	}

	machine := MakeEngine(Client, Message)
	defer machine.Close()
	code, err := command.GetRawScript(argSubcommand)
	if err != nil {
		Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("**"+argSubcommand+"**", " does not exist."))
	}
	err = machine.DoString(code)
	if err != nil {
		Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("Error!\n```", err, "```"))
	}
}

// AddTag takes a name and the contents of a script and stores it in a database.
func (command ScriptCommand) AddTag(name string, script string, author string) error {
	transaction, err := command.Database.Begin()
	if err != nil {
		return err
	}

	statement, err := transaction.Prepare(`
		INSERT INTO scripts(name, content, author, date, uses) VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(name, script, author, time.Now(), 0)
	if err != nil {
		return err
	}

	transaction.Commit()
	return nil
}

func (command ScriptCommand) DeleteTag(name string) error {
	statement, err := command.Database.Prepare("DELETE FROM scripts WHERE name = ?")
	if err != nil {
		return err
	}

	_, err = statement.Exec(name)
	if err != nil {
		return err
	}

	return nil
}

func (command ScriptCommand) GetRawScript(name string) (string, error) {
	statement, err := command.Database.Prepare("SELECT content FROM scripts WHERE name = ?")
	if err != nil {
		return "", err
	}
	content := ""
	err = statement.QueryRow(name).Scan(&content)
	if err != nil {
		return "", err
	}

	return content, nil
}

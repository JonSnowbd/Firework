package commands

import (
	"fmt"
	"strings"

	"github.com/JonSnowbd/Firework/scripting"
	"github.com/bwmarrin/discordgo"
)

func (command ScriptCommand) handleAdding(Client *discordgo.Session, Message *discordgo.MessageCreate, args []string) error {
	// Detect the tag data
	argName := args[2]
	argScript := strings.Join(args[3:], " ")
	argActualContent := strings.TrimPrefix(argScript, "```lua")
	argActualContent = strings.TrimSuffix(argActualContent, "```")
	argActualContent = strings.Trim(argActualContent, "\n")

	// Then add it to the database.
	err := scripting.AddTag(argName, argActualContent, Message.Author.Username, command.Database)
	if err != nil {
		Client.ChannelMessageSend(Message.ChannelID, "There was an error creating that script. **See your latest log file.**")
		fmt.Println(err)
		return err
	}

	Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("Succesfully created ", "**"+argName+"**"))
	return nil
}

func (command ScriptCommand) handleRemoving(Client *discordgo.Session, Message *discordgo.MessageCreate, args []string) error {
	target := args[2]
	err := scripting.DeleteTag(target, command.Database)
	if err != nil {
		fmt.Println(err)
		Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("Failed to delete ", "**"+target+"**"))
		return err
	}
	Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("Successfully deleted ", "**"+target+"**"))
	return nil
}

func (command ScriptCommand) handleViewing(Client *discordgo.Session, Message *discordgo.MessageCreate, args []string) error {
	target := args[2]
	s, err := scripting.GetRawScript(target, command.Database)
	if err != nil {
		fmt.Println(err)
		Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("**"+target+"**", " does not exist."))
		return err
	}
	Client.ChannelMessageSend(Message.ChannelID, "```lua\n"+s+"\n```")
	return nil
}

func (command ScriptCommand) handleExecution(Client *discordgo.Session, Message *discordgo.MessageCreate, args []string) error {
	target := args[1]

	// Create the scripting engine and defer its closing
	machine := scripting.MakeEngine(Client, Message)
	defer machine.Close()

	// First find it and get its script, making sure it exists.
	code, err := scripting.GetRawScript(target, command.Database)
	if err != nil {
		Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("**"+target+"**", " does not exist."))
		return err
	}
	err = machine.DoString(code)
	if err != nil {
		Client.ChannelMessageSend(Message.ChannelID, fmt.Sprint("Error!\n```", err, "```"))
		return err
	}

	return nil
}

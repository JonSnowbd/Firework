package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// State is how the bot is currently performing.
type State struct {
	// Running declares whether or not the bot is responding to input.
	Running bool
	// Prefix is the string that determines how commands will be triggered.
	Prefix string
	// Selfbot is whether or not this bot will connect as a user or as its own bot.
	Selfbot bool

	client   *discordgo.Session
	user     *discordgo.User
	commands []Command
}

// Start begins the bot and maybe returns an error. This does not block the program,
// So make sure to do so youreself when using this library.
func (b *State) Start(token string) error {
	if !b.Selfbot {
		token = "Bot " + token
	}
	// Connect to discord.
	dclient, err := discordgo.New(token)
	if err != nil {
		return err
	}

	dclient.AddHandler(b.messageHandler)

	// Establish discord Socket connection
	err = dclient.Open()
	if err != nil {
		return err
	}

	// Get user details.
	duser, err := dclient.User("@me")
	if err != nil {
		return err
	}

	b.client = dclient
	b.user = duser

	log.Info("Logged on as", duser.Username)

	b.Running = true

	return nil
}

func (b *State) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// If the bot is not running, discard anything.
	if !b.Running {
		return
	}
	// If its not a self bot, then discard own messages.
	if !b.Selfbot && m.Author.ID == b.user.ID {
		return
	}
	// If message does not start with the prefix, discard.
	if !strings.HasPrefix(m.Content, b.Prefix) {
		return
	}

	argCommand := strings.TrimPrefix(m.Content, b.Prefix)
	argCommand = strings.Split(argCommand, " ")[0]

	for _, cmd := range b.commands {
		if cmd.Match(argCommand, m.Author.ID == b.user.ID) {
			log.Info("Command ran:", argCommand)
			err := cmd.Run(s, m)
			if err != nil {
				log.Error(err)
			}
			return
		}
	}
	log.Info("Command doesn't exist:", argCommand)
}

// GetDefaultState returns a state that is set up to just work.
func GetDefaultState() State {
	return State{
		Selfbot:  false,
		Running:  false,
		commands: []Command{},
		Prefix:   "!",
	}
}

// AddCommand adds a command to the current bot.
func (b *State) AddCommand(command Command) {
	command.Init(*b)
	b.commands = append(b.commands, command)
}

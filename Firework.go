package main

import (
	"database/sql"
	"fmt"

	"github.com/JonSnowbd/Firework/bot"
	"github.com/JonSnowbd/Firework/commands"
	"github.com/JonSnowbd/Firework/icon"
	_ "github.com/mattn/go-sqlite3"

	"os"

	"github.com/getlantern/systray"
)

var (
	botState bot.State
)

func main() {
	botState = bot.GetDefaultState()
	botState.Selfbot = true
	botState.Prefix = "!!"

	db, err := sql.Open("sqlite3", "./Firework_data.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// Add functionality like so
	botState.AddCommand(commands.ScriptCommand{Database: db})
	// making sure the struct you supply
	// inherits everything from the command interface.
	// See _examples/blank_command.go for a barebones command.

	err = botState.Start(os.Getenv("Firework_token"))
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	go systray.Run(onReady, onExit)

	<-make(chan struct{}) // Simple blocking.
}

// Required as of the latest systray updated, nothing to do in it yet though.
func onExit() {}

// Runs the system tray icon.
func onReady() {
	systray.SetIcon(icon.DataOn)
	systray.SetTitle("Bot")
	systray.SetTooltip("Firework")

	mToggle := systray.AddMenuItem("Toggle", "Toggle the bot off.")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		select {
		// When clicked exit, close everything
		case <-mQuit.ClickedCh:
			systray.Quit()
			os.Exit(0)
			return
		// Otherwise toggle the bot.
		case <-mToggle.ClickedCh:
			if botState.Running {
				botState.Running = false
				systray.SetIcon(icon.DataOff)
			} else {
				botState.Running = true
				systray.SetIcon(icon.DataOn)
			}
		}
	}
}

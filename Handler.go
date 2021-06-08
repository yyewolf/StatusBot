package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func defineCommands() {
	commands["status"] = status
	commands["addip"] = addIP
	commands["remip"] = remIP
}

func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Check that sender is not a bot
	if m.Author.Bot {
		return
	}
	//Check that user is willing to use this bot
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	args := strings.Split(m.Content, " ")
	//Check that user is using one of our commands
	if len(args) < 1 {
		return
	}

	typedCmd := strings.Replace(args[0], prefix, "", 1)

	if val, ok := commands[typedCmd]; ok {
		//If command has been found in map
		val(s, m)
	}
}

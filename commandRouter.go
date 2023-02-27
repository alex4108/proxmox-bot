package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func startsWith(prefix string, content string) bool {
	return (strings.Split(content, " ")[0] == prefix)
}

func commandRouter(s *discordgo.Session, m *discordgo.MessageCreate) {
	if startsWith("*startvm", m.Content) {
		go startCommand(s, m)
	} else if startsWith("*stopvm", m.Content) {
		go stopCommand(s, m)
	} else if m.Content == "*help" {
		go helpCommand(s, m)
	} else if m.Content == "*ping" {
		go pingCommand(s, m)
	}
}

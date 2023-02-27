package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// startCommand is the command handler to start an azure VM.
func startCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	logicalName, err := getLogicalNameFromCommand(m.Content)
	if err != nil {
		response := fmt.Sprintf("Failed to get VM name from command %v: %v", m.Content, err)
		respond(s, m.ChannelID, response)
		return
	}

	log.Infof("Received start command for VM: %v", logicalName)

	respond(s, m.ChannelID, "Starting VM "+logicalName+"...")
	// todo Request Start via Proxmox API
	if err != nil {
		log.Errorf("failed to obtain a response: %v", err)
		respondError(s, m.ChannelID)
		return
	}

	// todo verify VM started
	log.Infof("VM %v started.", logicalName)
	respond(s, m.ChannelID, "VM "+logicalName+" started.")
}

func getLogicalNameFromCommand(content string) (string, error) {
	if len(strings.Split(content, " ")) < 2 {
		return "", fmt.Errorf("no VM name provided")
	}

	if len(strings.Split(content, " ")) > 2 {
		return "", fmt.Errorf("too many arguments")
	}

	targetVM := strings.Split(content, " ")[1]
	return targetVM, nil
}

// stopCommand is the command handler to stop an azure VM.
func stopCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	logicalName, err := getLogicalNameFromCommand(m.Content)
	if err != nil {
		response := fmt.Sprintf("Failed to get VM name from command %v: %v", m.Content, err)
		respond(s, m.ChannelID, response)
		return
	}

	log.Infof("Received stop command for VM: %v", logicalName)

	respond(s, m.ChannelID, "Stopping VM "+logicalName+"...")
	// todo Request Start via Proxmox API
	if err != nil {
		log.Errorf("failed to obtain a response: %v", err)
		respondError(s, m.ChannelID)
		return
	}

	// todo verify VM started
	log.Infof("VM %v stopped.", logicalName)
	respond(s, m.ChannelID, "VM "+logicalName+" stopped.")
}

// pingCommand is the command handler to ping the bot.
func pingCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	now := time.Now()
	latency := ""
	if timestampFieldExists(m) {
		diff := m.Timestamp.Sub(now)
		latency = "(" + strconv.Itoa(int(diff.Milliseconds())) + " ms)"
	}
	respond(s, m.ChannelID, "Pong! "+latency)
}

// helpCommand is the command handler to show the help message.
func helpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	messageContent := `proxmox-bot, an open source Discord Bot.

Available Commands:
	*stopvm <vm_name> Stops a VM
	*startvm <vm_name> Starts a VM
	*ping
	
Proudly maintained by Alex https://github.com/alex4108/proxmox-bot
`

	respond(s, m.ChannelID, messageContent)
}

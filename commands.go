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

	ce, err := getConfigEntry(logicalName)
	if err != nil {
		log.Errorf("Failed getConfigEntry for %v: %v", logicalName, err)
		respondError(s, m.ChannelID)
		return
	}

	pxmVm, err := getVmById(ce)
	if err != nil {
		log.Errorf("failed to obtain get VM by ID %v: %v", ce.VMId, err)
		respondError(s, m.ChannelID)
		return
	}

	task, err := pxmVm.Start()
	if err != nil {
		log.Errorf("failed to start VM: %v", err)
		respondError(s, m.ChannelID)
		return
	}
	respond(s, m.ChannelID, fmt.Sprintf("VM Start Requested, Task ID: %v", task.ID))

	for task.IsRunning {
		log.Infof("Waiting for VM to start...")
		time.Sleep(1 * time.Second)
	}

	if task.IsFailed {
		log.Errorf("VM Start TaskID %v FAILED: %v", task.ID, task)
		respond(s, m.ChannelID, fmt.Sprintf("VM Start Task ID %v FAILED!, %v", task.ID, task))
		return
	} else {
		log.Infof("VM Start TaskID %v SUCCEEDED: %v", task.ID, task)
		respond(s, m.ChannelID, fmt.Sprintf("VM %v started!", ce.LogicalName))
		return
	}
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

	ce, err := getConfigEntry(logicalName)
	if err != nil {
		log.Errorf("Failed getConfigEntry for %v: %v", logicalName, err)
		respondError(s, m.ChannelID)
		return
	}

	pxmVm, err := getVmById(ce)
	if err != nil {
		log.Errorf("failed to obtain get VM by ID %v: %v", ce.VMId, err)
		respondError(s, m.ChannelID)
		return
	}

	task, err := pxmVm.Shutdown()
	if err != nil {
		log.Errorf("failed to stop VM: %v", err)
		respondError(s, m.ChannelID)
		return
	}
	respond(s, m.ChannelID, fmt.Sprintf("VM Stop Requested, Task ID: %v", task.ID))

	for task.IsRunning {
		log.Infof("Waiting for VM to stop...")
		time.Sleep(1 * time.Second)
	}

	if task.IsFailed {
		log.Errorf("VM Stop TaskID %v FAILED: %v", task.ID, task)
		respond(s, m.ChannelID, fmt.Sprintf("VM Stop Task ID %v FAILED!, %v", task.ID, task))
		return
	} else {
		log.Infof("VM Stop TaskID %v SUCCEEDED: %v", task.ID, task)
		respond(s, m.ChannelID, fmt.Sprintf("VM %v stopped!", ce.LogicalName))
		return
	}
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

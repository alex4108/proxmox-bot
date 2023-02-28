package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/bwmarrin/discordgo"
	"github.com/luthermonson/go-proxmox"
	log "github.com/sirupsen/logrus"
)

// respondError is a quick way to respond with an error message.
func respondError(s *discordgo.Session, channelID string) {
	respond(s, channelID, "An internal error occured.  Please raise a bug on the github repository for further investigation.")
}

func respond(s *discordgo.Session, channelID string, response string) {
	_, err := s.ChannelMessageSend(channelID, response)
	if err != nil {
		log.Errorf("Failed to respond command. %v", err)
	}
}

func timestampFieldExists(obj *discordgo.MessageCreate) bool {
	metaValue := reflect.ValueOf(obj).Elem()
	field := metaValue.FieldByName("Timestamp")
	return field != (reflect.Value{})
}

func makeProxmoxClient(url string, name string) (*proxmox.Client, error) {
	insecureHTTPClient := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client := proxmox.NewClient(fmt.Sprintf("https://%s/api2/json", url),
		proxmox.WithClient(&insecureHTTPClient),
	)

	username, exist := os.LookupEnv(fmt.Sprintf("%s_USERNAME", name))
	if !exist {
		errMsg := fmt.Sprintf("Failed to find username variable for %v", name)
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}

	password, exist := os.LookupEnv(fmt.Sprintf("%s_PASSWORD", name))
	if !exist {
		errMsg := fmt.Sprintf("Failed to find password variable for %v", name)
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}

	if err := client.Login(username, password); err != nil {
		errMsg := fmt.Sprintf("Failed to login to proxmox @ %v: %v", url, err)
		log.Errorf(errMsg)
		return nil, errors.New(errMsg)
	}
	_, err := client.Version()
	if err != nil {
		log.Errorf("%v", err)
		return nil, err
	}
	return client, nil
}

func getConfigEntry(targetVM string) (ConfigEntry, error) {
	for _, v := range Cfg.Servers {
		if v.LogicalName == targetVM {
			return v, nil
		}
	}
	return ConfigEntry{}, errors.New("failed to find a VM with that Name")
}

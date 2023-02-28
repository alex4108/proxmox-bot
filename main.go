package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Cfg *Config

type ConfigEntry struct {
	VMId        string `yaml:"proxmox_vm_id"`
	LogicalName string `yaml:"logical_name"`
	VMHostUrl   string `yaml:"proxmox_host_url"`
	VMHostName  string `yaml:"proxmox_host_name"`
	VMHostUser  string
	VMHostPass  string
}

type Config struct {
	Servers []ConfigEntry `yaml:"vms"`
}

func main() {
	log.Info("Initializing...")
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	InCI, InCIExist := os.LookupEnv("CI")
	if InCIExist && InCI == "true" {
		log.Info("Running in CI.  This proves functionality?")
		os.Exit(0)
	}

	Token, tokenExists := os.LookupEnv("PROXMOX_BOT_DISCORD_TOKEN")
	if !tokenExists {
		log.Error("PROXMOX_BOT_DISCORD_TOKEN is not set.  Exiting.")
		os.Exit(2)
	}

	CfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	err = NewConfig(CfgPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Error("error creating Discord session,", err)
		os.Exit(3)
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers | discordgo.IntentMessageContent

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Error("Error opening websocket connection,", err)
		os.Exit(4)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Info("proxmox-bot is online!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// messageCreate handles new message events from Discord and routes them to the appropriate handler.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	commandRouter(s, m)
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) error {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return err
	}

	Cfg = config
	return nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "/workspace/proxmox-bot-config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	envPath, envPathExist := os.LookupEnv("CONFIG_PATH")
	if envPathExist {
		configPath = envPath
	}

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

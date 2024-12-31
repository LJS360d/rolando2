package main

import (
	"fmt"
	"os"
	"os/signal"
	"rolando/config"
	"strconv"
	"syscall"

	"rolando/cmd/handlers"
	"rolando/cmd/log"
	"rolando/cmd/repositories"
	"rolando/cmd/services"

	"github.com/bwmarrin/discordgo"
)

// LDFLAGS
var (
	Version string
	Env     string
	Build   string
)

const (
	DB_PATH = "rolando.db"
)

func main() {
	config.Version = Version
	config.Build = Build
	config.Env = Env
	fmt.Println("Version: ", config.Version)
	fmt.Println("Build: ", config.Build)
	fmt.Println("Env: ", config.Env)
	ds, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Log.Fatalf("error creating Discord session,", err)
	}

	ds.Identify.Intents = config.Intents

	// Open a websocket connection to Discord and begin listening.
	err = ds.Open()
	if err != nil {
		log.Log.Fatalln("error opening connection,", err)
	}

	err = ds.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Type: discordgo.ActivityTypeWatching,
				Name: strconv.Itoa(len(ds.State.Guilds)) + " servers",
			},
		},
		Status:    "online",
		AFK:       false,
		IdleSince: nil,
	})
	if err != nil {
		log.Log.Fatalf("error setting bot presence: %v", err)
	}
	// DI
	messagesRepo, err := repositories.NewMessagesRepository(DB_PATH)
	if err != nil {
		log.Log.Fatalf("error creating messages repository: %v", err)
	}
	chainsRepo, err := repositories.NewChainsRepository(DB_PATH)
	if err != nil {
		log.Log.Fatalf("error creating chains repository: %v", err)
	}
	chainsService := services.NewChainsService(ds, *chainsRepo, *messagesRepo)
	dataFetchService := services.NewDataFetchService(ds, chainsService, messagesRepo)
	// Handlers
	messagesHandler := handlers.NewMessageHandler(ds, chainsService)
	commandsHandler := handlers.NewSlashCommandsHandler(ds, chainsService)
	buttonsHandler := handlers.NewButtonsHandler(ds, dataFetchService, chainsService)
	// Register
	chainsService.LoadChains()
	ds.AddHandler(commandsHandler.OnSlashCommandInteraction)
	ds.AddHandler(messagesHandler.OnMessageCreate)
	ds.AddHandler(buttonsHandler.OnButtonInteraction)
	// Wait here until SIGINT or other term signal is received.
	log.Log.Infof("Logged in as %s!", ds.State.User.String())
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	ds.Close()
}

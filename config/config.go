package config

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	Token         string
	ApplicationID string
	Intents       discordgo.Intent
)

func init() {
	log.Println("Initializing configuration...")
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Assign the environment variables to package-level variables
	Token = os.Getenv("TOKEN")
	if Token == "" {
		log.Fatalf("TOKEN not set in the environment")
	}
	ApplicationID = os.Getenv("APPLICATION_ID")
	if ApplicationID == "" {
		log.Fatalf("APPLICATION_ID not set in the environment")
	}

	Intents = (discordgo.IntentDirectMessageReactions |
		discordgo.IntentDirectMessageTyping |
		discordgo.IntentDirectMessages |
		// discordgo.IntentAutoModerationConfiguration |
		// discordgo.IntentAutoModerationExecution |
		// discordgo.IntentDirectMessageReactions |
		// discordgo.IntentGuildEmojisAndStickers |
		// discordgo.IntentGuildIntegrations |
		discordgo.IntentGuildInvites |
		// discordgo.IntentGuildMembers |
		discordgo.IntentGuildMessageReactions |
		discordgo.IntentGuildMessageTyping |
		discordgo.IntentGuildMessages |
		// discordgo.IntentGuildModeration |
		// discordgo.IntentGuildPresences |
		// discordgo.IntentGuildScheduledEvents |
		// discordgo.IntentGuildVoiceStates |
		// discordgo.IntentGuildWebhooks |
		discordgo.IntentGuilds |
		discordgo.IntentMessageContent)

}

package handlers

import (
	"rolando/app/services"

	"github.com/bwmarrin/discordgo"
)

type SlashCommandsHandler struct {
	Client           *discordgo.Session
	ChainsService    *services.ChainsService
	DataFetchService *services.DataFetchService
	Commands         map[string]SlashCommandHandler
}

type SlashCommandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

type SlashCommand struct {
	Command *discordgo.ApplicationCommand
	Handler SlashCommandHandler
}

// Constructor function for SlashCommandsHandler
func NewSlashCommandsHandler(
	client *discordgo.Session,
	chainsService *services.ChainsService,
	dataFetchService *services.DataFetchService,
) *SlashCommandsHandler {
	handler := &SlashCommandsHandler{
		Client:           client,
		ChainsService:    chainsService,
		DataFetchService: dataFetchService,
		Commands:         make(map[string]SlashCommandHandler),
	}

	// Initialize commands
	handler.registerCommands([]SlashCommand{
		{
			Command: &discordgo.ApplicationCommand{
				Name:        "train",
				Description: "Fetches all available messages in the server to be used as training data",
			},
			Handler: handler.trainCommand,
		},
	})

	return handler
}

// Register commands and build the command map
func (h *SlashCommandsHandler) registerCommands(commands []SlashCommand) {
	for _, cmd := range commands {
		h.Commands[cmd.Command.Name] = cmd.Handler
		h.Client.ApplicationCommandCreate(h.Client.State.User.ID, "", cmd.Command)
	}
}

// Entry point for handling slash command interactions
func (h *SlashCommandsHandler) OnSlashCommandInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name
	if handler, exists := h.Commands[commandName]; exists {
		handler(s, i) // Call the function bound to this command
	}
}

// Implementation for the "train" command
func (h *SlashCommandsHandler) trainCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	guildID := i.GuildID

	// Acknowledge the interaction
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		return
	}

	// Fetch data using the DataFetchService
	trainingData, err := h.DataFetchService.FetchAllGuildMessages(guildID)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: "Failed to fetch training data.",
		})
		return
	}

	// Train the chain using ChainsService
	_, err = h.ChainsService.UpdateChainState(guildID, trainingData)
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
			Content: "Failed to train the chain.",
		})
		return
	}

	// Respond with success
	s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
		Content: "Training data successfully fetched and processed.",
	})
}

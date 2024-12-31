package handlers

import (
	"rolando/cmd/log"
	"rolando/cmd/services"

	"github.com/bwmarrin/discordgo"
)

type ButtonsHandler struct {
	Client           *discordgo.Session
	ChainsService    *services.ChainsService
	DataFetchService *services.DataFetchService
	Handlers         map[string]ButtonHandler
}

type ButtonHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

// Constructor for ButtonsHandler
func NewButtonsHandler(client *discordgo.Session, dataFetchService *services.DataFetchService, chainsService *services.ChainsService) *ButtonsHandler {
	handler := &ButtonsHandler{
		Client:           client,
		ChainsService:    chainsService,
		DataFetchService: dataFetchService,
		Handlers:         make(map[string]ButtonHandler),
	}

	// Register button handlers
	handler.registerButtons()

	return handler
}

// Register button handlers in the map
func (h *ButtonsHandler) registerButtons() {
	h.Handlers["confirm-train"] = h.onConfirmTrain
	h.Handlers["cancel-train"] = h.onCancelTrain
}

// Entry point for handling button interactions
func (h *ButtonsHandler) OnButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// Check if there's a handler for the button ID
	if handler, exists := h.Handlers[i.MessageComponentData().CustomID]; exists {
		handler(s, i) // Call the function bound to the button ID
	}
}

// Handle 'confirm-train' button interaction
func (h *ButtonsHandler) onConfirmTrain(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Defer the update
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	// Check if training is already completed
	chain, err := h.ChainsService.GetChainDocument(i.GuildID)
	if err != nil {
		log.Log.Errorf("Failed to fetch chain for guild %s: %v", i.GuildID, err)
		return
	}

	if chain.Trained {
		s.ChannelMessageSend(i.ChannelID, "Training already completed for this server.")
		return
	}

	// Start the training process
	// Send confirmation message
	s.ChannelMessageSend(i.ChannelID, "Training started! Fetching data...")

	go func() {
		_, err := h.DataFetchService.FetchAllGuildMessages(i.GuildID)
		if err != nil {
			log.Log.Errorf("Failed to fetch messages for guild %s: %v", i.GuildID, err)
			return
		}

		// Update chain status
		chain.Trained = true
		_, err = h.ChainsService.UpdateChainDocument(i.GuildID, map[string]any{"trained": true})
		if err != nil {
			log.Log.Errorf("Failed to update chain document for guild %s: %v", i.GuildID, err)
			return
		}

		// Send completion message
		s.ChannelMessageSend(i.ChannelID, "Training completed! Training data fetched.")
	}()
}

// Handle 'cancel-train' button interaction
func (h *ButtonsHandler) onCancelTrain(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Respond to cancel interaction
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Training process has been canceled.",
		},
	})
}

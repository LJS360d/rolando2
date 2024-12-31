package handlers

import (
	"fmt"
	"rolando/app/log"
	"rolando/app/model"
	"rolando/app/services"
	"rolando/app/utils"
	"rolando/config"

	"github.com/bwmarrin/discordgo"
)

type SlashCommandsHandler struct {
	Client        *discordgo.Session
	ChainsService *services.ChainsService
	Commands      map[string]SlashCommandHandler
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
) *SlashCommandsHandler {
	handler := &SlashCommandsHandler{
		Client:        client,
		ChainsService: chainsService,
		Commands:      make(map[string]SlashCommandHandler),
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
		{
			Command: &discordgo.ApplicationCommand{
				Name:        "gif",
				Description: "Returns a gif from the ones it knows",
			},
			Handler: handler.gifCommand,
		},
		{
			Command: &discordgo.ApplicationCommand{
				Name:        "image",
				Description: "Returns an image from the ones it knows",
			},
			Handler: handler.imageCommand,
		},
		{
			Command: &discordgo.ApplicationCommand{
				Name:        "video",
				Description: "Returns a video from the ones it knows",
			},
			Handler: handler.videoCommand,
		},
		{
			Command: &discordgo.ApplicationCommand{
				Name:        "analytics",
				Description: "Returns analytics about the bot",
			},
			Handler: handler.analyticsCommand,
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

func (h *SlashCommandsHandler) trainCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !h.checkAdmin(i) {
		return
	}
	h.ChainsService.GetChain(i.GuildID)
	chainDoc, err := h.ChainsService.GetChainDocument(i.GuildID)
	if err != nil {
		log.Log.Errorf("Failed to fetch chain document for guild %s: %v", i.GuildID, err)
		return
	}
	if chainDoc.Trained {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Training already completed for this server.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Create buttons for confirmation and cancellation
	confirmButton := &discordgo.Button{
		Label:    "Confirm",
		Style:    discordgo.PrimaryButton,
		CustomID: "confirm-train",
	}
	cancelButton := &discordgo.Button{
		Label:    "Cancel",
		Style:    discordgo.DangerButton,
		CustomID: "cancel-train",
	}

	// Create an action row with the buttons
	actionRow := &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			confirmButton,
			cancelButton,
		},
	}

	// Send the reply with buttons
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: `Are you sure you want to use **ALL SERVER MESSAGES** as training data for me?
This will fetch data in all accessible channels and delete all previous training data for this server.
If you wish to exclude specific channels, revoke my typing permissions in those channels.
`,
			Components: []discordgo.MessageComponent{*actionRow},
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}

func (h *SlashCommandsHandler) gifCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	chain, err := h.ChainsService.GetChain(i.GuildID)
	if err != nil {
		return
	}
	gif, err := chain.MediaStorage.GetMedia("gif")
	if err != nil || gif == "" {
		gif = "No valid gif found."
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: gif,
		},
	})
}

func (h *SlashCommandsHandler) imageCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	chain, err := h.ChainsService.GetChain(i.GuildID)
	if err != nil {
		return
	}
	image, err := chain.MediaStorage.GetMedia("image")
	if err != nil || image == "" {
		image = "No valid image found."
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: image,
		},
	})
}

func (h *SlashCommandsHandler) videoCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	chain, err := h.ChainsService.GetChain(i.GuildID)
	if err != nil {
		return
	}
	video, err := chain.MediaStorage.GetMedia("video")
	if err != nil || video == "" {
		video = "No valid video found."
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: video,
		},
	})
}

func (h *SlashCommandsHandler) analyticsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Fetch the chain data for the given guild
	chain, err := h.ChainsService.GetChain(i.GuildID)
	if err != nil {
		log.Log.Errorf("Failed to fetch chain for guild %s: %v", i.GuildID, err)
		return
	}
	analytics := model.NewMarkovChainAnalyzer(chain).GetRawAnalytics()
	chainDoc, err := h.ChainsService.GetChainDocument(i.GuildID)
	if err != nil {
		log.Log.Errorf("Failed to fetch chain document for guild %s: %v", i.GuildID, err)
		return
	}
	// Constructing the embed
	embed := &discordgo.MessageEmbed{
		Title:       "Analytics",
		Description: "",
		Color:       0xFFD700, // Gold color
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Complexity Score",
				Value:  fmt.Sprintf("```%d```", analytics.ComplexityScore),
				Inline: true,
			},
			{
				Name:   "Vocabulary",
				Value:  fmt.Sprintf("```%d words```", analytics.Words),
				Inline: true,
			},
			{
				Name:   "\t", // Empty field for spacing
				Value:  "\t",
				Inline: false,
			},
			{
				Name:   "Gifs",
				Value:  fmt.Sprintf("```%d```", analytics.Gifs),
				Inline: true,
			},
			{
				Name:   "Videos",
				Value:  fmt.Sprintf("```%d```", analytics.Videos),
				Inline: true,
			},
			{
				Name:   "Images",
				Value:  fmt.Sprintf("```%d```", analytics.Images),
				Inline: true,
			},
			{
				Name:   "\t", // Empty field for spacing
				Value:  "\t",
				Inline: false,
			},
			{
				Name:   "Processed Messages",
				Value:  fmt.Sprintf("```%d```", analytics.Messages),
				Inline: true,
			},
			{
				Name:   "Size",
				Value:  fmt.Sprintf("```%s / %s```", utils.FormatBytes(uint64(analytics.Size)), utils.FormatBytes(uint64(chainDoc.MaxSizeMb*1024*1024))),
				Inline: true,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    fmt.Sprintf("Version: %s", config.Version),
			IconURL: s.State.User.AvatarURL("256"),
		},
	}

	// Send the response with the embed
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		log.Log.Errorf("Failed to send analytics embed: %v", err)
	}
}

func (h *SlashCommandsHandler) checkAdmin(i *discordgo.InteractionCreate, msg ...string) bool {
	// Check if the user is in the list of owner IDs
	for _, ownerID := range config.OwnerIDs {
		if i.Member.User.ID == ownerID {
			return true
		}
	}
	/*
		// Check if the user has Administrator permissions
		perms, err := h.Client.check(i.Member.User.ID, i.GuildID)
		if err != nil {
			// Log the error if necessary
			return false
		}

		if perms&discordgo.PermissionAdministrator != 0 {
			return true
		}

		// Respond with a message if the user is not authorized
		if len(msg) > 0 {
			_ = h.Client.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg[0],
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		} else {
			_ = h.Client.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You are not authorized to use this command.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
	*/
	return false
}

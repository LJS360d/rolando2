package handlers

import (
	"github.com/bwmarrin/discordgo"
)

// SlashCommandsHandler handles incoming slash commands and dispatches them to the appropriate handler.
func ButtonInteractionsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	_ = i.ApplicationCommandData().Name

}

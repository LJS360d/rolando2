package services

import (
	"log"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type DataFetchService struct {
	Session        *discordgo.Session
	MessageLimit   int
	MaxFetchErrors int
	ChainService   *ChainsService
}

func NewDataFetchService(session *discordgo.Session, chainService *ChainsService) *DataFetchService {
	return &DataFetchService{
		Session:        session,
		MessageLimit:   750000,
		MaxFetchErrors: 5,
		ChainService:   chainService,
	}
}

// FetchAllGuildMessages fetches messages from all accessible channels in the guild.
func (d *DataFetchService) FetchAllGuildMessages(guildID string) ([]string, error) {
	guild, err := d.Session.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	messageCh := make(chan []string, len(guild.Channels))
	for _, channel := range guild.Channels {
		if !d.hasChannelAccess(channel) {
			continue
		}

		wg.Add(1)
		go func(channel *discordgo.Channel) {
			defer wg.Done()
			messages, err := d.fetchChannelMessages(channel.ID)
			if err != nil {
				log.Printf("Failed to fetch messages for channel %s: %v", channel.Name, err)
				return
			}
			messageCh <- messages
		}(channel)
	}

	wg.Wait()
	close(messageCh)

	var allMessages []string
	for msgs := range messageCh {
		allMessages = append(allMessages, msgs...)
	}

	log.Printf("Fetched %d messages in guild %s", len(allMessages), guildID)
	return allMessages, nil
}

func (d *DataFetchService) fetchChannelMessages(channelID string) ([]string, error) {
	var messages []string
	var lastMessageID string
	errorCount := 0

	for len(messages) < d.MessageLimit {
		batch, err := d.getMessageBatch(channelID, lastMessageID)
		if err != nil {
			errorCount++
			if errorCount > d.MaxFetchErrors {
				log.Printf("Error limit reached for channel %s: %v", channelID, err)
				break
			}
			continue
		}

		if len(batch) == 0 {
			break
		}

		for _, msg := range batch {
			messages = append(messages, msg.Content)
			// Update chain state (assumes ChainService has an UpdateChainState method)
			go d.ChainService.UpdateChainState(msg.GuildID, []string{msg.Content})
		}
		lastMessageID = batch[len(batch)-1].ID
	}

	log.Printf("Fetched %d messages from channel %s", len(messages), channelID)
	return messages, nil
}

func (d *DataFetchService) getMessageBatch(channelID, lastMessageID string) ([]*discordgo.Message, error) {

	messages, err := d.Session.ChannelMessages(channelID, 100, lastMessageID, "", "")
	if err != nil {
		return nil, err
	}

	var cleanMessages []*discordgo.Message
	for _, msg := range messages {
		if len(strings.Fields(msg.Content)) > 1 || d.containsURL(msg.Content) {
			cleanMessages = append(cleanMessages, msg)
		}
	}
	return cleanMessages, nil
}

func (d *DataFetchService) hasChannelAccess(channel *discordgo.Channel) bool {
	permissions, err := d.Session.State.UserChannelPermissions(d.Session.State.User.ID, channel.ID)
	if err != nil {
		return false
	}
	return permissions&discordgo.PermissionViewChannel != 0
}

func (d *DataFetchService) containsURL(content string) bool {
	return strings.Contains(content, "http://") || strings.Contains(content, "https://")
}

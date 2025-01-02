package bot

import (
	"rolando/cmd/log"
	"rolando/cmd/services"
	"rolando/config"
	"rolando/server/auth"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

type BotController struct {
	chainsService *services.ChainsService
	ds            *discordgo.Session
}

func NewController(chainsService *services.ChainsService, ds *discordgo.Session) *BotController {
	return &BotController{
		chainsService: chainsService,
		ds:            ds,
	}
}

type BroadcastRequest struct {
	Content string                   `json:"content"`
	Guilds  []*BroadcastGuildRequest `json:"guilds"`
}

type BroadcastGuildRequest struct {
	Id        string `json:"id"`
	ChannelId string `json:"channel_id"`
}

// POST /bot/broadcast, requires owner authorization
func (s *BotController) Broadcast(c *gin.Context) {
	errCode, err := auth.EnsureOwner(c, s.ds)
	if err != nil {
		c.JSON(errCode, gin.H{"error": err.Error()})
		return
	}
	req := &BroadcastRequest{}
	err = c.ShouldBindBodyWithJSON(req)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(req.Guilds))

	for _, g := range req.Guilds {
		wg.Add(1)

		go func(g *BroadcastGuildRequest) {
			defer wg.Done()
			channelId := g.ChannelId
			if channelId == "" {
				guild, err := s.ds.Guild(g.Id)
				if err != nil {
					errCh <- err
					return
				}
				channelId = guild.SystemChannelID
			}
			log.Log.Infof("Broadcasting message in guild: %s, channel: %s", g.Id, channelId)
			_, err := s.ds.ChannelMessageSend(channelId, req.Content)
			if err != nil {
				log.Log.Errorf("could not send message in guild: %s, channel: %s: %v", g.Id, channelId, err)
				errCh <- err
			}
		}(g)
	}

	wg.Wait()
	close(errCh)

	// Collect errors, if any
	for err := range errCh {
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(200, gin.H{"content": req.Content})
}

// GET /bot/guilds, requires owner authorization
func (s *BotController) GetBotGuilds(c *gin.Context) {
	errCode, err := auth.EnsureOwner(c, s.ds)
	if err != nil {
		c.JSON(errCode, gin.H{"error": err.Error()})
		return
	}
	guilds, err := s.ds.UserGuilds(200, "", "", true)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"guilds": guilds})
}

// GET /bot/user, public
func (s *BotController) GetBotUser(c *gin.Context) {
	botUser := s.ds.State.User
	commands, err := s.ds.ApplicationCommands(s.ds.State.User.ID, "")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"id":             botUser.ID,
		"username":       botUser.Username,
		"global_name":    botUser.Username + "#" + botUser.Discriminator,
		"avatar_url":     botUser.AvatarURL(""),
		"discriminator":  botUser.Discriminator,
		"verified":       botUser.Verified,
		"accent_color":   int32(botUser.AccentColor),
		"invite_url":     config.InviteUrl,
		"slash_commands": commands,
	})
}

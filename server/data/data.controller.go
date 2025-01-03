package data

import (
	"rolando/cmd/repositories"
	"rolando/server/auth"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

type DataController struct {
	messagesRepo *repositories.MessagesRepository
	ds           *discordgo.Session
}

func NewController(ds *discordgo.Session, messagesRepo *repositories.MessagesRepository) *DataController {
	return &DataController{
		messagesRepo: messagesRepo,
		ds:           ds,
	}
}

// GET /data/:chain, requires guild member authorization
func (s *DataController) GetData(c *gin.Context) {
	chainId := c.Param("chain")
	errCode, err := auth.EnsureGuildMember(c, s.ds, chainId)
	if err != nil {
		c.JSON(errCode, gin.H{"error": err.Error()})
		return
	}
	messages, err := s.messagesRepo.GetAllGuildMessages(chainId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	content := make([]string, len(messages))
	for i, message := range messages {
		content[i] = message.Content
	}
	c.JSON(200, content)
}

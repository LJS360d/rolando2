package analytics

import (
	"rolando/cmd/model"
	"rolando/cmd/services"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	chainsService *services.ChainsService
	ds            *discordgo.Session
}

func NewController(chainsService *services.ChainsService, ds *discordgo.Session) *AnalyticsController {
	return &AnalyticsController{
		chainsService: chainsService,
		ds:            ds,
	}
}

func (s *AnalyticsController) GetChainAnalytics(c *gin.Context) {
	chainId := c.Param("chain")
	chainDoc, err := s.chainsService.GetChainDocument(chainId)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	chain, err := s.chainsService.GetChain(chainDoc.ID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	analyzer := model.NewMarkovChainAnalyzer(chain)
	rawAnalytics := analyzer.GetRawAnalytics()

	c.JSON(200, gin.H{
		"complexity_score": rawAnalytics.ComplexityScore,
		"gifs":             rawAnalytics.Gifs,
		"images":           rawAnalytics.Images,
		"videos":           rawAnalytics.Videos,
		"reply_rate":       rawAnalytics.ReplyRate,
		"words":            rawAnalytics.Words,
		"messages":         rawAnalytics.Messages,
		"bytes":            rawAnalytics.Size,
		"id":               chain.ID,
		"name":             chainDoc.Name,
		"max_size_mb":      chainDoc.MaxSizeMb,
	})
}

func (s *AnalyticsController) GetAllChainsAnalytics(c *gin.Context) {
	// Retrieve chains
	chains, err := s.chainsService.GetAllChains()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	allAnalytics := make([]gin.H, 0)
	for _, chain := range chains {
		analyzer := model.NewMarkovChainAnalyzer(chain)
		chainDoc, err := s.chainsService.GetChainDocument(chain.ID)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Get the analytics data
		rawAnalytics := analyzer.GetRawAnalytics()
		chainAnalytics := gin.H{
			"complexity_score": rawAnalytics.ComplexityScore,
			"gifs":             rawAnalytics.Gifs,
			"images":           rawAnalytics.Images,
			"videos":           rawAnalytics.Videos,
			"reply_rate":       rawAnalytics.ReplyRate,
			"words":            rawAnalytics.Words,
			"messages":         rawAnalytics.Messages,
			"bytes":            rawAnalytics.Size,
			"id":               chain.ID,
			"name":             chainDoc.Name,
			"max_size_mb":      chainDoc.MaxSizeMb,
		}
		allAnalytics = append(allAnalytics, chainAnalytics)
	}
	c.JSON(200, allAnalytics)
}

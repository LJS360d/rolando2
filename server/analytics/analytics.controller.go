package analytics

import (
	"rolando/cmd/model"
	"rolando/cmd/services"

	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	chainsService *services.ChainsService
}

func NewAnalyticsController(chainsService *services.ChainsService) *AnalyticsController {
	return &AnalyticsController{
		chainsService: chainsService,
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
		"complexityScore": rawAnalytics.ComplexityScore,
		"gifs":            rawAnalytics.Gifs,
		"images":          rawAnalytics.Images,
		"videos":          rawAnalytics.Videos,
		"replyRate":       rawAnalytics.ReplyRate,
		"words":           rawAnalytics.Words,
		"messages":        rawAnalytics.Messages,
		"bytes":           rawAnalytics.Size,
		"id":              chain.ID,
		"name":            chainDoc.Name,
	})
}

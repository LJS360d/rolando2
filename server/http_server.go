package server

import (
	"rolando/cmd/log"
	"rolando/cmd/services"
	"rolando/config"
	"rolando/server/analytics"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	ChainsService  *services.ChainsService
	DiscordSession *discordgo.Session
	engine         *gin.Engine
}

func NewHttpServer(chainsService *services.ChainsService, discordSession *discordgo.Session) *HttpServer {
	return &HttpServer{
		ChainsService:  chainsService,
		DiscordSession: discordSession,
	}
}

func (s *HttpServer) Start() {
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	s.engine = r

	analyticsController := analytics.NewAnalyticsController(s.ChainsService)
	// Routes
	r.GET("/analytics/:chain", analyticsController.GetChainAnalytics)

	// Start the server
	log.Log.Infof("Server listening at %v", config.ServerAddress)
	if err := r.Run(config.ServerAddress); err != nil {
		log.Log.Fatalf("failed to start server: %v", err)
	}
}

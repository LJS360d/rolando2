package server

import (
	"rolando/cmd/log"
	"rolando/cmd/services"
	"rolando/config"
	"rolando/server/analytics"
	"rolando/server/auth"
	"rolando/server/bot"

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

	analyticsController := analytics.NewController(s.ChainsService, s.DiscordSession)
	botController := bot.NewController(s.ChainsService, s.DiscordSession)
	authController := auth.NewController(s.DiscordSession)
	// Routes
	r.GET("/auth/@me", authController.GetUser)

	r.GET("/analytics/:chain", analyticsController.GetChainAnalytics)
	r.GET("/analytics", analyticsController.GetAllChainsAnalytics)

	r.GET("/bot/user", botController.GetBotUser)
	r.GET("/bot/guilds", botController.GetBotGuilds)
	r.POST("/bot/broadcast", botController.Broadcast)

	// Start the server
	log.Log.Infof("Server listening at %v", config.ServerAddress)
	if err := r.Run(config.ServerAddress); err != nil {
		log.Log.Fatalf("failed to start server: %v", err)
	}
}

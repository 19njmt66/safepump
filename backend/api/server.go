package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StartServer sets up the Gin router and starts the HTTP REST & WebSocket gateway
func StartServer(addr string, hub *Hub) error {
	// Set Gin to release mode in production, but let default handle logging for development
	router := gin.Default()

	// Robust CORS middleware to enable connection from decentralized frontend (IPFS/Fleek)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// Serve static uploads
	router.Static("/uploads", "./uploads")

	// API Routes Group
	v1 := router.Group("/api/v1")
	{
		v1.GET("/tokens", GetTokens)
		v1.GET("/tokens/:address", GetTokenDetails)
		v1.GET("/tokens/:address/audit", GetTokenAudit)
		v1.POST("/tokens/:address/metadata", UpdateTokenMetadata)
		v1.GET("/tokens/:address/trades", GetTokenTrades)
		v1.GET("/tokens/:address/holders", GetTokenHolders)
		v1.GET("/tokens/:address/candles", GetTokenCandles)
		v1.GET("/creator/:creator/tokens", GetCreatorTokens)
		v1.POST("/upload", UploadFile)
		v1.GET("/users/:address", GetUserProfile)
		v1.GET("/users/:address/portfolio", GetUserPortfolio)
		v1.GET("/users/:address/trades", GetUserTrades)
		v1.GET("/trades", GetRecentTrades)
		v1.POST("/users/profile", UpdateUserProfile)
		v1.GET("/config", GetConfig)
	}

	// WebSocket endpoint
	router.GET("/ws", func(c *gin.Context) {
		ServeWs(hub, c.Writer, c.Request)
	})

	log.Printf("Starting API and WebSocket gateway on %s...", addr)
	return router.Run(addr)
}

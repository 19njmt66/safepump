package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend/api"
	"backend/config"
	"backend/database"
	"backend/indexer"
)

func main() {
	log.Println("Initializing SafePump backend services...")

	// 1. Load configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Critical error during database initialization: %v", err)
	}

	// 3. Initialize WebSocket Hub
	hub := api.NewHub()
	api.GlobalHub = hub
	go hub.Run()
	log.Println("WebSocket event distribution hub initialized.")

	// 4. Initialize Block Indexer
	idx, err := indexer.NewIndexer(cfg.RPCURL, db, cfg.FactoryAddress, cfg.StartBlock)
	if err != nil {
		log.Fatalf("Critical error during indexer initialization: %v", err)
	}

	// Bind Indexer callbacks to WebSocket notifications to enable live streaming feeds
	idx.OnTokenCreated = func(token database.Token) {
		log.Printf("[Live Feed] New Token Launched: %s (%s)", token.Name, token.Address)
		hub.BroadcastObject(map[string]interface{}{
			"type": "token_created",
			"data": token,
		})
	}

	idx.OnTrade = func(trade database.Trade) {
		log.Printf("[Live Feed] Trade Logged: %s | IsBuy=%t | ETH=%s | Token=%s",
			trade.TxHash[:8], trade.IsBuy, trade.EthAmount, trade.TokenAddress[:8])
		hub.BroadcastObject(map[string]interface{}{
			"type": "trade_processed",
			"data": trade,
		})
	}

	idx.OnMigrated = func(token database.Token) {
		log.Printf("[Live Feed] Graduation / Migration Complete: Token=%s | Pair=%s",
			token.Address[:8], token.PairAddress[:8])
		hub.BroadcastObject(map[string]interface{}{
			"type": "token_migrated",
			"data": token,
		})
	}

	// Start Indexer loop in the background
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go idx.Start(ctx)

	// 5. Start API Server (blocking call in this goroutine or running separately)
	// We run it in a goroutine and wait for OS signals for graceful shutdown
	go func() {
		err := api.StartServer(cfg.ServerAddr, hub)
		if err != nil {
			log.Fatalf("Critical error during API server execution: %v", err)
		}
	}()

	// 6. Wait for termination signals (Ctrl+C, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown signal received. Terminating backend services...")
}

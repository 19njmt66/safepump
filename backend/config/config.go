package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	DatabaseURL    string
	RPCURL         string
	FactoryAddress string
	StartBlock     uint64
	ServerAddr     string
}

// LoadConfig loads variables from .env or system environment
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Note: .env file not found, loading directly from system variables.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "host=localhost user=postgres password=postgres dbname=safepump port=5432 sslmode=disable"
	}

	rpcURL := os.Getenv("RPC_URL")
	if rpcURL == "" {
		// Default to public Base Sepolia L2 RPC
		rpcURL = "https://sepolia.base.org"
	}

	factoryAddr := os.Getenv("FACTORY_ADDRESS")
	if factoryAddr == "" {
		factoryAddr = "0x0000000000000000000000000000000000000000"
	}

	startBlockStr := os.Getenv("START_BLOCK")
	startBlock := uint64(0)
	if startBlockStr != "" {
		parsed, err := strconv.ParseUint(startBlockStr, 10, 64)
		if err == nil {
			startBlock = parsed
		}
	}

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	return &Config{
		DatabaseURL:    dbURL,
		RPCURL:         rpcURL,
		FactoryAddress: factoryAddr,
		StartBlock:     startBlock,
		ServerAddr:     serverAddr,
	}
}

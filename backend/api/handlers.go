package api

import (
	"math/big"
	"net/http"
	"time"

	"backend/database"

	"github.com/gin-gonic/gin"
)

// GetTokens retrieves all tokens sorted by order (default: bonding curve progress)
// Query params:
// - sort: "progress" (default), "new"
func GetTokens(c *gin.Context) {
	sort := c.DefaultQuery("sort", "progress")

	var tokens []database.Token
	query := database.DB

	if sort == "new" {
		query = query.Order("created_at desc")
	} else {
		// Sorting by progress (tokens sold)
		// We show active bonding curves first (migrated=false) ordered by progress, then migrated ones
		query = query.Order("migrated asc, (tokens_sold + 0) desc")
	}

	if err := query.Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query tokens"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// GetTokenDetails retrieves details of a specific token
func GetTokenDetails(c *gin.Context) {
	address := c.Param("address")

	var token database.Token
	if err := database.DB.Where("address = ?", address).First(&token).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	c.JSON(http.StatusOK, token)
}

// UpdateTokenMetadataRequest represents the request body to update token metadata
type UpdateTokenMetadataRequest struct {
	CreatorAddress string `json:"creator_address"`
	Description    string `json:"description"`
	ImageUrl       string `json:"image_url"`
}

// UpdateTokenMetadata updates the description and image_url of a token (only by the creator, in simple test mode)
func UpdateTokenMetadata(c *gin.Context) {
	address := c.Param("address")

	var req UpdateTokenMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var token database.Token
	err := database.DB.Where("address = ?", address).First(&token).Error
	if err != nil {
		// Not found yet! Create a new record with the address, creator, description, and image_url
		token = database.Token{
			Address:     address,
			Creator:     req.CreatorAddress,
			Description: req.Description,
			ImageUrl:    req.ImageUrl,
			TokensSold:  "0",
			EthRaised:   "0",
			Migrated:    false,
		}
	} else {
		// Simple local security validation: verify that creator_address matches the DB creator address
		if req.CreatorAddress != "" && token.Creator != req.CreatorAddress {
			c.JSON(http.StatusForbidden, gin.H{"error": "Only the creator can update metadata"})
			return
		}
		token.Description = req.Description
		token.ImageUrl = req.ImageUrl
	}

	if err := database.DB.Save(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token metadata"})
		return
	}

	c.JSON(http.StatusOK, token)
}

// GetTokenTrades retrieves historical trades for a specific token
func GetTokenTrades(c *gin.Context) {
	address := c.Param("address")

	var trades []database.Trade
	if err := database.DB.Where("token_address = ?", address).Order("timestamp desc").Limit(200).Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query trades"})
		return
	}

	c.JSON(http.StatusOK, trades)
}

// GetCreatorTokens retrieves all tokens launched by a creator address
func GetCreatorTokens(c *gin.Context) {
	creator := c.Param("creator")

	var tokens []database.Token
	if err := database.DB.Where("creator = ?", creator).Order("created_at desc").Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query creator tokens"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Candle represents an OHLC candlestick item for financial charts
type Candle struct {
	Time  int64   `json:"time"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

// GetTokenCandles retrieves aggregated candlestick chart data for a specific token
// Query params:
// - resolution: "1m" (default), "5m", "15m", "1h", "1d"
func GetTokenCandles(c *gin.Context) {
	address := c.Param("address")
	resolution := c.DefaultQuery("resolution", "1m")

	var duration time.Duration
	switch resolution {
	case "1m":
		duration = time.Minute
	case "5m":
		duration = 5 * time.Minute
	case "15m":
		duration = 15 * time.Minute
	case "1h":
		duration = time.Hour
	case "1d":
		duration = 24 * time.Hour
	default:
		duration = time.Minute
	}

	var trades []database.Trade
	// Query trades in ascending order (oldest first) to build the historical candles
	if err := database.DB.Where("token_address = ?", address).Order("timestamp asc").Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query trades"})
		return
	}

	if len(trades) == 0 {
		c.JSON(http.StatusOK, []Candle{})
		return
	}

	var candles []Candle
	var currentCandle *Candle

	for _, trade := range trades {
		// Calculate the instantaneous trade price = ethAmount / tokenAmount
		ethVal, ok1 := new(big.Float).SetString(trade.EthAmount)
		tokenVal, ok2 := new(big.Float).SetString(trade.TokenAmount)
		if !ok1 || !ok2 || tokenVal.Sign() == 0 {
			continue
		}

		price, _ := new(big.Float).Quo(ethVal, tokenVal).Float64()

		// Truncate the trade timestamp to the resolution bucket
		bucketTime := trade.Timestamp.Truncate(duration).Unix()

		if currentCandle == nil {
			currentCandle = &Candle{
				Time:  bucketTime,
				Open:  price,
				High:  price,
				Low:   price,
				Close: price,
			}
		} else if currentCandle.Time == bucketTime {
			// Update the current candle stats
			if price > currentCandle.High {
				currentCandle.High = price
			}
			if price < currentCandle.Low {
				currentCandle.Low = price
			}
			currentCandle.Close = price
		} else {
			// Store previous candle and initialize new bucket
			candles = append(candles, *currentCandle)
			currentCandle = &Candle{
				Time:  bucketTime,
				Open:  price,
				High:  price,
				Low:   price,
				Close: price,
			}
		}
	}

	if currentCandle != nil {
		candles = append(candles, *currentCandle)
	}

	c.JSON(http.StatusOK, candles)
}

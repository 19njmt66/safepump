package api

import (
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"backend/config"
	"backend/database"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
)

var defaultMemeImages = []string{
	"https://images.unsplash.com/photo-1543466835-00a7907e9de1?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1537151608828-ea2b117b6f86?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1583511655857-d19b40a7a54e?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1573865526739-10659fec78a5?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1533738363-b7f9aef128ce?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1561948955-570b270e7c36?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1592194996308-7b43878e84a6?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1472491235688-bdc81a63246e?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1507666480-1a4b9dfcd90f?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1591880911720-410d352270d1?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1618005182384-a83a8bd57fbe?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1614741118887-7a4ee193a5fa?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1506318137071-a8e063b4bec0?w=300&h=300&fit=crop",
	"https://images.unsplash.com/photo-1610296669228-602fa827fc1f?w=300&h=300&fit=crop",
}

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

	for i := range tokens {
		addr := strings.ToLower(tokens[i].Address)
		imgIdx := 0
		if len(addr) > 0 {
			imgIdx = int(addr[len(addr)-1]) % len(defaultMemeImages)
		}
		if tokens[i].ImageUrl == "" {
			tokens[i].ImageUrl = defaultMemeImages[imgIdx]
		}
		if tokens[i].Description == "" {
			tokens[i].Description = fmt.Sprintf("Bienvenue sur %s ($%s). Rejoignez notre communauté de détenteurs passionnés et aidez-nous à propulser ce projet vers la lune !", tokens[i].Name, tokens[i].Symbol)
		}
		if tokens[i].Website == "" {
			tokens[i].Website = fmt.Sprintf("https://www.%s.io", strings.ToLower(tokens[i].Symbol))
		}
		if tokens[i].Twitter == "" {
			tokens[i].Twitter = fmt.Sprintf("https://x.com/%s_coin", strings.ToLower(tokens[i].Symbol))
		}
		if tokens[i].Telegram == "" {
			tokens[i].Telegram = fmt.Sprintf("https://t.me/%s_chat", strings.ToLower(tokens[i].Symbol))
		}
	}

	c.JSON(http.StatusOK, tokens)
}

// GetTokenDetails retrieves details of a specific token
func GetTokenDetails(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))

	var token database.Token
	if err := database.DB.Where("address = ?", address).First(&token).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	addr := strings.ToLower(token.Address)
	imgIdx := 0
	if len(addr) > 0 {
		imgIdx = int(addr[len(addr)-1]) % len(defaultMemeImages)
	}
	if token.ImageUrl == "" {
		token.ImageUrl = defaultMemeImages[imgIdx]
	}
	if token.Description == "" {
		token.Description = fmt.Sprintf("Bienvenue sur %s ($%s). Rejoignez notre communauté de détenteurs passionnés et aidez-nous à propulser ce projet vers la lune !", token.Name, token.Symbol)
	}
	if token.Website == "" {
		token.Website = fmt.Sprintf("https://www.%s.io", strings.ToLower(token.Symbol))
	}
	if token.Twitter == "" {
		token.Twitter = fmt.Sprintf("https://x.com/%s_coin", strings.ToLower(token.Symbol))
	}
	if token.Telegram == "" {
		token.Telegram = fmt.Sprintf("https://t.me/%s_chat", strings.ToLower(token.Symbol))
	}

	c.JSON(http.StatusOK, token)
}

// UpdateTokenMetadataRequest represents the request body to update token metadata
type UpdateTokenMetadataRequest struct {
	CreatorAddress string `json:"creator_address"`
	Description    string `json:"description"`
	ImageUrl       string `json:"image_url"`
	Website        string `json:"website"`
	Twitter        string `json:"twitter"`
	Telegram       string `json:"telegram"`
}

// UpdateTokenMetadata updates the description, image_url, and social links of a token (only by the creator, in simple test mode)
func UpdateTokenMetadata(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))

	var req UpdateTokenMetadataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.CreatorAddress = strings.ToLower(req.CreatorAddress)

	var token database.Token
	err := database.DB.Where("address = ?", address).First(&token).Error
	if err != nil {
		// Not found yet! Create a new record with the address, creator, description, and image_url
		token = database.Token{
			Address:     address,
			Creator:     req.CreatorAddress,
			Description: req.Description,
			ImageUrl:    req.ImageUrl,
			Website:     req.Website,
			Twitter:     req.Twitter,
			Telegram:    req.Telegram,
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
		token.Website = req.Website
		token.Twitter = req.Twitter
		token.Telegram = req.Telegram
	}

	if err := database.DB.Save(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update token metadata"})
		return
	}

	c.JSON(http.StatusOK, token)
}

// GetTokenTrades retrieves historical trades for a specific token
func GetTokenTrades(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))

	var trades []database.Trade
	if err := database.DB.Where("token_address = ?", address).Order("timestamp desc").Limit(200).Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query trades"})
		return
	}

	c.JSON(http.StatusOK, trades)
}

// GetCreatorTokens retrieves all tokens launched by a creator address
func GetCreatorTokens(c *gin.Context) {
	creator := strings.ToLower(c.Param("creator"))

	var tokens []database.Token
	if err := database.DB.Where("creator = ?", creator).Order("created_at desc").Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query creator tokens"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Candle represents an OHLC candlestick item for financial charts
type Candle struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

// GetTokenCandles retrieves aggregated candlestick chart data for a specific token
// Query params:
// - resolution: "1m" (default), "5m", "15m", "1h", "1d"
func GetTokenCandles(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))
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
		tokenAmountFloat, _ := new(big.Float).Quo(tokenVal, new(big.Float).SetFloat64(1e18)).Float64()

		// Truncate the trade timestamp to the resolution bucket
		bucketTime := trade.Timestamp.Truncate(duration).Unix()

		if currentCandle == nil {
			currentCandle = &Candle{
				Time:   bucketTime,
				Open:   price,
				High:   price,
				Low:    price,
				Close:  price,
				Volume: tokenAmountFloat,
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
			currentCandle.Volume += tokenAmountFloat
		} else {
			// Store previous candle and initialize new bucket
			candles = append(candles, *currentCandle)
			currentCandle = &Candle{
				Time:   bucketTime,
				Open:   price,
				High:   price,
				Low:    price,
				Close:  price,
				Volume: tokenAmountFloat,
			}
		}
	}

	if currentCandle != nil {
		candles = append(candles, *currentCandle)
	}

	c.JSON(http.StatusOK, candles)
}

// UploadFile handles image uploads and returns the file URL
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("./uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create uploads directory"})
		return
	}

	// Generate a unique filename
	filename := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(file.Filename))
	dst := filepath.Join("./uploads", filename)

	// Save the file
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Construct absolute URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	fileURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, c.Request.Host, filename)

	c.JSON(http.StatusOK, gin.H{
		"url": fileURL,
	})
}

// GetUserProfile retrieves or creates a user profile by wallet address
func GetUserProfile(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))

	var user database.User
	err := database.DB.Where("LOWER(address) = ?", address).First(&user).Error
	if err == nil {
		c.JSON(http.StatusOK, user)
		return
	}

	// Create user with default values if not found
	short := address
	if len(address) > 8 {
		short = address[len(address)-6:]
	}
	user = database.User{
		Address:   address,
		Username:  "Degen-" + strings.ToUpper(short),
		Bio:       "Based SafePump Trader",
		AvatarUrl: "",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create default user profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfileReq defines payload for updating a user profile with Web3 personal signature validation
type UpdateProfileReq struct {
	Address   string `json:"address" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
	Message   string `json:"message" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

// UpdateUserProfile updates user profile details after validating signature
func UpdateUserProfile(c *gin.Context) {
	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	address := strings.ToLower(req.Address)

	// Cryptographic signature verification
	if !verifySignature(address, req.Message, req.Signature) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature verification"})
		return
	}

	// Retrieve or create user record
	var user database.User
	err := database.DB.Where("LOWER(address) = ?", address).First(&user).Error
	if err != nil {
		// Create new profile
		user = database.User{
			Address: address,
		}
	}

	// Update fields
	user.Username = req.Username
	user.Bio = req.Bio
	user.AvatarUrl = req.AvatarUrl

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save profile changes"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// verifySignature checks personal_sign signatures from Web3 wallets
func verifySignature(addressStr string, message string, signatureHex string) bool {
	signatureHex = strings.TrimPrefix(signatureHex, "0x")
	sigBytes, err := hexutil.Decode("0x" + signatureHex)
	if err != nil || len(sigBytes) != 65 {
		return false
	}

	// V adjustment: in Go-ethereum, recovery ID (V) must be 0 or 1, not 27 or 28
	if sigBytes[64] == 27 || sigBytes[64] == 28 {
		sigBytes[64] -= 27
	}

	// Add personal_sign message prefix
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(message), message)
	hash := crypto.Keccak256Hash([]byte(prefix))

	pubKey, err := crypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*pubKey)
	return strings.ToLower(recoveredAddress.Hex()) == strings.ToLower(addressStr)
}

// GetConfig retrieves the platform configuration (such as the factory address)
func GetConfig(c *gin.Context) {
	cfg := config.LoadConfig()
	c.JSON(http.StatusOK, gin.H{
		"factory_address": cfg.FactoryAddress,
	})
}

// GetUserPortfolio calculates portfolio stats and positions for a specific user
func GetUserPortfolio(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))

	var trades []database.Trade
	if err := database.DB.Where("buyer_or_seller = ?", address).Order("timestamp asc").Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user trades"})
		return
	}

	if len(trades) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"realized_pnl_eth":   0.0,
			"unrealized_pnl_eth": 0.0,
			"total_volume_eth":   0.0,
			"total_fees_eth":     0.0,
			"positions":          []interface{}{},
		})
		return
	}

	tokenAddressesMap := make(map[string]bool)
	for _, t := range trades {
		tokenAddressesMap[strings.ToLower(t.TokenAddress)] = true
	}
	var tokenAddresses []string
	for addr := range tokenAddressesMap {
		tokenAddresses = append(tokenAddresses, addr)
	}

	var tokens []database.Token
	if err := database.DB.Where("address IN ?", tokenAddresses).Find(&tokens).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query token metadata"})
		return
	}
	tokenMap := make(map[string]database.Token)
	for _, t := range tokens {
		tokenMap[strings.ToLower(t.Address)] = t
	}

	type LatestTradeResult struct {
		TokenAddress string
		EthAmount    string
		TokenAmount  string
	}
	var latestTrades []LatestTradeResult
	err := database.DB.Raw(`
		SELECT DISTINCT ON (token_address) token_address, eth_amount, token_amount 
		FROM trades 
		WHERE token_address IN ?
		ORDER BY token_address, timestamp DESC
	`, tokenAddresses).Scan(&latestTrades).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query current prices"})
		return
	}

	priceMap := make(map[string]float64)
	for _, lt := range latestTrades {
		ethVal, ok1 := new(big.Float).SetString(lt.EthAmount)
		tokenVal, ok2 := new(big.Float).SetString(lt.TokenAmount)
		if ok1 && ok2 && tokenVal.Sign() > 0 {
			price, _ := new(big.Float).Quo(ethVal, tokenVal).Float64()
			priceMap[strings.ToLower(lt.TokenAddress)] = price
		}
	}

	type PositionState struct {
		Holdings       *big.Int
		TotalCostEth   float64
		RealizedPnLEth float64
		TotalBoughtEth float64
		TotalSoldEth   float64
		TotalFeesEth   float64
	}
	states := make(map[string]*PositionState)

	for _, t := range trades {
		tokAddr := strings.ToLower(t.TokenAddress)
		state, exists := states[tokAddr]
		if !exists {
			state = &PositionState{
				Holdings: new(big.Int),
			}
			states[tokAddr] = state
		}

		tokenAmt, ok1 := new(big.Int).SetString(t.TokenAmount, 10)
		ethAmtFloat, ok2 := new(big.Float).SetString(t.EthAmount)
		feeFloat, ok3 := new(big.Float).SetString(t.Fee)

		if !ok1 || !ok2 || !ok3 {
			continue
		}

		ethAmt, _ := ethAmtFloat.Float64()
		ethAmt = ethAmt / 1e18
		fee, _ := feeFloat.Float64()
		fee = fee / 1e18

		state.TotalFeesEth += fee

		if t.IsBuy {
			state.Holdings.Add(state.Holdings, tokenAmt)
			state.TotalCostEth += ethAmt
			state.TotalBoughtEth += ethAmt
		} else {
			// Sell
			holdingsFloat, _ := new(big.Float).SetInt(state.Holdings).Float64()
			soldFloat, _ := new(big.Float).SetInt(tokenAmt).Float64()

			if holdingsFloat > 0 {
				costBasis := state.TotalCostEth / holdingsFloat
				if soldFloat > holdingsFloat {
					soldFloat = holdingsFloat
				}
				costOfSold := soldFloat * costBasis
				state.RealizedPnLEth += ethAmt - costOfSold

				state.Holdings.Sub(state.Holdings, tokenAmt)
				if state.Holdings.Sign() < 0 {
					state.Holdings.SetInt64(0)
				}
				state.TotalCostEth -= costOfSold
				if state.TotalCostEth < 0 {
					state.TotalCostEth = 0
				}
			} else {
				state.RealizedPnLEth += ethAmt
			}
			state.TotalSoldEth += ethAmt
		}
	}

	type UserPositionResponse struct {
		TokenAddress   string  `json:"token_address"`
		TokenName      string  `json:"token_name"`
		TokenSymbol    string  `json:"token_symbol"`
		ImageUrl       string  `json:"image_url"`
		Balance        string  `json:"balance"`
		AvgBuyPrice    float64 `json:"avg_buy_price"`
		CurrentPrice   float64 `json:"current_price"`
		RealizedPnL    float64 `json:"realized_pnl"`
		UnrealizedPnL  float64 `json:"unrealized_pnl"`
		TotalBoughtEth float64 `json:"total_bought_eth"`
		TotalSoldEth   float64 `json:"total_sold_eth"`
	}

	var positions []UserPositionResponse
	var totalRealized float64
	var totalUnrealized float64
	var totalVolume float64
	var totalFees float64
	var wins int
	var losses int

	for tokAddr, state := range states {
		tok, exists := tokenMap[tokAddr]
		if !exists {
			continue
		}

		currentPrice := priceMap[tokAddr]
		holdingsFloat, _ := new(big.Float).SetInt(state.Holdings).Float64()
		holdingsFloatVal := holdingsFloat / 1e18

		currentValuation := holdingsFloatVal * currentPrice
		unrealizedPnL := 0.0
		if holdingsFloatVal > 0 {
			unrealizedPnL = currentValuation - state.TotalCostEth
		}

		avgBuyPrice := 0.0
		if holdingsFloatVal > 0 {
			avgBuyPrice = state.TotalCostEth / holdingsFloatVal
		}

		pos := UserPositionResponse{
			TokenAddress:   tokAddr,
			TokenName:      tok.Name,
			TokenSymbol:    tok.Symbol,
			ImageUrl:       tok.ImageUrl,
			Balance:        state.Holdings.String(),
			AvgBuyPrice:    avgBuyPrice,
			CurrentPrice:   currentPrice,
			RealizedPnL:    state.RealizedPnLEth,
			UnrealizedPnL:  unrealizedPnL,
			TotalBoughtEth: state.TotalBoughtEth,
			TotalSoldEth:   state.TotalSoldEth,
		}

		positions = append(positions, pos)

		totalRealized += state.RealizedPnLEth
		totalUnrealized += unrealizedPnL
		totalVolume += state.TotalBoughtEth + state.TotalSoldEth
		totalFees += state.TotalFeesEth

		// Win/Loss calculation for Winrate
		posTotalPnL := state.RealizedPnLEth + unrealizedPnL
		if posTotalPnL > 0.00001 { // threshold to avoid tiny precision wins
			wins++
		} else if posTotalPnL < -0.00001 {
			losses++
		}
	}

	winrate := 0.0
	if wins+losses > 0 {
		winrate = (float64(wins) / float64(wins+losses)) * 100.0
	}

	c.JSON(http.StatusOK, gin.H{
		"realized_pnl_eth":   totalRealized,
		"unrealized_pnl_eth": totalUnrealized,
		"total_volume_eth":   totalVolume,
		"total_fees_eth":     totalFees,
		"winrate":            winrate,
		"positions":          positions,
	})
}

// GetUserTrades returns the last 50 trades of a specific user with token name and symbol details
func GetUserTrades(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))

	type Result struct {
		ID            uint      `json:"id"`
		TokenAddress  string    `json:"token_address"`
		TxHash        string    `json:"tx_hash"`
		BlockNumber   uint64    `json:"block_number"`
		Timestamp     time.Time `json:"timestamp"`
		IsBuy         bool      `json:"is_buy"`
		BuyerOrSeller string    `json:"buyer_or_seller"`
		TokenAmount   string    `json:"token_amount"`
		EthAmount     string    `json:"eth_amount"`
		Fee           string    `json:"fee"`
		TokenName     string    `json:"token_name"`
		TokenSymbol   string    `json:"token_symbol"`
	}

	var results []Result
	err := database.DB.Table("trades").
		Select("trades.*, tokens.name as token_name, tokens.symbol as token_symbol").
		Joins("left join tokens on tokens.address = trades.token_address").
		Where("trades.buyer_or_seller = ?", address).
		Order("trades.timestamp desc").
		Limit(50).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query user trades"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetRecentTrades returns the last 50 trades of the entire platform with token name and symbol details
func GetRecentTrades(c *gin.Context) {
	type Result struct {
		ID            uint      `json:"id"`
		TokenAddress  string    `json:"token_address"`
		TxHash        string    `json:"tx_hash"`
		BlockNumber   uint64    `json:"block_number"`
		Timestamp     time.Time `json:"timestamp"`
		IsBuy         bool      `json:"is_buy"`
		BuyerOrSeller string    `json:"buyer_or_seller"`
		TokenAmount   string    `json:"token_amount"`
		EthAmount     string    `json:"eth_amount"`
		Fee           string    `json:"fee"`
		TokenName     string    `json:"token_name"`
		TokenSymbol   string    `json:"token_symbol"`
	}

	var results []Result
	err := database.DB.Table("trades").
		Select("trades.*, tokens.name as token_name, tokens.symbol as token_symbol").
		Joins("left join tokens on tokens.address = trades.token_address").
		Order("trades.timestamp desc").
		Limit(50).
		Scan(&results).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query recent trades"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetTokenAudit calculates safety and audit metrics for a specific token
func GetTokenAudit(c *gin.Context) {
	tokenAddress := strings.ToLower(c.Param("address"))

	// 1. Get the token from database to find the creator
	var token database.Token
	if err := database.DB.Where("address = ?", tokenAddress).First(&token).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	cfg := config.LoadConfig()
	factoryAddr := strings.ToLower(cfg.FactoryAddress)

	// 2. Fetch all trades for this token to calculate balances
	var trades []database.Trade
	if err := database.DB.Where("token_address = ?", tokenAddress).Order("block_number asc, timestamp asc").Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query trades"})
		return
	}

	// 3. Compute balances per wallet
	balances := make(map[string]*big.Int)
	var totalFeesFloat float64

	// Track first block for sniper/bundler calculation
	var firstBlock uint64
	hasTrades := len(trades) > 0
	if hasTrades {
		firstBlock = trades[0].BlockNumber
	}

	var bundlersTokens *big.Int = big.NewInt(0)
	var snipersTokens *big.Int = big.NewInt(0)

	for _, t := range trades {
		buyerOrSeller := strings.ToLower(t.BuyerOrSeller)
		
		// Accumulate fees
		if f, ok := new(big.Float).SetString(t.Fee); ok {
			fFloat, _ := f.Float64()
			totalFeesFloat += fFloat / 1e18
		}

		// Calculate balance change
		amt, ok := new(big.Int).SetString(t.TokenAmount, 10)
		if !ok {
			continue
		}

		bal, exists := balances[buyerOrSeller]
		if !exists {
			bal = big.NewInt(0)
			balances[buyerOrSeller] = bal
		}

		if t.IsBuy {
			bal.Add(bal, amt)
			
			// Sniper / Bundler check (only on buys)
			if t.BlockNumber == firstBlock {
				// Bundler: bought in the exact first block of trading
				bundlersTokens.Add(bundlersTokens, amt)
			} else if t.BlockNumber > firstBlock && t.BlockNumber <= firstBlock+3 {
				// Sniper: bought in blocks +1, +2, +3
				snipersTokens.Add(snipersTokens, amt)
			}
		} else {
			bal.Sub(bal, amt)
			if bal.Sign() < 0 {
				bal.SetInt64(0)
			}
		}
	}

	// Exclude factory and zero/negative balances
	var activeBalances []float64
	var devBalanceFloat float64 = 0.0

	for addr, bal := range balances {
		if addr == factoryAddr || addr == "" {
			continue
		}
		
		balFloat, _ := new(big.Float).SetInt(bal).Float64()
		if balFloat <= 0 {
			continue
		}

		activeBalances = append(activeBalances, balFloat)

		if addr == strings.ToLower(token.Creator) {
			devBalanceFloat = balFloat
		}
	}

	holdersCount := len(activeBalances)

	// Sort active balances descending to find top 10
	sort.Slice(activeBalances, func(i, j int) bool {
		return activeBalances[i] > activeBalances[j]
	})

	var top10Sum float64
	for i := 0; i < len(activeBalances) && i < 10; i++ {
		top10Sum += activeBalances[i]
	}

	// Total supply of standard pump.fun tokens is 1 Billion (10^9 * 10^18 wei)
	totalSupply := 1000000000.0 // 1,000,000,000 tokens

	top10Percent := (top10Sum / 1e18 / totalSupply) * 100.0
	devPercent := (devBalanceFloat / 1e18 / totalSupply) * 100.0
	
	bundlersFloat, _ := new(big.Float).SetInt(bundlersTokens).Float64()
	bundlersPercent := (bundlersFloat / 1e18 / totalSupply) * 100.0

	snipersFloat, _ := new(big.Float).SetInt(snipersTokens).Float64()
	snipersPercent := (snipersFloat / 1e18 / totalSupply) * 100.0

	// Safe bounds check
	if top10Percent > 100.0 { top10Percent = 100.0 }
	if devPercent > 100.0 { devPercent = 100.0 }
	if bundlersPercent > 100.0 { bundlersPercent = 100.0 }
	if snipersPercent > 100.0 { snipersPercent = 100.0 }

	c.JSON(http.StatusOK, gin.H{
		"holders_count":    holdersCount,
		"top_10_percent":   top10Percent,
		"dev_percent":      devPercent,
		"snipers_percent":  snipersPercent,
		"bundlers_percent": bundlersPercent,
		"total_fees_eth":   totalFeesFloat,
	})
}

// GetTokenHolders calculates and returns the token holders sorted by balance
func GetTokenHolders(c *gin.Context) {
	tokenAddress := strings.ToLower(c.Param("address"))

	// 1. Get the token from database to find the creator
	var token database.Token
	if err := database.DB.Where("address = ?", tokenAddress).First(&token).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Token not found"})
		return
	}

	cfg := config.LoadConfig()
	factoryAddr := strings.ToLower(cfg.FactoryAddress)

	// 2. Fetch all trades for this token to calculate balances
	var trades []database.Trade
	if err := database.DB.Where("token_address = ?", tokenAddress).Order("block_number asc, timestamp asc").Find(&trades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query trades"})
		return
	}

	// 3. Compute balances per wallet
	balances := make(map[string]*big.Int)
	for _, t := range trades {
		buyerOrSeller := strings.ToLower(t.BuyerOrSeller)
		amt, ok := new(big.Int).SetString(t.TokenAmount, 10)
		if !ok {
			continue
		}

		bal, exists := balances[buyerOrSeller]
		if !exists {
			bal = big.NewInt(0)
			balances[buyerOrSeller] = bal
		}

		if t.IsBuy {
			bal.Add(bal, amt)
		} else {
			bal.Sub(bal, amt)
			if bal.Sign() < 0 {
				bal.SetInt64(0)
			}
		}
	}

	// 4. Fetch user profiles to map addresses to custom usernames
	var users []database.User
	if err := database.DB.Find(&users).Error; err != nil {
		fmt.Printf("Failed to query user profiles: %v\n", err)
	}
	usernames := make(map[string]string)
	avatarUrls := make(map[string]string)
	for _, u := range users {
		usernames[strings.ToLower(u.Address)] = u.Username
		avatarUrls[strings.ToLower(u.Address)] = u.AvatarUrl
	}

	type HolderInfo struct {
		Address    string  `json:"address"`
		Username   string  `json:"username"`
		AvatarUrl  string  `json:"avatar_url"`
		Balance    string  `json:"balance"`
		Percentage float64 `json:"percentage"`
		IsCreator  bool    `json:"is_creator"`
	}

	var holders []HolderInfo
	totalSupplyFloat := 1000000000.0 // 1 Billion tokens

	for addr, bal := range balances {
		if bal.Sign() <= 0 {
			continue
		}

		// Skip factory address (bonding curve reserves)
		if addr == factoryAddr {
			continue
		}

		balFloat, _ := new(big.Float).SetInt(bal).Float64()
		balEth := balFloat / 1e18
		percentage := (balEth / totalSupplyFloat) * 100.0

		username := usernames[addr]
		avatarUrl := avatarUrls[addr]
		if username == "" {
			username = addr
		}

		holders = append(holders, HolderInfo{
			Address:    addr,
			Username:   username,
			AvatarUrl:  avatarUrl,
			Balance:    bal.String(),
			Percentage: percentage,
			IsCreator:  addr == strings.ToLower(token.Creator),
		})
	}

	// Sort holders descending by balance
	sort.Slice(holders, func(i, j int) bool {
		balI, _ := new(big.Int).SetString(holders[i].Balance, 10)
		balJ, _ := new(big.Int).SetString(holders[j].Balance, 10)
		return balI.Cmp(balJ) > 0
	})

	c.JSON(http.StatusOK, holders)
}



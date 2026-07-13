package indexer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"backend/contracts"
	"backend/database"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

// Factory Event Topics Hex
const (
	TopicBuy          = "0x00f93dbdb72854b6b6fb35433086556f2635fc83c37080c667496fecfa650fb4"
	TopicSell         = "0x01fbb57444511e3de5b26ac09ad6bec45c3f9a1e59dd4a0f2b13a240d18476ce"
	TopicMigrated     = "0xf2ea3ee6d4d03a11390cbbcd09097d9fe2d7efb1b2825c2b509415d2fb95a7ba"
	TopicTokenCreated = "0xd5d05a8421149c74fd223cfc823befb883babf9bf0b0e4d6bf9c8fdb70e59bb4"
)

type Indexer struct {
	client         *ethclient.Client
	db             *gorm.DB
	factoryAddress common.Address
	filterer       *contracts.SafePumpFactoryFilterer
	startBlock     uint64

	// Callbacks for live notifications
	OnTokenCreated func(token database.Token)
	OnTrade        func(trade database.Trade)
	OnMigrated     func(token database.Token)

	// Live controls whether to trigger live callbacks (only after catch-up is complete)
	Live bool
}

// NewIndexer instantiates a new indexer connected to Base L2 RPC and database
func NewIndexer(rpcUrl string, db *gorm.DB, factoryAddrStr string, startBlock uint64) (*Indexer, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC client: %w", err)
	}

	factoryAddress := common.HexToAddress(factoryAddrStr)
	filterer, err := contracts.NewSafePumpFactoryFilterer(factoryAddress, client)
	if err != nil {
		return nil, fmt.Errorf("failed to bind filterer: %w", err)
	}

	return &Indexer{
		client:         client,
		db:             db,
		factoryAddress: factoryAddress,
		filterer:       filterer,
		startBlock:     startBlock,
	}, nil
}

// Start begins the indexing loop
func (idx *Indexer) Start(ctx context.Context) {
	log.Printf("Starting block indexer for Factory at %s from block %d...", idx.factoryAddress.Hex(), idx.startBlock)

	// Run initial catch-up batch silently (Live is false by default)
	log.Println("Running initial sync catch-up with blockchain history...")
	err := idx.processBatch(ctx)
	if err != nil {
		log.Println("Indexer initial catch-up execution warning:", err)
	}

	// Mark indexer as live to enable WebSockets broadcasts
	idx.Live = true
	log.Println("Indexer is now synced and listening for live events.")

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping indexer...")
			return
		case <-ticker.C:
			err := idx.processBatch(ctx)
			if err != nil {
				log.Println("Indexer batch execution error:", err)
			}
		}
	}
}

// processBatch processes blocks in chunks to stay within RPC limits
func (idx *Indexer) processBatch(ctx context.Context) error {
	// 1. Get last processed block from database
	var lastBlockConfig database.SystemConfig
	err := idx.db.Where("key = ?", "last_processed_block").First(&lastBlockConfig).Error
	var lastBlock uint64
	if errors.Is(err, gorm.ErrRecordNotFound) {
		lastBlock = idx.startBlock
	} else if err != nil {
		return err
	} else {
		var parsed uint64
		_, err := fmt.Sscanf(lastBlockConfig.Value, "%d", &parsed)
		if err != nil {
			lastBlock = idx.startBlock
		} else {
			lastBlock = parsed
		}
	}

	// 2. Query current block on-chain
	currentBlock, err := idx.client.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get latest block number: %w", err)
	}

	if currentBlock < lastBlock {
		log.Printf("WARNING: Chain reset detected! Latest block on-chain is %d, but last processed block was %d. Resetting indexer to block 0...", currentBlock, lastBlock)
		lastBlock = 0
		if lastBlockConfig.Key != "" {
			lastBlockConfig.Value = "0"
			idx.db.Save(&lastBlockConfig)
		} else {
			idx.db.Create(&database.SystemConfig{Key: "last_processed_block", Value: "0"})
		}

		log.Println("Clearing database tables for clean sync on chain reset...")
		idx.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&database.Trade{})
		idx.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&database.Token{})
	}

	if lastBlock >= currentBlock {
		return nil
	}

	// Limit chunks to 2000 blocks to prevent RPC timeouts
	toBlock := currentBlock
	if currentBlock-lastBlock > 2000 {
		toBlock = lastBlock + 2000
	}

	log.Printf("Syncing logs from block %d to %d...", lastBlock+1, toBlock)

	// 3. Query logs for Factory
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(lastBlock + 1)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{idx.factoryAddress},
	}

	logs, err := idx.client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to filter logs: %w", err)
	}

	// Cache timestamps for blocks queried in this batch to optimize RPC calls
	blockTimeCache := make(map[uint64]time.Time)

	// 4. Process logs
	for _, l := range logs {
		if len(l.Topics) == 0 {
			continue
		}

		topic := l.Topics[0].Hex()

		// Get block timestamp
		blockTime, exists := blockTimeCache[l.BlockNumber]
		if !exists {
			header, err := idx.client.HeaderByNumber(ctx, big.NewInt(int64(l.BlockNumber)))
			if err != nil {
				return fmt.Errorf("failed to fetch block header for block %d: %w", l.BlockNumber, err)
			}
			blockTime = time.Unix(int64(header.Time), 0)
			blockTimeCache[l.BlockNumber] = blockTime
		}

		var createdToken *database.Token
		var trade *database.Trade
		var migratedToken *database.Token

		// Handle database write inside a single transaction to maintain consistency
		err = idx.db.Transaction(func(tx *gorm.DB) error {
			switch topic {
			case TopicTokenCreated:
				ev, err := idx.filterer.ParseTokenCreated(l)
				if err != nil {
					return err
				}
				// Pick a deterministic image based on token address bytes
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
				tokenAddrHex := ev.Token.Hex()
				imgIdx := 0
				if len(tokenAddrHex) > 0 {
					imgIdx = int(tokenAddrHex[len(tokenAddrHex)-1]) % len(defaultMemeImages)
				}
				imgUrl := defaultMemeImages[imgIdx]
				description := fmt.Sprintf("Bienvenue sur %s ($%s). Rejoignez notre communauté de détenteurs passionnés et aidez-nous à propulser ce projet vers la lune !", ev.Name, ev.Symbol)
				website := fmt.Sprintf("https://www.%s.io", strings.ToLower(ev.Symbol))
				twitter := fmt.Sprintf("https://x.com/%s_coin", strings.ToLower(ev.Symbol))
				telegram := fmt.Sprintf("https://t.me/%s_chat", strings.ToLower(ev.Symbol))

				token := database.Token{
					Address:     strings.ToLower(ev.Token.Hex()),
					Creator:     strings.ToLower(ev.Creator.Hex()),
					Name:        ev.Name,
					Symbol:      ev.Symbol,
					Description: description,
					ImageUrl:    imgUrl,
					Website:     website,
					Twitter:     twitter,
					Telegram:    telegram,
					TokensSold:  "0",
					EthRaised:   "0",
					Migrated:    false,
					PairAddress: "",
					CreatedAt:   blockTime,
					UpdatedAt:   blockTime,
				}

				// Check if there was an initial buy on contract
				caller, err := contracts.NewSafePumpFactoryCaller(idx.factoryAddress, idx.client)
				var chainTokensSold *big.Int
				var chainEthRaised *big.Int
				if err == nil {
					chainInfo, err := caller.Tokens(nil, ev.Token)
					if err == nil {
						token.TokensSold = chainInfo.TokensSold.String()
						token.EthRaised = chainInfo.EthRaised.String()
						chainTokensSold = chainInfo.TokensSold
						chainEthRaised = chainInfo.EthRaised
					} else {
						log.Println("Indexer error calling Tokens on-chain:", err)
					}
				} else {
					log.Println("Indexer error creating factory caller:", err)
				}

				// INSERT/UPDATE TOKEN FIRST
				var existing database.Token
				if err := tx.Where("address = ?", strings.ToLower(ev.Token.Hex())).First(&existing).Error; err == nil {
					// Record exists! Update on-chain fields but preserve metadata description and image_url
					existing.Creator = strings.ToLower(ev.Creator.Hex())
					existing.Name = ev.Name
					existing.Symbol = ev.Symbol
					existing.TokensSold = token.TokensSold
					existing.EthRaised = token.EthRaised
					existing.Migrated = false
					existing.PairAddress = ""
					existing.CreatedAt = blockTime
					existing.UpdatedAt = blockTime
					if existing.Description == "" {
						existing.Description = description
					}
					if existing.ImageUrl == "" {
						existing.ImageUrl = imgUrl
					}
					if existing.Website == "" {
						existing.Website = website
					}
					if existing.Twitter == "" {
						existing.Twitter = twitter
					}
					if existing.Telegram == "" {
						existing.Telegram = telegram
					}
					if err := tx.Save(&existing).Error; err != nil {
						return err
					}
					createdToken = &existing
				} else {
					// Record doesn't exist, create it
					if err := tx.Create(&token).Error; err != nil {
						return err
					}
					createdToken = &token
				}

				// INSERT TRADE SECOND (Safe from Foreign Key constraint!)
				if chainTokensSold != nil && chainTokensSold.Cmp(big.NewInt(0)) > 0 {
					var existingTrade database.Trade
					// Try to see if this trade is already logged to prevent duplicates
					if err := tx.Where("tx_hash = ?", strings.ToLower(l.TxHash.Hex())).First(&existingTrade).Error; err != nil {
						t := database.Trade{
							TokenAddress:   strings.ToLower(ev.Token.Hex()),
							TxHash:         strings.ToLower(l.TxHash.Hex()),
							BlockNumber:    l.BlockNumber,
							Timestamp:      blockTime,
							IsBuy:          true,
							BuyerOrSeller:  strings.ToLower(ev.Creator.Hex()),
							TokenAmount:    chainTokensSold.String(),
							EthAmount:      chainEthRaised.String(),
							Fee:            "0",
							CreatedAt:      blockTime,
						}
						if err := tx.Create(&t).Error; err != nil {
							return err
						}
						trade = &t
						log.Printf("[Live Feed] Logged Initial Buy for Creator during creation: Token=%s Creator=%s ETH=%s Tokens=%s", token.Address, token.Creator, token.EthRaised, token.TokensSold)
					}
				}
				return nil

			case TopicBuy:
				ev, err := idx.filterer.ParseBuy(l)
				if err != nil {
					return err
				}
				t := database.Trade{
					TokenAddress:  strings.ToLower(ev.Token.Hex()),
					TxHash:        strings.ToLower(l.TxHash.Hex()),
					BlockNumber:   l.BlockNumber,
					Timestamp:     blockTime,
					IsBuy:         true,
					BuyerOrSeller: strings.ToLower(ev.Buyer.Hex()),
					TokenAmount:   ev.TokenAmount.String(),
					EthAmount:     ev.EthAmount.String(),
					Fee:           ev.Fee.String(),
					CreatedAt:     blockTime,
				}
				if err := tx.Create(&t).Error; err != nil {
					return err
				}

				var token database.Token
				if err := tx.Where("address = ?", strings.ToLower(ev.Token.Hex())).First(&token).Error; err != nil {
					return err
				}

				sold, _ := new(big.Int).SetString(token.TokensSold, 10)
				raised, _ := new(big.Int).SetString(token.EthRaised, 10)

				sold.Add(sold, ev.TokenAmount)
				raised.Add(raised, ev.EthAmount)

				token.TokensSold = sold.String()
				token.EthRaised = raised.String()
				token.UpdatedAt = blockTime

				if err := tx.Save(&token).Error; err != nil {
					return err
				}
				trade = &t
				return nil

			case TopicSell:
				ev, err := idx.filterer.ParseSell(l)
				if err != nil {
					return err
				}
				t := database.Trade{
					TokenAddress:  strings.ToLower(ev.Token.Hex()),
					TxHash:        strings.ToLower(l.TxHash.Hex()),
					BlockNumber:   l.BlockNumber,
					Timestamp:     blockTime,
					IsBuy:         false,
					BuyerOrSeller: strings.ToLower(ev.Seller.Hex()),
					TokenAmount:   ev.TokenAmount.String(),
					EthAmount:     ev.EthAmount.String(),
					Fee:           ev.Fee.String(),
					CreatedAt:     blockTime,
				}
				if err := tx.Create(&t).Error; err != nil {
					return err
				}

				var token database.Token
				if err := tx.Where("address = ?", strings.ToLower(ev.Token.Hex())).First(&token).Error; err != nil {
					return err
				}

				sold, _ := new(big.Int).SetString(token.TokensSold, 10)
				raised, _ := new(big.Int).SetString(token.EthRaised, 10)

				sold.Sub(sold, ev.TokenAmount)
				raised.Sub(raised, ev.EthAmount)

				token.TokensSold = sold.String()
				token.EthRaised = raised.String()
				token.UpdatedAt = blockTime

				if err := tx.Save(&token).Error; err != nil {
					return err
				}
				trade = &t
				return nil

			case TopicMigrated:
				ev, err := idx.filterer.ParseMigrated(l)
				if err != nil {
					return err
				}

				var token database.Token
				if err := tx.Where("address = ?", strings.ToLower(ev.Token.Hex())).First(&token).Error; err != nil {
					return err
				}

				token.Migrated = true
				token.PairAddress = strings.ToLower(ev.Pair.Hex())
				token.UpdatedAt = blockTime

				if err := tx.Save(&token).Error; err != nil {
					return err
				}
				migratedToken = &token
				return nil
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("transaction execution failed for tx %s: %w", l.TxHash.Hex(), err)
		}

		// Trigger callbacks after successful transaction commit if indexer is live
		if idx.Live {
			if createdToken != nil && idx.OnTokenCreated != nil {
				idx.OnTokenCreated(*createdToken)
			}
			if trade != nil && idx.OnTrade != nil {
				idx.OnTrade(*trade)
			}
			if migratedToken != nil && idx.OnMigrated != nil {
				idx.OnMigrated(*migratedToken)
			}
		}
	}

	// 5. Update indexer progress
	config := database.SystemConfig{
		Key:   "last_processed_block",
		Value: fmt.Sprintf("%d", toBlock),
	}
	err = idx.db.Save(&config).Error
	if err != nil {
		return fmt.Errorf("failed to save indexer progress: %w", err)
	}

	log.Printf("Successfully synced blocks up to %d.", toBlock)
	return nil
}

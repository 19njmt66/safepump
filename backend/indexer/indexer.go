package indexer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
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
				token := database.Token{
					Address:     ev.Token.Hex(),
					Creator:     ev.Creator.Hex(),
					Name:        ev.Name,
					Symbol:      ev.Symbol,
					TokensSold:  "0",
					EthRaised:   "0",
					Migrated:    false,
					PairAddress: "",
					CreatedAt:   blockTime,
					UpdatedAt:   blockTime,
				}
				var existing database.Token
				if err := tx.Where("address = ?", ev.Token.Hex()).First(&existing).Error; err == nil {
					// Record exists! Update on-chain fields but preserve metadata description and image_url
					existing.Creator = ev.Creator.Hex()
					existing.Name = ev.Name
					existing.Symbol = ev.Symbol
					existing.TokensSold = "0"
					existing.EthRaised = "0"
					existing.Migrated = false
					existing.PairAddress = ""
					existing.CreatedAt = blockTime
					existing.UpdatedAt = blockTime
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
				return nil

			case TopicBuy:
				ev, err := idx.filterer.ParseBuy(l)
				if err != nil {
					return err
				}
				t := database.Trade{
					TokenAddress:  ev.Token.Hex(),
					TxHash:        l.TxHash.Hex(),
					BlockNumber:   l.BlockNumber,
					Timestamp:     blockTime,
					IsBuy:         true,
					BuyerOrSeller: ev.Buyer.Hex(),
					TokenAmount:   ev.TokenAmount.String(),
					EthAmount:     ev.EthAmount.String(),
					Fee:           ev.Fee.String(),
					CreatedAt:     blockTime,
				}
				if err := tx.Create(&t).Error; err != nil {
					return err
				}

				var token database.Token
				if err := tx.Where("address = ?", ev.Token.Hex()).First(&token).Error; err != nil {
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
					TokenAddress:  ev.Token.Hex(),
					TxHash:        l.TxHash.Hex(),
					BlockNumber:   l.BlockNumber,
					Timestamp:     blockTime,
					IsBuy:         false,
					BuyerOrSeller: ev.Seller.Hex(),
					TokenAmount:   ev.TokenAmount.String(),
					EthAmount:     ev.EthAmount.String(),
					Fee:           ev.Fee.String(),
					CreatedAt:     blockTime,
				}
				if err := tx.Create(&t).Error; err != nil {
					return err
				}

				var token database.Token
				if err := tx.Where("address = ?", ev.Token.Hex()).First(&token).Error; err != nil {
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
				if err := tx.Where("address = ?", ev.Token.Hex()).First(&token).Error; err != nil {
					return err
				}

				token.Migrated = true
				token.PairAddress = ev.Pair.Hex()
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

		// Trigger callbacks after successful transaction commit
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

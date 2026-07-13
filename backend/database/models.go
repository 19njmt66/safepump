package database

import (
	"time"
)

// SystemConfig stores persistent key-value configuration and state (e.g. last indexed block)
type SystemConfig struct {
	Key       string    `gorm:"primaryKey;type:varchar(100)" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Token represents a meme coin deployed via the SafePump launchpad
type Token struct {
	Address      string    `gorm:"primaryKey;type:varchar(42)" json:"address"`
	Creator      string    `gorm:"type:varchar(42);index" json:"creator"`
	Name         string    `gorm:"type:varchar(255)" json:"name"`
	Symbol       string    `gorm:"type:varchar(50)" json:"symbol"`
	TokensSold   string    `gorm:"type:numeric" json:"tokens_sold"`
	EthRaised    string    `gorm:"type:numeric" json:"eth_raised"`
	Migrated     bool      `gorm:"default:false;index" json:"migrated"`
	PairAddress  string    `gorm:"type:varchar(42)" json:"pair_address"`
	Description  string    `gorm:"type:text" json:"description"`
	ImageUrl     string    `gorm:"type:text" json:"image_url"`
	Website      string    `gorm:"type:text" json:"website"`
	Twitter      string    `gorm:"type:text" json:"twitter"`
	Telegram     string    `gorm:"type:text" json:"telegram"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Trades       []Trade   `gorm:"foreignKey:TokenAddress" json:"-"`
}

// User represents a trader or creator profile on SafePump
type User struct {
	Address   string    `gorm:"primaryKey;type:varchar(42)" json:"address"`
	Username  string    `gorm:"type:varchar(100)" json:"username"`
	Bio       string    `gorm:"type:text" json:"bio"`
	AvatarUrl string    `gorm:"type:text" json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Trade represents a buy or sell transaction on the bonding curve
type Trade struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TokenAddress  string    `gorm:"type:varchar(42);index" json:"token_address"`
	TxHash        string    `gorm:"type:varchar(66)" json:"tx_hash"`
	BlockNumber   uint64    `json:"block_number"`
	Timestamp     time.Time `gorm:"index" json:"timestamp"`
	IsBuy         bool      `json:"is_buy"`
	BuyerOrSeller string    `gorm:"type:varchar(42);index" json:"buyer_or_seller"`
	TokenAmount   string    `gorm:"type:numeric" json:"token_amount"`
	EthAmount     string    `gorm:"type:numeric" json:"eth_amount"`
	Fee           string    `gorm:"type:numeric" json:"fee"`
	CreatedAt     time.Time `json:"created_at"`
}

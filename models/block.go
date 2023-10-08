package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Block struct {
	Height    int
	BlockHash string
	Timestamp time.Time
}

type BlockService struct {
	DB *sql.DB
}

func (bs *BlockService) Index(height int, hash string, timestamp time.Time) (*Block, error) {
	block := Block{
		Height:    height,
		BlockHash: hash,
		Timestamp: timestamp,
	}
	row := bs.DB.QueryRow(`
		INSERT INTO blocks (height, block_hash, timestamp)
		VALUES ($1, $2, $3) RETURNING block_hash`, height, hash, timestamp)
	err := row.Scan(&block.BlockHash)
	if err != nil {
		return nil, fmt.Errorf("error indexing block: %w", err)
	}
	return &block, nil
}

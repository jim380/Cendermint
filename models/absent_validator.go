package models

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

type AbsentValidator struct {
	BlockHeight    int
	ConsAddrBase64 string
}

type AbsentValidatorService struct {
	DB *sql.DB
}

func (abs *AbsentValidatorService) Index(height int, consAddrBase64 string) (*AbsentValidator, error) {
	zap.L().Info("Indexing absent validator", zap.Int("height", height), zap.String("consAddrBase64", consAddrBase64))
	abscentValidator := AbsentValidator{
		BlockHeight:    height,
		ConsAddrBase64: consAddrBase64,
	}
	row := abs.DB.QueryRow(`
		INSERT INTO absent_validators (block_height, cons_pub_address)
		VALUES ($1, $2) RETURNING cons_pub_address`, height, consAddrBase64)
	err := row.Scan(&abscentValidator.ConsAddrBase64)
	if err != nil {
		return nil, fmt.Errorf("error indexing absent validator: %w", err)
	}
	return &abscentValidator, nil
}

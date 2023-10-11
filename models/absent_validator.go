package models

import (
	"database/sql"
	"fmt"
)

type AbsentValidator struct {
	BlockHeight    int
	ConsHexAddress string
}

type AbsentValidatorService struct {
	DB *sql.DB
}

func (abs *AbsentValidatorService) Index(height int, consHex string) (*AbsentValidator, error) {
	abscentValidator := AbsentValidator{
		BlockHeight:    height,
		ConsHexAddress: consHex,
	}
	row := abs.DB.QueryRow(`
		INSERT INTO absent_validators (block_height, validator_cons_hex_addr)
		VALUES ($1, $2) RETURNING validator_cons_hex_addr`, height, consHex)
	err := row.Scan(&abscentValidator.ConsHexAddress)
	if err != nil {
		return nil, fmt.Errorf("error indexing absent validator: %w", err)
	}
	return &abscentValidator, nil
}

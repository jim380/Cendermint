package models

import (
	"database/sql"
	"fmt"
)

type Validator struct {
	ConsHexAddress string
	Moniker        string
}

type ValidatorService struct {
	DB *sql.DB
}

func (vs *ValidatorService) Index(consHexAddr, moniker string) (*Validator, error) {
	validator := Validator{
		ConsHexAddress: consHexAddr,
		Moniker:        moniker,
	}
	row := vs.DB.QueryRow(`
		INSERT INTO validators (cons_hex_address, moniker)
		VALUES ($1, $2) RETURNING cons_hex_address`, consHexAddr, moniker)
	err := row.Scan(&validator.ConsHexAddress)
	if err != nil {
		return nil, fmt.Errorf("error indexing validator: %w", err)
	}
	return &validator, nil
}

package controllers

import (
	"github.com/jim380/Cendermint/models"
)

type RestServices struct {
	BlockService           *models.BlockService
	ValidatorService       *models.ValidatorService
	AbsentValidatorService *models.AbsentValidatorService
}

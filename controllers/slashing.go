package controllers

import (
	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/types"
)

func (rs RestServices) GetSlashingInfo(cfg config.Config, rd *types.RESTData) {
	rs.SlashingService.GetSlashingParams(cfg, rd)
}

func (rs RestServices) GetSigningInfo(cfg config.Config, consAddr string, rd *types.RESTData) {
	rs.SlashingService.GetSigningInfo(cfg, consAddr, rd)
}

func (rs RestServices) GetCommitInfo(rd *types.RESTData, blockData types.Blocks, consHexAddr string) {
	rs.SlashingService.GetCommitInfo(rd, blockData, consHexAddr)
}

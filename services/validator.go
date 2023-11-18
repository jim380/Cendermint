package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"go.uber.org/zap"
)

type validatorsetsLegacy struct {
	Height string `json:"height"`

	Result struct {
		Block_Height string `json:"block_height"`
		Validators   []struct {
			ConsAddr         string                 `json:"address"`
			ConsPubKey       consPubKeyValSetLegacy `json:"pub_key"`
			ProposerPriority string                 `json:"proposer_priority"`
			VotingPower      string                 `json:"voting_power"`
		}
	}
}

type validatorsets struct {
	Block_Height string `json:"block_height"`
	Validators   []struct {
		ConsAddr         string           `json:"address"`
		ConsPubKey       consPubKeyValSet `json:"pub_key"`
		ProposerPriority string           `json:"proposer_priority"`
		VotingPower      string           `json:"voting_power"`
	} `json:"validators"`
}

type consPubKeyValSetLegacy struct {
	Type string `json:"type"`
	Key  string `json:"value"`
}

type consPubKeyValSet struct {
	Type string `json:"@type"`
	Key  string `json:"key"`
}

type Validator struct {
	ConsHexAddress string
	Moniker        string
}

type ValidatorService struct {
	DB *sql.DB
}

func (vs *ValidatorService) Init(db *sql.DB) {
	vs.DB = db
}

func (vs *ValidatorService) Index(consHexAddr, moniker string) (*Validator, error) {
	validator := Validator{
		ConsHexAddress: consHexAddr,
		Moniker:        moniker,
	}

	_, err := vs.DB.Exec(`
		INSERT INTO validators (cons_pub_address, moniker)
		VALUES ($1, $2) ON CONFLICT (cons_pub_address) DO NOTHING`, consHexAddr, moniker)
	if err != nil {
		return nil, fmt.Errorf("error indexing validator: %w", err)
	}

	return &validator, nil
}

func (vs *ValidatorService) GetValidatorInfo(cfg config.Config, currentBlockHeight int64, rd *types.RESTData) []string {
	var vSetsResultFinal map[string][]string

	if cfg.IsLegacySDKVersion() {
		var vSets, vSets2, vsetTest validatorsetsLegacy
		var vSetsResult map[string][]string = make(map[string][]string)
		var vSetsResult2 map[string][]string = make(map[string][]string)

		shouldRunPages := testPageLimit(cfg, currentBlockHeight, &vsetTest, 3)

		if shouldRunPages {
			runPages(cfg, currentBlockHeight, &vSets, vSetsResult, 1)
			runPages(cfg, currentBlockHeight, &vSets2, vSetsResult2, 2)

			for _, value := range vSets.Result.Validators {
				// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
				vSetsResult[value.ConsPubKey.Key] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
			}

			for _, value := range vSets2.Result.Validators {
				// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
				vSetsResult2[value.ConsPubKey.Key] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
			}

			vSetsResultFinal = mergeMap(vSetsResult, vSetsResult2)
			zap.L().Info("", zap.Bool("Success", true), zap.String("Active validators", fmt.Sprint(len(vSetsResultFinal))))
		} else {
			runPages(cfg, currentBlockHeight, &vSets, vSetsResult, 1)
			for _, value := range vSets.Result.Validators {
				// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
				vSetsResult[value.ConsPubKey.Key] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
			}
			vSetsResultFinal = vSetsResult
		}
	} else {
		var vSets validatorsets
		var vSetsResult map[string][]string = make(map[string][]string)

		route := rest.GetValidatorSetByHeightRoute(cfg)
		res, err := utils.HttpQuery(constants.RESTAddr + route + fmt.Sprint(currentBlockHeight))
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}

		json.Unmarshal(res, &vSets)

		if strings.Contains(string(res), "not found") {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
		}

		for _, value := range vSets.Validators {
			// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
			vSetsResult[value.ConsPubKey.Key] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
		}

		vSetsResultFinal = vSetsResult
		zap.L().Info("", zap.Bool("Success", true), zap.String("Active validators", fmt.Sprint(len(vSets.Validators))))
	}

	rd.Validatorsets = utils.Sort(vSetsResultFinal, 2) // sort by ProposerPriority
	for key, value := range rd.Validatorsets {
		zap.L().Debug("", zap.Bool("Success", true), zap.String(key, strings.Join(value, ", ")))
	}

	if len(rd.Validatorsets) == 0 {
		zap.L().Warn("", zap.Bool("Success", false), zap.String("err", "Validator set is empty"))
	}

	getValidator(cfg, rd)
	valInfo := locateValidatorInActiveSet(rd)
	return valInfo
}

// TO-DO if consumer chain, use cosmoshub's ConsPubKey
func locateValidatorInActiveSet(rd *types.RESTData) []string {
	valInfo, found := rd.Validatorsets[rd.Validator.ConsPubKey.Key]
	if !found {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", "Validator not found in the active set"))
	}
	return valInfo
}

func mergeMap(a map[string][]string, b map[string][]string) map[string][]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}

func runPages(cfg config.Config, currentBlockHeight int64, vSets *validatorsetsLegacy, vSetsResult map[string][]string, pages int) {
	route := rest.GetValidatorSetByHeightRoute(cfg)

	res, err := utils.HttpQuery(constants.RESTAddr + route + fmt.Sprint(currentBlockHeight))

	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	json.Unmarshal(res, &vSets)

	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	for _, value := range vSets.Result.Validators {
		// populate the validatorset map => [ConsPubKey][]string{ConsAddr, VotingPower, ProposerPriority}
		vSetsResult[value.ConsPubKey.Key] = []string{value.ConsAddr, value.VotingPower, value.ProposerPriority, "0"}
	}
}

func testPageLimit(cfg config.Config, currentBlockHeight int64, vSets *validatorsetsLegacy, maxPageNumber int64) bool {
	multiPagesSupported := true

	route := rest.GetValidatorSetByHeightRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + fmt.Sprint(currentBlockHeight) + "?page=2")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	json.Unmarshal(res, &vSets)

	if strings.Contains(string(res), "Internal error: page should be within") {
		zap.L().Info("", zap.String("warn", string(res)))
		multiPagesSupported = false
	}

	return multiPagesSupported
}

func getValidator(cfg config.Config, rd *types.RESTData) {
	var v types.Validators

	route := rest.GetValidatorByAddressRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + constants.OperAddr)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &v)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	rd.Validator = v.Validator
}

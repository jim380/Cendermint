package rest

import (
	"encoding/json"
	"strings"

	"go.uber.org/zap"

	utils "github.com/jim380/Cosmos-IE/utils"
)

type govInfo struct {
	TotalProposalCount  float64
	VotingProposalCount float64
}

type gov struct {
	Proposals  []proposal
	Pagination struct {
		Total string
	}
}

type proposal struct {
	Status string
}

func getGovInfo(log *zap.Logger) govInfo {
	var g gov
	var gi govInfo

	votingCount := 0

	res, _ := runRESTCommand("/cosmos/gov/v1beta1/proposals")
	json.Unmarshal(res, &g)
	if strings.Contains(string(res), "not found") {
		log.Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else {
		log.Info("\t", zap.Bool("Success", true), zap.String("Total Proposal Count", g.Pagination.Total))
	}

	for _, value := range g.Proposals {
		if value.Status == "PROPOSAL_STATUS_VOTING_PERIOD" {
			votingCount++
		}
	}

	gi.TotalProposalCount = utils.StringToFloat64(g.Pagination.Total)
	gi.VotingProposalCount = float64(votingCount)

	return gi
}

package rest

import (
	"encoding/json"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/jim380/Cendermint/config"
	utils "github.com/jim380/Cendermint/utils"
)

type govInfo struct {
	TotalProposalCount      float64
	VotingProposalCount     float64
	InVotingVotedCount      float64
	InVotingDidNotVoteCount float64
}

type gov struct {
	Proposals  []proposal `json:"proposals"`
	Pagination struct {
		Total string `json:"total"`
	} `json:"pagination"`
}

type proposal struct {
	ProposalID string `json:"proposal_id"`
	Status     string `json:"status"`
}

type voteInfo struct {
	Votes struct {
		Option string `json:"option"`
	} `json:"vote"`
}

func (rd *RESTData) getGovInfo(cfg config.Config) {
	var (
		g                  gov
		gi                 govInfo
		totalProposals     []string
		proposalsInVoting  []string
		inVotingVoted      int
		inVotingDidNotVote int
	)

	route := getProposalsRoute(cfg)
	res, err := HttpQuery(RESTAddr + route + "?pagination.limit=1000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &g)
	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	for _, value := range g.Proposals {
		totalProposals = append(totalProposals, value.ProposalID)
		if value.Status == "PROPOSAL_STATUS_VOTING_PERIOD" {
			proposalsInVoting = append(proposalsInVoting, value.ProposalID)
		}
	}
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Total proposals count", totalProposals[len(totalProposals)-1]))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Proposals in voting", strconv.Itoa(len(proposalsInVoting))))

	for _, value := range proposalsInVoting {
		var voteInfo voteInfo
		res, err := HttpQuery(RESTAddr + route + value + "/votes/" + utils.GetAccAddrFromOperAddr(OperAddr))
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		json.Unmarshal(res, &voteInfo)
		if voteInfo.Votes.Option != "" {
			inVotingVoted++
		} else {
			inVotingDidNotVote++
		}
	}
	totalProposalsCount, _ := strconv.ParseFloat(totalProposals[len(totalProposals)-1], 64)
	gi.TotalProposalCount = totalProposalsCount
	gi.VotingProposalCount = float64(len(proposalsInVoting))
	gi.InVotingVotedCount = float64(inVotingVoted)
	gi.InVotingDidNotVoteCount = float64(inVotingDidNotVote)

	rd.Gov = gi
}

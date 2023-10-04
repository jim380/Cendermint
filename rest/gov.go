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
		proposalsInVoting  []string
		inVotingVoted      int
		inVotingDidNotVote int
	)

	route := getProposalsRoute(cfg)
	res, err := HttpQuery(RESTAddr + route + "?pagination.limit=2000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	// Unmarshal the JSON response and check for errors
	if err := json.Unmarshal(res, &g); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}

	if strings.Contains(string(res), "not found") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	} else if strings.Contains(string(res), "error:") || strings.Contains(string(res), "error\\\":") {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", string(res)))
	}

	totalProposals := g.Proposals
	for _, value := range totalProposals {
		if value.Status == "PROPOSAL_STATUS_VOTING_PERIOD" {
			proposalsInVoting = append(proposalsInVoting, value.ProposalID)
		}
	}

	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Total proposals count", strconv.Itoa(len(totalProposals))))
	zap.L().Info("\t", zap.Bool("Success", true), zap.String("Proposals in voting", strconv.Itoa(len(proposalsInVoting))))

	for _, value := range proposalsInVoting {
		var voteInfo voteInfo
		res, err := HttpQuery(RESTAddr + route + value + "/votes/" + utils.GetAccAddrFromOperAddr(OperAddr))
		if err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}

		// Unmarshal the JSON response and check for errors
		if err := json.Unmarshal(res, &voteInfo); err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}

		if voteInfo.Votes.Option != "" {
			inVotingVoted++
		} else {
			inVotingDidNotVote++
		}
	}
	gi.TotalProposalCount = float64(len(totalProposals))
	gi.VotingProposalCount = float64(len(proposalsInVoting))
	gi.InVotingVotedCount = float64(inVotingVoted)
	gi.InVotingDidNotVoteCount = float64(inVotingDidNotVote)

	rd.Gov = gi
}

package services

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/jim380/Cendermint/config"
	"github.com/jim380/Cendermint/constants"
	"github.com/jim380/Cendermint/rest"
	"github.com/jim380/Cendermint/types"
	"github.com/jim380/Cendermint/utils"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
)

type GovService struct {
	DB *sql.DB
}

func (gs *GovService) Init(db *sql.DB) {
	gs.DB = db
}

func (rs *GovService) GetInfo(cfg config.Config, rd *types.RESTData) {
	var (
		g                  types.Gov
		gi                 types.GovInfo
		proposalsInVoting  []string
		inVotingVoted      int
		inVotingDidNotVote int
	)

	route := rest.GetProposalsRoute(cfg)
	res, err := utils.HttpQuery(constants.RESTAddr + route + "?pagination.limit=2000")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	if err := json.Unmarshal(res, &g); err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
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
		var voteInfo types.Vote
		accAddr, err := utils.GetAccAddrFromOperAddr(constants.OperAddr)
		if err != nil {
			log.Info().Msgf("Error getting account address from operator address: %v", err)
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		res, err := utils.HttpQuery(constants.RESTAddr + route + value + "/votes/" + accAddr)
		if err != nil {
			log.Info().Msgf("Error getting vote info: %v", err)
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		if err := json.Unmarshal(res, &voteInfo); err != nil {
			zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
		}
		if voteInfo.Vote.Options[0].Option != "" {
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

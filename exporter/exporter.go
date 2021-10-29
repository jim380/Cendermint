package exporter

import (
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/jim380/Cendermint/rest"
	utils "github.com/jim380/Cendermint/utils"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultGauges []prometheus.Gauge
	gaugesDenom   []prometheus.Gauge
)

func Run(chain string, log *zap.Logger) {
	denomList := getDenomList(chain)

	defaultGauges = make([]prometheus.Gauge, len(gaugesNamespaceList))
	gaugesDenom = make([]prometheus.Gauge, len(denomList)*3)

	// register nomal guages
	for i := 0; i < len(gaugesNamespaceList); i++ {
		defaultGauges[i] = utils.NewGauge("exporter", gaugesNamespaceList[i], "")
		prometheus.MustRegister(defaultGauges[i])
	}

	// register denom guages
	count := 0
	for i := 0; i < len(denomList)*3; i += 3 {

		gaugesDenom[i] = utils.NewGauge("exporter_balances", denomList[count], "")
		gaugesDenom[i+1] = utils.NewGauge("exporter_commission", denomList[count], "")
		gaugesDenom[i+2] = utils.NewGauge("exporter_rewards", denomList[count], "")
		prometheus.MustRegister(gaugesDenom[i])
		prometheus.MustRegister(gaugesDenom[i+1])
		prometheus.MustRegister(gaugesDenom[i+2])

		count++
	}

	// register label guages
	labels := []string{"chainId", "moniker", "operatorAddress", "accountAddress", "consHexAddress"}
	//	labels := []string{"chainId", "moniker", "operatorAddress", "accountAddress"}
	gaugesForLabel := utils.NewCounterVec("exporter", "labels", "", labels)

	prometheus.MustRegister(gaugesForLabel)

	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					//Error Log
					//panic("oops...something bad happened")
				}
				time.Sleep(5000 * time.Millisecond)
			}()

			var block rest.Blocks
			block.GetInfo()

			currentBlockHeight, _ := strconv.ParseInt(block.Block.Header.Height, 10, 64)
			if previousBlockHeight != currentBlockHeight {
				fmt.Println("")
				zap.L().Info("\t", zap.Bool("Success", true), zap.String("Block Height", fmt.Sprint(currentBlockHeight)))

				// fetch info from REST
				restData := rest.GetData(chain, currentBlockHeight, block, denomList[0])

				SetMetric(currentBlockHeight, restData, log)
				metricData := GetMetric()

				// set values for denom guages
				count := 0
				for i := 0; i < len(denomList)*3; i += 3 {
					for _, value := range metricData.Validator.Account.Balances {
						if value.Denom == denomList[count] {
							gaugesDenom[i].Set(utils.StringToFloat64(value.Amount))
						}
					}
					for _, value := range metricData.Validator.Account.Commission {
						if value.Denom == denomList[count] {
							gaugesDenom[i+1].Set(utils.StringToFloat64(value.Amount))
						}
					}
					for _, value := range metricData.Validator.Account.Rewards {
						if value.Denom == denomList[count] {
							gaugesDenom[i+2].Set(utils.StringToFloat64(value.Amount))
						}
					}
					count++
				}

				// set values for normal guages
				gaugesValue := [...]float64{
					float64(metricData.Network.BlockHeight),

					metricData.Network.Staking.NotBondedTokens,
					metricData.Network.Staking.BondedTokens,
					metricData.Network.Staking.TotalSupply,
					metricData.Network.Staking.BondedRatio,

					metricData.Network.Gov.TotalProposalCount,
					metricData.Network.Gov.VotingProposalCount,

					metricData.Validator.VotingPower,
					metricData.Validator.MinSelfDelegation,
					metricData.Validator.JailStatus,

					// metricData.Validator.Proposer.Ranking,
					// metricData.Validator.Proposer.Status,

					metricData.Validator.Delegation.Shares,
					metricData.Validator.Delegation.Ratio,
					metricData.Validator.Delegation.DelegatorCount,
					metricData.Validator.Delegation.Self,

					metricData.Validator.Commission.Rate,
					metricData.Validator.Commission.MaxRate,
					metricData.Validator.Commission.MaxChangeRate,
					// metricData.Validator.Commit.VoteType,
					metricData.Validator.Commit.PrecommitStatus,

					metricData.Network.Minting.Inflation,
					metricData.Network.Minting.ActualInflation,
				}

				for i := 0; i < len(gaugesNamespaceList); i++ {
					defaultGauges[i].Set(gaugesValue[i])
				}

				gaugesForLabel.WithLabelValues(metricData.Network.ChainID,
					metricData.Validator.Moniker,
					metricData.Validator.Address.Operator,
					metricData.Validator.Address.Account,
					metricData.Validator.Address.ConsensusHex,
				).Add(0)

			}

			previousBlockHeight = currentBlockHeight
		}()
	}
}

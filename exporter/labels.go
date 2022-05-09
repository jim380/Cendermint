package exporter

import (
	utils "github.com/jim380/Cendermint/utils"
	"github.com/prometheus/client_golang/prometheus"
)

type labelData struct {
	name   string
	labels []string
}

var counterVecs []prometheus.CounterVec

func registerLabels() []prometheus.CounterVec {
	for _, v := range labels {
		counterVec := utils.NewCounterVec("cendermint", v.name, "", v.labels)
		counterVecs = append(counterVecs, counterVec)
		prometheus.MustRegister(counterVec)
	}

	return counterVecs
}

// {"chain_id", "node_moniker", "node_id", "tm_version", "app_name", "binary_name", "app_version", "git_commit", "go_version", "sdk_version"}
func (metricData *metric) setNodeLabels(labels prometheus.CounterVec) {
	labels.WithLabelValues(
		metricData.Network.ChainID,
		metricData.Network.NodeInfo.Moniker,
		metricData.Network.NodeInfo.NodeID,
		metricData.Network.NodeInfo.TMVersion,
		metricData.Network.NodeInfo.AppName,
		metricData.Network.NodeInfo.Name,
		metricData.Network.NodeInfo.Version,
		metricData.Network.NodeInfo.GitCommit,
		metricData.Network.NodeInfo.GoVersion,
		metricData.Network.NodeInfo.SDKVersion,
	).Add(0)
}

func (metricData *metric) setAddrLabels(labels prometheus.CounterVec) {
	labels.WithLabelValues(
		metricData.Validator.Address.Operator,
		metricData.Validator.Address.Account,
		metricData.Validator.Address.ConsensusHex,
	).Add(0)
}

func (metricData *metric) setUpgradeLabels(labels prometheus.CounterVec) {
	labels.WithLabelValues(
		metricData.Upgrade.Name,
		metricData.Upgrade.Time,
		metricData.Upgrade.Height,
		metricData.Upgrade.Info,
	).Add(0)
}

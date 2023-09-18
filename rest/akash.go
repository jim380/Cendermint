package rest

import (
	"encoding/json"
	"strconv"

	"github.com/jim380/Cendermint/config"
	"go.uber.org/zap"
)

type akashInfo struct {
	// Deployments       akashDeployments
	TotalDeployments  int
	ActiveDeployments int
	ClosedDeployments int
}

type akashDeployments struct {
	Deployments []akashDeployment `json:"deployments"`
	Pagination  struct {
		Total string `json:"total"`
	} `json:"pagination"`
}

type akashDeployment struct {
	Deployment    `json:"deployment"`
	Groups        []Group `json:"groups"`
	EscrowAccount struct {
		Id struct {
			Scope string `json:"scope"`
			Xid   string `json:"xid"`
		} `json:"id"`
		Owner   string `json:"owner"`
		State   string `json:"state"`
		Balance struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balance"`
		Transferred struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"transferred"`
		SettledAt string `json:"settled_at"`
		Depositor string `json:"depositor"`
		Funds     struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"funds"`
	} `json:"escrow_account"`
}

type Deployment struct {
	DeploymentId `json:"deployment_id"`
	State        string `json:"state"`
	Version      string `json:"version"`
	CreatedAt    string `json:"created_at"`
}

type DeploymentId struct {
	Owner string `json:"owner"`
	Dseq  string `json:"dseq"`
}

type Group struct {
	GroupId   `json:"group_id"`
	State     string `json:"state"`
	GroupSpec `json:"group_spec"`
	CreatedAt string `json:"created_at"`
}

type GroupId struct {
	Owner string `json:"owner"`
	Dseq  string `json:"dseq"`
	Gseq  string `json:"gseq"`
}

type GroupSpec struct {
	Name string `json:"name"`
	// Requirements `json:"requirements"`
	Resources []struct {
		Resources struct {
			Cpu struct {
				Units struct {
					Val string `json:"Val"`
				} `json:"units"`
			} `json:"cpu"`
			Memory struct {
				Quantity struct {
					Val string `json:"Val"`
				} `json:"quantity"`
			} `json:"memory"`
			Storage struct {
				Name     string `json:"name"`
				Quantity struct {
					Val string `json:"Val"`
				} `json:"quantity"`
			} `json:"storage"`
			Endpoints []struct {
				Kind            string `json:"kind"`
				Sequence_number int    `json:"sequence_number"`
			} `json:"endpoints"`
		} `json:"resources"`
		Count string `json:"count"`
		Price struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"price"`
	} `json:"resources"`
}

func (rd *RESTData) getAkashDeployments(cfg config.Config) {
	if cfg.Chain.Chain != "akash" {
		return
	}
	var deployments, activeDeployments akashDeployments

	route := getDeploymentsRoute()
	res, err := HttpQuery(RESTAddr + route)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(res, &deployments)

	// get total deployments count
	totalDeploymentsCount, err := strconv.Atoi(deployments.Pagination.Total)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	rd.AkashInfo.TotalDeployments = totalDeploymentsCount

	// get active deployments count
	resActive, err := HttpQuery(RESTAddr + route + "?filters.state=active")
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	json.Unmarshal(resActive, &activeDeployments)

	activeDeploymentsCount, err := strconv.Atoi(activeDeployments.Pagination.Total)
	if err != nil {
		zap.L().Fatal("", zap.Bool("Success", false), zap.String("err", err.Error()))
	}
	rd.AkashInfo.ActiveDeployments = activeDeploymentsCount

	// get closed deployments count
	rd.AkashInfo.ClosedDeployments = totalDeploymentsCount - activeDeploymentsCount
}

package types

type AkashInfo struct {
	// Deployments       akashDeployments
	TotalDeployments  int
	ActiveDeployments int
	ClosedDeployments int
}

type AkashDeployments struct {
	Deployments []akashDeployment `json:"deployments"`
	Pagination  struct {
		Total string `json:"total"`
	} `json:"pagination"`
}

type akashDeployment struct {
	Deployment    `json:"deployment"`
	Groups        []Group `json:"groups"`
	EscrowAccount struct {
		ID struct {
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
	DeploymentID `json:"deployment_id"`
	State        string `json:"state"`
	Version      string `json:"version"`
	CreatedAt    string `json:"created_at"`
}

type DeploymentID struct {
	Owner string `json:"owner"`
	Dseq  string `json:"dseq"`
}

type Group struct {
	GroupID   `json:"group_id"`
	State     string `json:"state"`
	GroupSpec `json:"group_spec"`
	CreatedAt string `json:"created_at"`
}

type GroupID struct {
	Owner string `json:"owner"`
	Dseq  string `json:"dseq"`
	Gseq  string `json:"gseq"`
}

type GroupSpec struct {
	Name string `json:"name"`
	// Requirements `json:"requirements"`
	Resources []struct {
		Resources struct {
			CPU struct {
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
				Kind           string `json:"kind"`
				SequenceNumber int    `json:"sequence_number"`
			} `json:"endpoints"`
		} `json:"resources"`
		Count string `json:"count"`
		Price struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"price"`
	} `json:"resources"`
}

package akash

// for http metrics
type Deployments struct {
	TotalDeployments  int
	ActiveDeployments int
	ClosedDeployments int
}

// for http
type DeploymentsResponse struct {
	Deployments []Deployment `json:"deployments"`
	Pagination  Pagination   `json:"pagination"`
}

type Deployment struct {
	DeploymentDetails DeploymentDetails `json:"deployment"`
	Groups            []Group           `json:"groups"`
	EscrowAccount     EscrowAccount     `json:"escrow_account"`
}

type DeploymentDetails struct {
	DeploymentId DeploymentId `json:"deployment_id"`
	State        string       `json:"state"`
	Version      string       `json:"version"`
	CreatedAt    string       `json:"created_at"`
}

type DeploymentId struct {
	Owner string `json:"owner"`
	Dseq  string `json:"dseq"`
}

type Group struct {
	GroupId   GroupId   `json:"group_id"`
	State     string    `json:"state"`
	GroupSpec GroupSpec `json:"group_spec"`
	CreatedAt string    `json:"created_at"`
}

type GroupId struct {
	Owner string `json:"owner"`
	Dseq  string `json:"dseq"`
	Gseq  int    `json:"gseq"`
}

type GroupSpec struct {
	Name         string       `json:"name"`
	Requirements Requirements `json:"requirements"`
	Resources    []Resource   `json:"resources"`
}

type Requirements struct {
	SignedBy   SignedBy    `json:"signed_by"`
	Attributes []Attribute `json:"attributes"`
}

type SignedBy struct {
	AllOf []string `json:"all_of"`
	AnyOf []string `json:"any_of"`
}

type Resource struct {
	ResourceDetails ResourceDetails `json:"resource"`
	Count           string          `json:"count"`
	Price           Price           `json:"price"`
}

type ResourceDetails struct {
	ID        int        `json:"id"`
	CPU       CPU        `json:"cpu"`
	Memory    Memory     `json:"memory"`
	Storage   []Storage  `json:"storage"`
	GPU       GPU        `json:"gpu"`
	Endpoints []Endpoint `json:"endpoints"`
}

type CPU struct {
	Units Units `json:"units"`
}

type Memory struct {
	Quantity Units `json:"quantity"`
}

type Storage struct {
	Name     string `json:"name"`
	Quantity Units  `json:"quantity"`
}

type GPU struct {
	Units Units `json:"units"`
}

type Units struct {
	Val string `json:"val"`
}

type Endpoint struct {
	Kind            string `json:"kind"`
	Sequence_number int    `json:"sequence_number"`
}

type Price struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type EscrowAccount struct {
	ID          ID     `json:"id"`
	Owner       string `json:"owner"`
	State       string `json:"state"`
	Balance     Amount `json:"balance"`
	Transferred Amount `json:"transferred"`
	SettledAt   string `json:"settled_at"`
	Depositor   string `json:"depositor"`
	Funds       Amount `json:"funds"`
}

type ID struct {
	Scope string `json:"scope"`
	Xid   string `json:"xid"`
}

type Amount struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

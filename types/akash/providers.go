package akash

// for http
type ProvidersResponse struct {
	Providers  []Provider `json:"providers"`
	Pagination Pagination `json:"pagination"`
}

type Provider struct {
	Owner      string      `json:"owner"`
	HostURI    string      `json:"host_uri"`
	Attributes []Attribute `json:"attributes"`
	Info       Info        `json:"info"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Info struct {
	Email   string `json:"email"`
	Website string `json:"website"`
}

type Pagination struct {
	NextKey string `json:"next_key"`
	Total   string `json:"total"`
}

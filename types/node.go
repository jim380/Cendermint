package types

type NodeInfo struct {
	Default     DefaultInfo `json:"default_node_info"`
	Application appVersion  `json:"application_version"`
}

type DefaultInfo struct {
	NodeID    string `json:"default_node_id"`
	TMVersion string `json:"version"`
	Moniker   string `json:"moniker"`
}

type appVersion struct {
	AppName    string `json:"name"`
	Name       string `json:"app_name"`
	Version    string `json:"version"`
	GitCommit  string `json:"git_commit"`
	GoVersion  string `json:"go_version"`
	SDKVersion string `json:"cosmos_sdk_version"`
}

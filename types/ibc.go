package types

type IbcInfo struct {
	IbcChannelInfo
	IbcConnectionInfo
}

type IbcChannelInfo struct {
	OpenChannels int
	IBCChannels  ibcChannels `json:"channels"`
	Pagination   struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
	Height struct {
		RevisionNumber string `json:"revision_number"`
		RevisionHeight string `json:"revision_height"`
	}
}

type IbcConnectionInfo struct {
	OpenConnections int
	IBConnections   ibcConnections `json:"connections"`
	Pagination      struct {
		NextKey string `json:"next_key"`
		Total   string `json:"total"`
	}
	Height struct {
		RevisionNumber string `json:"revision_number"`
		RevisionHeight string `json:"revision_height"`
	}
}

type ibcChannels []struct {
	State          string              `json:"state"`
	Ordering       string              `json:"ordering"`
	Counterparty   counterpartyChannel `json:"counterparty"`
	ConnectionHops []struct {
		string
	}
	Version   string `json:"version"`
	PortID    string `json:"port_id"`
	ChannelID string `json:"channel_id"`
}

type ibcConnections []struct {
	ID           string                 `json:"id"`
	ClientID     string                 `json:"client_id"`
	Versions     connectionVersions     `json:"versions"`
	State        string                 `json:"state"`
	Counterparty counterpartyConnection `json:"counterparty"`
	DelayPeriod  string                 `json:"delay_period"`
}

type connectionVersions []struct {
	Identifier string   `json:"identifier"`
	Features   []string `json:"features"`
}

type counterpartyChannel struct {
	PortID    string `json:"port_id"`
	ChannelID string `json:"channel_id"`
}

type counterpartyConnection struct {
	ClientID     string `json:"client_id"`
	ConnectionID string `json:"connection_id"`
	Prefix       struct {
		KeyPrefix string `json:"key_prefix"`
	}
}

package akash

type AuditorsResponse struct {
	Providers []struct {
		Auditor string `json:"auditor"`
	} `json:"providers"`
	Pagination Pagination `json:"pagination"`
}

package retailcrm

// VersionResponse return available API versions
type VersionResponse struct {
	Success  bool     `json:"success,omitempty"`
	Versions []string `json:"versions,brackets,omitempty"`
}

// CredentialResponse return available API methods
type CredentialResponse struct {
	Success        bool     `json:"success,omitempty"`
	Credentials    []string `json:"credentials,brackets,omitempty"`
	SiteAccess     string   `json:"siteAccess,omitempty"`
	SitesAvailable []string `json:"sitesAvailable,brackets,omitempty"`
}

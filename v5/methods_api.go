package v5

import (
	"encoding/json"
	"fmt"
)

// ApiVersions get available API versions
func (c *Client) ApiVersions() (*VersionResponse, int, error) {
	var resp VersionResponse
	data, status, err := c.GetRequest(fmt.Sprintf("%s/api-versions", unversionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ApiCredentials get available API methods
func (c *Client) ApiCredentials() (*CredentialResponse, int, error) {
	var resp CredentialResponse
	data, status, err := c.GetRequest(fmt.Sprintf("%s/credentials", unversionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

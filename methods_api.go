package retailcrm

import (
	"fmt"
	"encoding/json"
	"net/http"
	"errors"
)

// ApiVersions get available API versions
func (c *Client) ApiVersions() (*VersionResponse, int, error) {
	var resp VersionResponse
	data, status, err := c.getRequest(fmt.Sprintf("/api-versions"))

	if err != nil {
		return &resp, status, err
	}

	if status >= http.StatusBadRequest {
		return &resp, status, errors.New(fmt.Sprintf("HTTP request error. Status code: %d.\n", status))
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ApiCredentials get available API methods
func (c *Client) ApiCredentials() (*CredentialResponse, int, error) {
	var resp CredentialResponse
	data, status, err := c.getRequest(fmt.Sprintf("/credentials"))

	if err != nil {
		return &resp, status, err
	}

	if status >= http.StatusBadRequest {
		return &resp, status, errors.New(fmt.Sprintf("HTTP request error. Status code: %d.\n", status))
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}



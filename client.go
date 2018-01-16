package retailcrm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiPrefix = "/api"
)

// Client type
type Client struct {
	Url        string
	Key        string
	Version    string
	httpClient *http.Client
}

// ErrorResponse type
type ErrorResponse struct {
	ErrorMsg string            `json:"errorMsg,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

// New initalize client
func New(url string, key string, version string) *Client {
	return &Client {
		url,
		key,
		version,
		&http.Client{Timeout: 20 * time.Second},
	}
}

func buildUrl(url string, version string) string {
	var versionPlaceholder string = "/" + version
	if version == "" {
		versionPlaceholder = ""
	}

	var requestUrl = fmt.Sprintf("%s%s%s", url, apiPrefix, versionPlaceholder)

	return requestUrl
}

func (c *Client) getRequest(urlWithParameters string) ([]byte, int, error) {
	var res []byte
    var reqUrl = buildUrl(c.Url, c.Version)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", reqUrl , urlWithParameters), nil)
	if err != nil {
		return res, 0, err
	}

	req.Header.Set("X-API-KEY", c.Key)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return res, 0, err
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return res, resp.StatusCode, errors.New(fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode))
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		return res, 0, err
	}

	return res, resp.StatusCode, nil
}

func (c *Client) postRequest(url string, postParams url.Values) ([]byte, int, error) {
	var res []byte
	var reqUrl = buildUrl(c.Url, c.Version)

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", reqUrl, url),
		strings.NewReader(postParams.Encode()),
	)
	if err != nil {
		return res, 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-API-KEY", c.Key)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return res, 0, err
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return res, resp.StatusCode, errors.New(fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode))
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		return res, 0, err
	}

	return res, resp.StatusCode, nil
}

func buildRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	return res, nil
}

// ErrorResponse method
func (c *Client) ErrorResponse(data []byte) (*ErrorResponse, error) {
	var resp ErrorResponse
	err := json.Unmarshal(data, &resp)

	return &resp, err
}

// CheckBy select identifier type
func CheckBy(by string) string  {
	var context = "id"

	if by != "id" {
		context = "externalId"
	}

	return context
}

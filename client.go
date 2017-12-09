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
	ApiPrefix = "/api/v5"
)

type Client struct {
	Url        string
	apiKey     string
	httpClient *http.Client
}

type ErrorResponse struct {
	ErrorMsg string            `json:"errorMsg,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

func New(url string, apiKey string) *Client {
	return &Client{
		url,
		apiKey,
		&http.Client{Timeout: 20 * time.Second},
	}
}

func (r *Client) GetRequest(urlWithParameters string) ([]byte, int, error) {
	var res []byte

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", r.Url, ApiPrefix, urlWithParameters), nil)
	if err != nil {
		return res, 0, err
	}

	req.Header.Set("X-API-KEY", r.apiKey)

	resp, err := r.httpClient.Do(req)
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

func (r *Client) PostRequest(url string, postParams url.Values) ([]byte, int, error) {
	var res []byte

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s%s", r.Url, ApiPrefix, url),
		strings.NewReader(postParams.Encode()),
	)
	if err != nil {
		return res, 0, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-API-KEY", r.apiKey)

	resp, err := r.httpClient.Do(req)
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

func (r *Client) ErrorResponse(data []byte) (*ErrorResponse, error) {
	var resp ErrorResponse
	err := json.Unmarshal(data, &resp)

	return &resp, err
}

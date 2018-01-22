package v5

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	versionedPrefix   = "/api/v5"
	unversionedPrefix = "/api"
)

// Client type
type Client struct {
	Url        string
	Key        string
	httpClient *http.Client
}

// ErrorResponse type
type ErrorResponse struct {
	ErrorMsg string            `json:"errorMsg,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

// New initalize client
func New(url string, key string) *Client {
	return &Client{
		url,
		key,
		&http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) getRequest(urlWithParameters string) ([]byte, int, error) {
	var res []byte

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.Url, urlWithParameters), nil)
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

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", c.Url, url),
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

// checkBy select identifier type
func checkBy(by string) string {
	var context = "id"

	if by != "id" {
		context = "externalId"
	}

	return context
}

// ApiVersions get available API versions
func (c *Client) ApiVersions() (*VersionResponse, int, error) {
	var resp VersionResponse
	data, status, err := c.getRequest(fmt.Sprintf("%s/api-versions", unversionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ApiCredentials get available API methods
func (c *Client) ApiCredentials() (*CredentialResponse, int, error) {
	var resp CredentialResponse
	data, status, err := c.getRequest(fmt.Sprintf("%s/credentials", unversionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Customer get method
func (c *Client) Customer(id, by, site string) (*CustomerResponse, int, error) {
	var resp CustomerResponse
	var context = checkBy(by)

	fw := CustomerGetFilter{context, site}
	params, _ := query.Values(fw)
	data, status, err := c.getRequest(fmt.Sprintf("%s/customers/%s?%s", versionedPrefix, id, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Customers list method
func (c *Client) Customers(filter CustomersFilter, limit, page int) (*CustomersResponse, int, error) {
	var resp CustomersResponse

	if limit == 0 {
		limit = 20
	}

	if page == 0 {
		page = 1
	}

	fw := CustomersParameters{filter, limit, page}
	params, _ := query.Values(fw)

	data, status, err := c.getRequest(fmt.Sprintf("%s/customers?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

func (c *Client) CustomerCreate(customer Customer, site ...string) (*CustomerChangeResponse, int, error) {
	var resp CustomerChangeResponse
	customerJson, _ := json.Marshal(&customer)

	p := url.Values{
		"customer": {string(customerJson[:])},
	}

	if len(site) > 0 {
		s := site[0]

		if s != "" {
			p.Add("site", s)
		}
	}

	data, status, err := c.postRequest(fmt.Sprintf("%s/customers/create", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

func (c *Client) CustomerEdit(customer Customer, by string, site ...string) (*CustomerChangeResponse, int, error) {
	var resp CustomerChangeResponse
	var uid = strconv.Itoa(customer.Id)
	var context = checkBy(by)

	if context == "externalId" {
		uid = customer.ExternalId
	}

	customerJson, _ := json.Marshal(&customer)

	p := url.Values{
		"by":       {string(context)},
		"customer": {string(customerJson[:])},
	}

	if len(site) > 0 {
		s := site[0]

		if s != "" {
			p.Add("site", s)
		}
	}

	data, status, err := c.postRequest(fmt.Sprintf("%s/customers/%s/edit", versionedPrefix, uid), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

func (c *Client) CustomersUpload() {

}

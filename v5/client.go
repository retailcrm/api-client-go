package v5

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	versionedPrefix   = "/api/v5"
	unversionedPrefix = "/api"
)

// New initalize client
func New(url string, key string) *Client {
	return &Client{
		url,
		key,
		&http.Client{Timeout: 20 * time.Second},
	}
}

func (c *Client) GetRequest(urlWithParameters string) ([]byte, int, error) {
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

func (c *Client) PostRequest(url string, postParams url.Values) ([]byte, int, error) {
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

// checkBy select identifier type
func checkBy(by string) string {
	var context = "id"

	if by != "id" {
		context = "externalId"
	}

	return context
}

// fillSite add site code to parameters if present
func fillSite(p *url.Values, site []string) {
	if len(site) > 0 {
		s := site[0]

		if s != "" {
			p.Add("site", s)
		}
	}
}

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

// Customer get method
func (c *Client) Customer(id, by, site string) (*CustomerResponse, int, error) {
	var resp CustomerResponse
	var context = checkBy(by)

	fw := CustomerRequest{context, site}
	params, _ := query.Values(fw)
	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers/%s?%s", versionedPrefix, id, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Customers list method
func (c *Client) Customers(parameters CustomersRequest) (*CustomersResponse, int, error) {
	var resp CustomersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomerCreate method
func (c *Client) CustomerCreate(customer Customer, site ...string) (*CustomerChangeResponse, int, error) {
	var resp CustomerChangeResponse
	customerJson, _ := json.Marshal(&customer)

	p := url.Values{
		"customer": {string(customerJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/create", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomerEdit method
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

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/%s/edit", versionedPrefix, uid), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersUpload method
func (c *Client) CustomersUpload(customers []Customer, site ...string) (*CustomersUploadResponse, int, error) {
	var resp CustomersUploadResponse

	uploadJson, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(uploadJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/upload", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersFixExternalIds method
func (c *Client) CustomersFixExternalIds(customers []IdentifiersPair) (*SucessfulResponse, int, error) {
	var resp SucessfulResponse

	customersJson, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(customersJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/fix-external-ids", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersHistory method
func (c *Client) CustomersHistory(parameters CustomersHistoryRequest) (*CustomersHistoryResponse, int, error) {
	var resp CustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers/history?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Order get method
func (c *Client) Order(id, by, site string) (*OrderResponse, int, error) {
	var resp OrderResponse
	var context = checkBy(by)

	fw := OrderRequest{context, site}
	params, _ := query.Values(fw)
	data, status, err := c.GetRequest(fmt.Sprintf("%s/orders/%s?%s", versionedPrefix, id, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Orders list method
func (c *Client) Orders(parameters OrdersRequest) (*OrdersResponse, int, error) {
	var resp OrdersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/orders?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrderCreate method
func (c *Client) OrderCreate(order Order, site ...string) (*OrderChangeResponse, int, error) {
	var resp OrderChangeResponse
	orderJson, _ := json.Marshal(&order)

	p := url.Values{
		"order": {string(orderJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/create", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomerEdit method
func (c *Client) OrderEdit(order Order, by string, site ...string) (*OrderChangeResponse, int, error) {
	var resp OrderChangeResponse
	var uid = strconv.Itoa(order.Id)
	var context = checkBy(by)

	if context == "externalId" {
		uid = order.ExternalId
	}

	orderJson, _ := json.Marshal(&order)

	p := url.Values{
		"by":    {string(context)},
		"order": {string(orderJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/%s/edit", versionedPrefix, uid), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrdersUpload method
func (c *Client) OrdersUpload(orders []Order, site ...string) (*OrdersUploadResponse, int, error) {
	var resp OrdersUploadResponse

	uploadJson, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(uploadJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/upload", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrdersFixExternalIds method
func (c *Client) OrdersFixExternalIds(orders []IdentifiersPair) (*SucessfulResponse, int, error) {
	var resp SucessfulResponse

	ordersJson, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(ordersJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/fix-external-ids", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrdersHistory method
func (c *Client) OrdersHistory(parameters OrdersHistoryRequest) (*CustomersHistoryResponse, int, error) {
	var resp CustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers/history?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// User get method
func (c *Client) User(id int) (*UserResponse, int, error) {
	var resp UserResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/users/%d", versionedPrefix, id))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Users list method
func (c *Client) Users(parameters UsersRequest) (*UsersResponse, int, error) {
	var resp UsersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/users?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// UserGroups list method
func (c *Client) UserGroups(parameters UserGroupsRequest) (*UserGroupsResponse, int, error) {
	var resp UserGroupsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/user-groups", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// UserStatus update method
func (c *Client) UserStatus(id int, status string) (*SucessfulResponse, int, error) {
	var resp SucessfulResponse

	p := url.Values{
		"status": {string(status)},
	}

	data, st, err := c.PostRequest(fmt.Sprintf("%s/users/%d/status", versionedPrefix, id), p)
	if err != nil {
		return &resp, st, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, st, err
}

// Task get method
func (c *Client) Task(id int) (*TaskResponse, int, error) {
	var resp TaskResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/tasks/%d", versionedPrefix, id))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Tasks list method
func (c *Client) Tasks(parameters TasksRequest) (*TasksResponse, int, error) {
	var resp TasksResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/tasks?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// TaskCreate method
func (c *Client) TaskCreate(task Task, site ...string) (*TaskChangeResponse, int, error) {
	var resp TaskChangeResponse
	taskJson, _ := json.Marshal(&task)

	p := url.Values{
		"task": {string(taskJson[:])},
	}

	if len(site) > 0 {
		s := site[0]

		if s != "" {
			p.Add("site", s)
		}
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/tasks/create", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// TaskEdit method
func (c *Client) TaskEdit(task Task, site ...string) (*SucessfulResponse, int, error) {
	var resp SucessfulResponse
	var uid = strconv.Itoa(task.Id)

	taskJson, _ := json.Marshal(&task)

	p := url.Values{
		"task": {string(taskJson[:])},
	}

	if len(site) > 0 {
		s := site[0]

		if s != "" {
			p.Add("site", s)
		}
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/tasks/%s/edit", versionedPrefix, uid), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

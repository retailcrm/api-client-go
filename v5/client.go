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

// GetRequest implements GET Request
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

// PostRequest implements POST Request
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

// StaticticUpdate update statistic
func (c *Client) StaticticUpdate() (*SucessfulResponse, int, error) {
	var resp SucessfulResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/statistic/update", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Segments get segments
func (c *Client) Segments(parameters SegmentsRequest) (*SegmentsResponse, int, error) {
	var resp SegmentsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/segments?%s", versionedPrefix, params.Encode()))
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
func (c *Client) OrderCreate(order Order, site ...string) (*CreateResponse, int, error) {
	var resp CreateResponse
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

// OrderEdit method
func (c *Client) OrderEdit(order Order, by string, site ...string) (*CreateResponse, int, error) {
	var resp CreateResponse
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
func (c *Client) TaskCreate(task Task, site ...string) (*CreateResponse, int, error) {
	var resp CreateResponse
	taskJson, _ := json.Marshal(&task)

	p := url.Values{
		"task": {string(taskJson[:])},
	}

	fillSite(&p, site)

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

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/tasks/%s/edit", versionedPrefix, uid), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Notes list method
func (c *Client) Notes(parameters NotesRequest) (*NotesResponse, int, error) {
	var resp NotesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers/notes?%s", versionedPrefix, params.Encode()))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// NoteCreate method
func (c *Client) NoteCreate(note Note, site ...string) (*CreateResponse, int, error) {
	var resp CreateResponse

	noteJson, _ := json.Marshal(&note)

	p := url.Values{
		"note": {string(noteJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/notes/create", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// NoteDelete method
func (c *Client) NoteDelete(id int) (*SucessfulResponse, int, error) {
	var resp SucessfulResponse

	p := url.Values{
		"id": {string(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/notes/%d/delete", versionedPrefix, id), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentCreate method
func (c *Client) PaymentCreate(payment Payment, site ...string) (*CreateResponse, int, error) {
	var resp CreateResponse

	paymentJson, _ := json.Marshal(&payment)

	p := url.Values{
		"payment": {string(paymentJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/payments/create", versionedPrefix), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentDelete method
func (c *Client) PaymentDelete(id int) (*SucessfulResponse, int, error) {
	var resp SucessfulResponse

	p := url.Values{
		"id": {string(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/payments/%d/delete", versionedPrefix, id), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentEdit method
func (c *Client) PaymentEdit(payment Payment, by string, site ...string) (*SucessfulResponse, int, error) {
	var resp SucessfulResponse
	var uid = strconv.Itoa(payment.Id)
	var context = checkBy(by)

	if context == "externalId" {
		uid = payment.ExternalId
	}

	paymentJson, _ := json.Marshal(&payment)

	p := url.Values{
		"payment": {string(paymentJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/payments/%s/edit", versionedPrefix, uid), p)
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Countries method
func (c *Client) Countries() (*CountriesResponse, int, error) {
	var resp CountriesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/countries", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CostGroups method
func (c *Client) CostGroups() (*CostGroupsResponse, int, error) {
	var resp CostGroupsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/cost-groups", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CostItems method
func (c *Client) CostItems() (*CostItemsResponse, int, error) {
	var resp CostItemsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/cost-items", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Couriers method
func (c *Client) Couriers() (*CouriersResponse, int, error) {
	var resp CouriersResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/couriers", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryService method
func (c *Client) DeliveryService() (*DeliveryServiceResponse, int, error) {
	var resp DeliveryServiceResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/delivery-services", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryTypes method
func (c *Client) DeliveryTypes() (*DeliveryTypesResponse, int, error) {
	var resp DeliveryTypesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/delivery-types", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// LegalEntities method
func (c *Client) LegalEntities() (*LegalEntitiesResponse, int, error) {
	var resp LegalEntitiesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/legal-entities", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrderMethods method
func (c *Client) OrderMethods() (*OrderMethodsResponse, int, error) {
	var resp OrderMethodsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/order-methods", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrderTypes method
func (c *Client) OrderTypes() (*OrderTypesResponse, int, error) {
	var resp OrderTypesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/order-types", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentStatuses method
func (c *Client) PaymentStatuses() (*PaymentStatusesResponse, int, error) {
	var resp PaymentStatusesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/payment-statuses", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentTypes method
func (c *Client) PaymentTypes() (*PaymentTypesResponse, int, error) {
	var resp PaymentTypesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/payment-types", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PriceTypes method
func (c *Client) PriceTypes() (*PriceTypesResponse, int, error) {
	var resp PriceTypesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/price-types", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ProductStatuses method
func (c *Client) ProductStatuses() (*ProductStatusesResponse, int, error) {
	var resp ProductStatusesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/product-statuses", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Statuses method
func (c *Client) Statuses() (*StatusesResponse, int, error) {
	var resp StatusesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/statuses", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// StatusGroups method
func (c *Client) StatusGroups() (*StatusGroupsResponse, int, error) {
	var resp StatusGroupsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/status-groups", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Sites method
func (c *Client) Sites() (*SitesResponse, int, error) {
	var resp SitesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/sites", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Stores method
func (c *Client) Stores() (*StoresResponse, int, error) {
	var resp StoresResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/stores", versionedPrefix))
	if err != nil {
		return &resp, status, err
	}

	err = json.Unmarshal(data, &resp)

	return &resp, status, err
}

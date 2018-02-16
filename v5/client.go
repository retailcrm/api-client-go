package v5

import (
	"encoding/json"
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
func (c *Client) GetRequest(urlWithParameters string) ([]byte, int, ErrorResponse) {
	var res []byte
	var bug ErrorResponse

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.Url, urlWithParameters), nil)
	if err != nil {
		bug.ErrorMsg = err.Error()
		return res, 0, bug
	}

	req.Header.Set("X-API-KEY", c.Key)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		bug.ErrorMsg = err.Error()
		return res, 0, bug
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		bug.ErrorMsg = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, bug
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		bug.ErrorMsg = err.Error()
	}

	eresp, _ := c.ErrorResponse(res)
	if eresp.ErrorMsg != "" {
		return res, resp.StatusCode, eresp
	}

	return res, resp.StatusCode, bug
}

// PostRequest implements POST Request
func (c *Client) PostRequest(url string, postParams url.Values) ([]byte, int, ErrorResponse) {
	var res []byte
	var bug ErrorResponse

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s", c.Url, url),
		strings.NewReader(postParams.Encode()),
	)
	if err != nil {
		bug.ErrorMsg = err.Error()
		return res, 0, bug
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-API-KEY", c.Key)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		bug.ErrorMsg = err.Error()
		return res, 0, bug
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		bug.ErrorMsg = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, bug
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		bug.ErrorMsg = err.Error()
		return res, 0, bug
	}

	eresp, _ := c.ErrorResponse(res)
	if eresp.ErrorMsg != "" {
		return res, resp.StatusCode, eresp
	}

	return res, resp.StatusCode, bug
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
func (c *Client) ApiVersions() (*VersionResponse, int, ErrorResponse) {
	var resp VersionResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/api-versions", unversionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ApiCredentials get available API methods
func (c *Client) ApiCredentials() (*CredentialResponse, int, ErrorResponse) {
	var resp CredentialResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/credentials", unversionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// StaticticUpdate update statistic
func (c *Client) StaticticUpdate() (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/statistic/update", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Segments get segments
func (c *Client) Segments(parameters SegmentsRequest) (*SegmentsResponse, int, ErrorResponse) {
	var resp SegmentsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/segments?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Customer get method
func (c *Client) Customer(id, by, site string) (*CustomerResponse, int, ErrorResponse) {
	var resp CustomerResponse
	var context = checkBy(by)

	fw := CustomerRequest{context, site}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers/%s?%s", versionedPrefix, id, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Customers list method
func (c *Client) Customers(parameters CustomersRequest) (*CustomersResponse, int, ErrorResponse) {
	var resp CustomersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomerCreate method
func (c *Client) CustomerCreate(customer Customer, site ...string) (*CustomerChangeResponse, int, ErrorResponse) {
	var resp CustomerChangeResponse

	customerJson, _ := json.Marshal(&customer)

	p := url.Values{
		"customer": {string(customerJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/create", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomerEdit method
func (c *Client) CustomerEdit(customer Customer, by string, site ...string) (*CustomerChangeResponse, int, ErrorResponse) {
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
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersUpload method
func (c *Client) CustomersUpload(customers []Customer, site ...string) (*CustomersUploadResponse, int, ErrorResponse) {
	var resp CustomersUploadResponse

	uploadJson, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(uploadJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/upload", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersCombine method
func (c *Client) CustomersCombine(customers []Customer, resultCustomer Customer) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	combineJsonIn, _ := json.Marshal(&customers)
	combineJsonOut, _ := json.Marshal(&resultCustomer)

	p := url.Values{
		"customers":      {string(combineJsonIn[:])},
		"resultCustomer": {string(combineJsonOut[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/combine", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersFixExternalIds method
func (c *Client) CustomersFixExternalIds(customers []IdentifiersPair) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	customersJson, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(customersJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/fix-external-ids", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CustomersHistory method
func (c *Client) CustomersHistory(parameters CustomersHistoryRequest) (*CustomersHistoryResponse, int, ErrorResponse) {
	var resp CustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers/history?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Order get method
func (c *Client) Order(id, by, site string) (*OrderResponse, int, ErrorResponse) {
	var resp OrderResponse
	var context = checkBy(by)

	fw := OrderRequest{context, site}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/orders/%s?%s", versionedPrefix, id, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Orders list method
func (c *Client) Orders(parameters OrdersRequest) (*OrdersResponse, int, ErrorResponse) {
	var resp OrdersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/orders?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrderCreate method
func (c *Client) OrderCreate(order Order, site ...string) (*CreateResponse, int, ErrorResponse) {
	var resp CreateResponse
	orderJson, _ := json.Marshal(&order)

	p := url.Values{
		"order": {string(orderJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/create", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrderEdit method
func (c *Client) OrderEdit(order Order, by string, site ...string) (*CreateResponse, int, ErrorResponse) {
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
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrdersUpload method
func (c *Client) OrdersUpload(orders []Order, site ...string) (*OrdersUploadResponse, int, ErrorResponse) {
	var resp OrdersUploadResponse

	uploadJson, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(uploadJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/upload", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrdersCombine method
func (c *Client) OrdersCombine(technique string, order, resultOrder Order) (*OperationResponse, int, ErrorResponse) {
	var resp OperationResponse

	combineJsonIn, _ := json.Marshal(&order)
	combineJsonOut, _ := json.Marshal(&resultOrder)

	p := url.Values{
		"technique":   {technique},
		"order":       {string(combineJsonIn[:])},
		"resultOrder": {string(combineJsonOut[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/combine", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrdersFixExternalIds method
func (c *Client) OrdersFixExternalIds(orders []IdentifiersPair) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	ordersJson, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(ordersJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/fix-external-ids", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrdersHistory method
func (c *Client) OrdersHistory(parameters OrdersHistoryRequest) (*CustomersHistoryResponse, int, ErrorResponse) {
	var resp CustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/orders/history?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Pack get method
func (c *Client) Pack(id int) (*PackResponse, int, ErrorResponse) {
	var resp PackResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/orders/packs/%d", versionedPrefix, id))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Packs list method
func (c *Client) Packs(parameters PacksRequest) (*PacksResponse, int, ErrorResponse) {
	var resp PacksResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/orders/packs?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PackCreate method
func (c *Client) PackCreate(pack Pack) (*CreateResponse, int, ErrorResponse) {
	var resp CreateResponse
	packJson, _ := json.Marshal(&pack)

	p := url.Values{
		"pack": {string(packJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/packs/create", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PackEdit method
func (c *Client) PackEdit(pack Pack) (*CreateResponse, int, ErrorResponse) {
	var resp CreateResponse

	packJson, _ := json.Marshal(&pack)

	p := url.Values{
		"pack": {string(packJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/packs/%d/edit", versionedPrefix, pack.Id), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PackDelete method
func (c *Client) PackDelete(id int) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/packs/%d/delete", versionedPrefix, id), url.Values{})
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PacksHistory method
func (c *Client) PacksHistory(parameters PacksHistoryRequest) (*PacksHistoryResponse, int, ErrorResponse) {
	var resp PacksHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/orders/packs/history?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// User get method
func (c *Client) User(id int) (*UserResponse, int, ErrorResponse) {
	var resp UserResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/users/%d", versionedPrefix, id))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Users list method
func (c *Client) Users(parameters UsersRequest) (*UsersResponse, int, ErrorResponse) {
	var resp UsersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/users?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// UserGroups list method
func (c *Client) UserGroups(parameters UserGroupsRequest) (*UserGroupsResponse, int, ErrorResponse) {
	var resp UserGroupsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/user-groups", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// UserStatus update method
func (c *Client) UserStatus(id int, status string) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	p := url.Values{
		"status": {string(status)},
	}

	data, st, err := c.PostRequest(fmt.Sprintf("%s/users/%d/status", versionedPrefix, id), p)
	if err.ErrorMsg != "" {
		return &resp, st, err
	}

	json.Unmarshal(data, &resp)

	return &resp, st, err
}

// Task get method
func (c *Client) Task(id int) (*TaskResponse, int, ErrorResponse) {
	var resp TaskResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/tasks/%d", versionedPrefix, id))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Tasks list method
func (c *Client) Tasks(parameters TasksRequest) (*TasksResponse, int, ErrorResponse) {
	var resp TasksResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/tasks?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// TaskCreate method
func (c *Client) TaskCreate(task Task, site ...string) (*CreateResponse, int, ErrorResponse) {
	var resp CreateResponse
	taskJson, _ := json.Marshal(&task)

	p := url.Values{
		"task": {string(taskJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/tasks/create", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// TaskEdit method
func (c *Client) TaskEdit(task Task, site ...string) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse
	var uid = strconv.Itoa(task.Id)

	taskJson, _ := json.Marshal(&task)

	p := url.Values{
		"task": {string(taskJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/tasks/%s/edit", versionedPrefix, uid), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Notes list method
func (c *Client) Notes(parameters NotesRequest) (*NotesResponse, int, ErrorResponse) {
	var resp NotesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/customers/notes?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// NoteCreate method
func (c *Client) NoteCreate(note Note, site ...string) (*CreateResponse, int, ErrorResponse) {
	var resp CreateResponse

	noteJson, _ := json.Marshal(&note)

	p := url.Values{
		"note": {string(noteJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/notes/create", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// NoteDelete method
func (c *Client) NoteDelete(id int) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	p := url.Values{
		"id": {string(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/customers/notes/%d/delete", versionedPrefix, id), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentCreate method
func (c *Client) PaymentCreate(payment Payment, site ...string) (*CreateResponse, int, ErrorResponse) {
	var resp CreateResponse

	paymentJson, _ := json.Marshal(&payment)

	p := url.Values{
		"payment": {string(paymentJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/payments/create", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentDelete method
func (c *Client) PaymentDelete(id int) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	p := url.Values{
		"id": {string(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/orders/payments/%d/delete", versionedPrefix, id), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentEdit method
func (c *Client) PaymentEdit(payment Payment, by string, site ...string) (*SucessfulResponse, int, ErrorResponse) {
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
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Countries method
func (c *Client) Countries() (*CountriesResponse, int, ErrorResponse) {
	var resp CountriesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/countries", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CostGroups method
func (c *Client) CostGroups() (*CostGroupsResponse, int, ErrorResponse) {
	var resp CostGroupsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/cost-groups", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CostItems method
func (c *Client) CostItems() (*CostItemsResponse, int, ErrorResponse) {
	var resp CostItemsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/cost-items", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Couriers method
func (c *Client) Couriers() (*CouriersResponse, int, ErrorResponse) {
	var resp CouriersResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/couriers", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryServices method
func (c *Client) DeliveryServices() (*DeliveryServiceResponse, int, ErrorResponse) {
	var resp DeliveryServiceResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/delivery-services", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryTypes method
func (c *Client) DeliveryTypes() (*DeliveryTypesResponse, int, ErrorResponse) {
	var resp DeliveryTypesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/delivery-types", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// LegalEntities method
func (c *Client) LegalEntities() (*LegalEntitiesResponse, int, ErrorResponse) {
	var resp LegalEntitiesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/legal-entities", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrderMethods method
func (c *Client) OrderMethods() (*OrderMethodsResponse, int, ErrorResponse) {
	var resp OrderMethodsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/order-methods", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrderTypes method
func (c *Client) OrderTypes() (*OrderTypesResponse, int, ErrorResponse) {
	var resp OrderTypesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/order-types", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentStatuses method
func (c *Client) PaymentStatuses() (*PaymentStatusesResponse, int, ErrorResponse) {
	var resp PaymentStatusesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/payment-statuses", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentTypes method
func (c *Client) PaymentTypes() (*PaymentTypesResponse, int, ErrorResponse) {
	var resp PaymentTypesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/payment-types", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PriceTypes method
func (c *Client) PriceTypes() (*PriceTypesResponse, int, ErrorResponse) {
	var resp PriceTypesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/price-types", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ProductStatuses method
func (c *Client) ProductStatuses() (*ProductStatusesResponse, int, ErrorResponse) {
	var resp ProductStatusesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/product-statuses", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Statuses method
func (c *Client) Statuses() (*StatusesResponse, int, ErrorResponse) {
	var resp StatusesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/statuses", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// StatusGroups method
func (c *Client) StatusGroups() (*StatusGroupsResponse, int, ErrorResponse) {
	var resp StatusGroupsResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/status-groups", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Sites method
func (c *Client) Sites() (*SitesResponse, int, ErrorResponse) {
	var resp SitesResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/sites", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Stores method
func (c *Client) Stores() (*StoresResponse, int, ErrorResponse) {
	var resp StoresResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/reference/stores", versionedPrefix))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CostGroupEdit method
func (c *Client) CostGroupEdit(costGroup CostGroup) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&costGroup)

	p := url.Values{
		"costGroup": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/cost-groups/%s/edit", versionedPrefix, costGroup.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CostItemEdit method
func (c *Client) CostItemEdit(costItem CostItem) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&costItem)

	p := url.Values{
		"costItem": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/cost-items/%s/edit", versionedPrefix, costItem.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CourierCreate method
func (c *Client) CourierCreate(courier Courier) (*CreateResponse, int, ErrorResponse) {
	var resp CreateResponse

	objJson, _ := json.Marshal(&courier)

	p := url.Values{
		"courier": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/couriers/create", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// CourierEdit method
func (c *Client) CourierEdit(courier Courier) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&courier)

	p := url.Values{
		"courier": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/couriers/%d/edit", versionedPrefix, courier.Id), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryServiceEdit method
func (c *Client) DeliveryServiceEdit(deliveryService DeliveryService) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&deliveryService)

	p := url.Values{
		"deliveryService": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/delivery-services/%s/edit", versionedPrefix, deliveryService.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryTypeEdit method
func (c *Client) DeliveryTypeEdit(deliveryType DeliveryType) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&deliveryType)

	p := url.Values{
		"deliveryType": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/delivery-types/%s/edit", versionedPrefix, deliveryType.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// LegalEntityEditorderMe method
func (c *Client) LegalEntityEdit(legalEntity LegalEntity) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&legalEntity)

	p := url.Values{
		"legalEntity": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/legal-entities/%s/edit", versionedPrefix, legalEntity.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrderMethodEdit method
func (c *Client) OrderMethodEdit(orderMethod OrderMethod) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&orderMethod)

	p := url.Values{
		"orderMethod": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/order-methods/%s/edit", versionedPrefix, orderMethod.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// OrderTypeEdit method
func (c *Client) OrderTypeEdit(orderType OrderType) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&orderType)

	p := url.Values{
		"orderType": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/order-types/%s/edit", versionedPrefix, orderType.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentStatusEdit method
func (c *Client) PaymentStatusEdit(paymentStatus PaymentStatus) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&paymentStatus)

	p := url.Values{
		"paymentStatus": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/payment-statuses/%s/edit", versionedPrefix, paymentStatus.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PaymentTypeEdit method
func (c *Client) PaymentTypeEdit(paymentType PaymentType) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&paymentType)

	p := url.Values{
		"paymentType": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/payment-types/%s/edit", versionedPrefix, paymentType.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PriceTypeEdit method
func (c *Client) PriceTypeEdit(priceType PriceType) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&priceType)

	p := url.Values{
		"priceType": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/price-types/%s/edit", versionedPrefix, priceType.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ProductStatusEdit method
func (c *Client) ProductStatusEdit(productStatus ProductStatus) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&productStatus)

	p := url.Values{
		"productStatus": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/product-statuses/%s/edit", versionedPrefix, productStatus.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// StatusEdit method
func (c *Client) StatusEdit(st Status) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&st)

	p := url.Values{
		"status": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/statuses/%s/edit", versionedPrefix, st.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// SiteEdit method
func (c *Client) SiteEdit(site Site) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&site)

	p := url.Values{
		"site": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/sites/%s/edit", versionedPrefix, site.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// StoreEdit method
func (c *Client) StoreEdit(store Store) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	objJson, _ := json.Marshal(&store)

	p := url.Values{
		"store": {string(objJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/reference/stores/%s/edit", versionedPrefix, store.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// Inventories method
func (c *Client) Inventories(parameters InventoriesRequest) (*InventoriesResponse, int, ErrorResponse) {
	var resp InventoriesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/store/inventories?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// InventoriesUpload method
func (c *Client) InventoriesUpload(inventories []InventoryUpload, site ...string) (*StoreUploadResponse, int, ErrorResponse) {
	var resp StoreUploadResponse

	uploadJson, _ := json.Marshal(&inventories)

	p := url.Values{
		"offers": {string(uploadJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/store/inventories/upload", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// PricesUpload method
func (c *Client) PricesUpload(prices []OfferPriceUpload) (*StoreUploadResponse, int, ErrorResponse) {
	var resp StoreUploadResponse

	uploadJson, _ := json.Marshal(&prices)

	p := url.Values{
		"prices": {string(uploadJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/store/prices/upload", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ProductsGroup method
func (c *Client) ProductsGroup(parameters ProductsGroupsRequest) (*ProductsGroupsResponse, int, ErrorResponse) {
	var resp ProductsGroupsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/store/product-groups?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ProductsGroup method
func (c *Client) Products(parameters ProductsRequest) (*ProductsResponse, int, ErrorResponse) {
	var resp ProductsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/store/products?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// ProductsProperties method
func (c *Client) ProductsProperties(parameters ProductsPropertiesRequest) (*ProductsPropertiesResponse, int, ErrorResponse) {
	var resp ProductsPropertiesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/store/products/properties?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryTracking method
func (c *Client) DeliveryTracking(parameters DeliveryTrackingRequest, subcode string) (*SucessfulResponse, int, ErrorResponse) {
	var resp SucessfulResponse

	updateJson, _ := json.Marshal(&parameters)

	p := url.Values{
		"statusUpdate": {string(updateJson[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/delivery/generic/%s/tracking", versionedPrefix, subcode), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryShipments method
func (c *Client) DeliveryShipments(parameters DeliveryShipmentsRequest) (*DeliveryShipmentsResponse, int, ErrorResponse) {
	var resp DeliveryShipmentsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("%s/delivery/shipments?%s", versionedPrefix, params.Encode()))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryShipments method
func (c *Client) DeliveryShipment(id int) (*DeliveryShipmentResponse, int, ErrorResponse) {
	var resp DeliveryShipmentResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/delivery/shipments/%s", versionedPrefix, id))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryShipmentEdit method
func (c *Client) DeliveryShipmentEdit(shipment DeliveryShipment, site ...string) (*DeliveryShipmentUpdateResponse, int, ErrorResponse) {
	var resp DeliveryShipmentUpdateResponse
	updateJson, _ := json.Marshal(&shipment)

	p := url.Values{
		"deliveryShipment": {string(updateJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/delivery/shipments/%s/edit", versionedPrefix, strconv.Itoa(shipment.Id)), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// DeliveryShipmentCreate method
func (c *Client) DeliveryShipmentCreate(shipment DeliveryShipment, deliveryType string, site ...string) (*DeliveryShipmentUpdateResponse, int, ErrorResponse) {
	var resp DeliveryShipmentUpdateResponse
	updateJson, _ := json.Marshal(&shipment)

	p := url.Values{
		"deliveryType":     {string(deliveryType)},
		"deliveryShipment": {string(updateJson[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("%s/delivery/shipments/create", versionedPrefix), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// IntegrationModule method
func (c *Client) IntegrationModule(code string) (*IntegrationModuleResponse, int, ErrorResponse) {
	var resp IntegrationModuleResponse

	data, status, err := c.GetRequest(fmt.Sprintf("%s/integration-modules/%s", versionedPrefix, code))
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

// IntegrationModuleEdit method
func (c *Client) IntegrationModuleEdit(integrationModule IntegrationModule) (*IntegrationModuleEditResponse, int, ErrorResponse) {
	var resp IntegrationModuleEditResponse
	updateJson, _ := json.Marshal(&integrationModule)

	p := url.Values{"integrationModule": {string(updateJson[:])}}

	data, status, err := c.PostRequest(fmt.Sprintf("%s/integration-modules/%s/edit", versionedPrefix, integrationModule.Code), p)
	if err.ErrorMsg != "" {
		return &resp, status, err
	}

	json.Unmarshal(data, &resp)

	return &resp, status, err
}

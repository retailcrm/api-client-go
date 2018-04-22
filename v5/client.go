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
	"github.com/retailcrm/api-client-go/errs"
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
func (c *Client) GetRequest(urlWithParameters string, versioned ...bool) ([]byte, int, errs.Failure) {
	var res []byte
	var cerr errs.Failure
	var prefix = "/api/v5"

	if len(versioned) > 0 {
		s := versioned[0]

		if s == false {
			prefix = "/api"
		}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.URL, prefix, urlWithParameters), nil)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	req.Header.Set("X-API-KEY", c.Key)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		cerr.ApiErr = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, cerr
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		cerr.RuntimeErr = err
	}

	return res, resp.StatusCode, cerr
}

// PostRequest implements POST Request
func (c *Client) PostRequest(url string, postParams url.Values) ([]byte, int, errs.Failure) {
	var res []byte
	var cerr errs.Failure
	var prefix = "/api/v5"

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", c.URL, prefix, url), strings.NewReader(postParams.Encode()))
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-API-KEY", c.Key)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		cerr.ApiErr = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, cerr
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		cerr.RuntimeErr = err
		return res, 0, cerr
	}

	return res, resp.StatusCode, cerr
}

func buildRawResponse(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	return res, nil
}

func buildErr(data []byte) errs.Failure {
	var err = errs.Failure{}

	eresp, errr := errs.ErrorResponse(data)
	err.RuntimeErr = errr
	err.ApiErr = eresp.ErrorMsg

	if eresp.Errors != nil {
		err.ApiErrs = eresp.Errors
	}

	return err
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

// APIVersions get all available API versions for exact account
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.APIVersions()
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.ErrorMsg)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ErrorMsg)
// 	}
//
// 	for _, value := range data.versions {
// 		fmt.Printf("%v\n", value)
// 	}
func (c *Client) APIVersions() (VersionResponse, int, errs.Failure) {
	var resp VersionResponse

	data, status, err := c.GetRequest("/api-versions", false)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// APICredentials get all available API methods for exact account
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.APICredentials()
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.ErrorMsg)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ErrorMsg)
// 	}
//
// 	for _, value := range data.credentials {
// 		fmt.Printf("%v\n", value)
// 	}
func (c *Client) APICredentials() (CredentialResponse, int, errs.Failure) {
	var resp CredentialResponse

	data, status, err := c.GetRequest("/credentials", false)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Customers list method
func (c *Client) Customers(parameters CustomersRequest) (CustomersResponse, int, errs.Failure) {
	var resp CustomersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CustomersCombine method
func (c *Client) CustomersCombine(customers []Customer, resultCustomer Customer) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	combineJSONIn, _ := json.Marshal(&customers)
	combineJSONOut, _ := json.Marshal(&resultCustomer)

	p := url.Values{
		"customers":      {string(combineJSONIn[:])},
		"resultCustomer": {string(combineJSONOut[:])},
	}

	data, status, err := c.PostRequest("/customers/combine", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CustomerCreate method
func (c *Client) CustomerCreate(customer Customer, site ...string) (CustomerChangeResponse, int, errs.Failure) {
	var resp CustomerChangeResponse

	customerJSON, _ := json.Marshal(&customer)

	p := url.Values{
		"customer": {string(customerJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers/create", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CustomersFixExternalIds method
func (c *Client) CustomersFixExternalIds(customers []IdentifiersPair) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	customersJSON, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(customersJSON[:])},
	}

	data, status, err := c.PostRequest("/customers/fix-external-ids", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CustomersHistory method
func (c *Client) CustomersHistory(parameters CustomersHistoryRequest) (CustomersHistoryResponse, int, errs.Failure) {
	var resp CustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers/history?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CustomerNotes list method
func (c *Client) CustomerNotes(parameters NotesRequest) (NotesResponse, int, errs.Failure) {
	var resp NotesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers/notes?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CustomerNoteCreate method
func (c *Client) CustomerNoteCreate(note Note, site ...string) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse

	noteJSON, _ := json.Marshal(&note)

	p := url.Values{
		"note": {string(noteJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers/notes/create", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CustomerNoteDelete method
func (c *Client) CustomerNoteDelete(id int) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	p := url.Values{
		"id": {string(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/customers/notes/%d/delete", id), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CustomersUpload method
func (c *Client) CustomersUpload(customers []Customer, site ...string) (CustomersUploadResponse, int, errs.Failure) {
	var resp CustomersUploadResponse

	uploadJSON, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(uploadJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers/upload", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Customer get method
func (c *Client) Customer(id, by, site string) (CustomerResponse, int, errs.Failure) {
	var resp CustomerResponse
	var context = checkBy(by)

	fw := CustomerRequest{context, site}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers/%s?%s", id, params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CustomerEdit method
func (c *Client) CustomerEdit(customer Customer, by string, site ...string) (CustomerChangeResponse, int, errs.Failure) {
	var resp CustomerChangeResponse
	var uid = strconv.Itoa(customer.ID)
	var context = checkBy(by)

	if context == "externalId" {
		uid = customer.ExternalID
	}

	customerJSON, _ := json.Marshal(&customer)

	p := url.Values{
		"by":       {string(context)},
		"customer": {string(customerJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/customers/%s/edit", uid), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeliveryTracking method
func (c *Client) DeliveryTracking(parameters DeliveryTrackingRequest, subcode string) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	updateJSON, _ := json.Marshal(&parameters)

	p := url.Values{
		"statusUpdate": {string(updateJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/delivery/generic/%s/tracking", subcode), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeliveryShipments method
func (c *Client) DeliveryShipments(parameters DeliveryShipmentsRequest) (DeliveryShipmentsResponse, int, errs.Failure) {
	var resp DeliveryShipmentsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/delivery/shipments?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeliveryShipmentCreate method
func (c *Client) DeliveryShipmentCreate(shipment DeliveryShipment, deliveryType string, site ...string) (DeliveryShipmentUpdateResponse, int, errs.Failure) {
	var resp DeliveryShipmentUpdateResponse
	updateJSON, _ := json.Marshal(&shipment)

	p := url.Values{
		"deliveryType":     {string(deliveryType)},
		"deliveryShipment": {string(updateJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/delivery/shipments/create", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeliveryShipment method
func (c *Client) DeliveryShipment(id int) (DeliveryShipmentResponse, int, errs.Failure) {
	var resp DeliveryShipmentResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/delivery/shipments/%d", id))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeliveryShipmentEdit method
func (c *Client) DeliveryShipmentEdit(shipment DeliveryShipment, site ...string) (DeliveryShipmentUpdateResponse, int, errs.Failure) {
	var resp DeliveryShipmentUpdateResponse
	updateJSON, _ := json.Marshal(&shipment)

	p := url.Values{
		"deliveryShipment": {string(updateJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/delivery/shipments/%s/edit", strconv.Itoa(shipment.ID)), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// IntegrationModule method
func (c *Client) IntegrationModule(code string) (IntegrationModuleResponse, int, errs.Failure) {
	var resp IntegrationModuleResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/integration-modules/%s", code))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// IntegrationModuleEdit method
func (c *Client) IntegrationModuleEdit(integrationModule IntegrationModule) (IntegrationModuleEditResponse, int, errs.Failure) {
	var resp IntegrationModuleEditResponse
	updateJSON, _ := json.Marshal(&integrationModule)

	p := url.Values{"integrationModule": {string(updateJSON[:])}}

	data, status, err := c.PostRequest(fmt.Sprintf("/integration-modules/%s/edit", integrationModule.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Orders list method
func (c *Client) Orders(parameters OrdersRequest) (OrdersResponse, int, errs.Failure) {
	var resp OrdersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrdersCombine method
func (c *Client) OrdersCombine(technique string, order, resultOrder Order) (OperationResponse, int, errs.Failure) {
	var resp OperationResponse

	combineJSONIn, _ := json.Marshal(&order)
	combineJSONOut, _ := json.Marshal(&resultOrder)

	p := url.Values{
		"technique":   {technique},
		"order":       {string(combineJSONIn[:])},
		"resultOrder": {string(combineJSONOut[:])},
	}

	data, status, err := c.PostRequest("/orders/combine", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrderCreate method
func (c *Client) OrderCreate(order Order, site ...string) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse
	orderJSON, _ := json.Marshal(&order)

	p := url.Values{
		"order": {string(orderJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/create", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrdersFixExternalIds method
func (c *Client) OrdersFixExternalIds(orders []IdentifiersPair) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	ordersJSON, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(ordersJSON[:])},
	}

	data, status, err := c.PostRequest("/orders/fix-external-ids", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrdersHistory method
func (c *Client) OrdersHistory(parameters OrdersHistoryRequest) (CustomersHistoryResponse, int, errs.Failure) {
	var resp CustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/history?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrderPaymentCreate method
func (c *Client) OrderPaymentCreate(payment Payment, site ...string) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse

	paymentJSON, _ := json.Marshal(&payment)

	p := url.Values{
		"payment": {string(paymentJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/payments/create", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrderPaymentDelete method
func (c *Client) OrderPaymentDelete(id int) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	p := url.Values{
		"id": {string(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/payments/%d/delete", id), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrderPaymentEdit method
func (c *Client) OrderPaymentEdit(payment Payment, by string, site ...string) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse
	var uid = strconv.Itoa(payment.ID)
	var context = checkBy(by)

	if context == "externalId" {
		uid = payment.ExternalID
	}

	paymentJSON, _ := json.Marshal(&payment)

	p := url.Values{
		"payment": {string(paymentJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/payments/%s/edit", uid), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrdersUpload method
func (c *Client) OrdersUpload(orders []Order, site ...string) (OrdersUploadResponse, int, errs.Failure) {
	var resp OrdersUploadResponse

	uploadJSON, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(uploadJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/upload", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Order get method
func (c *Client) Order(id, by, site string) (OrderResponse, int, errs.Failure) {
	var resp OrderResponse
	var context = checkBy(by)

	fw := OrderRequest{context, site}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/%s?%s", id, params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrderEdit method
func (c *Client) OrderEdit(order Order, by string, site ...string) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse
	var uid = strconv.Itoa(order.ID)
	var context = checkBy(by)

	if context == "externalId" {
		uid = order.ExternalID
	}

	orderJSON, _ := json.Marshal(&order)

	p := url.Values{
		"by":    {string(context)},
		"order": {string(orderJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/%s/edit", uid), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Packs list method
func (c *Client) Packs(parameters PacksRequest) (PacksResponse, int, errs.Failure) {
	var resp PacksResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/packs?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PackCreate method
func (c *Client) PackCreate(pack Pack) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse
	packJSON, _ := json.Marshal(&pack)

	p := url.Values{
		"pack": {string(packJSON[:])},
	}

	data, status, err := c.PostRequest("/orders/packs/create", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PacksHistory method
func (c *Client) PacksHistory(parameters PacksHistoryRequest) (PacksHistoryResponse, int, errs.Failure) {
	var resp PacksHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/packs/history?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Pack get method
func (c *Client) Pack(id int) (PackResponse, int, errs.Failure) {
	var resp PackResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/packs/%d", id))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PackDelete method
func (c *Client) PackDelete(id int) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/packs/%d/delete", id), url.Values{})
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PackEdit method
func (c *Client) PackEdit(pack Pack) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse

	packJSON, _ := json.Marshal(&pack)

	p := url.Values{
		"pack": {string(packJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/packs/%d/edit", pack.ID), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Countries method
func (c *Client) Countries() (CountriesResponse, int, errs.Failure) {
	var resp CountriesResponse

	data, status, err := c.GetRequest("/reference/countries")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CostGroups method
func (c *Client) CostGroups() (CostGroupsResponse, int, errs.Failure) {
	var resp CostGroupsResponse

	data, status, err := c.GetRequest("/reference/cost-groups")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CostGroupEdit method
func (c *Client) CostGroupEdit(costGroup CostGroup) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&costGroup)

	p := url.Values{
		"costGroup": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/cost-groups/%s/edit", costGroup.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CostItems method
func (c *Client) CostItems() (CostItemsResponse, int, errs.Failure) {
	var resp CostItemsResponse

	data, status, err := c.GetRequest("/reference/cost-items")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CostItemEdit method
func (c *Client) CostItemEdit(costItem CostItem) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&costItem)

	p := url.Values{
		"costItem": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/cost-items/%s/edit", costItem.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Couriers method
func (c *Client) Couriers() (CouriersResponse, int, errs.Failure) {
	var resp CouriersResponse

	data, status, err := c.GetRequest("/reference/couriers")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CourierCreate method
func (c *Client) CourierCreate(courier Courier) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse

	objJSON, _ := json.Marshal(&courier)

	p := url.Values{
		"courier": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest("/reference/couriers/create", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// CourierEdit method
func (c *Client) CourierEdit(courier Courier) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&courier)

	p := url.Values{
		"courier": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/couriers/%d/edit", courier.ID), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeliveryServices method
func (c *Client) DeliveryServices() (DeliveryServiceResponse, int, errs.Failure) {
	var resp DeliveryServiceResponse

	data, status, err := c.GetRequest("/reference/delivery-services")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeliveryServiceEdit method
func (c *Client) DeliveryServiceEdit(deliveryService DeliveryService) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&deliveryService)

	p := url.Values{
		"deliveryService": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/delivery-services/%s/edit", deliveryService.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeliveryTypes method
func (c *Client) DeliveryTypes() (DeliveryTypesResponse, int, errs.Failure) {
	var resp DeliveryTypesResponse

	data, status, err := c.GetRequest("/reference/delivery-types")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// DeliveryTypeEdit method
func (c *Client) DeliveryTypeEdit(deliveryType DeliveryType) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&deliveryType)

	p := url.Values{
		"deliveryType": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/delivery-types/%s/edit", deliveryType.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// LegalEntities method
func (c *Client) LegalEntities() (LegalEntitiesResponse, int, errs.Failure) {
	var resp LegalEntitiesResponse

	data, status, err := c.GetRequest("/reference/legal-entities")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// LegalEntityEdit method
func (c *Client) LegalEntityEdit(legalEntity LegalEntity) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&legalEntity)

	p := url.Values{
		"legalEntity": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/legal-entities/%s/edit", legalEntity.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrderMethods method
func (c *Client) OrderMethods() (OrderMethodsResponse, int, errs.Failure) {
	var resp OrderMethodsResponse

	data, status, err := c.GetRequest("/reference/order-methods")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrderMethodEdit method
func (c *Client) OrderMethodEdit(orderMethod OrderMethod) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&orderMethod)

	p := url.Values{
		"orderMethod": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/order-methods/%s/edit", orderMethod.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrderTypes method
func (c *Client) OrderTypes() (OrderTypesResponse, int, errs.Failure) {
	var resp OrderTypesResponse

	data, status, err := c.GetRequest("/reference/order-types")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// OrderTypeEdit method
func (c *Client) OrderTypeEdit(orderType OrderType) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&orderType)

	p := url.Values{
		"orderType": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/order-types/%s/edit", orderType.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PaymentStatuses method
func (c *Client) PaymentStatuses() (PaymentStatusesResponse, int, errs.Failure) {
	var resp PaymentStatusesResponse

	data, status, err := c.GetRequest("/reference/payment-statuses")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PaymentStatusEdit method
func (c *Client) PaymentStatusEdit(paymentStatus PaymentStatus) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&paymentStatus)

	p := url.Values{
		"paymentStatus": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/payment-statuses/%s/edit", paymentStatus.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PaymentTypes method
func (c *Client) PaymentTypes() (PaymentTypesResponse, int, errs.Failure) {
	var resp PaymentTypesResponse

	data, status, err := c.GetRequest("/reference/payment-types")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PaymentTypeEdit method
func (c *Client) PaymentTypeEdit(paymentType PaymentType) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&paymentType)

	p := url.Values{
		"paymentType": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/payment-types/%s/edit", paymentType.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PriceTypes method
func (c *Client) PriceTypes() (PriceTypesResponse, int, errs.Failure) {
	var resp PriceTypesResponse

	data, status, err := c.GetRequest("/reference/price-types")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PriceTypeEdit method
func (c *Client) PriceTypeEdit(priceType PriceType) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&priceType)

	p := url.Values{
		"priceType": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/price-types/%s/edit", priceType.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// ProductStatuses method
func (c *Client) ProductStatuses() (ProductStatusesResponse, int, errs.Failure) {
	var resp ProductStatusesResponse

	data, status, err := c.GetRequest("/reference/product-statuses")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// ProductStatusEdit method
func (c *Client) ProductStatusEdit(productStatus ProductStatus) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&productStatus)

	p := url.Values{
		"productStatus": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/product-statuses/%s/edit", productStatus.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Sites method
func (c *Client) Sites() (SitesResponse, int, errs.Failure) {
	var resp SitesResponse

	data, status, err := c.GetRequest("/reference/sites")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// SiteEdit method
func (c *Client) SiteEdit(site Site) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&site)

	p := url.Values{
		"site": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/sites/%s/edit", site.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// StatusGroups method
func (c *Client) StatusGroups() (StatusGroupsResponse, int, errs.Failure) {
	var resp StatusGroupsResponse

	data, status, err := c.GetRequest("/reference/status-groups")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Statuses method
func (c *Client) Statuses() (StatusesResponse, int, errs.Failure) {
	var resp StatusesResponse

	data, status, err := c.GetRequest("/reference/statuses")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// StatusEdit method
func (c *Client) StatusEdit(st Status) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&st)

	p := url.Values{
		"status": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/statuses/%s/edit", st.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Stores method
func (c *Client) Stores() (StoresResponse, int, errs.Failure) {
	var resp StoresResponse

	data, status, err := c.GetRequest("/reference/stores")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// StoreEdit method
func (c *Client) StoreEdit(store Store) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&store)

	p := url.Values{
		"store": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/stores/%s/edit", store.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Segments get segments
func (c *Client) Segments(parameters SegmentsRequest) (SegmentsResponse, int, errs.Failure) {
	var resp SegmentsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/segments?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Inventories method
func (c *Client) Inventories(parameters InventoriesRequest) (InventoriesResponse, int, errs.Failure) {
	var resp InventoriesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/inventories?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// InventoriesUpload method
func (c *Client) InventoriesUpload(inventories []InventoryUpload, site ...string) (StoreUploadResponse, int, errs.Failure) {
	var resp StoreUploadResponse

	uploadJSON, _ := json.Marshal(&inventories)

	p := url.Values{
		"offers": {string(uploadJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/store/inventories/upload", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// PricesUpload method
func (c *Client) PricesUpload(prices []OfferPriceUpload) (StoreUploadResponse, int, errs.Failure) {
	var resp StoreUploadResponse

	uploadJSON, _ := json.Marshal(&prices)

	p := url.Values{
		"prices": {string(uploadJSON[:])},
	}

	data, status, err := c.PostRequest("/store/prices/upload", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// ProductsGroup method
func (c *Client) ProductsGroup(parameters ProductsGroupsRequest) (ProductsGroupsResponse, int, errs.Failure) {
	var resp ProductsGroupsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/product-groups?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Products method
func (c *Client) Products(parameters ProductsRequest) (ProductsResponse, int, errs.Failure) {
	var resp ProductsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/products?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// ProductsProperties method
func (c *Client) ProductsProperties(parameters ProductsPropertiesRequest) (ProductsPropertiesResponse, int, errs.Failure) {
	var resp ProductsPropertiesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/products/properties?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Tasks list method
func (c *Client) Tasks(parameters TasksRequest) (TasksResponse, int, errs.Failure) {
	var resp TasksResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/tasks?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// TaskCreate method
func (c *Client) TaskCreate(task Task, site ...string) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse
	taskJSON, _ := json.Marshal(&task)

	p := url.Values{
		"task": {string(taskJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/tasks/create", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Task get method
func (c *Client) Task(id int) (TaskResponse, int, errs.Failure) {
	var resp TaskResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/tasks/%d", id))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// TaskEdit method
func (c *Client) TaskEdit(task Task, site ...string) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse
	var uid = strconv.Itoa(task.ID)

	taskJSON, _ := json.Marshal(&task)

	p := url.Values{
		"task": {string(taskJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/tasks/%s/edit", uid), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// UserGroups list method
func (c *Client) UserGroups(parameters UserGroupsRequest) (UserGroupsResponse, int, errs.Failure) {
	var resp UserGroupsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/user-groups?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// Users list method
func (c *Client) Users(parameters UsersRequest) (UsersResponse, int, errs.Failure) {
	var resp UsersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/users?%s", params.Encode()))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// User get method
func (c *Client) User(id int) (UserResponse, int, errs.Failure) {
	var resp UserResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/users/%d", id))
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

// UserStatus update method
func (c *Client) UserStatus(id int, status string) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	p := url.Values{
		"status": {string(status)},
	}

	data, st, err := c.PostRequest(fmt.Sprintf("/users/%d/status", id), p)
	if err.RuntimeErr != nil {
		return resp, st, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, st, buildErr(data)
	}

	return resp, st, err
}

// StaticticsUpdate update statistic
func (c *Client) StaticticsUpdate() (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	data, status, err := c.GetRequest("/statistic/update")
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	if resp.Success == false {
		return resp, status, buildErr(data)
	}

	return resp, status, err
}

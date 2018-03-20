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

// New initalize client
func New(url string, key string) *Client {
	return &Client{
		url,
		key,
		&http.Client{Timeout: 20 * time.Second},
	}
}

// GetRequest implements GET Request
func (c *Client) GetRequest(urlWithParameters string, versioned ...bool) ([]byte, int, ErrorResponse) {
	var res []byte
	var bug ErrorResponse
	var prefix = "/api/v5"

	if len(versioned) > 0 {
		s := versioned[0]

		if s == false {
			prefix = "/api"
		}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.URL, prefix, urlWithParameters), nil)
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
	var prefix = "/api/v5"

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s%s%s", c.URL, prefix, url),
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

// APIVersions get all available API versions for exact account
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.APIVersions()
//
// 	if err.ErrorMsg != "" {
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
func (c *Client) APIVersions() (VersionResponse, int, ErrorResponse) {
	var resp VersionResponse

	data, status, err := c.GetRequest("/api-versions", false)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

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
// 	if err.ErrorMsg != "" {
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
func (c *Client) APICredentials() (CredentialResponse, int, ErrorResponse) {
	var resp CredentialResponse

	data, status, err := c.GetRequest("/credentials", false)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Customers list method
func (c *Client) Customers(parameters CustomersRequest) (CustomersResponse, int, ErrorResponse) {
	var resp CustomersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomersCombine method
func (c *Client) CustomersCombine(customers []Customer, resultCustomer Customer) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	combineJSONIn, _ := json.Marshal(&customers)
	combineJSONOut, _ := json.Marshal(&resultCustomer)

	p := url.Values{
		"customers":      {string(combineJSONIn[:])},
		"resultCustomer": {string(combineJSONOut[:])},
	}

	data, status, err := c.PostRequest("/customers/combine", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomerCreate method
func (c *Client) CustomerCreate(customer Customer, site ...string) (CustomerChangeResponse, int, ErrorResponse) {
	var resp CustomerChangeResponse

	customerJSON, _ := json.Marshal(&customer)

	p := url.Values{
		"customer": {string(customerJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers/create", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomersFixExternalIds method
func (c *Client) CustomersFixExternalIds(customers []IdentifiersPair) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	customersJSON, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(customersJSON[:])},
	}

	data, status, err := c.PostRequest("/customers/fix-external-ids", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomersHistory method
func (c *Client) CustomersHistory(parameters CustomersHistoryRequest) (CustomersHistoryResponse, int, ErrorResponse) {
	var resp CustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers/history?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomerNotes list method
func (c *Client) CustomerNotes(parameters NotesRequest) (NotesResponse, int, ErrorResponse) {
	var resp NotesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers/notes?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomerNoteCreate method
func (c *Client) CustomerNoteCreate(note Note, site ...string) (CreateResponse, int, ErrorResponse) {
	var resp CreateResponse

	noteJSON, _ := json.Marshal(&note)

	p := url.Values{
		"note": {string(noteJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers/notes/create", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomerNoteDelete method
func (c *Client) CustomerNoteDelete(id int) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	p := url.Values{
		"id": {string(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/customers/notes/%d/delete", id), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomersUpload method
func (c *Client) CustomersUpload(customers []Customer, site ...string) (CustomersUploadResponse, int, ErrorResponse) {
	var resp CustomersUploadResponse

	uploadJSON, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(uploadJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/customers/upload", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Customer get method
func (c *Client) Customer(id, by, site string) (CustomerResponse, int, ErrorResponse) {
	var resp CustomerResponse
	var context = checkBy(by)

	fw := CustomerRequest{context, site}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("/customers/%s?%s", id, params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomerEdit method
func (c *Client) CustomerEdit(customer Customer, by string, site ...string) (CustomerChangeResponse, int, ErrorResponse) {
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
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// DeliveryTracking method
func (c *Client) DeliveryTracking(parameters DeliveryTrackingRequest, subcode string) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	updateJSON, _ := json.Marshal(&parameters)

	p := url.Values{
		"statusUpdate": {string(updateJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/delivery/generic/%s/tracking", subcode), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// DeliveryShipments method
func (c *Client) DeliveryShipments(parameters DeliveryShipmentsRequest) (DeliveryShipmentsResponse, int, ErrorResponse) {
	var resp DeliveryShipmentsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/delivery/shipments?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// DeliveryShipmentCreate method
func (c *Client) DeliveryShipmentCreate(shipment DeliveryShipment, deliveryType string, site ...string) (DeliveryShipmentUpdateResponse, int, ErrorResponse) {
	var resp DeliveryShipmentUpdateResponse
	updateJSON, _ := json.Marshal(&shipment)

	p := url.Values{
		"deliveryType":     {string(deliveryType)},
		"deliveryShipment": {string(updateJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/delivery/shipments/create", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// DeliveryShipment method
func (c *Client) DeliveryShipment(id int) (DeliveryShipmentResponse, int, ErrorResponse) {
	var resp DeliveryShipmentResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/delivery/shipments/%d", id))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// DeliveryShipmentEdit method
func (c *Client) DeliveryShipmentEdit(shipment DeliveryShipment, site ...string) (DeliveryShipmentUpdateResponse, int, ErrorResponse) {
	var resp DeliveryShipmentUpdateResponse
	updateJSON, _ := json.Marshal(&shipment)

	p := url.Values{
		"deliveryShipment": {string(updateJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/delivery/shipments/%s/edit", strconv.Itoa(shipment.ID)), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// IntegrationModule method
func (c *Client) IntegrationModule(code string) (IntegrationModuleResponse, int, ErrorResponse) {
	var resp IntegrationModuleResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/integration-modules/%s", code))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// IntegrationModuleEdit method
func (c *Client) IntegrationModuleEdit(integrationModule IntegrationModule) (IntegrationModuleEditResponse, int, ErrorResponse) {
	var resp IntegrationModuleEditResponse
	updateJSON, _ := json.Marshal(&integrationModule)

	p := url.Values{"integrationModule": {string(updateJSON[:])}}

	data, status, err := c.PostRequest(fmt.Sprintf("/integration-modules/%s/edit", integrationModule.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Orders list method
func (c *Client) Orders(parameters OrdersRequest) (OrdersResponse, int, ErrorResponse) {
	var resp OrdersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrdersCombine method
func (c *Client) OrdersCombine(technique string, order, resultOrder Order) (OperationResponse, int, ErrorResponse) {
	var resp OperationResponse

	combineJSONIn, _ := json.Marshal(&order)
	combineJSONOut, _ := json.Marshal(&resultOrder)

	p := url.Values{
		"technique":   {technique},
		"order":       {string(combineJSONIn[:])},
		"resultOrder": {string(combineJSONOut[:])},
	}

	data, status, err := c.PostRequest("/orders/combine", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrderCreate method
func (c *Client) OrderCreate(order Order, site ...string) (CreateResponse, int, ErrorResponse) {
	var resp CreateResponse
	orderJSON, _ := json.Marshal(&order)

	p := url.Values{
		"order": {string(orderJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/create", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrdersFixExternalIds method
func (c *Client) OrdersFixExternalIds(orders []IdentifiersPair) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	ordersJSON, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(ordersJSON[:])},
	}

	data, status, err := c.PostRequest("/orders/fix-external-ids", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrdersHistory method
func (c *Client) OrdersHistory(parameters OrdersHistoryRequest) (CustomersHistoryResponse, int, ErrorResponse) {
	var resp CustomersHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/history?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrderPaymentCreate method
func (c *Client) OrderPaymentCreate(payment Payment, site ...string) (CreateResponse, int, ErrorResponse) {
	var resp CreateResponse

	paymentJSON, _ := json.Marshal(&payment)

	p := url.Values{
		"payment": {string(paymentJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/payments/create", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrderPaymentDelete method
func (c *Client) OrderPaymentDelete(id int) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	p := url.Values{
		"id": {string(id)},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/payments/%d/delete", id), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrderPaymentEdit method
func (c *Client) OrderPaymentEdit(payment Payment, by string, site ...string) (SuccessfulResponse, int, ErrorResponse) {
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
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrdersUpload method
func (c *Client) OrdersUpload(orders []Order, site ...string) (OrdersUploadResponse, int, ErrorResponse) {
	var resp OrdersUploadResponse

	uploadJSON, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(uploadJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/orders/upload", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Order get method
func (c *Client) Order(id, by, site string) (OrderResponse, int, ErrorResponse) {
	var resp OrderResponse
	var context = checkBy(by)

	fw := OrderRequest{context, site}
	params, _ := query.Values(fw)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/%s?%s", id, params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrderEdit method
func (c *Client) OrderEdit(order Order, by string, site ...string) (CreateResponse, int, ErrorResponse) {
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
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Packs list method
func (c *Client) Packs(parameters PacksRequest) (PacksResponse, int, ErrorResponse) {
	var resp PacksResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/packs?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PackCreate method
func (c *Client) PackCreate(pack Pack) (CreateResponse, int, ErrorResponse) {
	var resp CreateResponse
	packJSON, _ := json.Marshal(&pack)

	p := url.Values{
		"pack": {string(packJSON[:])},
	}

	data, status, err := c.PostRequest("/orders/packs/create", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PacksHistory method
func (c *Client) PacksHistory(parameters PacksHistoryRequest) (PacksHistoryResponse, int, ErrorResponse) {
	var resp PacksHistoryResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/packs/history?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Pack get method
func (c *Client) Pack(id int) (PackResponse, int, ErrorResponse) {
	var resp PackResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/orders/packs/%d", id))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PackDelete method
func (c *Client) PackDelete(id int) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/packs/%d/delete", id), url.Values{})
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PackEdit method
func (c *Client) PackEdit(pack Pack) (CreateResponse, int, ErrorResponse) {
	var resp CreateResponse

	packJSON, _ := json.Marshal(&pack)

	p := url.Values{
		"pack": {string(packJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/orders/packs/%d/edit", pack.ID), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Countries method
func (c *Client) Countries() (CountriesResponse, int, ErrorResponse) {
	var resp CountriesResponse

	data, status, err := c.GetRequest("/reference/countries")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CostGroups method
func (c *Client) CostGroups() (CostGroupsResponse, int, ErrorResponse) {
	var resp CostGroupsResponse

	data, status, err := c.GetRequest("/reference/cost-groups")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CostGroupEdit method
func (c *Client) CostGroupEdit(costGroup CostGroup) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&costGroup)

	p := url.Values{
		"costGroup": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/cost-groups/%s/edit", costGroup.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CostItems method
func (c *Client) CostItems() (CostItemsResponse, int, ErrorResponse) {
	var resp CostItemsResponse

	data, status, err := c.GetRequest("/reference/cost-items")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CostItemEdit method
func (c *Client) CostItemEdit(costItem CostItem) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&costItem)

	p := url.Values{
		"costItem": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/cost-items/%s/edit", costItem.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Couriers method
func (c *Client) Couriers() (CouriersResponse, int, ErrorResponse) {
	var resp CouriersResponse

	data, status, err := c.GetRequest("/reference/couriers")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CourierCreate method
func (c *Client) CourierCreate(courier Courier) (CreateResponse, int, ErrorResponse) {
	var resp CreateResponse

	objJSON, _ := json.Marshal(&courier)

	p := url.Values{
		"courier": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest("/reference/couriers/create", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CourierEdit method
func (c *Client) CourierEdit(courier Courier) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&courier)

	p := url.Values{
		"courier": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/couriers/%d/edit", courier.ID), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// DeliveryServices method
func (c *Client) DeliveryServices() (DeliveryServiceResponse, int, ErrorResponse) {
	var resp DeliveryServiceResponse

	data, status, err := c.GetRequest("/reference/delivery-services")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// DeliveryServiceEdit method
func (c *Client) DeliveryServiceEdit(deliveryService DeliveryService) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&deliveryService)

	p := url.Values{
		"deliveryService": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/delivery-services/%s/edit", deliveryService.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// DeliveryTypes method
func (c *Client) DeliveryTypes() (DeliveryTypesResponse, int, ErrorResponse) {
	var resp DeliveryTypesResponse

	data, status, err := c.GetRequest("/reference/delivery-types")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// DeliveryTypeEdit method
func (c *Client) DeliveryTypeEdit(deliveryType DeliveryType) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&deliveryType)

	p := url.Values{
		"deliveryType": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/delivery-types/%s/edit", deliveryType.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// LegalEntities method
func (c *Client) LegalEntities() (LegalEntitiesResponse, int, ErrorResponse) {
	var resp LegalEntitiesResponse

	data, status, err := c.GetRequest("/reference/legal-entities")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// LegalEntityEdit method
func (c *Client) LegalEntityEdit(legalEntity LegalEntity) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&legalEntity)

	p := url.Values{
		"legalEntity": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/legal-entities/%s/edit", legalEntity.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrderMethods method
func (c *Client) OrderMethods() (OrderMethodsResponse, int, ErrorResponse) {
	var resp OrderMethodsResponse

	data, status, err := c.GetRequest("/reference/order-methods")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrderMethodEdit method
func (c *Client) OrderMethodEdit(orderMethod OrderMethod) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&orderMethod)

	p := url.Values{
		"orderMethod": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/order-methods/%s/edit", orderMethod.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrderTypes method
func (c *Client) OrderTypes() (OrderTypesResponse, int, ErrorResponse) {
	var resp OrderTypesResponse

	data, status, err := c.GetRequest("/reference/order-types")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// OrderTypeEdit method
func (c *Client) OrderTypeEdit(orderType OrderType) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&orderType)

	p := url.Values{
		"orderType": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/order-types/%s/edit", orderType.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PaymentStatuses method
func (c *Client) PaymentStatuses() (PaymentStatusesResponse, int, ErrorResponse) {
	var resp PaymentStatusesResponse

	data, status, err := c.GetRequest("/reference/payment-statuses")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PaymentStatusEdit method
func (c *Client) PaymentStatusEdit(paymentStatus PaymentStatus) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&paymentStatus)

	p := url.Values{
		"paymentStatus": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/payment-statuses/%s/edit", paymentStatus.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PaymentTypes method
func (c *Client) PaymentTypes() (PaymentTypesResponse, int, ErrorResponse) {
	var resp PaymentTypesResponse

	data, status, err := c.GetRequest("/reference/payment-types")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PaymentTypeEdit method
func (c *Client) PaymentTypeEdit(paymentType PaymentType) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&paymentType)

	p := url.Values{
		"paymentType": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/payment-types/%s/edit", paymentType.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PriceTypes method
func (c *Client) PriceTypes() (PriceTypesResponse, int, ErrorResponse) {
	var resp PriceTypesResponse

	data, status, err := c.GetRequest("/reference/price-types")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PriceTypeEdit method
func (c *Client) PriceTypeEdit(priceType PriceType) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&priceType)

	p := url.Values{
		"priceType": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/price-types/%s/edit", priceType.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// ProductStatuses method
func (c *Client) ProductStatuses() (ProductStatusesResponse, int, ErrorResponse) {
	var resp ProductStatusesResponse

	data, status, err := c.GetRequest("/reference/product-statuses")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// ProductStatusEdit method
func (c *Client) ProductStatusEdit(productStatus ProductStatus) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&productStatus)

	p := url.Values{
		"productStatus": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/product-statuses/%s/edit", productStatus.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Sites method
func (c *Client) Sites() (SitesResponse, int, ErrorResponse) {
	var resp SitesResponse

	data, status, err := c.GetRequest("/reference/sites")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// SiteEdit method
func (c *Client) SiteEdit(site Site) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&site)

	p := url.Values{
		"site": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/sites/%s/edit", site.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// StatusGroups method
func (c *Client) StatusGroups() (StatusGroupsResponse, int, ErrorResponse) {
	var resp StatusGroupsResponse

	data, status, err := c.GetRequest("/reference/status-groups")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Statuses method
func (c *Client) Statuses() (StatusesResponse, int, ErrorResponse) {
	var resp StatusesResponse

	data, status, err := c.GetRequest("/reference/statuses")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// StatusEdit method
func (c *Client) StatusEdit(st Status) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&st)

	p := url.Values{
		"status": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/statuses/%s/edit", st.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Stores method
func (c *Client) Stores() (StoresResponse, int, ErrorResponse) {
	var resp StoresResponse

	data, status, err := c.GetRequest("/reference/stores")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// StoreEdit method
func (c *Client) StoreEdit(store Store) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	objJSON, _ := json.Marshal(&store)

	p := url.Values{
		"store": {string(objJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/reference/stores/%s/edit", store.Code), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Segments get segments
func (c *Client) Segments(parameters SegmentsRequest) (SegmentsResponse, int, ErrorResponse) {
	var resp SegmentsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/segments?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Inventories method
func (c *Client) Inventories(parameters InventoriesRequest) (InventoriesResponse, int, ErrorResponse) {
	var resp InventoriesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/inventories?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// InventoriesUpload method
func (c *Client) InventoriesUpload(inventories []InventoryUpload, site ...string) (StoreUploadResponse, int, ErrorResponse) {
	var resp StoreUploadResponse

	uploadJSON, _ := json.Marshal(&inventories)

	p := url.Values{
		"offers": {string(uploadJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/store/inventories/upload", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// PricesUpload method
func (c *Client) PricesUpload(prices []OfferPriceUpload) (StoreUploadResponse, int, ErrorResponse) {
	var resp StoreUploadResponse

	uploadJSON, _ := json.Marshal(&prices)

	p := url.Values{
		"prices": {string(uploadJSON[:])},
	}

	data, status, err := c.PostRequest("/store/prices/upload", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// ProductsGroup method
func (c *Client) ProductsGroup(parameters ProductsGroupsRequest) (ProductsGroupsResponse, int, ErrorResponse) {
	var resp ProductsGroupsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/product-groups?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Products method
func (c *Client) Products(parameters ProductsRequest) (ProductsResponse, int, ErrorResponse) {
	var resp ProductsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/products?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// ProductsProperties method
func (c *Client) ProductsProperties(parameters ProductsPropertiesRequest) (ProductsPropertiesResponse, int, ErrorResponse) {
	var resp ProductsPropertiesResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/store/products/properties?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Tasks list method
func (c *Client) Tasks(parameters TasksRequest) (TasksResponse, int, ErrorResponse) {
	var resp TasksResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/tasks?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// TaskCreate method
func (c *Client) TaskCreate(task Task, site ...string) (CreateResponse, int, ErrorResponse) {
	var resp CreateResponse
	taskJSON, _ := json.Marshal(&task)

	p := url.Values{
		"task": {string(taskJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/tasks/create", p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Task get method
func (c *Client) Task(id int) (TaskResponse, int, ErrorResponse) {
	var resp TaskResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/tasks/%d", id))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// TaskEdit method
func (c *Client) TaskEdit(task Task, site ...string) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse
	var uid = strconv.Itoa(task.ID)

	taskJSON, _ := json.Marshal(&task)

	p := url.Values{
		"task": {string(taskJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/tasks/%s/edit", uid), p)
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// UserGroups list method
func (c *Client) UserGroups(parameters UserGroupsRequest) (UserGroupsResponse, int, ErrorResponse) {
	var resp UserGroupsResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/user-groups?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Users list method
func (c *Client) Users(parameters UsersRequest) (UsersResponse, int, ErrorResponse) {
	var resp UsersResponse

	params, _ := query.Values(parameters)

	data, status, err := c.GetRequest(fmt.Sprintf("/users?%s", params.Encode()))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// User get method
func (c *Client) User(id int) (UserResponse, int, ErrorResponse) {
	var resp UserResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/users/%d", id))
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// UserStatus update method
func (c *Client) UserStatus(id int, status string) (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	p := url.Values{
		"status": {string(status)},
	}

	data, st, err := c.PostRequest(fmt.Sprintf("/users/%d/status", id), p)
	if err.ErrorMsg != "" {
		return resp, st, err
	}

	json.Unmarshal(data, &resp)

	return resp, st, err
}

// StaticticsUpdate update statistic
func (c *Client) StaticticsUpdate() (SuccessfulResponse, int, ErrorResponse) {
	var resp SuccessfulResponse

	data, status, err := c.GetRequest("/statistic/update")
	if err.ErrorMsg != "" {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

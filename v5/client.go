package v5

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		URL:        url,
		Key:        key,
		httpClient: &http.Client{Timeout: 20 * time.Second},
	}
}

// GetRequest implements GET Request
func (c *Client) GetRequest(urlWithParameters string, versioned ...bool) ([]byte, int, errs.Failure) {
	var res []byte
	var failure errs.Failure
	var prefix = "/api/v5"

	if len(versioned) > 0 {
		s := versioned[0]

		if s == false {
			prefix = "/api"
		}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", c.URL, prefix, urlWithParameters), nil)
	if err != nil {
		failure.RuntimeErr = err
		return res, 0, failure
	}

	req.Header.Set("X-API-KEY", c.Key)

	if c.Debug {
		log.Printf("API Request: %s %s", fmt.Sprintf("%s%s%s", c.URL, prefix, urlWithParameters), c.Key)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		failure.RuntimeErr = err
		return res, 0, failure
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		failure.ApiErr = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, failure
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		failure.RuntimeErr = err
	}

	if c.Debug {
		log.Printf("API Response: %s", res)
	}

	return res, resp.StatusCode, failure
}

// PostRequest implements POST Request
func (c *Client) PostRequest(url string, postParams url.Values) ([]byte, int, errs.Failure) {
	var res []byte
	var failure errs.Failure
	var prefix = "/api/v5"

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", c.URL, prefix, url), strings.NewReader(postParams.Encode()))
	if err != nil {
		failure.RuntimeErr = err
		return res, 0, failure
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-API-KEY", c.Key)

	if c.Debug {
		log.Printf("API Request: %s %s %s", url, c.Key, postParams.Encode())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		failure.RuntimeErr = err
		return res, 0, failure
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		failure.ApiErr = fmt.Sprintf("HTTP request error. Status code: %d.\n", resp.StatusCode)
		return res, resp.StatusCode, failure
	}

	res, err = buildRawResponse(resp)
	if err != nil {
		failure.RuntimeErr = err
		return res, 0, failure
	}

	if c.Debug {
		log.Printf("API Response: %s", res)
	}

	return res, resp.StatusCode, failure
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

	resp, e := errs.ErrorResponse(data)
	err.RuntimeErr = e
	err.ApiErr = resp.ErrorMsg

	if resp.Errors != nil {
		err.ApiErrs = resp.Errors
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
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-api-versions
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.APIVersions()
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
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
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-api-versions
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.APICredentials()
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
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

// Customers returns list of customers matched the specified filter
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-customers
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Customers(v5.CustomersRequest{
//		Filter: CustomersFilter{
//			City: "Moscow",
//		},
//		Page: 3,
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	for _, value := range data.Customers {
// 		fmt.Printf("%v\n", value)
// 	}
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

// CustomersCombine combine given customers
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-customers-combine
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CustomersCombine([]v5.Customer{{ID: 1}, {ID: 2}}, Customer{ID: 3})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// CustomerCreate creates customer
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-customers-create
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CustomersCombine(v5.Customer{
//		FirstName:  "Ivan",
//		LastName:   "Ivanov",
//		Patronymic: "Ivanovich",
//		ExternalID: 1,
//		Email:      "ivanov@example.com",
//		Address: &v5.Address{
//			City:     "Moscow",
//			Street:   "Kutuzovsky",
//		},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	if data.Success == true {
//		fmt.Printf("%v", err.Id)
//	}
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

// CustomersFixExternalIds customers external ID
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-customers-fix-external-ids
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CustomersFixExternalIds([]v5.IdentifiersPair{{
//		ID:         1,
//		ExternalID: 12,
//	}})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// CustomersHistory returns customer's history
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-customers-history
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CustomersHistory(v5.CustomersHistoryRequest{
//		Filter: v5.CustomersHistoryFilter{
//			SinceID: 20,
//		},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	for _, value := range data.History {
// 		fmt.Printf("%v\n", value)
// 	}
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

// CustomerNotes returns customer related notes
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-customers-notes
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CustomerNotes(v5.NotesRequest{
//		Filter: v5.NotesFilter{
//			CustomerIds: []int{1,2,3}
// 		},
//		Page: 1,
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	for _, value := range data.Notes {
// 		fmt.Printf("%v\n", value)
// 	}
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

// CustomerNoteCreate note creation
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-customers-notes-create
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CustomerNoteCreate(v5.Note{
//		Text:      "some text",
//		ManagerID: 12,
//		Customer: &v5.Customer{
//			ID: 1,
//		},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	if data.Success == true {
// 		fmt.Printf("%v\n", data.ID)
// 	}
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

// CustomerNoteDelete remove customer related note
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-customers-notes-id-delete
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CustomerNoteDelete(12)
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// CustomersUpload customers batch upload
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-customers-upload
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CustomersUpload([]v5.Customer{
//		{
//			FirstName:  "Ivan",
//			LastName:   "Ivanov",
//			Patronymic: "Ivanovich",
//			ExternalID: 1,
//			Email:      "ivanov@example.com",
//		},
//		{
//			FirstName:  "Petr",
//			LastName:   "Petrov",
//			Patronymic: "Petrovich",
//			ExternalID: 2,
//			Email:      "petrov@example.com",
//		},
//	}}
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	if data.Success == true {
// 		fmt.Printf("%v\n", data.UploadedCustomers)
// 	}
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

// Customer returns information about customer
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-customers-externalId
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Customer(12, "externalId", "")
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	if data.Success == true {
// 		fmt.Printf("%v\n", data.Customer)
// 	}
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

// CustomerEdit edit exact customer
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-customers-externalId-edit
//
// Example:
//
//	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomerEdit(
//		v5.Customer{
//			FirstName:  "Ivan",
//			LastName:   "Ivanov",
//			Patronymic: "Ivanovich",
//			ID: 		1,
//			Email:      "ivanov@example.com",
//		},
//		"id",
//	)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	if data.Success == true {
// 		fmt.Printf("%v\n", data.Customer)
// 	}
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

// DeliveryTracking updates tracking data
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-delivery-generic-subcode-tracking
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.DeliveryTracking(
//		v5.DeliveryTrackingRequest{
//			DeliveryID: "1",
//			TrackNumber "123",
//			History: []DeliveryHistoryRecord{
//				{
//					Code: "cancel",
//					UpdatedAt: "2012-12-12 12:12:12"
//				},
//			}
//		},
//		"delivery-1",
//	)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
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

// DeliveryShipments returns list of shipments
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-delivery-shipments
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.DeliveryShipments(v5.DeliveryShipmentsRequest{
//		Limit: 12,
//		Filter: v5.ShipmentFilter{
//			DateFrom: "2012-12-12",
//		},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	for _, value := range data.DeliveryShipments {
// 		fmt.Printf("%v\n", value)
// 	}
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

// DeliveryShipmentCreate creates shipment
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-delivery-shipments-create
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.DeliveryShipmentCreate(
//		v5.DeliveryShipment{
//			Date: "2012-12-12",
//			Time: v5.DeliveryTime{
//				From: "18:00",
//				To: "20:00",
//			},
//			Orders: []v5.Order{{Number: "12"}},
//		},
//		"sdek",
//	)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	if data.Success == true {
// 		fmt.Printf("%v\n", data.ID)
// 	}
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

// DeliveryShipment get information about shipment
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-delivery-shipments-id
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.DeliveryShipment(12)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	if data.Success == true {
// 		fmt.Printf("%v\n", data.DeliveryShipment)
// 	}
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

// DeliveryShipmentEdit shipment editing
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-delivery-shipments-id-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.DeliveryShipmentEdit(v5.DeliveryShipment{
//		ID: "12",
//		Time: v5.DeliveryTime{
//			From: "14:00",
//			To: "18:00",
//		},
// 	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// IntegrationModule returns integration module
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-integration-modules-code
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.IntegrationModule("moysklad3")
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	if data.Success == true {
// 		fmt.Printf("%v\n", data.IntegrationModule)
// 	}
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

// IntegrationModuleEdit integration module create/edit
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-integration-modules-code
//
// Example:
//
//	var client = v5.New("https://demo.url", "09jIJ")
//
//	name := "MS"
//	code := "moysklad3"
//
//	data, status, err := client.IntegrationModuleEdit(v5.IntegrationModule{
//		Code:            code,
//		IntegrationCode: code,
//		Active:          false,
//		Name:            fmt.Sprintf("Integration module %s", name),
//		AccountURL:      fmt.Sprintf("http://example.com/%s/account", name),
//		BaseURL:         fmt.Sprintf("http://example.com/%s", name),
//		ClientID:        "123",
//		Logo:            "https://cdn.worldvectorlogo.com/logos/github-icon.svg",
//	})
//
//	if err.RuntimeErr != nil {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.ApiErr())
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v\n", data.Info)
//	}
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

// Orders returns list of orders matched the specified filters
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-orders
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Orders(v5.OrdersRequest{Filter: v5.OrdersFilter{City: "Moscow"}, Page: 1})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	for _, value := range data.Orders {
// 		fmt.Printf("%v\n", value)
// 	}
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

// OrdersCombine combines given orders
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-combine
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrdersCombine("ours", v5.Order{ID: 1}, v5.Order{ID: 1})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// OrderCreate creates an order
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-create
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrderCreate(v5.Order{
//		FirstName:  "Ivan",
//		LastName:   "Ivanov",
//		Patronymic: "Ivanovich",
//		Email:      "ivanov@example.com",
//		Items:      []v5.OrderItem{{Offer: v5.Offer{ID: 12}, Quantity: 5}},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	if data.Success == true {
// 		fmt.Printf("%v\n", data.ID)
// 	}
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

// OrdersFixExternalIds set order external ID
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-fix-external-ids
//
// Example:
//
//	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.OrdersFixExternalIds(([]v5.IdentifiersPair{{
//		ID:         1,
//		ExternalID: 12,
//	}})
//
//	if err.RuntimeErr != nil {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.ApiErr())
//	}
//
//	if data.Success == true {
//		fmt.Printf("%v\n", data.ID)
//	}
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

// OrdersHistory returns orders history
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-orders-history
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrdersHistory(v5.OrdersHistoryRequest{Filter: v5.OrdersHistoryFilter{SinceID: 20}})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	for _, value := range data.History {
// 		fmt.Printf("%v\n", value)
// 	}
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

// OrderPaymentCreate creates payment
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-payments-create
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrderPaymentCreate(v5.Payment{
//		Order: &v5.Order{
//			ID: 12,
//		},
//		Amount: 300,
//		Type:   "cash",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	if data.Success == true {
// 		fmt.Printf("%v\n", data.ID)
// 	}
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

// OrderPaymentDelete payment removing
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-payments-id-delete
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrderPaymentDelete(12)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// OrderPaymentEdit edit payment
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-payments-id-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrderPaymentEdit(
// 		v5.Payment{
//			ID:     12,
//			Amount: 500,
//		},
//		"id",
//	)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// OrdersUpload batch orders uploading
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-upload
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrdersUpload([]v5.Order{
//		{
//			FirstName:  "Ivan",
//			LastName:   "Ivanov",
//			Patronymic: "Ivanovich",
//			Email:      "ivanov@example.com",
//			Items:      []v5.OrderItem{{Offer: v5.Offer{ID: 12}, Quantity: 5}},
//		},
//		{
//			FirstName:  "Pert",
//			LastName:   "Petrov",
//			Patronymic: "Petrovich",
//			Email:      "petrov@example.com",
//			Items:      []v5.OrderItem{{Offer: v5.Offer{ID: 13}, Quantity: 1}},
//		}
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	if data.Success == true {
// 		fmt.Printf("%v\n", data.UploadedOrders)
// 	}
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

// Order returns information about order
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-orders-externalId
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Order(12, "externalId", "")
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	if data.Success == true {
// 		fmt.Printf("%v\n", data.Order)
// 	}
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

// OrderEdit edit an order
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-externalId-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrderEdit(
// 		v5.Order{
//			ID:    12,
//			Items: []v5.OrderItem{{Offer: v5.Offer{ID: 13}, Quantity: 6}},
//		},
// 		"id",
// 	)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// Packs returns list of packs matched the specified filters
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-orders-packs
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Packs(v5.PacksRequest{Filter: v5.PacksFilter{OrderID: 12}})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	for _, value := range data.Packs {
// 		fmt.Printf("%v\n", value)
// 	}
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

// PackCreate creates a pack
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-packs-create
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.PackCreate(Pack{
//		Store:    "store-1",
//		ItemID:   12,
//		Quantity: 1,
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	if data.Success == true {
// 		fmt.Printf("%v\n", data.ID)
// 	}
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

// PacksHistory returns a history of order packing
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-orders-packs-history
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.PacksHistory(v5.PacksHistoryRequest{Filter: v5.OrdersHistoryFilter{SinceID: 5}})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	for _, value := range data.History {
// 		fmt.Printf("%v\n", value)
// 	}
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

// Pack returns a pack info
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-orders-packs-id
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Pack(112)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	if data.Success == true {
// 		fmt.Printf("%v\n", data.Pack)
// 	}
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

// PackDelete removes a pack
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-packs-id-delete
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.PackDelete(112)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// PackEdit edit a pack
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-orders-packs-id-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.PackEdit(Pack{ID: 12, Quantity: 2})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// Countries returns list of available country codes
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-countries
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

// CostGroups returns costs groups list
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-cost-groups
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

// CostGroupEdit edit costs groups
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-cost-groups-code-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CostGroupEdit(v5.CostGroup{
//		Code:   "group-1",
//		Color:  "#da5c98",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// CostItems returns costs items list
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-cost-items
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

// CostItemEdit edit costs items
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-cost-items-code-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CostItemEdit(v5.CostItem{
//		Code:            "seo",
//		Active:          false,
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// Couriers returns list of couriers
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-couriers
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

// CourierCreate creates a courier
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-couriers
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CostItemEdit(v5.Courier{
//		Active:    true,
//		Email:     "courier1@example.com",
//		FirstName: "Ivan",
//		LastName:  "Ivanov",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	if data.Success == true {
// 		fmt.Printf("%v", data.ID)
// 	}
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

// CourierEdit edit a courier
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-couriers-id-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.CostItemEdit(v5.Courier{
//		ID:    	  1,
//		Patronymic: "Ivanovich",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// DeliveryServices returns list of delivery services
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-delivery-services
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

// DeliveryServiceEdit delivery service create/edit
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-delivery-services-code-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.DeliveryServiceEdit(v5.DeliveryService{
//		Active: false,
//		Code:   "delivery-1",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// DeliveryTypes returns list of delivery types
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-delivery-types
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

// DeliveryTypeEdit delivery type create/edit
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-delivery-types-code-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.DeliveryTypeEdit(v5.DeliveryType{
//		Active:        false,
//		Code:          "type-1",
//		DefaultCost:   300,
//		DefaultForCrm: false,
//	}
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// LegalEntities returns list of legal entities
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-legal-entities
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

// LegalEntityEdit change information about legal entity
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-legal-entities-code-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.LegalEntityEdit(v5.LegalEntity{
//		Code:          "legal-entity-1",
//		CertificateDate:   "2012-12-12",
//	}
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// OrderMethods returns list of order methods
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-order-methods
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

// OrderMethodEdit order method create/edit
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-order-methods-code-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrderMethodEdit(v5.OrderMethod{
//		Code:          "method-1",
//		Active:        false,
//		DefaultForCRM: false,
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// OrderTypes return list of order types
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-order-types
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

// OrderTypeEdit create/edit order type
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-order-methods-code-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.OrderTypeEdit(v5.OrderType{
//		Code:          "order-type-1",
//		Active:        false,
//		DefaultForCRM: false,
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// PaymentStatuses returns list of payment statuses
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-payment-statuses
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

// PaymentStatusEdit payment status creation/editing
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-payment-statuses-code-edit
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

// PaymentTypes returns list of payment types
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-payment-types
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

// PaymentTypeEdit payment type create/edit
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-payment-types-code-edit
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

// PriceTypes returns list of price types
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-price-types
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

// PriceTypeEdit price type create/edit
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-price-types-code-edit
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

// ProductStatuses returns list of item statuses in order
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-product-statuses
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

// ProductStatusEdit order item status create/edit
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-product-statuses-code-edit
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

// Sites returns the sites list
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-sites
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

// SiteEdit site create/edit
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-sites-code-edit
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

// StatusGroups returns list of order status groups
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-status-groups
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

// Statuses returns list of order statuses
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-statuses
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

// StatusEdit order status create/edit
//
// For more information see www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-sites-code-edit
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

// Stores returns list of warehouses
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-reference-stores
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

// StoreEdit warehouse create/edit
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-reference-stores-code-edit
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

// Segments returns segments
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-segments
//
// Example:
//
//	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Segments(SegmentsRequest{
//		Filter: v5.SegmentsFilter{
//			Ids: []int{1,2,3}
//		}
//	})
//
//	if err.RuntimeErr != nil {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.ApiErr())
//	}
//
//	for _, value := range data.Segments {
//		fmt.Printf("%v\n", value)
//	}
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

// Inventories returns leftover stocks and purchasing prices
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-store-inventories
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Inventories(v5.InventoriesRequest{Filter: v5.InventoriesFilter{Details: 1, ProductActive: 1}, Page: 1})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	for _, value := range data.Offers {
// 		fmt.Printf("%v\n", value)
// 	}
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

// InventoriesUpload updates the leftover stocks and purchasing prices
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-store-inventories-upload
//
// Example:
//
//	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.InventoriesUpload(
//	   []v5.InventoryUpload{
//		   {
//			   XMLID: "pT22K9YzX21HTdzFCe1",
//			   Stores: []InventoryUploadStore{
//				   {Code: "test-store-v5", Available: 10, PurchasePrice: 1500},
//				   {Code: "test-store-v4", Available: 20, PurchasePrice: 1530},
//				   {Code: "test-store", Available: 30, PurchasePrice: 1510},
//			   },
//		   },
//		   {
//			   XMLID: "JQICtiSpOV3AAfMiQB3",
//			   Stores: []InventoryUploadStore{
//				   {Code: "test-store-v5", Available: 45, PurchasePrice: 1500},
//				   {Code: "test-store-v4", Available: 32, PurchasePrice: 1530},
//				   {Code: "test-store", Available: 46, PurchasePrice: 1510},
//			   },
//		   },
//	   },
//	)
//
//	if err.RuntimeErr != nil {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.ApiErr())
//	}
//
//	fmt.Printf("%v\n", data.NotFoundOffers)
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

// PricesUpload updates prices
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-store-prices-upload
//
// Example:
//
//	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.PricesUpload([]v5.OfferPriceUpload{
//		{
//			ID         1
//			Site       "store-1"
//			Prices     []PriceUpload{{Code:  "price-1"}}
//		},
//		{
//			ID         2
//			Site       "store-1"
//			Prices     []PriceUpload{{Code:  "price-2"}}
//		},
//	})
//
//	if err.RuntimeErr != nil {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.ApiErr())
//	}
//
//	fmt.Printf("%v\n", data.NotFoundOffers)
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

// ProductsGroup returns list of product groups
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-store-product-groups
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.ProductsGroup(v5.ProductsGroupsRequest{
//		Filter: v5.ProductsGroupsFilter{
//			Active: 1,
//		},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	for _, value := range data.ProductGroup {
// 		fmt.Printf("%v\n", value)
// 	}
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

// Products returns list of products and SKU
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-store-products
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Products(v5.ProductsRequest{
//		Filter: v5.ProductsFilter{
//			Active:   1,
//			MinPrice: 1000,
//		},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	for _, value := range data.Products {
// 		fmt.Printf("%v\n", value)
// 	}
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

// ProductsProperties returns list of item properties, matching the specified filters
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-store-products-properties
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.ProductsProperties(v5.ProductsPropertiesRequest{
//		Filter: v5.ProductsPropertiesFilter{
//			Sites: []string["store"],
//		},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	for _, value := range data.Properties {
// 		fmt.Printf("%v\n", value)
// 	}
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

// Tasks returns task list
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-tasks
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Tasks(v5.TasksRequest{
//		Filter: TasksFilter{
//			DateFrom: "2012-12-12",
//		},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	for _, value := range data.Tasks {
// 		fmt.Printf("%v\n", value)
// 	}
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

// TaskCreate create a task
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-tasks-create
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Tasks(v5.Task{
//		Text:        "task №1",
//		PerformerID: 12,
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	if data.Success == true {
// 		fmt.Printf("%v\n", data.ID)
// 	}
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

// Task returns task
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-tasks-id
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Task(12)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	if data.Success == true {
// 		fmt.Printf("%v\n", data.Task)
// 	}
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

// TaskEdit edit a task
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-tasks-id-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Task(v5.Task{
//		ID:   12
//		Text: "task №2",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// UserGroups returns list of user groups
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-user-groups
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.UserGroups(v5.UserGroupsRequest{Page: 1})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	for _, value := range data.Groups {
// 		fmt.Printf("%v\n", value)
// 	}
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

// Users returns list of users matched the specified filters
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-users
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.Users(v5.UsersRequest{Filter: v5.UsersFilter{Active: 1}, Page: 1})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	for _, value := range data.Users {
// 		fmt.Printf("%v\n", value)
// 	}
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

// User returns information about user
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-users-id
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.User(12)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
//	if data.Success == true {
// 		fmt.Printf("%v\n", data.User)
// 	}
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

// UserStatus change user status
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-users
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
// 	data, status, err := client.UserStatus(12, "busy")
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
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

// StaticticsUpdate updates statistics
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-statistic-update
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

// Costs returns costs list
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-costs
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Costs(CostsRequest{
//		Filter: CostsFilter{
//			Ids: []string{"1","2","3"},
//			MinSumm: "1000"
//		},
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	for _, value := range data.Costs {
// 		fmt.Printf("%v\n", value.Summ)
// 	}
func (c *Client) Costs(costs CostsRequest) (CostsResponse, int, errs.Failure) {
	var resp CostsResponse

	params, _ := query.Values(costs)

	data, status, err := c.GetRequest(fmt.Sprintf("/costs?%s", params.Encode()))

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CostCreate create an cost
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-costs-create
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostCreate(
// 		v5.CostRecord{
//			DateFrom:  "2012-12-12",
//			DateTo:    "2012-12-12",
//			Summ:      12,
//			CostItem:  "calculation-of-costs",
//			Order: Order{
// 				Number: "1"
// 			},
// 			Sites:    []string{"store"},
//		},
//		"store"
// 	)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	If data.Success == true {
// 		fmt.Printf("%v", data.ID)
// 	}
func (c *Client) CostCreate(cost CostRecord, site ...string) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse

	costJSON, _ := json.Marshal(&cost)

	p := url.Values{
		"cost": {string(costJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest("/costs/create", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CostsDelete removes a cost
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-costs-delete
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.client.CostsDelete([]int{1, 2, 3, 48, 49, 50})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	If data.Success == true {
// 		fmt.Printf("Not removed costs: %v", data.NotRemovedIds)
// 	}
func (c *Client) CostsDelete(ids []int) (CostsDeleteResponse, int, errs.Failure) {
	var resp CostsDeleteResponse

	costJSON, _ := json.Marshal(&ids)

	p := url.Values{
		"ids": {string(costJSON[:])},
	}

	data, status, err := c.PostRequest("/costs/delete", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CostsUpload batch costs upload
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-costs-upload
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostCreate([]v5.CostRecord{
//		{
//			DateFrom:  "2012-12-12",
//			DateTo:    "2012-12-12",
//			Summ:      12,
//			CostItem:  "calculation-of-costs",
//			Order: Order{
// 				Number: "1"
// 			},
// 			Sites:    []string{"store"},
//		},
//		{
//			DateFrom:  "2012-12-13",
//			DateTo:    "2012-12-13",
//			Summ:      13,
//			CostItem:  "seo",
//		}
// 	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	If data.Success == true {
// 		fmt.Printf("Uploaded costs: %v", data.UploadedCosts)
// 	}
func (c *Client) CostsUpload(cost []CostRecord) (CostsUploadResponse, int, errs.Failure) {
	var resp CostsUploadResponse

	costJSON, _ := json.Marshal(&cost)

	p := url.Values{
		"costs": {string(costJSON[:])},
	}

	data, status, err := c.PostRequest("/costs/upload", p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// Cost returns information about specified cost
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-costs-id
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.Cost(1)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	If data.Success == true {
// 		fmt.Printf("%v", data.Cost)
// 	}
func (c *Client) Cost(id int) (CostResponse, int, errs.Failure) {
	var resp CostResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/costs/%d", id))

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CostDelete removes cost
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-costs-id-delete
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostDelete(1)
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
func (c *Client) CostDelete(id int) (SuccessfulResponse, int, errs.Failure) {
	var resp SuccessfulResponse

	costJSON, _ := json.Marshal(&id)

	p := url.Values{
		"costs": {string(costJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/costs/%d/delete", id), p)

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CostEdit edit a cost
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-costs-id-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CostEdit(1, v5.Cost{
//		DateFrom:  "2012-12-12",
//		DateTo:    "2018-12-13",
//		Summ:      321,
//		CostItem:  "seo",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	If data.Success == true {
// 		fmt.Printf("%v", data.Id)
// 	}
func (c *Client) CostEdit(id int, cost CostRecord, site ...string) (CreateResponse, int, errs.Failure) {
	var resp CreateResponse

	costJSON, _ := json.Marshal(&cost)

	p := url.Values{
		"cost": {string(costJSON[:])},
	}

	fillSite(&p, site)

	data, status, err := c.PostRequest(fmt.Sprintf("/costs/%d/edit", id), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomFields returns list of custom fields
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-custom-fields
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomFields(v5.CustomFieldsRequest{
//		Type: "string",
// 		Entity: "customer",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	for _, value := range data.CustomFields {
//		fmt.Printf("%v\n", value)
//	}
func (c *Client) CustomFields(customFields CustomFieldsRequest) (CustomFieldsResponse, int, errs.Failure) {
	var resp CustomFieldsResponse

	params, _ := query.Values(customFields)

	data, status, err := c.GetRequest(fmt.Sprintf("/custom-fields?%s", params.Encode()))

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomDictionaries returns list of custom directory
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-custom-fields-dictionaries
//
// Example:
//
//	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomDictionaries(v5.CustomDictionariesRequest{
//		Filter: v5.CustomDictionariesFilter{
//			Name: "Dictionary-1",
//		},
//	})
//
//	if err.RuntimeErr != nil {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	for _, value := range data.CustomDictionaries {
//		fmt.Printf("%v\n", value.Elements)
//	}
func (c *Client) CustomDictionaries(customDictionaries CustomDictionariesRequest) (CustomDictionariesResponse, int, errs.Failure) {
	var resp CustomDictionariesResponse

	params, _ := query.Values(customDictionaries)

	data, status, err := c.GetRequest(fmt.Sprintf("/custom-fields/dictionaries?%s", params.Encode()))

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomDictionariesCreate creates a custom dictionary
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-custom-fields-dictionaries-create
//
// Example:
//
//	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomDictionariesCreate(v5.CustomDictionary{
//		Name: "Courier profiles",
//		Code: "courier-profiles",
//		Elements: []Element{
//			{
//				Name: "Name",
//				Code: "name",
//			},
//			{
//				Name: "Lastname",
//				Code: "lastname",
//			}
//		},
//	})
//
//	if err.RuntimeErr != nil {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	If data.Success == true {
//		fmt.Printf("%v", data.Code)
//	}
func (c *Client) CustomDictionariesCreate(customDictionary CustomDictionary) (CustomResponse, int, errs.Failure) {
	var resp CustomResponse

	costJSON, _ := json.Marshal(&customDictionary)

	p := url.Values{
		"customDictionary": {string(costJSON[:])},
	}

	data, status, err := c.PostRequest("/custom-fields/dictionaries/create", p)

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomDictionary returns information about dictionary
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-custom-fields-entity-code
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomDictionary("courier-profiles")
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	If data.Success == true {
// 		fmt.Printf("%v", data.CustomDictionary.Name)
// 	}
func (c *Client) CustomDictionary(code string) (CustomDictionaryResponse, int, errs.Failure) {
	var resp CustomDictionaryResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/custom-fields/dictionaries/%s", code))

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomDictionaryEdit edit custom dictionary
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-custom-fields-dictionaries-code-edit
//
// Example:
//
//	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomDictionaryEdit(v5.CustomDictionary{
//	Name: "Courier profiles",
//		Code: "courier-profiles",
//		Elements: []Element{
//			{
//				Name: "Name",
//				Code: "name",
//			},
//			{
//				Name: "Lastname",
//				Code: "lastname",
//			}
//		},
//	})
//
//	if err.RuntimeErr != nil {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	if status >= http.StatusBadRequest {
//		fmt.Printf("%v", err.RuntimeErr)
//	}
//
//	If data.Success == true {
//		fmt.Printf("%v", data.Code)
//	}
func (c *Client) CustomDictionaryEdit(customDictionary CustomDictionary) (CustomResponse, int, errs.Failure) {
	var resp CustomResponse

	costJSON, _ := json.Marshal(&customDictionary)

	p := url.Values{
		"customDictionary": {string(costJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/custom-fields/dictionaries/%s/edit", customDictionary.Code), p)
	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomFieldsCreate creates custom field
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-custom-fields-entity-create
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomFieldsCreate(CustomFields{
//		Name:        "First order",
//		Code:        "first-order",
//		Type:        "bool",
//		Entity:      "order",
//		DisplayArea: "customer",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	If data.Success == true {
// 		fmt.Printf("%v", data.Code)
// 	}
func (c *Client) CustomFieldsCreate(customFields CustomFields) (CustomResponse, int, errs.Failure) {
	var resp CustomResponse

	costJSON, _ := json.Marshal(&customFields)

	p := url.Values{
		"customField": {string(costJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/custom-fields/%s/create", customFields.Entity), p)

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomField returns information about custom fields
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#get--api-v5-custom-fields-entity-code
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomField("order", "first-order")
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	If data.Success == true {
// 		fmt.Printf("%v", data.CustomField)
// 	}
func (c *Client) CustomField(entity, code string) (CustomFieldResponse, int, errs.Failure) {
	var resp CustomFieldResponse

	data, status, err := c.GetRequest(fmt.Sprintf("/custom-fields/%s/%s", entity, code))

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

// CustomFieldEdit list method
//
// For more information see http://www.retailcrm.pro/docs/Developers/ApiVersion5#post--api-v5-custom-fields-entity-code-edit
//
// Example:
//
// 	var client = v5.New("https://demo.url", "09jIJ")
//
//	data, status, err := client.CustomFieldEdit(CustomFields{
//		Code:        "first-order",
//		Entity:      "order",
//		DisplayArea: "delivery",
//	})
//
// 	if err.RuntimeErr != nil {
// 		fmt.Printf("%v", err.RuntimeErr)
// 	}
//
// 	if status >= http.StatusBadRequest {
// 		fmt.Printf("%v", err.ApiErr())
// 	}
//
// 	If data.Success == true {
// 		fmt.Printf("%v", data.Code)
// 	}
func (c *Client) CustomFieldEdit(customFields CustomFields) (CustomResponse, int, errs.Failure) {
	var resp CustomResponse

	costJSON, _ := json.Marshal(&customFields)

	p := url.Values{
		"customField": {string(costJSON[:])},
	}

	data, status, err := c.PostRequest(fmt.Sprintf("/custom-fields/%s/%s/edit", customFields.Entity, customFields.Code), p)

	if err.RuntimeErr != nil {
		return resp, status, err
	}

	json.Unmarshal(data, &resp)

	return resp, status, err
}

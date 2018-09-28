package v5

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/h2non/gock"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if os.Getenv("DEVELOPER_NODE") == "1" {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		os.Exit(m.Run())
	}
}

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

var (
	r        *rand.Rand // Rand for this package.
	user, _  = strconv.Atoi(os.Getenv("RETAILCRM_USER"))
	statuses = map[int]bool{http.StatusOK: true, http.StatusCreated: true}
	crmURL   = os.Getenv("RETAILCRM_URL")
	badURL   = "https://qwertypoiu.retailcrm.ru"

	statusFail  = "FailTest: status < 400"
	successFail = "FailTest: Success == true"
	codeFail    = "test-12345"
	iCodeFail   = 123123
)

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)

	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}

	return string(result)
}

func client() *Client {
	c := New(crmURL, os.Getenv("RETAILCRM_KEY"))
	c.Debug = true
	return c
}

func badurlclient() *Client {
	return New(badURL, os.Getenv("RETAILCRM_KEY"))
}

func badkeyclient() *Client {
	return New(os.Getenv("RETAILCRM_URL"), "1234567890")
}

func TestGetRequest(t *testing.T) {
	c := client()
	_, status, _ := c.GetRequest("/fake-method")

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestPostRequest(t *testing.T) {
	c := client()
	_, status, _ := c.PostRequest("/fake-method", url.Values{})

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestClient_ApiVersionsVersions(t *testing.T) {
	c := client()

	data, status, err := c.APIVersions()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_ApiVersionsVersionsBadUrl(t *testing.T) {
	c := badurlclient()

	defer gock.Off()

	gock.New(badURL).
		Get("/api/api-versions").
		Reply(200).
		BodyString(`{"success": false, "errorMsg" : "Account does not exist"}`)

	data, status, err := c.APIVersions()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != false {
		t.Logf("%v", err.ApiError())
	}
}

func TestClient_ApiCredentialsCredentialsBadKey(t *testing.T) {
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/credentials").
		Reply(403).
		BodyString(`{"success": false, "errorMsg": "Wrong \"apiKey\" value"}`)

	c := badkeyclient()

	data, status, err := c.APICredentials()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Logf("%v", err.ApiError())
	}

	if data.Success != false {
		t.Logf("%v", err.ApiError())
	}
}

func TestClient_ApiCredentialsCredentials(t *testing.T) {
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/credentials").
		Reply(200).
		BodyString(`{"success": true}`)

	c := client()

	data, status, err := c.APICredentials()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomersCustomers(t *testing.T) {
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/v5/customers").
		MatchParam("filter[city]", "Москва").
		MatchParam("page", "3").
		Reply(200).
		BodyString(`{"success": true}`)

	c := client()

	data, status, err := c.Customers(CustomersRequest{
		Filter: CustomersFilter{
			City: "Москва",
		},
		Page: 3,
	})

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Logf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Logf("%v", err.ApiError())
	}
}

func TestClient_CustomersCustomers_Fail(t *testing.T) {
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/v5/customers").
		MatchParam("filter[ids][]", codeFail).
		Reply(400).
		BodyString(`{"success": false,"errorMsg": "Internal Server Error"}`)

	c := client()

	data, status, err := c.Customers(CustomersRequest{
		Filter: CustomersFilter{
			Ids: []string{codeFail},
		},
	})

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CustomerChange(t *testing.T) {
	c := client()

	f := Customer{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalID: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		Address: &Address{
			City:     "Москва",
			Street:   "Кутузовский",
			Building: "14",
		},
	}

	defer gock.Off()

	str, _ := json.Marshal(f)

	p := url.Values{
		"customer": {string(str)},
	}

	gock.New(crmURL).
		Post("/api/v5/customers/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 123}`)

	cr, sc, err := c.CustomerCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if cr.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	f.ID = cr.ID
	f.Vip = true

	str, _ = json.Marshal(f)

	p = url.Values{
		"by":       {string("id")},
		"customer": {string(str)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/customers/%v/edit", cr.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	ed, se, err := c.CustomerEdit(f, "id")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if se != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if ed.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/api/v5/customers/%v", f.ExternalID)).
		MatchParam("by", "externalId").
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.Customer(f.ExternalID, "externalId", "")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomerChange_Fail(t *testing.T) {
	c := client()

	f := Customer{
		FirstName: "Понтелей",
	}

	defer gock.Off()

	str, _ := json.Marshal(f)

	p := url.Values{
		"customer": {string(str)},
	}

	gock.New(crmURL).
		Post("/api/v5/customers/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Parameter 'externalId' in 'customer' is missing"}`)

	cr, sc, err := c.CustomerCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sc < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if cr.Success != false {
		t.Error(successFail)
	}

	f.ID = cr.ID
	f.Vip = true

	str, _ = json.Marshal(f)

	p = url.Values{
		"by":       {string("id")},
		"customer": {string(str)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/customers/%v/edit", cr.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	ed, se, err := c.CustomerEdit(f, "id")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if se < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if ed.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/api/v5/customers/%v", codeFail)).
		MatchParam("by", "externalId").
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	data, status, err := c.Customer(codeFail, "externalId", "")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CustomersUpload(t *testing.T) {
	c := client()
	customers := make([]Customer, 3)

	for i := range customers {
		customers[i] = Customer{
			FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
			LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
			ExternalID: RandomString(8),
			Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		}
	}

	defer gock.Off()

	str, _ := json.Marshal(customers)

	p := url.Values{
		"customers": {string(str)},
	}

	gock.New(crmURL).
		Post("/api/v5/customers/upload").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.CustomersUpload(customers)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomersUpload_Fail(t *testing.T) {
	c := client()

	customers := []Customer{{ExternalID: strconv.Itoa(iCodeFail)}}

	defer gock.Off()

	str, _ := json.Marshal(customers)
	p := url.Values{
		"customers": {string(str)},
	}

	gock.New(crmURL).
		Post("/api/v5/customers/upload").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(460).
		BodyString(`{"success": false, "errorMsg": "Customers are loaded with errors"}`)

	data, status, err := c.CustomersUpload(customers)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CustomersCombine(t *testing.T) {
	c := client()

	defer gock.Off()

	customers := []Customer{{ID: 1}, {ID: 2}}
	resultCustomer := Customer{ID: 3}

	jr, _ := json.Marshal(&customers)
	combineJSONOut, _ := json.Marshal(&resultCustomer)

	p := url.Values{
		"customers":      {string(jr[:])},
		"resultCustomer": {string(combineJSONOut[:])},
	}

	gock.New(crmURL).
		Post("/customers/combine").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.CustomersCombine(customers, resultCustomer)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomersCombine_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	customers := []Customer{{}, {}}
	resultCustomer := Customer{}

	jr, _ := json.Marshal(&customers)
	combineJSONOut, _ := json.Marshal(&resultCustomer)

	p := url.Values{
		"customers":      {string(jr[:])},
		"resultCustomer": {string(combineJSONOut[:])},
	}

	gock.New(crmURL).
		Post("/customers/combine").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Invalid input parameters"}`)

	data, status, err := c.CustomersCombine([]Customer{{}, {}}, Customer{})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CustomersFixExternalIds(t *testing.T) {
	c := client()

	customers := []IdentifiersPair{{
		ID:         123,
		ExternalID: RandomString(8),
	}}

	defer gock.Off()

	jr, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/customers/fix-external-ids").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	fx, fe, err := c.CustomersFixExternalIds(customers)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if fe != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if fx.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomersFixExternalIds_Fail(t *testing.T) {
	c := client()

	customers := []IdentifiersPair{{ExternalID: RandomString(8)}}

	defer gock.Off()

	jr, _ := json.Marshal(&customers)

	p := url.Values{
		"customers": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/customers/fix-external-ids").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "errors": {"id": "ID must be an integer"}}`)

	data, status, err := c.CustomersFixExternalIds(customers)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CustomersHistory(t *testing.T) {
	c := client()
	f := CustomersHistoryRequest{
		Filter: CustomersHistoryFilter{
			SinceID: 20,
		},
	}
	defer gock.Off()

	gock.New(crmURL).
		Get("/customers/history").
		MatchParam("filter[sinceId]", "20").
		Reply(200).
		BodyString(`{"success": true, "history": [{"id": 1}]}`)

	data, status, err := c.CustomersHistory(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	if len(data.History) == 0 {
		t.Errorf("%v", "Empty history")
	}
}

func TestClient_CustomersHistory_Fail(t *testing.T) {
	c := client()
	f := CustomersHistoryRequest{
		Filter: CustomersHistoryFilter{
			StartDate: "2020-13-12",
		},
	}

	defer gock.Off()

	gock.New(crmURL).
		Get("/customers/history").
		MatchParam("filter[startDate]", "2020-13-12").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "errors": {"children[startDate]": "Значение недопустимо."}}`)

	data, status, err := c.CustomersHistory(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_NotesNotes(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/customers/notes").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true, "notes": [{"id": 1}]}`)

	data, status, err := c.CustomerNotes(NotesRequest{Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	if len(data.Notes) == 0 {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_NotesNotes_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/customers/notes").
		MatchParam("filter[createdAtFrom]", "2020-13-12").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "errors": {"children[createdAtFrom]": "This value is not valid."}}`)

	data, status, err := c.CustomerNotes(NotesRequest{
		Filter: NotesFilter{CreatedAtFrom: "2020-13-12"},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_NotesCreateDelete(t *testing.T) {
	c := client()

	note := Note{
		Text:      "some text",
		ManagerID: user,
		Customer: &Customer{
			ID: 123,
		},
	}

	defer gock.Off()

	jr, _ := json.Marshal(&note)

	p := url.Values{
		"note": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/customers/notes/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 1}`)

	noteCreateResponse, noteCreateStatus, err := c.CustomerNoteCreate(note)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if noteCreateStatus != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if noteCreateResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	p = url.Values{
		"id": {string(1)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/api/v5/customers/notes/%d/delete", 1)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	noteDeleteResponse, noteDeleteStatus, err := c.CustomerNoteDelete(noteCreateResponse.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if noteDeleteStatus != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if noteDeleteResponse.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_NotesCreateDelete_Fail(t *testing.T) {
	c := client()
	defer gock.Off()

	note := Note{
		Text:      "some text",
		ManagerID: user,
		Customer:  &Customer{},
	}

	jr, _ := json.Marshal(&note)

	p := url.Values{
		"note": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/customers/notes/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the entity format", "errors": {"customer": "Set one of the following fields: id, externalId"}}`)

	data, status, err := c.CustomerNoteCreate(note)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}

	p = url.Values{
		"id": {string(iCodeFail)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/customers/notes/%d/delete", iCodeFail)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Note not found map"}`)

	noteDeleteResponse, noteDeleteStatus, err := c.CustomerNoteDelete(iCodeFail)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if noteDeleteStatus < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if noteDeleteResponse.Success != false {
		t.Error(successFail)
	}
}

func TestClient_OrdersOrders(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders").
		MatchParam("filter[city]", "Москва").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true, "orders":  [{"id": 1}]}`)

	data, status, err := c.Orders(OrdersRequest{Filter: OrdersFilter{City: "Москва"}, Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersOrders_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders").
		MatchParam("filter[attachments]", "7").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "errors": {"children[attachments]": "SThis value is not valid."}}`)

	data, status, err := c.Orders(OrdersRequest{Filter: OrdersFilter{Attachments: 7}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_OrderChange(t *testing.T) {
	c := client()

	random := RandomString(8)

	f := Order{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalID: random,
		Email:      fmt.Sprintf("%s@example.com", random),
	}

	defer gock.Off()

	jr, _ := json.Marshal(&f)

	p := url.Values{
		"order": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/orders/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 1}`)

	cr, sc, err := c.OrderCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if cr.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	f.ID = cr.ID
	f.CustomerComment = "test comment"

	jr, _ = json.Marshal(&f)

	p = url.Values{
		"by":    {string("id")},
		"order": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/%d/edit", f.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	ed, se, err := c.OrderEdit(f, "id")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if se != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if ed.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/orders/%s", f.ExternalID)).
		MatchParam("by", "externalId").
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.Order(f.ExternalID, "externalId", "")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderChange_Fail(t *testing.T) {
	c := client()

	random := RandomString(8)

	f := Order{
		FirstName:       "Понтелей",
		LastName:        "Турбин",
		Patronymic:      "Аристархович",
		ExternalID:      random,
		Email:           fmt.Sprintf("%s@example.com", random),
		CustomerComment: "test comment",
	}

	defer gock.Off()

	jr, _ := json.Marshal(&f)

	p := url.Values{
		"by":    {string("id")},
		"order": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/%d/edit", f.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found map"}`)

	data, status, err := c.OrderEdit(f, "id")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_OrdersUpload(t *testing.T) {
	c := client()
	orders := make([]Order, 3)

	for i := range orders {
		orders[i] = Order{
			FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
			LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
			ExternalID: RandomString(8),
			Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		}
	}

	defer gock.Off()

	jr, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/orders/upload").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.OrdersUpload(orders)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersUpload_Fail(t *testing.T) {
	c := client()

	orders := []Order{
		{
			FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
			LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
			ExternalID: strconv.Itoa(iCodeFail),
			Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		},
	}

	defer gock.Off()

	jr, _ := json.Marshal(&orders)

	p := url.Values{
		"orders": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/orders/upload").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(460).
		BodyString(`{"success": false, "errorMsg": "Orders are loaded with errors"}`)

	data, status, err := c.OrdersUpload(orders)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_OrdersCombine(t *testing.T) {
	c := client()

	defer gock.Off()

	jr1, _ := json.Marshal(&Order{ID: 1})
	jr2, _ := json.Marshal(&Order{ID: 2})
	p := url.Values{
		"technique":   {"ours"},
		"order":       {string(jr1)},
		"resultOrder": {string(jr2)},
	}

	gock.New(crmURL).
		Post("/orders/combine").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.OrdersCombine("ours", Order{ID: 1}, Order{ID: 2})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersCombine_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	jr, _ := json.Marshal(&Order{})
	p := url.Values{
		"technique":   {"ours"},
		"order":       {string(jr)},
		"resultOrder": {string(jr)},
	}

	gock.New(crmURL).
		Post("/orders/combine").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Invalid input parameters"}`)

	data, status, err := c.OrdersCombine("ours", Order{}, Order{})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_OrdersFixExternalIds(t *testing.T) {
	c := client()

	orders := []IdentifiersPair{{
		ID:         123,
		ExternalID: RandomString(8),
	}}

	defer gock.Off()

	jr, _ := json.Marshal(orders)
	p := url.Values{
		"orders": {string(jr)},
	}

	gock.New(crmURL).
		Post("/orders/fix-external-ids").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	fx, fe, err := c.OrdersFixExternalIds(orders)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if fe != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if fx.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersFixExternalIds_Fail(t *testing.T) {
	c := client()
	orders := []IdentifiersPair{{ExternalID: RandomString(8)}}

	defer gock.Off()

	jr, _ := json.Marshal(orders)
	p := url.Values{
		"orders": {string(jr)},
	}

	gock.New(crmURL).
		Post("/orders/fix-external-ids").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Invalid input parameters"}`)

	data, status, err := c.OrdersFixExternalIds(orders)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_OrdersHistory(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders/history").
		MatchParam("filter[sinceId]", "20").
		Reply(200).
		BodyString(`{"success": true, "history": [{"id": 1}]}`)

	data, status, err := c.OrdersHistory(OrdersHistoryRequest{Filter: OrdersHistoryFilter{SinceID: 20}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	if len(data.History) == 0 {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrdersHistory_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders/history").
		MatchParam("filter[startDate]", "2020-13-12").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "errors": {"children[startDate]": "Значение недопустимо."}}`)

	data, status, err := c.OrdersHistory(OrdersHistoryRequest{Filter: OrdersHistoryFilter{StartDate: "2020-13-12"}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_PaymentCreateEditDelete(t *testing.T) {
	c := client()

	f := Payment{
		Order: &Order{
			ID: 123,
		},
		Amount: 300,
		Type:   "cash",
	}

	defer gock.Off()

	jr, _ := json.Marshal(f)
	p := url.Values{
		"payment": {string(jr)},
	}

	gock.New(crmURL).
		Post("/orders/payments/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 1}`)

	paymentCreateResponse, status, err := c.OrderPaymentCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if paymentCreateResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	k := Payment{
		ID:     paymentCreateResponse.ID,
		Amount: 500,
	}

	jr, _ = json.Marshal(k)
	p = url.Values{
		"payment": {string(jr)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/payments/%d/edit", paymentCreateResponse.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	paymentEditResponse, status, err := c.OrderPaymentEdit(k, "id")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if paymentEditResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	p = url.Values{
		"id": {string(paymentCreateResponse.ID)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/payments/%d/delete", paymentCreateResponse.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	paymentDeleteResponse, status, err := c.OrderPaymentDelete(paymentCreateResponse.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if paymentDeleteResponse.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentCreateEditDelete_Fail(t *testing.T) {
	c := client()

	f := Payment{
		Order:  &Order{},
		Amount: 300,
		Type:   "cash",
	}

	defer gock.Off()

	jr, _ := json.Marshal(f)
	p := url.Values{
		"payment": {string(jr)},
	}

	gock.New(crmURL).
		Post("/orders/payments/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the entity format", "errors": {"order": "Set one of the following fields: id, externalId, number"}}`)

	data, status, err := c.OrderPaymentCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}

	k := Payment{
		ID:     iCodeFail,
		Amount: 500,
	}

	jr, _ = json.Marshal(k)
	p = url.Values{
		"payment": {string(jr)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/payments/%d/edit", iCodeFail)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Payment not found"}`)

	paymentEditResponse, status, err := c.OrderPaymentEdit(k, "id")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if paymentEditResponse.Success != false {
		t.Error(successFail)
	}

	p = url.Values{
		"id": {string(iCodeFail)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/payments/%d/delete", iCodeFail)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Payment not found"}`)

	paymentDeleteResponse, status, err := c.OrderPaymentDelete(iCodeFail)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if paymentDeleteResponse.Success != false {
		t.Error(successFail)
	}
}

func TestClient_TasksTasks(t *testing.T) {
	c := client()

	f := TasksRequest{
		Filter: TasksFilter{
			Creators: []int{123},
		},
		Page: 1,
	}
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/v5/tasks").
		MatchParam("filter[creators][]", "123").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.Tasks(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_TasksTasks_Fail(t *testing.T) {
	c := client()

	f := TasksRequest{
		Filter: TasksFilter{
			Creators: []int{123123},
		},
	}
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/v5/tasks").
		MatchParam("filter[creators][]", "123123").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	data, status, err := c.Tasks(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_TaskChange(t *testing.T) {
	c := client()

	f := Task{
		Text:        RandomString(15),
		PerformerID: user,
	}
	defer gock.Off()

	jr, _ := json.Marshal(f)
	p := url.Values{
		"task": {string(jr)},
	}

	gock.New(crmURL).
		Post("/tasks/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 1}`)

	cr, sc, err := c.TaskCreate(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if cr.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	f.ID = cr.ID
	f.Commentary = RandomString(20)

	gock.New(crmURL).
		Get(fmt.Sprintf("/tasks/%d", f.ID)).
		Reply(200).
		BodyString(`{"success": true}`)

	gt, sg, err := c.Task(f.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if sg != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if gt.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	jr, _ = json.Marshal(f)
	p = url.Values{
		"task": {string(jr)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/tasks/%d/edit", f.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.TaskEdit(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_TaskChange_Fail(t *testing.T) {
	c := client()

	f := Task{
		Text:       RandomString(15),
		Commentary: RandomString(20),
	}

	defer gock.Off()

	jr, _ := json.Marshal(f)
	p := url.Values{
		"task": {string(jr)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/tasks/%d/edit", f.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Task is not loaded", "errors": {"performerId": "This value should not be blank."}}`)

	data, status, err := c.TaskEdit(f)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_UsersUsers(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/users").
		MatchParam("filter[active]", "1").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.Users(UsersRequest{Filter: UsersFilter{Active: 1}, Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_UsersUsers_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/users").
		MatchParam("filter[active]", "3").
		MatchParam("page", "1").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "errors": {"active": "he value you selected is not a valid choice."}}`)

	data, status, err := c.Users(UsersRequest{Filter: UsersFilter{Active: 3}, Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_UsersUser(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get(fmt.Sprintf("/users/%d", user)).
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.User(user)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_UsersUser_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get(fmt.Sprintf("/users/%d", iCodeFail)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	data, status, err := c.User(iCodeFail)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_UsersGroups(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/user-groups").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.UserGroups(UserGroupsRequest{Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_UsersUpdate(t *testing.T) {
	c := client()

	defer gock.Off()

	p := url.Values{
		"status": {string("busy")},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/users/%d/status", user)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.UserStatus(user, "busy")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_UsersUpdate_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	p := url.Values{
		"status": {string("busy")},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/users/%d/status", iCodeFail)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	data, status, err := c.UserStatus(iCodeFail, "busy")
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_StaticticsUpdate(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/statistic/update").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.StaticticsUpdate()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Countries(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/couriers").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.Couriers()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CostGroups(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/cost-groups").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.CostGroups()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CostItems(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/cost-items").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.CostItems()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Couriers(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/couriers").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.Couriers()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryService(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/delivery-services").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.DeliveryServices()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryTypes(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/delivery-types").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.DeliveryTypes()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_LegalEntities(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/legal-entities").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.LegalEntities()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderMethods(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/order-methods").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.OrderMethods()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderTypes(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/order-types").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.OrderTypes()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentStatuses(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/payment-statuses").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.PaymentStatuses()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentTypes(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/payment-types").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.PaymentTypes()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PriceTypes(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/price-types").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.PriceTypes()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_ProductStatuses(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/product-statuses").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.ProductStatuses()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Statuses(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/statuses").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.Statuses()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_StatusGroups(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/status-groups").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.StatusGroups()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Sites(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/sites").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.Sites()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Stores(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/reference/stores").
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.Stores()
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CostGroupItemEdit(t *testing.T) {
	c := client()

	uid := RandomString(4)

	costGroup := CostGroup{
		Code:   fmt.Sprintf("cost-gr-%s", uid),
		Active: false,
		Color:  "#da5c98",
		Name:   fmt.Sprintf("CostGroup-%s", uid),
	}

	defer gock.Off()

	jr, _ := json.Marshal(&costGroup)

	p := url.Values{
		"costGroup": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/api/v5/reference/cost-groups/%s/edit", costGroup.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	data, st, err := c.CostGroupEdit(costGroup)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	costItem := CostItem{
		Code:            fmt.Sprintf("cost-it-%s", uid),
		Name:            fmt.Sprintf("CostItem-%s", uid),
		Group:           fmt.Sprintf("cost-gr-%s", uid),
		Type:            "const",
		AppliesToOrders: true,
		Active:          false,
	}

	jr, _ = json.Marshal(&costItem)

	p = url.Values{
		"costItem": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/api/v5/reference/cost-items/%s/edit", costItem.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	idata, st, err := c.CostItemEdit(costItem)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if idata.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CostGroupItemEdit_Fail(t *testing.T) {
	c := client()

	uid := RandomString(4)
	costGroup := CostGroup{
		Code:   fmt.Sprintf("cost-gr-%s", uid),
		Active: false,
		Name:   fmt.Sprintf("CostGroup-%s", uid),
	}

	defer gock.Off()

	jr, _ := json.Marshal(&costGroup)

	p := url.Values{
		"costGroup": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/cost-groups/%s/edit", costGroup.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	data, status, err := c.CostGroupEdit(costGroup)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}

	costItem := CostItem{
		Code:            fmt.Sprintf("666cost-it-%s", uid),
		Name:            fmt.Sprintf("CostItem-%s", uid),
		Group:           fmt.Sprintf("cost-gr-%s", uid),
		Type:            "const",
		AppliesToOrders: true,
		Active:          false,
	}

	jr, _ = json.Marshal(&costItem)

	p = url.Values{
		"costItem": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/cost-items/%s/edit", costItem.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	idata, st, err := c.CostItemEdit(costItem)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if idata.Success != false {
		t.Error(successFail)
	}
}

func TestClient_Courier(t *testing.T) {
	c := client()

	cur := Courier{
		Active:    true,
		Email:     fmt.Sprintf("%s@example.com", RandomString(5)),
		FirstName: fmt.Sprintf("%s", RandomString(5)),
		LastName:  fmt.Sprintf("%s", RandomString(5)),
	}

	defer gock.Off()

	jr, _ := json.Marshal(&cur)

	p := url.Values{
		"courier": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/api/v5/reference/couriers/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 1}`)

	data, st, err := c.CourierCreate(cur)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	cur.ID = data.ID
	cur.Patronymic = fmt.Sprintf("%s", RandomString(5))

	jr, _ = json.Marshal(&cur)

	p = url.Values{
		"courier": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/couriers/%d/edit", cur.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	idata, st, err := c.CourierEdit(cur)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if idata.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Courier_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	jr, _ := json.Marshal(&Courier{})

	p := url.Values{
		"courier": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/api/v5/reference/couriers/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the entity format", "errors": {"firstName": "Specify the first name"}}`)

	data, st, err := c.CourierCreate(Courier{})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}

	cur := Courier{Patronymic: fmt.Sprintf("%s", RandomString(5))}
	jr, _ = json.Marshal(&cur)

	p = url.Values{
		"courier": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/couriers/%d/edit", cur.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "An attempt was made to edit a nonexistent record"}`)

	idata, st, err := c.CourierEdit(cur)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if idata.Success != false {
		t.Error(successFail)
	}
}

func TestClient_DeliveryServiceEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	deliveryService := DeliveryService{
		Active: false,
		Code:   RandomString(5),
		Name:   RandomString(5),
	}

	jr, _ := json.Marshal(&deliveryService)

	p := url.Values{
		"deliveryService": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/delivery-services/%s/edit", deliveryService.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true, "id": 1}`)

	data, st, err := c.DeliveryServiceEdit(deliveryService)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryServiceEdit_Fail(t *testing.T) {
	c := client()
	defer gock.Off()

	deliveryService := DeliveryService{
		Active: false,
		Name:   RandomString(5),
	}

	jr, _ := json.Marshal(&deliveryService)

	p := url.Values{
		"deliveryService": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/delivery-services/%s/edit", deliveryService.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.DeliveryServiceEdit(deliveryService)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_DeliveryTypeEdit(t *testing.T) {
	c := client()

	x := []string{"cash", "bank-card"}

	defer gock.Off()

	deliveryType := DeliveryType{
		Active:        false,
		Code:          RandomString(5),
		Name:          RandomString(5),
		DefaultCost:   300,
		PaymentTypes:  x,
		DefaultForCrm: false,
	}

	jr, _ := json.Marshal(&deliveryType)

	p := url.Values{
		"deliveryType": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/delivery-types/%s/edit", deliveryType.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.DeliveryTypeEdit(deliveryType)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryTypeEdit_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	deliveryType := DeliveryType{
		Active:        false,
		Name:          RandomString(5),
		DefaultCost:   300,
		PaymentTypes:  []string{"cash", "bank-card"},
		DefaultForCrm: false,
	}

	jr, _ := json.Marshal(&deliveryType)

	p := url.Values{
		"deliveryType": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/delivery-types/%s/edit", deliveryType.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.DeliveryTypeEdit(deliveryType)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_OrderMethodEdit(t *testing.T) {
	c := client()
	defer gock.Off()

	orderMethod := OrderMethod{
		Code:          RandomString(5),
		Name:          RandomString(5),
		Active:        false,
		DefaultForCRM: false,
	}

	jr, _ := json.Marshal(&orderMethod)

	p := url.Values{
		"orderMethod": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/order-methods/%s/edit", orderMethod.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.OrderMethodEdit(orderMethod)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderMethodEdit_Fail(t *testing.T) {
	c := client()
	defer gock.Off()

	orderMethod := OrderMethod{
		Name:          RandomString(5),
		Active:        false,
		DefaultForCRM: false,
	}

	jr, _ := json.Marshal(&orderMethod)

	p := url.Values{
		"orderMethod": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/order-methods/%s/edit", orderMethod.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.OrderMethodEdit(orderMethod)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_OrderTypeEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	orderType := OrderType{
		Code:          RandomString(5),
		Name:          RandomString(5),
		Active:        false,
		DefaultForCRM: false,
	}

	jr, _ := json.Marshal(&orderType)

	p := url.Values{
		"orderType": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/order-types/%s/edit", orderType.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.OrderTypeEdit(orderType)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_OrderTypeEdit_Fail(t *testing.T) {
	c := client()
	defer gock.Off()

	orderType := OrderType{
		Name:          RandomString(5),
		Active:        false,
		DefaultForCRM: false,
	}
	jr, _ := json.Marshal(&orderType)

	p := url.Values{
		"orderType": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/order-types/%s/edit", orderType.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.OrderTypeEdit(orderType)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_PaymentStatusEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	paymentStatus := PaymentStatus{
		Code:            RandomString(5),
		Name:            RandomString(5),
		Active:          false,
		DefaultForCRM:   false,
		PaymentTypes:    []string{"cash"},
		PaymentComplete: false,
	}

	jr, _ := json.Marshal(&paymentStatus)

	p := url.Values{
		"paymentStatus": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/payment-statuses/%s/edit", paymentStatus.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.PaymentStatusEdit(paymentStatus)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentStatusEdit_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	paymentStatus := PaymentStatus{
		Name:            RandomString(5),
		Active:          false,
		DefaultForCRM:   false,
		PaymentTypes:    []string{"cash"},
		PaymentComplete: false,
	}
	jr, _ := json.Marshal(&paymentStatus)

	p := url.Values{
		"paymentStatus": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/payment-statuses/%s/edit", paymentStatus.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.PaymentStatusEdit(paymentStatus)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_PaymentTypeEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	paymentType := PaymentType{
		Code:          RandomString(5),
		Name:          RandomString(5),
		Active:        false,
		DefaultForCRM: false,
	}
	jr, _ := json.Marshal(&paymentType)

	p := url.Values{
		"paymentType": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/payment-types/%s/edit", paymentType.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.PaymentTypeEdit(paymentType)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PaymentTypeEdit_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	paymentType := PaymentType{
		Name:          RandomString(5),
		Active:        false,
		DefaultForCRM: false,
	}
	jr, _ := json.Marshal(&paymentType)

	p := url.Values{
		"paymentType": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/payment-types/%s/edit", paymentType.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.PaymentTypeEdit(paymentType)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_PriceTypeEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	priceType := PriceType{
		Code:   RandomString(5),
		Name:   RandomString(5),
		Active: false,
	}
	jr, _ := json.Marshal(&priceType)

	p := url.Values{
		"priceType": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/price-types/%s/edit", priceType.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.PriceTypeEdit(priceType)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PriceTypeEdit_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	priceType := PriceType{
		Name:   RandomString(5),
		Active: false,
	}
	jr, _ := json.Marshal(&priceType)

	p := url.Values{
		"priceType": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/price-types/%s/edit", priceType.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.PriceTypeEdit(priceType)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_ProductStatusEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	productStatus := ProductStatus{
		Code:         RandomString(5),
		Name:         RandomString(5),
		Active:       false,
		CancelStatus: false,
	}
	jr, _ := json.Marshal(&productStatus)

	p := url.Values{
		"productStatus": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/product-statuses/%s/edit", productStatus.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.ProductStatusEdit(productStatus)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_ProductStatusEdit_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	productStatus := ProductStatus{
		Name:         RandomString(5),
		Active:       false,
		CancelStatus: false,
	}
	jr, _ := json.Marshal(&productStatus)

	p := url.Values{
		"productStatus": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/product-statuses/%s/edit", productStatus.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.ProductStatusEdit(productStatus)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_StatusEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	status := Status{
		Code:   RandomString(5),
		Name:   RandomString(5),
		Active: false,
		Group:  "new",
	}
	jr, _ := json.Marshal(&status)

	p := url.Values{
		"status": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/statuses/%s/edit", status.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.StatusEdit(status)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_StatusEdit_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	status := Status{
		Name:   RandomString(5),
		Active: false,
		Group:  "new",
	}

	jr, _ := json.Marshal(&status)

	p := url.Values{
		"status": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/statuses/%s/edit", status.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.StatusEdit(status)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_SiteEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	site := Site{
		Code:        RandomString(5),
		Name:        RandomString(5),
		URL:         fmt.Sprintf("https://%s.example.com", RandomString(4)),
		LoadFromYml: false,
	}
	jr, _ := json.Marshal(&site)

	p := url.Values{
		"site": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/sites/%s/edit", site.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, _, err := c.SiteEdit(site)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if data.Success == false {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_SiteEdit_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	site := Site{
		Code:        RandomString(5),
		Name:        RandomString(5),
		LoadFromYml: false,
	}

	jr, _ := json.Marshal(&site)

	p := url.Values{
		"site": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/sites/%s/edit", site.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(405).
		BodyString(`{"success": false, "errorMsg": "Method Not Allowed"}`)

	data, _, err := c.SiteEdit(site)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_StoreEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	store := Store{
		Code:          RandomString(5),
		Name:          RandomString(5),
		Active:        false,
		Type:          "store-type-warehouse",
		InventoryType: "integer",
	}

	jr, _ := json.Marshal(&store)

	p := url.Values{
		"store": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/stores/%s/edit", store.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.StoreEdit(store)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[st] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_StoreEdit_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	store := Store{
		Name:          RandomString(5),
		Active:        false,
		Type:          "store-type-warehouse",
		InventoryType: "integer",
	}

	jr, _ := json.Marshal(&store)

	p := url.Values{
		"store": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/reference/stores/%s/edit", store.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	data, st, err := c.StoreEdit(store)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_PackChange(t *testing.T) {
	c := client()
	defer gock.Off()

	pack := Pack{
		Store:    "test-store",
		ItemID:   123,
		Quantity: 1,
	}

	jr, _ := json.Marshal(&pack)

	pr := url.Values{
		"pack": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/orders/packs/create").
		MatchType("url").
		BodyString(pr.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 1}`)

	p, status, err := c.PackCreate(pack)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if p.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/orders/packs/%d", p.ID)).
		Reply(200).
		BodyString(`{"success": true}`)

	s, status, err := c.Pack(p.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if s.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	jr, _ = json.Marshal(&Pack{ID: p.ID, Quantity: 2})

	pr = url.Values{
		"pack": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/packs/%d/edit", p.ID)).
		MatchType("url").
		BodyString(pr.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	e, status, err := c.PackEdit(Pack{ID: p.ID, Quantity: 2})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if e.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/packs/%d/delete", p.ID)).
		MatchType("url").
		Reply(200).
		BodyString(`{"success": true}`)

	d, status, err := c.PackDelete(p.ID)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if d.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PackChange_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	pack := Pack{
		Store:    "test-store",
		ItemID:   iCodeFail,
		Quantity: 1,
	}

	jr, _ := json.Marshal(&pack)

	pr := url.Values{
		"pack": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/orders/packs/create").
		MatchType("url").
		BodyString(pr.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the entity format"}`)

	p, status, err := c.PackCreate(pack)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if p.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/orders/packs/%d", iCodeFail)).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the entity format"}`)

	s, status, err := c.Pack(iCodeFail)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if s.Success != false {
		t.Error(successFail)
	}

	jr, _ = json.Marshal(&Pack{ID: iCodeFail, Quantity: 2})

	pr = url.Values{
		"pack": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/packs/%d/edit", iCodeFail)).
		MatchType("url").
		BodyString(pr.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Pack with id 123123 not found"}`)

	e, status, err := c.PackEdit(Pack{ID: iCodeFail, Quantity: 2})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if e.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/packs/%d/delete", iCodeFail)).
		MatchType("url").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Pack not found"}`)

	d, status, err := c.PackDelete(iCodeFail)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if d.Success != false {
		t.Error(successFail)
	}
}

func TestClient_PacksHistory(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders/packs/history").
		MatchParam("filter[sinceId]", "5").
		Reply(200).
		BodyString(`{"success": true, "history": [{"id": 1}]}`)

	data, status, err := c.PacksHistory(PacksHistoryRequest{Filter: OrdersHistoryFilter{SinceID: 5}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	if len(data.History) == 0 {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_PacksHistory_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders/packs/history").
		MatchParam("filter[startDate]", "2020-13-12").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	data, status, err := c.PacksHistory(PacksHistoryRequest{Filter: OrdersHistoryFilter{StartDate: "2020-13-12"}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_Packs(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders/packs").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.Packs(PacksRequest{Filter: PacksFilter{}, Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Packs_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders/packs").
		MatchParam("filter[shipmentDateFrom]", "2020-13-12").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	data, status, err := c.Packs(PacksRequest{Filter: PacksFilter{ShipmentDateFrom: "2020-13-12"}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_Inventories(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/store/inventories").
		MatchParam("filter[details]", "1").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true, "id": 1}`)

	data, status, err := c.Inventories(InventoriesRequest{Filter: InventoriesFilter{Details: 1}, Page: 1})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Inventories_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/store/inventories").
		MatchParam("filter[sites][]", codeFail).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	data, status, err := c.Inventories(InventoriesRequest{Filter: InventoriesFilter{Sites: []string{codeFail}}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_Segments(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/segments").
		Reply(200).
		BodyString(`{"success": true, "id": 1}`)

	data, status, err := c.Segments(SegmentsRequest{})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Segments_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/segments").
		MatchParam("filter[active]", "3").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	data, status, err := c.Segments(SegmentsRequest{Filter: SegmentsFilter{Active: 3}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_IntegrationModule(t *testing.T) {
	c := client()

	name := RandomString(5)
	code := RandomString(8)

	defer gock.Off()

	integrationModule := IntegrationModule{
		Code:            code,
		IntegrationCode: code,
		Active:          false,
		Name:            fmt.Sprintf("Integration module %s", name),
		AccountURL:      fmt.Sprintf("http://example.com/%s/account", name),
		BaseURL:         fmt.Sprintf("http://example.com/%s", name),
		ClientID:        RandomString(10),
		Logo:            "https://cdn.worldvectorlogo.com/logos/github-icon.svg",
	}

	jr, _ := json.Marshal(&integrationModule)

	pr := url.Values{
		"integrationModule": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/integration-modules/%s/edit", integrationModule.Code)).
		MatchType("url").
		BodyString(pr.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	m, status, err := c.IntegrationModuleEdit(integrationModule)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err.ApiError())
	}

	if m.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/integration-modules/%s", code)).
		Reply(200).
		BodyString(`{"success": true}`)

	g, status, err := c.IntegrationModule(code)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_IntegrationModule_Fail(t *testing.T) {
	c := client()

	name := RandomString(5)
	code := RandomString(8)

	defer gock.Off()

	integrationModule := IntegrationModule{
		IntegrationCode: code,
		Active:          false,
		Name:            fmt.Sprintf("Integration module %s", name),
		AccountURL:      fmt.Sprintf("http://example.com/%s/account", name),
		BaseURL:         fmt.Sprintf("http://example.com/%s", name),
		ClientID:        RandomString(10),
		Logo:            "https://cdn.worldvectorlogo.com/logos/github-icon.svg",
	}

	jr, _ := json.Marshal(&integrationModule)

	pr := url.Values{
		"integrationModule": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/integration-modules/%s/edit", integrationModule.Code)).
		MatchType("url").
		BodyString(pr.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	m, status, err := c.IntegrationModuleEdit(integrationModule)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if m.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/integration-modules/%s", code)).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	g, status, err := c.IntegrationModule(code)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if g.Success != false {
		t.Error(successFail)
	}
}

func TestClient_ProductsGroup(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/store/product-groups").
		MatchParam("filter[active]", "1").
		Reply(200).
		BodyString(`{"success": true}`)

	g, status, err := c.ProductsGroup(ProductsGroupsRequest{
		Filter: ProductsGroupsFilter{
			Active: 1,
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_ProductsGroup_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/store/product-groups").
		MatchParam("filter[active]", "3").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	g, status, err := c.ProductsGroup(ProductsGroupsRequest{
		Filter: ProductsGroupsFilter{
			Active: 3,
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if g.Success != false {
		t.Error(successFail)
	}
}

func TestClient_Products(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/store/products").
		MatchParam("filter[active]", "1").
		MatchParam("filter[minPrice]", "1").
		Reply(200).
		BodyString(`{"success": true}`)

	g, status, err := c.Products(ProductsRequest{
		Filter: ProductsFilter{
			Active:   1,
			MinPrice: 1,
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Products_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/store/products").
		MatchParam("filter[active]", "3").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	g, status, err := c.Products(ProductsRequest{
		Filter: ProductsFilter{Active: 3},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if g.Success != false {
		t.Error(successFail)
	}
}

func TestClient_ProductsProperties(t *testing.T) {
	c := client()

	sites := make([]string, 1)
	sites[0] = os.Getenv("RETAILCRM_SITE")

	defer gock.Off()

	gock.New(crmURL).
		Get("/store/products").
		MatchParam("filter[sites][]", sites[0]).
		Reply(200).
		BodyString(`{"success": true}`)

	g, status, err := c.ProductsProperties(ProductsPropertiesRequest{
		Filter: ProductsPropertiesFilter{
			Sites: sites,
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryShipments(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/delivery/shipments").
		MatchParam("filter[dateFrom]", "2017-10-10").
		Reply(200).
		BodyString(`{"success": true}`)

	g, status, err := c.DeliveryShipments(DeliveryShipmentsRequest{
		Filter: ShipmentFilter{
			DateFrom: "2017-10-10",
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status != http.StatusOK {
		t.Errorf("%v", err.ApiError())
	}

	if g.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_DeliveryShipments_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/delivery/shipments").
		MatchParam("filter[stores][]", codeFail).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	g, status, err := c.DeliveryShipments(DeliveryShipmentsRequest{
		Filter: ShipmentFilter{
			Stores: []string{codeFail},
		},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if g.Success != false {
		t.Error(successFail)
	}
}

func TestClient_Cost(t *testing.T) {
	c := client()

	costRecord := CostRecord{
		DateFrom: "2018-04-02",
		DateTo:   "2018-04-02",
		Summ:     124,
		CostItem: "seo",
	}

	defer gock.Off()

	str, _ := json.Marshal(costRecord)

	p := url.Values{
		"cost": {string(str)},
	}

	gock.New(crmURL).
		Post("/costs/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 123}`)

	data, status, err := c.CostCreate(costRecord)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	id := data.ID

	gock.New(crmURL).
		Get("/costs").
		MatchParam("filter[ids][]", strconv.Itoa(id)).
		MatchParam("limit", "20").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true}`)

	costs, status, err := c.Costs(CostsRequest{
		Filter: CostsFilter{
			Ids: []string{strconv.Itoa(id)},
		},
		Limit: 20,
		Page:  1,
	})

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if costs.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/costs/%d", id)).
		Reply(200).
		BodyString(`{"success": true}`)

	cost, status, err := c.Cost(id)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if cost.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	costRecord.DateFrom = "2018-04-09"
	costRecord.DateTo = "2018-04-09"
	costRecord.Summ = 421

	str, _ = json.Marshal(costRecord)

	p = url.Values{
		"cost": {string(str)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/costs/%d/edit", id)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	costEdit, status, err := c.CostEdit(id, costRecord)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if costEdit.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	j, _ := json.Marshal(&id)

	p = url.Values{
		"costs": {string(j)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/costs/%d/delete", id)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	costDelete, status, err := c.CostDelete(id)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if costDelete.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_Cost_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	costRecord := CostRecord{
		DateFrom: "2018-13-13",
		DateTo:   "2012-04-02",
		Summ:     124,
		CostItem: "seo",
	}

	str, _ := json.Marshal(costRecord)

	p := url.Values{
		"cost": {string(str)},
	}

	gock.New(crmURL).
		Post("/costs/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Cost is not loaded"}`)

	data, status, err := c.CostCreate(costRecord)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	id := data.ID

	gock.New(crmURL).
		Get("/costs").
		MatchParam("filter[sites][]", codeFail).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	costs, status, err := c.Costs(CostsRequest{
		Filter: CostsFilter{Sites: []string{codeFail}},
	})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if costs.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/costs/%d", id)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	cost, status, err := c.Cost(id)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if cost.Success != false {
		t.Error(successFail)
	}

	costRecord.DateFrom = "2020-13-12"
	costRecord.DateTo = "2012-04-09"
	costRecord.Summ = 421
	costRecord.Sites = []string{codeFail}

	str, _ = json.Marshal(costRecord)

	p = url.Values{
		"cost": {string(str)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/costs/%d/edit", id)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Cost is not loaded"}`)

	costEdit, status, err := c.CostEdit(id, costRecord)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if costEdit.Success != false {
		t.Error(successFail)
	}

	j, _ := json.Marshal(&id)

	p = url.Values{
		"costs": {string(j)},
	}
	gock.New(crmURL).
		Post(fmt.Sprintf("/costs/%d/delete", id)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	costDelete, status, err := c.CostDelete(id)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if costDelete.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CostsUpload(t *testing.T) {
	c := client()

	defer gock.Off()

	costsUpload := []CostRecord{
		{
			Source:   nil,
			DateFrom: "2018-04-02",
			DateTo:   "2018-04-02",
			Summ:     124,
			CostItem: "seo",
			Order:    nil,
		},
		{
			Source:   nil,
			DateFrom: "2018-04-03",
			DateTo:   "2018-04-03",
			Summ:     125,
			CostItem: "seo",
			Order:    nil,
			Sites:    []string{"retailcrm-ru"},
		},
	}

	j, _ := json.Marshal(&costsUpload)

	p := url.Values{
		"costs": {string(j)},
	}

	gock.New(crmURL).
		Post("/costs/upload").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "uploadedCosts": [1, 2]}`)

	data, status, err := c.CostsUpload(costsUpload)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}

	ids := data.UploadedCosts

	j, _ = json.Marshal(&ids)

	p = url.Values{
		"ids": {string(j)},
	}

	gock.New(crmURL).
		Post("/costs/delete").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	costsDelete, status, err := c.CostsDelete(ids)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if costsDelete.Success != true {

		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CostsUpload_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	costsUpload := []CostRecord{
		{
			Source:   nil,
			DateFrom: "2018-04-03",
			DateTo:   "2018-04-03",
			Summ:     125,
			CostItem: "seo",
			Order:    nil,
			Sites:    []string{codeFail},
		},
	}

	j, _ := json.Marshal(&costsUpload)

	p := url.Values{
		"costs": {string(j)},
	}

	gock.New(crmURL).
		Post("/costs/upload").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(460).
		BodyString(`{"success": false, "errorMsg": "Costs are loaded with errors"}`)

	data, status, err := c.CostsUpload(costsUpload)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}

	ids := data.UploadedCosts

	j, _ = json.Marshal(&ids)

	p = url.Values{
		"ids": {string(j)},
	}

	gock.New(crmURL).
		Post("/costs/delete").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Expected array, but got NULL: null"}`)

	costsDelete, status, err := c.CostsDelete(ids)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if costsDelete.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CustomFields(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/custom-fields").
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.CustomFields(CustomFieldsRequest{})

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomFields_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/custom-fields").
		MatchParam("filter[type]", codeFail).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	data, status, err := c.CustomFields(CustomFieldsRequest{Filter: CustomFieldsFilter{Type: codeFail}})
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CustomDictionariesCreate(t *testing.T) {
	c := client()

	code := "test_" + RandomString(8)

	defer gock.Off()

	customDictionary := CustomDictionary{
		Name: "test2",
		Code: code,
		Elements: []Element{
			{
				Name: "test",
				Code: "test",
			},
		},
	}

	j, _ := json.Marshal(&customDictionary)

	p := url.Values{
		"customDictionary": {string(j)},
	}

	gock.New(crmURL).
		Post("/custom-fields/dictionaries/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 1}`)

	data, status, err := c.CustomDictionariesCreate(customDictionary)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	gock.New(crmURL).
		Get("/custom-fields/dictionaries").
		MatchParam("filter[name]", "test").
		MatchParam("limit", "10").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true}`)

	cds, status, err := c.CustomDictionaries(CustomDictionariesRequest{
		Filter: CustomDictionariesFilter{
			Name: "test",
		},
		Limit: 10,
		Page:  1,
	})

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {

		t.Errorf("%v", err.ApiError())
	}

	if cds.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/custom-fields/dictionaries/%s", code)).
		Reply(200).
		BodyString(`{"success": true}`)

	cd, status, err := c.CustomDictionary(code)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if cd.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	customDictionary.Name = "test223"
	customDictionary.Elements = []Element{
		{
			Name: "test3",
			Code: "test3",
		},
	}

	j, _ = json.Marshal(&customDictionary)

	p = url.Values{
		"customDictionary": {string(j)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/custom-fields/dictionaries/%s/edit", customDictionary.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	cde, status, err := c.CustomDictionaryEdit(customDictionary)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if cde.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomDictionariesCreate_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	customDictionary := CustomDictionary{
		Name: "test2",
		Code: RandomString(8),
		Elements: []Element{
			{
				Name: "test",
				Code: "test",
			},
		},
	}

	j, _ := json.Marshal(&customDictionary)

	p := url.Values{
		"customDictionary": {string(j)},
	}

	gock.New(crmURL).
		Post("/custom-fields/dictionaries/create").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	data, status, err := c.CustomDictionariesCreate(customDictionary)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if data.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/custom-fields/dictionaries/%s", codeFail)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	cd, status, err := c.CustomDictionary(codeFail)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if cd.Success != false {
		t.Error(successFail)
	}

	customDictionary.Name = "test223"
	customDictionary.Elements = []Element{
		{
			Name: "test3",
		},
	}

	j, _ = json.Marshal(&customDictionary)

	p = url.Values{
		"customDictionary": {string(j)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/custom-fields/dictionaries/%s/edit", customDictionary.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	cde, status, err := c.CustomDictionaryEdit(customDictionary)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Errorf("%v", err.ApiError())
	}

	if cde.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CustomFieldsCreate(t *testing.T) {
	c := client()
	codeCustomField := RandomString(8)

	defer gock.Off()

	customFields := CustomFields{
		Name:        codeCustomField,
		Code:        codeCustomField,
		Type:        "text",
		Entity:      "order",
		DisplayArea: "customer",
	}

	j, _ := json.Marshal(&customFields)

	p := url.Values{
		"customField": {string(j)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/custom-fields/%s/create", customFields.Entity)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true, "id": 1}`)

	data, status, err := c.CustomFieldsCreate(customFields)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if data.Success != true {

		t.Errorf("%v", err.ApiError())
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/custom-fields/%s/%s", "order", codeCustomField)).
		Reply(200).
		BodyString(`{"success": true}`)

	customField, status, err := c.CustomField("order", codeCustomField)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if customField.Success != true {
		t.Errorf("%v", err.ApiError())
	}

	customFields.DisplayArea = "delivery"

	j, _ = json.Marshal(&customFields)

	p = url.Values{
		"customField": {string(j)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/custom-fields/%s/%s/edit", customFields.Entity, customFields.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	customFieldEdit, status, err := c.CustomFieldEdit(customFields)

	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if !statuses[status] {
		t.Errorf("%v", err.ApiError())
	}

	if customFieldEdit.Success != true {
		t.Errorf("%v", err.ApiError())
	}
}

func TestClient_CustomFieldsCreate_Fail(t *testing.T) {
	c := client()

	codeCustomField := "test_" + RandomString(8)

	defer gock.Off()

	customFields := CustomFields{
		Name:        codeCustomField,
		Type:        "text",
		Entity:      "order",
		DisplayArea: "customer",
	}

	j, _ := json.Marshal(&customFields)

	p := url.Values{
		"customField": {string(j)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/custom-fields/%s/create", customFields.Entity)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters"}`)

	data, status, err := c.CustomFieldsCreate(customFields)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/custom-fields/%s/%s", "order", codeCustomField)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	customField, status, err := c.CustomField("order", codeCustomField)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if customField.Success != false {
		t.Error(successFail)
	}

	customFields.DisplayArea = "delivery"

	j, _ = json.Marshal(&customFields)

	p = url.Values{
		"customField": {string(j)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/custom-fields/%s/%s/edit", customFields.Entity, customFields.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "API method not found"}`)

	customFieldEdit, status, err := c.CustomFieldEdit(customFields)
	if err.RuntimeErr != nil {
		t.Errorf("%v", err.Error())
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if customFieldEdit.Success != false {
		t.Error(successFail)
	}
}

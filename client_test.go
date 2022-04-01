package retailcrm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/joho/godotenv"
	gock "gopkg.in/h2non/gock.v1"
)

func TestMain(m *testing.M) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	os.Exit(m.Run())
}

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

var (
	r        *rand.Rand // Rand for this package.
	user, _  = strconv.Atoi(os.Getenv("RETAILCRM_USER"))
	statuses = map[int]bool{http.StatusOK: true, http.StatusCreated: true}
	crmURL   = os.Getenv("RETAILCRM_URL")
	badURL   = "https://qwertypoiu.retailcrm.pro"

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

func TestGetRequest(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/api/v5/fake-method").
		Reply(404).
		BodyString(`{"success": false, "errorMsg" : "Method not found"}`)

	_, status, _ := c.GetRequest("/fake-method")

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestPostRequest(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/api/v5/fake-method").
		Reply(404).
		BodyString(`{"success": false, "errorMsg" : "Method not found"}`)

	_, status, _ := c.PostRequest("/fake-method", url.Values{})

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestClient_ApiVersionsVersions(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/api/api-versions").
		Reply(200).
		BodyString(`{"success": true, "versions" : ["1.0", "4.0", "3.0", "4.0", "5.0"]}`)

	data, _, err := c.APIVersions()
	if err != nil {
		t.Errorf("%v", err.Error())
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_ApiVersionsVersionsBadUrl(t *testing.T) {
	c := badurlclient()

	defer gock.Off()

	gock.New(badURL).
		Get("/api/api-versions").
		Reply(400).
		BodyString(`{"success": false, "errorMsg" : "Account does not exist"}`)

	data, _, err := c.APIVersions()
	if err == nil {
		t.Error("Error must be return")
	}

	if data.Success != false {
		t.Logf("%v", err)
	}
}

func TestClient_ApiCredentialsCredentialsBadKey(t *testing.T) {
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/credentials").
		Reply(403).
		BodyString(`{"success": false, "errorMsg": "Wrong \"apiKey\" value"}`)

	c := client()

	data, _, err := c.APICredentials()
	if err == nil {
		t.Error("Error must be return")
	}

	if data.Success != false {
		t.Logf("%v", err)
	}
}

func TestClient_ApiCredentialsCredentials(t *testing.T) {
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/credentials").
		Reply(200).
		BodyString(`{"success": true}`)

	c := client()

	data, _, err := c.APICredentials()
	if err != nil {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_APISystemInfo(t *testing.T) {
	defer gock.Off()

	r := SystemInfoResponse{
		Success:       true,
		SystemVersion: "8.7.77",
		PublicURL:     crmURL,
		TechnicalURL:  fmt.Sprintf("https://%s.io", RandomString(30)),
	}

	data, err := json.Marshal(r)
	if err != nil {
		t.Errorf("%v", err)
	}

	gock.New(crmURL).
		Get("/api/system-info").
		Reply(200).
		BodyString(string(data))

	c := client()
	res, s, e := c.APISystemInfo()
	if e != nil {
		t.Errorf("%v", e)
	}

	assert.Equal(t, s, 200)
	assert.Equal(t, crmURL, res.PublicURL)
	assert.Contains(t, res.TechnicalURL, ".io")
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Logf("%v", err)
	}

	if data.Success != true {
		t.Logf("%v", err)
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

	data, _, err := c.Customers(CustomersRequest{
		Filter: CustomersFilter{
			Ids: []string{codeFail},
		},
	})

	if err == nil {
		t.Error("Error must be return")
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
		Tags: []Tag{
			{"first", "#3e89b6", false},
			{"second", "#ffa654", false},
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err)
	}

	if cr.Success != true {
		t.Errorf("%v", err)
	}

	f.ID = cr.ID
	f.Vip = true

	str, _ = json.Marshal(f)

	p = url.Values{
		"by":       {string(ByID)},
		"customer": {string(str)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/customers/%v/edit", cr.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	ed, se, err := c.CustomerEdit(f, ByID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if se != http.StatusOK {
		t.Errorf("%v", err)
	}

	if ed.Success != true {
		t.Errorf("%v", err)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/api/v5/customers/%v", f.ExternalID)).
		MatchParam("by", ByExternalID).
		Reply(200).
		BodyString(`{
			"success": true, 
			"customer": {
				"tags": [
					{
						"name": "first",
						"color": "#3e89b6",
						"attached": false
					},
					{
						"name": "second",
						"color": "#ffa654",
						"attached": false
					}
				]
			}
		}`)

	data, _, err := c.Customer(f.ExternalID, ByExternalID, "")
	if err != nil {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(data.Customer.Tags, f.Tags) {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
		"by":       {string(ByID)},
		"customer": {string(str)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/customers/%v/edit", cr.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	ed, se, err := c.CustomerEdit(f, ByID)
	if err == nil {
		t.Error("Error must be return")
	}

	if se < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if ed.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/api/v5/customers/%v", codeFail)).
		MatchParam("by", ByExternalID).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	data, _, err := c.Customer(codeFail, ByExternalID, "")
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Customers are loaded with ErrorsList"}`)

	data, _, err := c.CustomersUpload(customers)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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

	data, _, err := c.CustomersCombine([]Customer{{}, {}}, Customer{})
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if fe != http.StatusOK {
		t.Errorf("%v", err)
	}

	if fx.Success != true {
		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "ErrorsList": {"id": "ID must be an integer"}}`)

	data, _, err := c.CustomersFixExternalIds(customers)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "ErrorsList": {"children[startDate]": "Значение недопустимо."}}`)

	data, _, err := c.CustomersHistory(f)
	if err == nil {
		t.Error("Error must be return")
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CorporateCustomersList(t *testing.T) {
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/v5/customers-corporate").
		MatchParam("filter[city]", "Москва").
		MatchParam("page", "3").
		Reply(200).
		BodyString(`{"success":true,"pagination":{"limit":20,"totalCount":1,"currentPage":3,"totalPageCount":1}}`)

	c := client()

	data, status, err := c.CorporateCustomers(CorporateCustomersRequest{
		Filter: CorporateCustomersFilter{
			City: "Москва",
		},
		Page: 3,
	})

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Logf("%v", err)
	}

	if data.Success != true {
		t.Logf("%v", err)
	}
}

func TestClient_CorporateCustomersCreate(t *testing.T) {
	defer gock.Off()

	gock.New(crmURL).
		Post("/api/v5/customers-corporate/create").
		Reply(201).
		BodyString(`{"success":true,"id":2344}`)

	c := client()
	customer := CorporateCustomer{
		ExternalID:         "ext-id",
		Nickname:           "Test Customer",
		Vip:                true,
		Bad:                false,
		CustomFields:       nil,
		PersonalDiscount:   10,
		DiscountCardNumber: "1234567890",
		Source: &Source{
			Source:   "source",
			Medium:   "medium",
			Campaign: "campaign",
			Keyword:  "keyword",
			Content:  "content",
		},
		Companies: []Company{
			{
				IsMain:     true,
				ExternalID: "company-ext-id",
				Active:     true,
				Name:       "name",
				Brand:      "brand",
				Site:       "https://retailcrm.pro",
				Contragent: &Contragent{
					ContragentType: "legal-entity",
					LegalName:      "Legal Name",
					LegalAddress:   "Legal Address",
					INN:            "000000000",
					OKPO:           "000000000",
					KPP:            "000000000",
					OGRN:           "000000000",
					BIK:            "000000000",
					Bank:           "bank",
					BankAddress:    "bankAddress",
					CorrAccount:    "corrAccount",
					BankAccount:    "bankAccount",
				},
				Address: &IdentifiersPair{
					ID:         0,
					ExternalID: "ext-addr-id",
				},
				CustomFields: nil,
			},
		},
		Addresses: []CorporateCustomerAddress{
			{
				Index:        "123456",
				CountryISO:   "RU",
				Region:       "Russia",
				RegionID:     0,
				City:         "Moscow",
				CityID:       0,
				CityType:     "city",
				Street:       "Pushkinskaya",
				StreetID:     0,
				StreetType:   "street",
				Building:     "",
				Flat:         "",
				IntercomCode: "",
				Floor:        0,
				Block:        0,
				House:        "",
				Housing:      "",
				Metro:        "",
				Notes:        "",
				Text:         "",
				ExternalID:   "ext-addr-id",
				Name:         "Main Address",
			},
		},
	}

	data, status, err := c.CorporateCustomerCreate(customer, "site")

	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_CorporateCustomersFixExternalIds(t *testing.T) {
	c := client()

	customers := []IdentifiersPair{{
		ID:         123,
		ExternalID: RandomString(8),
	}}

	defer gock.Off()

	jr, _ := json.Marshal(&customers)

	p := url.Values{
		"customersCorporate": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/customers-corporate/fix-external-ids").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	fx, fe, err := c.CorporateCustomersFixExternalIds(customers)
	if err != nil {
		t.Errorf("%v", err)
	}

	if fe != http.StatusOK {
		t.Errorf("%v", err)
	}

	if fx.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_CorporateCustomersFixExternalIds_Fail(t *testing.T) {
	c := client()

	customers := []IdentifiersPair{{ExternalID: RandomString(8)}}

	defer gock.Off()

	jr, _ := json.Marshal(&customers)

	p := url.Values{
		"customersCorporate": {string(jr[:])},
	}

	gock.New(crmURL).
		Post("/customers-corporate/fix-external-ids").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "ErrorsList": {"id": "ID must be an integer"}}`)

	data, _, err := c.CorporateCustomersFixExternalIds(customers)
	if err == nil {
		t.Error("Error must be return")
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CorporateCustomersHistory(t *testing.T) {
	c := client()
	f := CorporateCustomersHistoryRequest{
		Filter: CorporateCustomersHistoryFilter{
			SinceID: 20,
		},
	}
	defer gock.Off()

	gock.New(crmURL).
		Get("/customers-corporate/history").
		MatchParam("filter[sinceId]", "20").
		Reply(200).
		BodyString(`{"success": true, "history": [{"id": 1}]}`)

	data, status, err := c.CorporateCustomersHistory(f)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if len(data.History) == 0 {
		t.Errorf("%v", "Empty history")
	}
}

func TestClient_CorporateCustomersHistory_Fail(t *testing.T) {
	c := client()
	f := CorporateCustomersHistoryRequest{
		Filter: CorporateCustomersHistoryFilter{
			StartDate: "2020-13-12",
		},
	}

	defer gock.Off()

	gock.New(crmURL).
		Get("/customers-corporate/history").
		MatchParam("filter[startDate]", "2020-13-12").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "ErrorsList": {"children[startDate]": "Значение недопустимо."}}`)

	data, _, err := c.CorporateCustomersHistory(f)
	if err == nil {
		t.Error("Error must be return")
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CorporateCustomersNotes(t *testing.T) {
	defer gock.Off()

	gock.New(crmURL).
		Get("/api/v5/customers-corporate/notes").
		MatchParam("filter[text]", "sample").
		Reply(200).
		BodyString(`
		{
		  "success": true,
		  "pagination": {
		    "limit": 20,
		    "totalCount": 1,
		    "currentPage": 1,
		    "totalPageCount": 1
		  },
		  "notes": [
		    {
		      "customer": {
		        "site": "site",
		        "id": 2346,
		        "externalId": "ext-id",
		        "type": "customer_corporate"
		      },
		      "managerId": 24,
		      "id": 106,
		      "text": "<p>sample text</p>",
		      "createdAt": "2019-10-15 17:08:59"
		    }
		  ]
		}
		`)

	c := client()

	data, status, err := c.CorporateCustomersNotes(CorporateCustomersNotesRequest{
		Filter: CorporateCustomersNotesFilter{
			Text: "sample",
		},
	})

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if data.Notes[0].Text != "<p>sample text</p>" {
		t.Errorf("invalid note text")
	}
}

func TestClient_CorporateCustomerNoteCreate(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/customers-corporate/notes/create").
		Reply(201).
		BodyString(`{"success":true,"id":1}`)

	data, status, err := c.CorporateCustomerNoteCreate(CorporateCustomerNote{
		Text: "another note",
		Customer: &IdentifiersPair{
			ID: 1,
		},
	}, "site")

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if data.ID != 1 {
		t.Error("invalid note id")
	}
}

func TestClient_CorporateCustomerNoteDelete(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/customers-corporate/notes/1/delete").
		Reply(200).
		BodyString(`{"success":true}`)

	data, status, err := c.CorporateCustomerNoteDelete(1)

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_CorporateCustomersUpload(t *testing.T) {
	c := client()
	customers := make([]CorporateCustomer, 3)

	for i := range customers {
		customers[i] = CorporateCustomer{
			Nickname:   fmt.Sprintf("Name_%s", RandomString(8)),
			ExternalID: RandomString(8),
		}
	}

	defer gock.Off()

	str, _ := json.Marshal(customers)

	p := url.Values{
		"customersCorporate": {string(str)},
	}

	gock.New(crmURL).
		Post("/api/v5/customers-corporate/upload").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.CorporateCustomersUpload(customers)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_CorporateCustomersUpload_Fail(t *testing.T) {
	c := client()

	customers := []CorporateCustomer{{ExternalID: strconv.Itoa(iCodeFail)}}

	defer gock.Off()

	str, _ := json.Marshal(customers)
	p := url.Values{
		"customersCorporate": {string(str)},
	}

	gock.New(crmURL).
		Post("/api/v5/customers-corporate/upload").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(460).
		BodyString(`{"success": false, "errorMsg": "Customers are loaded with ErrorsList"}`)

	data, _, err := c.CorporateCustomersUpload(customers)
	if err == nil {
		t.Error("Error must be return")
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_CorporateCustomer(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/customers-corporate/ext-id").
		MatchParam("by", ByExternalID).
		MatchParam("site", "site").
		Reply(200).
		BodyString(`
		{
		  "success": true,
		  "customerCorporate": {
		    "type": "customer_corporate",
		    "id": 2346,
		    "externalId": "ext-id",
		    "nickName": "Test Customer 2",
		    "mainAddress": {
		      "id": 2034,
		      "externalId": "ext-addr-id223",
		      "name": "Main Address"
		    },
		    "createdAt": "2019-10-15 16:16:56",
		    "vip": false,
		    "bad": false,
		    "site": "site",
		    "marginSumm": 0,
		    "totalSumm": 0,
		    "averageSumm": 0,
		    "ordersCount": 0,
		    "costSumm": 0,
		    "customFields": [],
		    "personalDiscount": 10,
		    "mainCompany": {
		      "id": 26,
		      "externalId": "company-ext-id",
		      "name": "name"
		    }
		  }
		}
		`)

	data, status, err := c.CorporateCustomer("ext-id", ByExternalID, "site")
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_CorporateCustomerAddresses(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/customers-corporate/ext-id/addresses").
		MatchParams(map[string]string{
			"by":           ByExternalID,
			"filter[name]": "Main Address",
			"limit":        "20",
			"page":         "1",
			"site":         "site",
		}).
		Reply(200).
		BodyString(`
		{
		  "success": true,
		  "addresses": [
		    {
		      "id": 2034,
		      "index": "123456",
		      "countryIso": "RU",
		      "region": "Russia",
		      "city": "Moscow",
		      "cityType": "city",
		      "street": "Pushkinskaya",
		      "streetType": "street",
		      "text": "street Pushkinskaya",
		      "externalId": "ext-addr-id223",
		      "name": "Main Address"
		    }
		  ]
		}
		`)

	data, status, err := c.CorporateCustomerAddresses("ext-id", CorporateCustomerAddressesRequest{
		Filter: CorporateCustomerAddressesFilter{
			Name: "Main Address",
		},
		By:    ByExternalID,
		Site:  "site",
		Limit: 20,
		Page:  1,
	})

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if len(data.Addresses) == 0 {
		t.Error("data.Addresses must not be empty")
	}
}

func TestClient_CorporateCustomerAddressesCreate(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/customers-corporate/ext-id/addresses/create").
		Reply(201).
		BodyString(`{"success":true,"id":1}`)

	data, status, err := c.CorporateCustomerAddressesCreate("ext-id", ByExternalID, CorporateCustomerAddress{
		Text: "this is new address",
		Name: "New Address",
	}, "site")

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_CorporateCustomerAddressesEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/customers-corporate/customer-ext-id/addresses/addr-ext-id/edit").
		Reply(200).
		BodyString(`{"success":true,"id":1}`)

	data, status, err := c.CorporateCustomerAddressesEdit(
		"customer-ext-id",
		ByExternalID,
		ByExternalID,
		CorporateCustomerAddress{
			ExternalID: "addr-ext-id",
			Name:       "Main Address 2",
		},
		"site",
	)

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_CorporateCustomerCompanies(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/customers-corporate/ext-id/companies").
		MatchParams(map[string]string{
			"by":            ByExternalID,
			"filter[ids][]": "1",
			"limit":         "20",
			"page":          "1",
			"site":          "site",
		}).
		Reply(200).
		BodyString(`
		{
		  "success": true,
		  "companies": [
		    {
		      "isMain": true,
		      "id": 1,
		      "externalId": "company-ext-id",
		      "customer": {
		        "site": "site",
		        "id": 2346,
		        "externalId": "ext-id",
		        "type": "customer_corporate"
		      },
		      "active": true,
		      "name": "name",
		      "brand": "brand",
		      "site": "https://retailcrm.pro",
		      "createdAt": "2019-10-15 16:16:56",
		      "contragent": {
		        "contragentType": "legal-entity",
		        "legalName": "Legal Name",
		        "legalAddress": "Legal Address",
		        "INN": "000000000",
		        "OKPO": "000000000",
		        "KPP": "000000000",
		        "OGRN": "000000000",
		        "BIK": "000000000",
		        "bank": "bank",
		        "bankAddress": "bankAddress",
		        "corrAccount": "corrAccount",
		        "bankAccount": "bankAccount"
		      },
		      "address": {
		        "id": 2034,
		        "index": "123456",
		        "countryIso": "RU",
		        "region": "Russia",
		        "city": "Moscow",
		        "cityType": "city",
		        "street": "Pushkinskaya",
		        "streetType": "street",
		        "text": "street Pushkinskaya",
		        "externalId": "ext-addr-id",
		        "name": "Main Address 2"
		      },
		      "marginSumm": 0,
		      "totalSumm": 0,
		      "averageSumm": 0,
		      "ordersCount": 0,
		      "costSumm": 0,
		      "customFields": []
		    }
		  ]
		}
		`)

	data, status, err := c.CorporateCustomerCompanies("ext-id", IdentifiersPairRequest{
		Filter: IdentifiersPairFilter{
			Ids: []string{"1"},
		},
		By:    ByExternalID,
		Site:  "site",
		Limit: 20,
		Page:  1,
	})

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if len(data.Companies) == 0 {
		t.Error("data.Companies must not be empty")
	}
}

func TestClient_CorporateCustomerCompaniesCreate(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/customers-corporate/ext-id/companies/create").
		Reply(201).
		BodyString(`{"success":true,"id":1}`)

	data, status, err := c.CorporateCustomerCompaniesCreate("ext-id", ByExternalID, Company{
		Name: "New Company",
	}, "site")

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_CorporateCustomerCompaniesEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/customers-corporate/customer-ext-id/companies/company-ext-id/edit").
		Reply(200).
		BodyString(`{"success":true,"id":1}`)

	data, status, err := c.CorporateCustomerCompaniesEdit(
		"customer-ext-id",
		ByExternalID,
		ByExternalID,
		Company{
			ExternalID: "company-ext-id",
			Name:       "New Company Name 2",
		},
		"site",
	)

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_CorporateCustomerContacts(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/customers-corporate/ext-id/contacts").
		MatchParams(map[string]string{
			"by":    ByExternalID,
			"limit": "20",
			"page":  "1",
			"site":  "site",
		}).
		Reply(200).
		BodyString(`
		{
		  "success": true,
		  "contacts": [
		    {
		      "isMain": false,
		      "customer": {
		        "id": 2347,
		        "site": "site"
		      },
		      "companies": []
		    }
		  ]
		}
		`)

	data, status, err := c.CorporateCustomerContacts("ext-id", IdentifiersPairRequest{
		Filter: IdentifiersPairFilter{},
		By:     ByExternalID,
		Site:   "site",
		Limit:  20,
		Page:   1,
	})

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if len(data.Contacts) == 0 {
		t.Error("data.Contacts must not be empty")
	}
}

func TestClient_CorporateCustomerContactsCreate(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/customers/create").
		Reply(201).
		BodyString(`{"success":true,"id":2}`)

	gock.New(crmURL).
		Post("/customers-corporate/ext-id/contacts/create").
		Reply(201).
		BodyString(`{"success":true,"id":3}`)

	createResponse, createStatus, createErr := c.CustomerCreate(Customer{
		ExternalID: "test-customer-as-contact-person",
		FirstName:  "Contact",
		LastName:   "Person",
	}, "site")

	if createErr != nil {
		t.Errorf("%v", createErr)
	}

	if createStatus >= http.StatusBadRequest {
		t.Errorf("%v", createErr)
	}

	if createResponse.Success != true {
		t.Errorf("%v", createErr)
	}

	if createResponse.ID != 2 {
		t.Errorf("invalid createResponse.ID: should be `2`, got `%d`", createResponse.ID)
	}

	data, status, err := c.CorporateCustomerContactsCreate("ext-id", ByExternalID, CorporateCustomerContact{
		IsMain: false,
		Customer: CorporateCustomerContactCustomer{
			ExternalID: "test-customer-as-contact-person",
			Site:       "site",
		},
		Companies: []IdentifiersPair{},
	}, "site")

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if data.ID != 3 {
		t.Errorf("invalid data.ID: should be `3`, got `%d`", data.ID)
	}
}

func TestClient_CorporateCustomerContactsEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/customers-corporate/ext-id/contacts/2350/edit").
		Reply(200).
		BodyString(`{"success":true,"id":19}`)

	data, status, err := c.CorporateCustomerContactsEdit("ext-id", ByExternalID, ByID, CorporateCustomerContact{
		IsMain: false,
		Customer: CorporateCustomerContactCustomer{
			ID: 2350,
		},
	}, "site")

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("(%d) %v", status, err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if data.ID == 0 {
		t.Errorf("invalid data.ID: should be `19`, got `%d`", data.ID)
	}
}

func TestClient_CorporateCustomerEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Post("/customers-corporate/ext-id/edit").
		Reply(200).
		BodyString(`{"success":true,"id":2346}`)

	data, status, err := c.CorporateCustomerEdit(CorporateCustomer{
		ExternalID: "ext-id",
		Nickname:   "Another Nickname 2",
		Vip:        true,
	}, ByExternalID, "site")

	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("(%d) %v", status, err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if len(data.Notes) == 0 {
		t.Errorf("%v", err)
	}
}

func TestClient_NotesNotes_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/customers/notes").
		MatchParam("filter[createdAtFrom]", "2020-13-12").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "ErrorsList": {"children[createdAtFrom]": "This value is not valid."}}`)

	data, _, err := c.CustomerNotes(NotesRequest{
		Filter: NotesFilter{CreatedAtFrom: "2020-13-12"},
	})
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if noteCreateStatus != http.StatusCreated {
		t.Errorf("%v", err)
	}

	if noteCreateResponse.Success != true {
		t.Errorf("%v", err)
	}

	p = url.Values{
		"id": {"1"},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/api/v5/customers/notes/%d/delete", 1)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	noteDeleteResponse, noteDeleteStatus, err := c.CustomerNoteDelete(noteCreateResponse.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if noteDeleteStatus != http.StatusOK {
		t.Errorf("%v", err)
	}

	if noteDeleteResponse.Success != true {

		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Errors in the entity format", "ErrorsList": {"customer": "Set one of the following fields: id, externalId"}}`)

	data, _, err := c.CustomerNoteCreate(note)
	if err == nil {
		t.Error("Error must be return")
	}

	if data.Success != false {
		t.Error(successFail)
	}

	p = url.Values{
		"id": {strconv.Itoa(iCodeFail)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/customers/notes/%d/delete", iCodeFail)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Note not found map"}`)

	noteDeleteResponse, noteDeleteStatus, err := c.CustomerNoteDelete(iCodeFail)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
	}
}

func TestClient_OrdersOrders_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders").
		MatchParam("filter[attachments]", "7").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "ErrorsList": {"children[attachments]": "SThis value is not valid."}}`)

	data, _, err := c.Orders(OrdersRequest{Filter: OrdersFilter{Attachments: 7}})
	if err == nil {
		t.Error("Error must be return")
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
		BodyString(`{
    "success": true,
    "id": 1,
    "order": {
        "slug": 1,
        "id": 1,
        "number": "1A",
        "orderMethod": "shopping-cart",
        "countryIso": "RU",
        "createdAt": "2020-07-14 11:44:43",
        "statusUpdatedAt": "2020-07-14 11:44:43",
        "summ": 0,
        "totalSumm": 0,
        "prepaySum": 0,
        "purchaseSumm": 0,
        "markDatetime": "2020-07-14 11:44:43",
        "call": false,
        "expired": false,
        "customer": {
            "type": "customer",
            "id": 1,
            "isContact": false,
            "createdAt": "2020-07-14 11:44:43",
            "vip": false,
            "bad": false,
            "site": "site",
            "contragent": {
                "contragentType": "individual"
            },
            "tags": [],
            "marginSumm": 0,
            "totalSumm": 0,
            "averageSumm": 0,
            "ordersCount": 0,
            "personalDiscount": 0,
            "segments": [],
            "email": "",
            "phones": []
        },
        "contact": {
            "type": "customer",
            "id": 4512,
            "isContact": false,
            "createdAt": "2020-07-14 11:44:43",
            "vip": false,
            "bad": false,
            "site": "site",
            "contragent": {
                "contragentType": "individual"
            },
            "tags": [],
            "marginSumm": 0,
            "totalSumm": 0,
            "averageSumm": 0,
            "ordersCount": 0,
            "personalDiscount": 0,
            "segments": [],
            "email": "",
            "phones": []
        },
        "contragent": {
            "contragentType": "individual"
        },
        "delivery": {
            "cost": 0,
            "netCost": 0,
            "address": {}
        },
        "site": "site",
        "status": "new",
        "items": [],
        "payments": [],
        "fromApi": true,
        "shipmentStore": "main",
        "shipped": false
    }
}`)

	cr, sc, err := c.OrderCreate(f)
	if err != nil {
		t.Errorf("%v", err)
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err)
	}

	if cr.Success != true {
		t.Errorf("%v", err)
	}

	if cr.Order.Number != "1A" {
		t.Errorf("invalid order number: got %s want %s", cr.Order.Number, "1A")
	}

	f.ID = cr.ID
	f.CustomerComment = "test comment"

	jr, _ = json.Marshal(&f)

	p = url.Values{
		"by":    {string(ByID)},
		"order": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/%d/edit", f.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	ed, se, err := c.OrderEdit(f, ByID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if se != http.StatusOK {
		t.Errorf("%v", err)
	}

	if ed.Success != true {
		t.Errorf("%v", err)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/orders/%s", f.ExternalID)).
		MatchParam("by", ByExternalID).
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.Order(f.ExternalID, ByExternalID, "")
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
		"by":    {string(ByID)},
		"order": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/%d/edit", f.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found map"}`)

	data, _, err := c.OrderEdit(f, ByID)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Orders are loaded with ErrorsList"}`)

	data, _, err := c.OrdersUpload(orders)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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

	data, _, err := c.OrdersCombine("ours", Order{}, Order{})
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if fe != http.StatusOK {
		t.Errorf("%v", err)
	}

	if fx.Success != true {

		t.Errorf("%v", err)
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

	data, _, err := c.OrdersFixExternalIds(orders)
	if err == nil {
		t.Error("Error must be return")
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_OrdersStatuses(t *testing.T) {
	c := client()
	defer gock.Off()

	gock.New(crmURL).
		Get("/orders/statuses").
		MatchParam("ids[]", "1").
		MatchParam("externalIds[]", "2").
		Reply(200).
		BodyString(`
			{
				"success": true,
				"orders": [
					{
						"id": 1,
						"externalId": "123",
						"status": "New",
						"group": "new"
					},
					{
						"id": 123,
						"externalId": "2",
						"status": "New",
						"group": "new"
					}
				]
			}`)

	data, status, err := c.OrdersStatuses(OrdersStatusesRequest{
		IDs:         []int{1},
		ExternalIDs: []string{"2"},
	})
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if len(data.Orders) == 0 {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if len(data.History) == 0 {

		t.Errorf("%v", err)
	}
}

func TestClient_OrdersHistory_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/orders/history").
		MatchParam("filter[startDate]", "2020-13-12").
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "ErrorsList": {"children[startDate]": "Значение недопустимо."}}`)

	data, _, err := c.OrdersHistory(OrdersHistoryRequest{Filter: OrdersHistoryFilter{StartDate: "2020-13-12"}})
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err)
	}

	if paymentCreateResponse.Success != true {
		t.Errorf("%v", err)
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

	paymentEditResponse, status, err := c.OrderPaymentEdit(k, ByID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if paymentEditResponse.Success != true {
		t.Errorf("%v", err)
	}

	p = url.Values{
		"id": {strconv.Itoa(paymentCreateResponse.ID)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/payments/%d/delete", paymentCreateResponse.ID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	paymentDeleteResponse, status, err := c.OrderPaymentDelete(paymentCreateResponse.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if paymentDeleteResponse.Success != true {
		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Errors in the entity format", "ErrorsList": {"order": "Set one of the following fields: id, externalId, number"}}`)

	data, _, err := c.OrderPaymentCreate(f)
	if err == nil {
		t.Error("Error must be return")
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

	paymentEditResponse, status, err := c.OrderPaymentEdit(k, ByID)
	if err == nil {
		t.Error("Error must be return")
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if paymentEditResponse.Success != false {
		t.Error(successFail)
	}

	p = url.Values{
		"id": {strconv.Itoa(iCodeFail)},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/payments/%d/delete", iCodeFail)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Payment not found"}`)

	paymentDeleteResponse, status, err := c.OrderPaymentDelete(iCodeFail)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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

	data, _, err := c.Tasks(f)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if sc != http.StatusCreated {
		t.Errorf("%v", err)
	}

	if cr.Success != true {
		t.Errorf("%v", err)
	}

	f.ID = cr.ID
	f.Commentary = RandomString(20)

	gock.New(crmURL).
		Get(fmt.Sprintf("/tasks/%d", f.ID)).
		Reply(200).
		BodyString(`{"success": true}`)

	gt, sg, err := c.Task(f.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if sg != http.StatusOK {
		t.Errorf("%v", err)
	}

	if gt.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Task is not loaded", "ErrorsList": {"performerId": "This value should not be blank."}}`)

	data, _, err := c.TaskEdit(f)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Errors in the input parameters", "ErrorsList": {"active": "he value you selected is not a valid choice."}}`)

	data, _, err := c.Users(UsersRequest{Filter: UsersFilter{Active: 3}, Page: 1})
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
	}
}

func TestClient_UsersUser_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get(fmt.Sprintf("/users/%d", iCodeFail)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	data, _, err := c.User(iCodeFail)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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

	data, _, err := c.UserStatus(iCodeFail, "busy")
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_OrderPaymentEdit(t *testing.T) {
	c := client()
	payment := Payment{
		ExternalID: RandomString(8),
	}

	defer gock.Off()

	jr, _ := json.Marshal(&payment)
	p := url.Values{
		"by":      {"externalId"},
		"payment": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/payments/%s/edit", payment.ExternalID)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	data, status, err := c.OrderPaymentEdit(payment, "externalId")
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if idata.Success != true {

		t.Errorf("%v", err)
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

	data, _, err := c.CostGroupEdit(costGroup)
	if err == nil {
		t.Error("Error must be return")
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
	if err == nil {
		t.Error("Error must be return")
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
		FirstName: RandomString(5),
		LastName:  RandomString(5),
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	cur.ID = data.ID
	cur.Patronymic = RandomString(5)

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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if idata.Success != true {

		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Errors in the entity format", "ErrorsList": {"firstName": "Specify the first name"}}`)

	data, st, err := c.CourierCreate(Courier{})
	if err == nil {
		t.Error("Error must be return")
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}

	cur := Courier{Patronymic: RandomString(5)}
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if data.Success == false {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
	}

	if st < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if data.Success != false {
		t.Error(successFail)
	}
}

func TestClient_Units(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/api/v5/reference/units").
		Reply(200).
		BodyString(`{"success": true, "units": []}`)

	data, st, err := c.Units()
	if err != nil {
		t.Errorf("%v", err)
	}

	if st != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_UnitsEdit(t *testing.T) {
	c := client()

	defer gock.Off()

	unit := Unit{
		Code:    RandomString(5),
		Name:    RandomString(5),
		Sym:     RandomString(2),
		Default: false,
		Active:  true,
	}

	jr, _ := json.Marshal(&unit)

	p := url.Values{
		"unit": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/api/v5/reference/units/%s/edit", unit.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(201).
		BodyString(`{"success": true}`)

	data, st, err := c.UnitEdit(unit)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[st] {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_UnitEdit_Fail(t *testing.T) {
	c := client()

	defer gock.Off()

	unit := Unit{
		Name:   RandomString(5),
		Active: false,
	}

	jr, _ := json.Marshal(&unit)

	p := url.Values{
		"unit": {string(jr[:])},
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/api/v5/reference/units/%s/edit", unit.Code)).
		MatchType("url").
		BodyString(p.Encode()).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Method not found"}`)

	data, st, err := c.UnitEdit(unit)
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err)
	}

	if p.Success != true {
		t.Errorf("%v", err)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/orders/packs/%d", p.ID)).
		Reply(200).
		BodyString(`{"success": true}`)

	s, status, err := c.Pack(p.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if s.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if e.Success != true {
		t.Errorf("%v", err)
	}

	gock.New(crmURL).
		Post(fmt.Sprintf("/orders/packs/%d/delete", p.ID)).
		MatchType("url").
		Reply(200).
		BodyString(`{"success": true}`)

	d, status, err := c.PackDelete(p.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if d.Success != true {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err == nil {
		t.Error("Error must be return")
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
	if err == nil {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if len(data.History) == 0 {

		t.Errorf("%v", err)
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

	data, _, err := c.PacksHistory(PacksHistoryRequest{Filter: OrdersHistoryFilter{StartDate: "2020-13-12"}})
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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

	data, _, err := c.Packs(PacksRequest{Filter: PacksFilter{ShipmentDateFrom: "2020-13-12"}})
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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

	data, _, err := c.Inventories(InventoriesRequest{Filter: InventoriesFilter{Sites: []string{codeFail}}})
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
	}
}

func TestClient_Settings(t *testing.T) {
	c := client()

	defer gock.Off()

	gock.New(crmURL).
		Get("/settings").
		Reply(200).
		BodyString(`{
  "success": true,
  "settings": {
    "default_currency": {
      "value": "RUB",
      "updated_at": "2019-02-13 13:57:20"
    },
    "system_language": {
      "value": "RU",
      "updated_at": "2019-02-13 14:02:23"
    },
    "timezone": {
      "value": "Europe/Moscow",
      "updated_at": "2019-02-13 13:57:20"
    }
  }
}
`)

	data, status, err := c.Settings()
	if err != nil {
		t.Errorf("%v", err)
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if data.Settings.DefaultCurrency.Value != "RUB" {
		t.Errorf("Invalid default_currency value: %v", data.Settings.DefaultCurrency.Value)
	}

	if data.Settings.SystemLanguage.Value != "RU" {
		t.Errorf("Invalid system_language value: %v", data.Settings.SystemLanguage.Value)
	}

	if data.Settings.Timezone.Value != "Europe/Moscow" {
		t.Errorf("Invalid timezone value: %v", data.Settings.Timezone.Value)
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

	data, _, err := c.Segments(SegmentsRequest{Filter: SegmentsFilter{Active: 3}})
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusCreated {
		t.Errorf("%v", err)
	}

	if m.Success != true {
		t.Errorf("%v", err)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/integration-modules/%s", code)).
		Reply(200).
		BodyString(`{"success": true}`)

	g, status, err := c.IntegrationModule(code)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if g.Success != true {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if g.Success != true {

		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if g.Success != true {

		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if g.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_DeliveryTracking(t *testing.T) {
	c := client()

	defer gock.Off()

	p := url.Values{
		"statusUpdate": {`[{"deliveryId":"123","history":[{"code":"1","updatedAt":"2020-01-01T00:00:00:000"}]}]`},
	}

	gock.New(crmURL).
		Post("/delivery/generic/subcode/tracking").
		MatchType("url").
		BodyString(p.Encode()).
		Reply(200).
		BodyString(`{"success": true}`)

	g, status, err := c.DeliveryTracking([]DeliveryTrackingRequest{{
		DeliveryID: "123",
		History: []DeliveryHistoryRecord{{
			Code:      "1",
			UpdatedAt: "2020-01-01T00:00:00:000",
		}},
	}}, "subcode")

	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if g.Success != true {
		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if g.Success != true {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Error("Error must be return")
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if costs.Success != true {
		t.Errorf("%v", err)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/costs/%d", id)).
		Reply(200).
		BodyString(`{"success": true}`)

	cost, status, err := c.Cost(id)

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if cost.Success != true {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if costEdit.Success != true {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if costDelete.Success != true {
		t.Errorf("%v", err)
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

	data, _, err := c.CostCreate(costRecord)
	if err == nil {
		t.Error("Error must be return")
	}

	if err == nil {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
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
	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if costsDelete.Success != true {

		t.Errorf("%v", err)
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
		BodyString(`{"success": false, "errorMsg": "Costs are loaded with ErrorsList"}`)

	data, _, err := c.CostsUpload(costsUpload)
	if err == nil {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Errorf("%v", err)
	}

	if status < http.StatusBadRequest {
		t.Error(statusFail)
	}

	if costsDelete.Success != false {
		t.Error(successFail)
	}
}

func TestClient_Files(t *testing.T) {
	c := client()
	fileID := 14925

	defer gock.Off()

	gock.New(crmURL).
		Get("/files").
		MatchParam("filter[ids][]", strconv.Itoa(fileID)).
		MatchParam("limit", "20").
		MatchParam("page", "1").
		Reply(200).
		BodyString(`{"success": true,"pagination": {"limit": 20,"totalCount": 0,"currentPage": 1,"totalPageCount": 0},"files": []}`)

	_, status, err := c.Files(FilesRequest{
		Limit: 20,
		Page:  1,
		Filter: FilesFilter{
			Ids: []int{fileID},
		},
	})

	if status != 200 {
		t.Errorf("%v %v", err.Error(), err)
	}
}

func TestClient_FileUpload(t *testing.T) {
	c := client()
	file := strings.NewReader(`test file contents`)

	defer gock.Off()

	gock.New(crmURL).
		Post("/files/upload").
		Reply(200).
		BodyString(`{"success": true, "file": {"id": 1}}`)

	data, status, err := c.FileUpload(file)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if data.File.ID != 1 {
		t.Error("invalid file id")
	}
}

func TestClient_FileUploadFail(t *testing.T) {
	c := client()
	file := strings.NewReader(`test file contents`)

	defer gock.Off()

	gock.New(crmURL).
		Post("/files/upload").
		Reply(400).
		BodyString(`{"success":false,"errorMsg":"Your account doesn't have enough money to upload files."}`)

	_, status, err := c.FileUpload(file)
	if err == nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusBadRequest {
		t.Errorf("status should be `%d`, got `%d` instead", http.StatusBadRequest, status)
	}
}

func TestClient_File(t *testing.T) {
	c := client()
	invalidFile := 20
	fileResponse := &FileResponse{
		Success: true,
		File: &File{
			ID:         19,
			Filename:   "image.jpg",
			Type:       "image/jpeg",
			CreatedAt:  time.Now().String(),
			Size:       10000,
			Attachment: nil,
		},
	}
	respData, errr := json.Marshal(fileResponse)
	if errr != nil {
		t.Errorf("%v", errr.Error())
	}

	defer gock.Off()

	gock.New(crmURL).
		Get(fmt.Sprintf("/files/%d", fileResponse.File.ID)).
		Reply(200).
		BodyString(string(respData))

	gock.New(crmURL).
		Get(fmt.Sprintf("/files/%d", invalidFile)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not Found"}`)

	s, status, err := c.File(fileResponse.File.ID)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if s.Success != true {
		t.Errorf("%v", err)
	}

	if s.File.ID != fileResponse.File.ID {
		t.Error("invalid response data")
	}

	s, status, err = c.File(invalidFile)
	if err == nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusNotFound {
		t.Errorf("status should be `%d`, got `%d` instead", http.StatusNotFound, status)
	}
}

func TestClient_FileDelete(t *testing.T) {
	c := client()
	successful := 19
	badRequest := 20
	notFound := 21

	defer gock.Off()

	gock.New(crmURL).
		Post(fmt.Sprintf("/files/%d/delete", successful)).
		Reply(200).
		BodyString(`{"success": true}`)

	gock.New(crmURL).
		Post(fmt.Sprintf("/files/%d/delete", badRequest)).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Error"}`)

	gock.New(crmURL).
		Post(fmt.Sprintf("/files/%d/delete", notFound)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not Found"}`)

	data, status, err := c.FileDelete(successful)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	data, _, err = c.FileDelete(badRequest)
	if err == nil {
		t.Errorf("%v", err)
	}

	data, _, err = c.FileDelete(notFound)
	if err == nil {
		t.Errorf("%v", err)
	}
}

func TestClient_FileDownload(t *testing.T) {
	c := client()
	successful := 19
	fail := 20
	fileData := "file data"

	defer gock.Off()

	gock.New(crmURL).
		Get(fmt.Sprintf("/files/%d/download", successful)).
		Reply(200).
		BodyString(fileData)

	gock.New(crmURL).
		Get(fmt.Sprintf("/files/%d/download", fail)).
		Reply(400).
		BodyString("")

	data, status, err := c.FileDownload(successful)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	fetchedByte, err := ioutil.ReadAll(data)
	if err != nil {
		t.Error(err)
	}

	fetched := string(fetchedByte)
	if fetched != fileData {
		t.Error("file data mismatch")
	}

	data, status, err = c.FileDownload(fail)
	if err == nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusBadRequest {
		t.Errorf("status should be `%d`, got `%d` instead", http.StatusBadRequest, status)
	}
}

func TestClient_FileEdit(t *testing.T) {
	c := client()
	successful := 19
	fail := 20
	resp := FileResponse{
		Success: true,
		File:    &File{Filename: "image.jpg"},
	}
	respData, _ := json.Marshal(resp)

	defer gock.Off()

	gock.New(crmURL).
		Post(fmt.Sprintf("/files/%d/edit", successful)).
		Reply(200).
		BodyString(string(respData))

	gock.New(crmURL).
		Post(fmt.Sprintf("/files/%d/edit", fail)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not Found"}`)

	data, status, err := c.FileEdit(successful, *resp.File)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
	}

	if data.File.Filename != resp.File.Filename {
		t.Errorf("filename should be `%s`, got `%s` instead", resp.File.Filename, data.File.Filename)
	}

	data, _, err = c.FileEdit(fail, *resp.File)
	if err == nil {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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

	data, _, err := c.CustomFields(CustomFieldsRequest{Filter: CustomFieldsFilter{Type: codeFail}})
	if err == nil {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if data.Success != true {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {

		t.Errorf("%v", err)
	}

	if cds.Success != true {
		t.Errorf("%v", err)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/custom-fields/dictionaries/%s", code)).
		Reply(200).
		BodyString(`{"success": true}`)

	cd, status, err := c.CustomDictionary(code)

	if err != nil {
		t.Errorf("%v", err)
	}

	if cd.Success != true {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if cde.Success != true {
		t.Errorf("%v", err)
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

	data, _, err := c.CustomDictionariesCreate(customDictionary)
	if err == nil {
		t.Errorf("%v", err)
	}

	if data.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/custom-fields/dictionaries/%s", codeFail)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	cd, status, err := c.CustomDictionary(codeFail)
	if err == nil {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Errorf("%v", err)
	}

	if status < http.StatusBadRequest {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if data.Success != true {

		t.Errorf("%v", err)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/custom-fields/%s/%s", "order", codeCustomField)).
		Reply(200).
		BodyString(`{"success": true}`)

	customField, status, err := c.CustomField("order", codeCustomField)

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if customField.Success != true {
		t.Errorf("%v", err)
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

	if err != nil {
		t.Errorf("%v", err)
	}

	if !statuses[status] {
		t.Errorf("%v", err)
	}

	if customFieldEdit.Success != true {
		t.Errorf("%v", err)
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

	data, _, err := c.CustomFieldsCreate(customFields)
	if err == nil {
		t.Errorf("%v", err)
	}

	if data.Success != false {
		t.Error(successFail)
	}

	gock.New(crmURL).
		Get(fmt.Sprintf("/custom-fields/%s/%s", "order", codeCustomField)).
		Reply(404).
		BodyString(`{"success": false, "errorMsg": "Not found"}`)

	customField, status, err := c.CustomField("order", codeCustomField)
	if err == nil {
		t.Errorf("%v", err)
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
	if err == nil {
		t.Errorf("%v", err)
	}

	if customFieldEdit.Success != false {
		t.Error(successFail)
	}
}

func TestClient_UpdateScopes(t *testing.T) {
	c := client()

	code := RandomString(8)

	defer gock.Off()

	request := UpdateScopesRequest{Requires: ScopesRequired{Scopes: []string{"scope1", "scope2"}}}

	jr, _ := json.Marshal(&request)

	gock.New(crmURL).
		Post(fmt.Sprintf("/integration-modules/%s/update-scopes", code)).
		BodyString(string(jr[:])).
		Reply(200).
		BodyString(`{"success": true, "apiKey": "newApiKey"}`)

	m, status, err := c.UpdateScopes(code, request)
	if err != nil {
		t.Errorf("%v", err)
	}

	if status != http.StatusOK {
		t.Errorf("%v", err)
	}

	if m.Success != true {
		t.Errorf("%v", err)
	}
}

func TestClient_UpdateScopes_Fail(t *testing.T) {
	c := client()

	code := RandomString(8)

	defer gock.Off()

	request := UpdateScopesRequest{Requires: ScopesRequired{Scopes: []string{"scope1", "scope2"}}}

	jr, _ := json.Marshal(&request)

	gock.New(crmURL).
		Post(fmt.Sprintf("/integration-modules/%s/update-scopes", code)).
		BodyString(string(jr[:])).
		Reply(400).
		BodyString(`{"success": false, "errorMsg": "Not enabled simple connection"}`)

	m, status, err := c.UpdateScopes(code, request)
	if err == nil {
		t.Error("Error must be return")
	}

	if status != http.StatusBadRequest {
		t.Errorf("%v", err)
	}

	if m.Success != false {
		t.Error(successFail)
	}
}

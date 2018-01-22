package v5

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

type Configuration struct {
	Url string `json:"url"`
	Key string `json:"key"`
	Ver string `json:"version"`
}

func buildConfiguration() *Configuration {
	file, _ := os.Open("../config.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return &Configuration{
		configuration.Url,
		configuration.Key,
		configuration.Ver,
	}
}

var r *rand.Rand // Rand for this package.

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)

	for i := range result {
		result[i] = chars[r.Intn(len(chars))]
	}

	return string(result)
}

func client() *Client {
	configuration := buildConfiguration()
	return New(configuration.Url, configuration.Key)
}

func TestGetRequest(t *testing.T) {
	c := client()
	_, status, _ := c.getRequest("/fake-method")

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestPostRequest(t *testing.T) {
	c := client()
	_, status, _ := c.postRequest("/fake-method", url.Values{})

	if status != http.StatusNotFound {
		t.Fail()
	}
}

func TestClient_ApiVersionsVersions(t *testing.T) {
	c := client()
	data, status, err := c.ApiVersions()
	if err != nil {
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Fail()
	}

	if data.Success != true {
		t.Fail()
	}
}

func TestClient_ApiCredentialsCredentials(t *testing.T) {
	c := client()
	data, status, err := c.ApiCredentials()
	if err != nil {
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Fail()
	}

	if data.Success != true {
		t.Fail()
	}
}

func TestClient_CustomersCustomers(t *testing.T) {
	c := client()
	f := CustomersFilter{}
	f.City = "Москва"

	data, status, err := c.Customers(f, 20, 1)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestClient_CustomerChange(t *testing.T) {
	c := client()
	f := Customer{}

	random := RandomString(8)

	f.FirstName = "Понтелей"
	f.LastName = "Турбин"
	f.Patronymic = "Аристархович"
	f.ExternalId = random
	f.Email = fmt.Sprintf("%s@example.com", random)

	cr, sc, err := c.CustomerCreate(f)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if sc != http.StatusCreated {
		t.Errorf("%s", err)
		t.Fail()
	}

	if cr.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	f.Id = cr.Id
	f.Vip = true

	ed, se, err := c.CustomerEdit(f, "id")
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if se != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if ed.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	data, status, err := c.Customer(f.ExternalId, "externalId", "")
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status >= http.StatusBadRequest {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Customer.ExternalId != f.ExternalId {
		t.Errorf("%s", err)
		t.Fail()
	}
}

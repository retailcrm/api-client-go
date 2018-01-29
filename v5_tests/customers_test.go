package v5_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/retailcrm/api-client-go/v5"
)

func TestClient_CustomersCustomers(t *testing.T) {
	c := client()
	f := v5.CustomersRequest{
		Filter: v5.CustomersFilter{
			City: "Москва",
		},
		Page: 3,
	}

	data, status, err := c.Customers(f)
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

	random := RandomString(8)

	f := v5.Customer{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalId: random,
		Email:      fmt.Sprintf("%s@example.com", random),
	}

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

func TestClient_CustomersUpload(t *testing.T) {
	c := client()
	customers := make([]v5.Customer, 3)

	for i := range customers {
		customers[i] = v5.Customer{
			FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
			LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
			ExternalId: RandomString(8),
			Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		}
	}

	data, status, err := c.CustomersUpload(customers)
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

func TestClient_CustomersFixExternalIds(t *testing.T) {
	c := client()
	f := v5.Customer{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalId: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	cr, sc, err := c.CustomerCreate(f)
	if err != nil {
		t.Errorf("%s", sc)
		t.Fail()
	}

	if sc != http.StatusCreated {
		t.Errorf("%s", sc)
		t.Fail()
	}

	if cr.Success != true {
		t.Errorf("%s", sc)
		t.Fail()
	}

	customers := []v5.IdentifiersPair{{
		Id:         cr.Id,
		ExternalId: RandomString(8),
	}}

	fx, fe, err := c.CustomersFixExternalIds(customers)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if fe != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if fx.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestClient_CustomersHistory(t *testing.T) {
	c := client()
	f := v5.CustomersHistoryRequest{
		Filter: v5.CustomersHistoryFilter{
			SinceId: 20,
		},
	}

	data, status, err := c.CustomersHistory(f)
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

	if len(data.History) == 0 {
		t.Errorf("%s", err)
		t.Fail()
	}
}

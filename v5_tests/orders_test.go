package v5_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/retailcrm/api-client-go/v5"
)

func TestClient_OrdersOrders(t *testing.T) {
	c := client()
	f := v5.OrdersRequest{
		Filter: v5.OrdersFilter{
			City: "Москва",
		},
		Page: 1,
	}

	data, status, err := c.Orders(f)
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

func TestClient_OrderChange(t *testing.T) {
	c := client()

	random := RandomString(8)

	f := v5.Order{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalId: random,
		Email:      fmt.Sprintf("%s@example.com", random),
	}

	cr, sc, err := c.OrderCreate(f)
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
	f.CustomerComment = "test comment"

	ed, se, err := c.OrderEdit(f, "id")
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

	data, status, err := c.Order(f.ExternalId, "externalId", "")
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

	if data.Order.ExternalId != f.ExternalId {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestClient_OrdersUpload(t *testing.T) {
	c := client()
	orders := make([]v5.Order, 3)

	for i := range orders {
		orders[i] = v5.Order{
			FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
			LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
			ExternalId: RandomString(8),
			Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
		}
	}

	data, status, err := c.OrdersUpload(orders)
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

func TestClient_OrdersFixExternalIds(t *testing.T) {
	c := client()
	f := v5.Order{
		FirstName:  fmt.Sprintf("Name_%s", RandomString(8)),
		LastName:   fmt.Sprintf("Test_%s", RandomString(8)),
		ExternalId: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	cr, sc, err := c.OrderCreate(f)
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

	orders := []v5.IdentifiersPair{{
		Id:         cr.Id,
		ExternalId: RandomString(8),
	}}

	fx, fe, err := c.OrdersFixExternalIds(orders)
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

func TestClient_OrdersHistory(t *testing.T) {
	c := client()
	f := v5.OrdersHistoryRequest{
		Filter: v5.OrdersHistoryFilter{
			SinceId: 20,
		},
	}

	data, status, err := c.OrdersHistory(f)
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

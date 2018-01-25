package v5_tests

import (
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

func TestClient_OrdersOrder(t *testing.T) {
	c := client()

	data, status, err := c.Order("upload-b-1480333204", "externalId", "")
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

func TestClient_OrdersHistory(t *testing.T) {
	c := client()
	f := v5.OrdersHistoryRequest{
		Filter: v5.OrdersHistoryFilter{
			SinceId: 100,
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

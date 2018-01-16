package retailcrm

import (
	"testing"
	"net/http"
)

func TestCustomer(t *testing.T) {
	c := VersionedClient()
	data, status, err := c.Customer("163", "id", "")

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

	if data.Customer.Id != 163 {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestCustomers(t *testing.T) {
	c := VersionedClient()
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

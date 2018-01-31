package v5_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/retailcrm/api-client-go/v5"
)

func TestClient_PaymentCreateEditDelete(t *testing.T) {
	c := client()

	order := v5.Order{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalId: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	createOrderResponse, status, err := c.OrderCreate(order)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status != http.StatusCreated {
		t.Errorf("%s", err)
		t.Fail()
	}

	if createOrderResponse.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	f := v5.Payment{
		Order: &v5.Order{
			Id: createOrderResponse.Id,
		},
		Amount: 300,
		Type:   "cash",
	}

	paymentCreateResponse, status, err := c.PaymentCreate(f)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status != http.StatusCreated {
		t.Errorf("%s", err)
		t.Fail()
	}

	if paymentCreateResponse.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	k := v5.Payment{
		Id:     paymentCreateResponse.Id,
		Amount: 500,
	}

	paymentEditResponse, status, err := c.PaymentEdit(k, "id")
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if paymentEditResponse.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	paymentDeleteResponse, status, err := c.PaymentDelete(paymentCreateResponse.Id)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if paymentDeleteResponse.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}

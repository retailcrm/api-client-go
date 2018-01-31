package v5_tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/retailcrm/api-client-go/v5"
)

func TestClient_NotesNotes(t *testing.T) {
	c := client()
	f := v5.NotesRequest{
		Page: 1,
	}

	data, status, err := c.Notes(f)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if status != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if data.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}

func TestClient_NotesCreateDelete(t *testing.T) {
	c := client()

	customer := v5.Customer{
		FirstName:  "Понтелей",
		LastName:   "Турбин",
		Patronymic: "Аристархович",
		ExternalId: RandomString(8),
		Email:      fmt.Sprintf("%s@example.com", RandomString(8)),
	}

	createCustomerResponse, createCustomerStatus, err := c.CustomerCreate(customer)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if createCustomerStatus != http.StatusCreated {
		t.Errorf("%s", err)
		t.Fail()
	}

	if createCustomerResponse.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	f := v5.Note{
		Text:      "some text",
		ManagerId: GetUser(),
		Customer: v5.Customer{
			Id: createCustomerResponse.Id,
		},
	}

	noteCreateResponse, noteCreateStatus, err := c.NoteCreate(f)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if noteCreateStatus != http.StatusCreated {
		t.Errorf("%s", err)
		t.Fail()
	}

	if noteCreateResponse.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}

	noteDeleteResponse, noteDeleteStatus, err := c.NoteDelete(noteCreateResponse.Id)
	if err != nil {
		t.Errorf("%s", err)
		t.Fail()
	}

	if noteDeleteStatus != http.StatusOK {
		t.Errorf("%s", err)
		t.Fail()
	}

	if noteDeleteResponse.Success != true {
		t.Errorf("%s", err)
		t.Fail()
	}
}

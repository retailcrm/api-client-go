package v5

import (
	"errors"
	"testing"
)

func TestFailure_ApiErrorsSlice(t *testing.T) {
	b := []byte(`{"success": false, "errorMsg": "Failed to activate module", "errors": ["Your account has insufficient funds to activate integration module"]}`)
	expected := map[string]string{
		"0": "Your account has insufficient funds to activate integration module",
	}

	var expEr *APIError
	e := NewAPIError(b)

	if errors.As(e, &expEr) {
		if eq := expEr.Errors["0"] == expected["0"]; eq != true {
			t.Errorf("%+v", eq)
		}
	} else {
		t.Errorf("Error must be type of APIError: %v", e)
	}
}

func TestFailure_ApiErrorsMap(t *testing.T) {
	b := []byte(`{"success": false, "errorMsg": "Failed to activate module", "errors": {"id": "ID must be an integer"}}`)
	expected := map[string]string{
		"id": "ID must be an integer",
	}

	var expEr *APIError
	e := NewAPIError(b)

	if errors.As(e, &expEr) {
		if eq := expected["id"] == expEr.Errors["id"]; eq != true {
			t.Errorf("%+v", eq)
		}
	} else {
		t.Errorf("Error must be type of APIError: %v", e)
	}
}

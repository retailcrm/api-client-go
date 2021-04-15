package v5

import (
	"testing"
)

func TestFailure_ApiErrorsSlice(t *testing.T) {
	b := []byte(`{"success": false, "errorMsg": "Failed to activate module", "errors": ["Your account has insufficient funds to activate integration module"]}`)
	expected := map[string]string{
		"0": "Your account has insufficient funds to activate integration module",
	}

	e := NewApiError(b)
	if eq := e.(*ApiError).Errors["0"] == expected["0"]; eq != true {
		t.Errorf("%+v", eq)
	}
}

func TestFailure_ApiErrorsMap(t *testing.T) {
	b := []byte(`{"success": false, "errorMsg": "Failed to activate module", "errors": {"id": "ID must be an integer"}}`)
	expected := map[string]string{
		"id": "ID must be an integer",
	}

	e := NewApiError(b)
	if eq := expected["id"] == e.(*ApiError).Errors["id"]; eq != true {
		t.Errorf("%+v", eq)
	}
}

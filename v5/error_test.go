package v5

import (
	"reflect"
	"testing"

	"golang.org/x/xerrors"
)

func TestFailure_ApiErrorsSlice(t *testing.T) {
	b := []byte(`{"success": false,
				"errorMsg": "Failed to activate module",
				"errors": [
					"Your account has insufficient funds to activate integration module",
					"Test error"
				]}`)
	expected := APIErrorsList{
		"0": "Your account has insufficient funds to activate integration module",
		"1": "Test error",
	}

	var expEr *APIError
	e := NewAPIError(b)

	if xerrors.As(e, &expEr) {
		if eq := reflect.DeepEqual(expEr.Errors, expected); eq != true {
			t.Errorf("%+v", eq)
		}
	} else {
		t.Errorf("Error must be type of APIError: %v", e)
	}
}

func TestFailure_ApiErrorsMap(t *testing.T) {
	b := []byte(`{"success": false,
				"errorMsg": "Failed to activate module",
				"errors": {"id": "ID must be an integer", "test": "Test error"}}`,
	)
	expected := APIErrorsList{
		"id":   "ID must be an integer",
		"test": "Test error",
	}

	var expEr *APIError
	e := NewAPIError(b)

	if xerrors.As(e, &expEr) {
		if eq := reflect.DeepEqual(expEr.Errors, expected); eq != true {
			t.Errorf("%+v", eq)
		}
	} else {
		t.Errorf("Error must be type of APIError: %v", e)
	}
}

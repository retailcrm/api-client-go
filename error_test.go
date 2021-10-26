package retailcrm

import (
	"errors"
	"reflect"
	"testing"
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

	e := CreateAPIError(b)

	if errors.Is(e, ErrGeneric) {
		if eq := reflect.DeepEqual(expected, e.(APIError).Errors()); eq != true {
			t.Errorf("%+v", eq)
		}
	} else {
		t.Errorf("Error must be type of ErrGeneric: %v", e)
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

	e := CreateAPIError(b)
	if errors.Is(e, ErrGeneric) {
		if eq := reflect.DeepEqual(expected, e.(APIError).Errors()); eq != true {
			t.Errorf("%+v", eq)
		}
	} else {
		t.Errorf("Error must be type of ErrGeneric: %v", e)
	}
}

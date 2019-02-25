package errs

import (
	"reflect"
	"testing"
)

func TestFailure_ApiErrorsSlice(t *testing.T) {
	var err = Failure{}
	b := []byte(`{"success": false, "errorMsg": "Failed to activate module", "errors": ["Your account has insufficient funds to activate integration module"]}`)
	expected := map[string]string{
		"0": "Your account has insufficient funds to activate integration module",
	}

	resp, e := ErrorResponse(b)
	err.SetRuntimeError(e)
	err.SetApiError(resp.ErrorMsg)

	if resp.Errors != nil {
		err.SetApiErrors(ErrorsHandler(resp.Errors))
	}

	eq := reflect.DeepEqual(expected, err.ApiErrors())

	if eq != true {
		t.Errorf("%+v", eq)
	}
}

func TestFailure_ApiErrorsMap(t *testing.T) {
	var err = Failure{}
	b := []byte(`{"success": false, "errorMsg": "Failed to activate module", "errors": {"id": "ID must be an integer"}}`)
	expected := map[string]string{
		"id": "ID must be an integer",
	}

	resp, e := ErrorResponse(b)
	err.SetRuntimeError(e)
	err.SetApiError(resp.ErrorMsg)

	if resp.Errors != nil {
		err.SetApiErrors(ErrorsHandler(resp.Errors))
	}

	eq := reflect.DeepEqual(expected, err.ApiErrors())

	if eq != true {
		t.Errorf("%+v", eq)
	}
}

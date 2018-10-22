package errs

import (
	"testing"
)

func TestFailure_ApiErrorsSlice(t *testing.T) {
	var err = Failure{}
	b := []byte(`{"success": false, "errorMsg": "Failed to activate module", "errors": ["Your account has insufficient funds to activate integration module"]}`)

	resp, e := ErrorResponse(b)
	err.RuntimeErr = e
	err.ApiErr = resp.ErrorMsg

	if resp.Errors != nil {
		err.ApiErrs = resp.Errors
	}

	f, ok := resp.Errors.([]interface{})

	if !ok {
		t.Errorf("%+v", f)
	}

}

func TestFailure_ApiErrorsMap(t *testing.T) {
	var err = Failure{}
	b := []byte(`{"success": false, "errorMsg": "Failed to activate module", "errors": {"id": "ID must be an integer"}}`)

	resp, e := ErrorResponse(b)
	err.RuntimeErr = e
	err.ApiErr = resp.ErrorMsg

	if resp.Errors != nil {
		err.ApiErrs = resp.Errors
	}

	f, ok := resp.Errors.(map[string]interface{})

	if !ok {
		t.Errorf("%+v", f)
	}

}

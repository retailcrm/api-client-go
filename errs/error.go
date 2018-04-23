package errs

import (
	"encoding/json"
	"fmt"
)

// Error returns the string representation of the error and satisfies the error interface.
func (f *Failure) Error() string {
	return f.RuntimeErr.Error()
}

// ApiError returns formatted string representation of the API error
func (f *Failure) ApiError() string {
	return fmt.Sprintf("%v", f.ApiErr)
}

// ApiErrors returns array of formatted strings that represents API errors
func (f *Failure) ApiErrors() map[string]string {
	return f.ApiErrs
}

// ErrorResponse method
func ErrorResponse(data []byte) (FailureResponse, error) {
	var resp FailureResponse
	err := json.Unmarshal(data, &resp)

	return resp, err
}

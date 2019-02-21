package errs

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Error returns the string representation of the error and satisfies the error interface.
func (f *Failure) Error() string {
	if f != nil && f.runtimeErr != nil {
		return f.runtimeErr.Error()
	}

	return ""
}

// ApiError returns formatted string representation of the API error
func (f *Failure) ApiError() string {
	if f != nil && f.apiErr != "" {
		return fmt.Sprintf("%+v", f.apiErr)
	}

	return ""
}

// ApiErrors returns array of formatted strings that represents API errors
func (f *Failure) ApiErrors() map[string]string {
	if len(f.apiErrs) > 0 {
		return f.apiErrs
	}

	return nil
}

// SetRuntimeError set runtime error value
func (f *Failure) SetRuntimeError(e error) {
	f.runtimeErr = e
}

// SetApiError set api error value
func (f *Failure) SetApiError(e string) {
	f.apiErr = e
}

// SetApiErrors set api errors value
func (f *Failure) SetApiErrors(e map[string]string) {
	f.apiErrs = e
}

// ErrorsHandler returns map
func ErrorsHandler(errs interface{}) map[string]string {
	m := make(map[string]string)

	switch errs.(type) {
	case map[string]interface{}:
		for idx, val := range errs.(map[string]interface{}) {
			m[idx] = val.(string)
		}
	case []interface{}:
		for idx, val := range errs.([]interface{}) {
			m[strconv.Itoa(idx)] = val.(string)
		}
	}

	return m
}

// ErrorResponse method
func ErrorResponse(data []byte) (FailureResponse, error) {
	var resp FailureResponse
	err := json.Unmarshal(data, &resp)

	return resp, err
}

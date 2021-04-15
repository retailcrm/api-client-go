package v5

import (
	"encoding/json"
	"strconv"
)

// ApiErrorsList struct
type ApiErrorsList map[string]string

// ApiError struct
type ApiError struct {
	SuccessfulResponse
	ErrorMsg string `json:"errorMsg,omitempty"`
	Errors ApiErrorsList `json:"errors,omitempty"`
}

func (e *ApiError) Error() string {
	return e.ErrorMsg
}

func (e *ApiError) GetApiErrors() map[string]string {
	return e.Errors
}

func NewApiError (dataResponse []byte) error {
	a := &ApiError{}

	if err := json.Unmarshal(dataResponse, a); err != nil {
		return err
	}

	return a
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
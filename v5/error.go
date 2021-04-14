package v5

import (
	"strconv"
)

// ApiErrorsList struct
type ApiErrorsList struct {
	Success bool `json:"success"`
	ErrorsMsg string `json:"errorMsg,omitempty"`
	Errors	interface{} `json:"errors,omitempty"`
}

// ApiError struct
type ApiError struct {
	SuccessfulResponse
	ErrorMsg string
	Errors map[string]string
}

func (e *ApiError) Error() string {
	return e.ErrorMsg
}

func NewApiError (dataResponse []byte) error {
	a := &ApiError{}

	if err := a.UnmarshalJSON(dataResponse); err != nil {
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
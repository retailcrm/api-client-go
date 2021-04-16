package v5

import (
	"encoding/json"
	"strconv"
)

const ArrowHTML = 60

// APIErrorsList struct.
type APIErrorsList map[string]string

// APIError struct.
type APIError struct {
	SuccessfulResponse
	ErrorMsg string        `json:"errorMsg,omitempty"`
	Errors   APIErrorsList `json:"errors,omitempty"`
}

func (e *APIError) Error() string {
	return e.ErrorMsg
}

func NewAPIError(dataResponse []byte) error {
	a := &APIError{}

	if dataResponse[0] == ArrowHTML {
		a.ErrorMsg = "405 Not Allowed"
		return a
	}

	if err := json.Unmarshal(dataResponse, a); err != nil {
		return err
	}

	return a
}

// ErrorsHandler returns map.
func ErrorsHandler(errs interface{}) map[string]string {
	m := make(map[string]string)

	switch e := errs.(type) {
	case map[string]interface{}:
		for idx, val := range e {
			m[idx] = val.(string)
		}
	case []interface{}:
		for idx, val := range e {
			m[strconv.Itoa(idx)] = val.(string)
		}
	}

	return m
}

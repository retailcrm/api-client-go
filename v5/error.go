package v5

import "encoding/json"

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

	if len(dataResponse) > 0 && dataResponse[0] == '<' {
		a.ErrorMsg = "Account does not exist."
		return a
	}

	if err := json.Unmarshal(dataResponse, a); err != nil {
		return err
	}

	return a
}

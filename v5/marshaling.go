package v5

import (
	"encoding/json"
)

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Name)
}

func (a *ApiError) UnmarshalJSON(data []byte) error {
	var e ApiErrorsList

	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	a.SuccessfulResponse = SuccessfulResponse{e.Success}
	a.ErrorMsg = e.ErrorsMsg

	if e.Errors != nil {
		a.Errors = ErrorsHandler(e.Errors)
	}

	return nil
}
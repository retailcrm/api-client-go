package v5

import (
	"encoding/json"
)

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Name)
}

func (a *APIErrorsList) UnmarshalJSON(data []byte) error {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	*a = ErrorsHandler(i)
	return nil
}

package v5

import "encoding/json"

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Name)
}

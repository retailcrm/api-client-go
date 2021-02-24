package v5

import "encoding/json"

func (t Tag) MarshalJSON() ([]byte, error) {
	name, err := json.Marshal(t.Name)

	return name, err
}

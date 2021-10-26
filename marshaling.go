package retailcrm

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Name)
}

func (a *APIErrorsList) UnmarshalJSON(data []byte) error {
	var i interface{}
	var m map[string]string
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		m = make(map[string]string, len(e))
		for idx, val := range e {
			m[idx] = fmt.Sprint(val)
		}
	case []interface{}:
		m = make(map[string]string, len(e))
		for idx, val := range e {
			m[strconv.Itoa(idx)] = fmt.Sprint(val)
		}
	}

	*a = m
	return nil
}

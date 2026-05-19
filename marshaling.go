package retailcrm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Name)
}

func (v *StringOrNumber) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*v = ""
		return nil
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return err
	}

	switch value := value.(type) {
	case string:
		*v = StringOrNumber(value)
	case json.Number:
		*v = StringOrNumber(value.String())
	default:
		return fmt.Errorf("string or number: expected string or number, got %T", value)
	}

	return nil
}

func (a *APIErrorsList) UnmarshalJSON(data []byte) error {
	var i interface{}
	var m APIErrorsList
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		m = make(APIErrorsList, len(e))
		for idx, val := range e {
			m[idx] = fmt.Sprint(val)
		}
	case []interface{}:
		m = make(APIErrorsList, len(e))
		for idx, val := range e {
			m[strconv.Itoa(idx)] = fmt.Sprint(val)
		}
	}

	*a = m
	return nil
}

func (l *StringMap) UnmarshalJSON(data []byte) error {
	var i interface{}
	var m StringMap
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		m = make(StringMap, len(e))
		for idx, val := range e {
			m[idx] = fmt.Sprint(val)
		}
	case []interface{}:
		m = make(StringMap, len(e))
		for idx, val := range e {
			m[strconv.Itoa(idx)] = fmt.Sprint(val)
		}
	}

	*l = m
	return nil
}

func (l *CustomFieldMap) UnmarshalJSON(data []byte) error {
	var i interface{}
	var items CustomFieldMap
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		items = make(CustomFieldMap, len(e))
		for idx, val := range e {
			items[idx] = val
		}
	case []interface{}:
		items = make(CustomFieldMap, len(e))
		for idx, val := range e {
			items[strconv.Itoa(idx)] = val
		}
	}

	*l = items
	return nil
}

func (p *OrderPayments) UnmarshalJSON(data []byte) error {
	var i interface{}
	var m OrderPayments
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		m = make(OrderPayments, len(e))
		for idx, val := range e {
			var res OrderPayment
			err := unmarshalMap(val.(map[string]interface{}), &res)
			if err != nil {
				return err
			}
			m[idx] = res
		}
	case []interface{}:
		m = make(OrderPayments, len(e))
		for idx, val := range e {
			var res OrderPayment
			err := unmarshalMap(val.(map[string]interface{}), &res)
			if err != nil {
				return err
			}
			m[strconv.Itoa(idx)] = res
		}
	}

	*p = m
	return nil
}

func (p *Properties) UnmarshalJSON(data []byte) error {
	var i interface{}
	var m Properties
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}

	switch e := i.(type) {
	case map[string]interface{}:
		m = make(Properties, len(e))
		for idx, val := range e {
			var res Property
			err := unmarshalMap(val.(map[string]interface{}), &res)
			if err != nil {
				return err
			}
			m[idx] = res
		}
	case []interface{}:
		m = make(Properties, len(e))
		for idx, val := range e {
			var res Property
			err := unmarshalMap(val.(map[string]interface{}), &res)
			if err != nil {
				return err
			}
			m[strconv.Itoa(idx)] = res
		}
	}

	*p = m
	return nil
}

func unmarshalMap(m map[string]interface{}, v interface{}) (err error) {
	var data []byte
	data, err = json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

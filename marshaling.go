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

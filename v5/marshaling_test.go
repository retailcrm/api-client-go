package v5

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTag_MarshalJSON(t *testing.T) {
	tags := []Tag{
		{"first", "#3e89b6", false},
		{"second", "#ffa654", false},
	}
	names := []byte(`["first","second"]`)
	str, err := json.Marshal(tags)

	if err != nil {
		t.Errorf("%v", err.Error())
	}

	if !reflect.DeepEqual(str, names) {
		t.Errorf("Marshaled: %#v\nExpected: %#v\n", str, names)
	}
}

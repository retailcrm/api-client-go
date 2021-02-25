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
	str, _ := json.Marshal(tags)

	eq := reflect.DeepEqual(str, names)
	if eq != true {
		t.Errorf("Marshaled: %#v\nExpected: %#v\n", str, names)
	}
}

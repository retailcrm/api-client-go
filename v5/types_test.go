package v5

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestClient_OrderDeliveryData(t *testing.T) {
	d := OrderDeliveryData{
		OrderDeliveryDataBasic: OrderDeliveryDataBasic{
			"track",
			"status",
			"address",
			"type",
		},
	}

	data, _ := json.Marshal(d)
	expectedStr := `{"payerType":"type","pickuppointAddress":"address","status":"status","trackNumber":"track"}`
	if string(data) != expectedStr {
		t.Error("error")
	}

	d.AdditionalFields = map[string]interface{}{
		"customFirst":  "one",
		"customSecond": "two",
	}

	data, _ = json.Marshal(d)
	expectedStr = `{"customFirst":"one","customSecond":"two","payerType":"type","pickuppointAddress":"address","status":"status","trackNumber":"track"}`
	if string(data) != expectedStr {
		t.Error("error")
	}

	d = OrderDeliveryData{}
	json.Unmarshal(data, &d)
	expected := OrderDeliveryData{
		OrderDeliveryDataBasic: OrderDeliveryDataBasic{
			"track",
			"status",
			"address",
			"type",
		},
		AdditionalFields: map[string]interface{}{
			"customFirst":  "one",
			"customSecond": "two",
		},
	}

	eq := reflect.DeepEqual(expected, d)
	if eq != true {
		t.Errorf("Unmarshaled: %#v\nExpected: %#v\n", d, expected)
	}
}

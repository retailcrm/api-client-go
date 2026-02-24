package retailcrm

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
		t.Errorf("Marshaled: %s\nExpected: %s\n", data, expectedStr)
	}

	d.AdditionalFields = map[string]interface{}{
		"customFirst":  "one",
		"customSecond": "two",
	}

	data, _ = json.Marshal(d)
	expectedStr = `{"customFirst":"one","customSecond":"two","payerType":"type","pickuppointAddress":"address","status":"status","trackNumber":"track"}`
	if string(data) != expectedStr {
		t.Errorf("Marshaled: %s\nExpected: %s\n", data, expectedStr)
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

func TestCustomer_IsContactJSON(t *testing.T) {
	customer := Customer{ID: 1, IsContact: true}

	data, err := json.Marshal(customer)
	if err != nil {
		t.Fatalf("marshal customer: %v", err)
	}

	var marshaled map[string]interface{}
	if err := json.Unmarshal(data, &marshaled); err != nil {
		t.Fatalf("unmarshal marshaled payload: %v", err)
	}

	if value, ok := marshaled["isContact"]; !ok || value != true {
		t.Fatalf("expected isContact=true in marshaled payload, got %#v", marshaled["isContact"])
	}

	var decoded Customer
	if err := json.Unmarshal([]byte(`{"id":2,"isContact":true}`), &decoded); err != nil {
		t.Fatalf("unmarshal customer with isContact: %v", err)
	}

	if !decoded.IsContact {
		t.Fatalf("expected IsContact=true after unmarshal")
	}
}

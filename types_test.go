package retailcrm

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestClient_OrderDeliveryData(t *testing.T) {
	d := OrderDeliveryData{
		OrderDeliveryDataBasic: OrderDeliveryDataBasic{
			TrackNumber:        "track",
			Status:             "status",
			PickuppointAddress: "address",
			PayerType:          "type",
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
			TrackNumber:        "track",
			Status:             "status",
			PickuppointAddress: "address",
			PayerType:          "type",
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

func TestCustomer_MGCustomersExternalIDJSON(t *testing.T) {
	var fromNumber Customer
	if err := json.Unmarshal([]byte(`{"mgCustomers":[{"id":1,"externalId":123}]}`), &fromNumber); err != nil {
		t.Fatalf("unmarshal customer with numeric mgCustomers externalId: %v", err)
	}

	if len(fromNumber.MgCustomers) != 1 {
		t.Fatalf("expected one mgCustomer from numeric payload, got %d", len(fromNumber.MgCustomers))
	}

	if fromNumber.MgCustomers[0].ExternalID != 123 {
		t.Fatalf("expected mgCustomers[0].ExternalID=123, got %d", fromNumber.MgCustomers[0].ExternalID)
	}

	var fromString Customer
	if err := json.Unmarshal([]byte(`{"mgCustomers":[{"id":1,"externalId":"123"}]}`), &fromString); err == nil {
		t.Fatal("expected string mgCustomers externalId to be rejected")
	}
}

func TestAPIMethodDTOFieldsJSON(t *testing.T) {
	order := Order{
		Company:                &Company{ID: 10, ExternalID: "company-ext"},
		LoyaltyEventDiscountID: 55,
		IsFromCart:             true,
		Items: []OrderItem{{
			ExternalID:     "item-ext",
			ExternalIDs:    []CodeValueModel{{Code: "marketplace", Value: "item-1"}},
			MarkingCodes:   []string{"mark-1"},
			MarkingObjects: []MarkingObject{{Code: "mark-2", Provider: "chestny_znak"}},
			Ordering:       2,
		}},
		Delivery: &OrderDelivery{
			Service: &OrderDeliveryService{DeliveryType: "courier"},
			Data: &OrderDeliveryData{OrderDeliveryDataBasic: OrderDeliveryDataBasic{
				ExternalID: "delivery-ext",
				Cost:       100,
				ExtraData:  StringMap{"terminal": "A1"},
				ItemDeclaredValues: []DeliveryItemDeclaredValue{{
					OrderProduct: OrderProductIdentifier{ExternalID: "item-ext"},
					Value:        500,
				}},
				Packages: []DeliveryPackage{{
					PackageID: "box-1",
					Items: []DeliveryPackageItem{{
						OrderProduct: OrderProductIdentifier{ExternalID: "item-ext"},
						Quantity:     1,
					}},
				}},
			}},
		},
	}

	data, err := json.Marshal(order)
	if err != nil {
		t.Fatalf("marshal order: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("unmarshal order payload: %v", err)
	}

	if _, ok := payload["company"]; !ok {
		t.Fatalf("expected company in marshaled order: %s", data)
	}
	if payload["loyaltyEventDiscountId"].(float64) != 55 {
		t.Fatalf("expected loyaltyEventDiscountId=55, got %#v", payload["loyaltyEventDiscountId"])
	}
	if payload["isFromCart"] != true {
		t.Fatalf("expected isFromCart=true, got %#v", payload["isFromCart"])
	}

	var decodedPack Pack
	if err := json.Unmarshal([]byte(`{"item":{"id":1,"externalId":"item-ext","externalIds":[{"code":"marketplace","value":"item-1"}]}}`), &decodedPack); err != nil {
		t.Fatalf("unmarshal pack: %v", err)
	}
	if decodedPack.Item == nil || decodedPack.Item.ExternalIDs[0].Code != "marketplace" {
		t.Fatalf("expected pack item externalIds, got %#v", decodedPack.Item)
	}
}

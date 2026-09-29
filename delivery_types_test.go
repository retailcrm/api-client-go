package retailcrm

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Keep the original method signature compatible with existing client interfaces.
var _ interface {
	DeliveryTypes() (DeliveryTypesResponse, int, error)
} = (*Client)(nil)

func TestDeliveryTypesQuery(t *testing.T) {
	cases := []struct {
		name   string
		params *DeliveryTypesRequest
		query  url.Values
	}{
		{name: "legacy", query: url.Values{}},
		{name: "empty", params: &DeliveryTypesRequest{}, query: url.Values{}},
		{
			name: "filter",
			params: &DeliveryTypesRequest{
				Filter: DeliveryTypesFilter{Codes: []string{"courier", "pickup-point"}},
			},
			query: url.Values{"filter[codes][]": {"courier", "pickup-point"}},
		},
		{
			name: "cost settings",
			params: &DeliveryTypesRequest{
				Filter:              DeliveryTypesFilter{Codes: []string{"courier"}},
				IncludeCostSettings: true,
			},
			query: url.Values{"filter[codes][]": {"courier"}, "includeCostSettings": {"true"}},
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				assert.Equal(t, http.MethodGet, request.Method)
				assert.Equal(t, "/api/v5/reference/delivery-types", request.URL.Path)
				assert.Equal(t, testCase.query, request.URL.Query())
				assert.Equal(t, "test-key", request.Header.Get("X-API-KEY"))
				if len(testCase.query) == 0 {
					assert.False(t, request.URL.ForceQuery)
				}
				writer.Header().Set("Content-Type", "application/json")
				_, err := io.WriteString(writer, `{"success":true,"deliveryTypes":{}}`)
				assert.NoError(t, err)
			}))
			t.Cleanup(server.Close)
			client := New(server.URL, "test-key").WithHTTPClient(server.Client())
			var response DeliveryTypesResponse
			var status int
			var err error
			if testCase.params == nil {
				response, status, err = client.DeliveryTypes()
			} else {
				response, status, err = client.DeliveryTypesWithParams(*testCase.params)
			}
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, status)
			assert.True(t, response.Success)
		})
	}
}

func TestDeliveryCostSettings(t *testing.T) {
	const payload = `{
		"success": true,
		"deliveryTypes": {
			"courier": {
				"code": "courier",
				"costCalculationType": "dynamic",
				"netCostCalculationType": "static",
				"codMarkup": 1.5,
				"limitByRegions": true,
				"dateTimeCostConditions": [
					{"days":["1","2","3","4","5"],"timeStart":"10:00","timeEnd":"21:30","value":50,"netValue":10},
					{"days":[],"timeStart":"10:00","timeEnd":"22:00","value":0,"netValue":0}
				],
				"regionWeightCostConditions": [{
					"country":"RU","region":"Регион","regionId":1,"city":"Город","cityId":2,
					"weightStart":0,"weightEnd":5000,"summStart":null,"summEnd":5000,
					"value":500,"netValue":100,"netValueType":"subtract"
				}]
			}
		}
	}`
	var response DeliveryTypesResponse
	require.NoError(t, json.Unmarshal([]byte(payload), &response))
	delivery := response.DeliveryTypes["courier"]
	assert.Equal(t, "dynamic", delivery.CostCalculationType)
	assert.Equal(t, "static", delivery.NetCostCalculationType)
	assert.Equal(t, float32(1.5), delivery.CODMarkup)
	assert.True(t, delivery.LimitByRegions)
	require.Len(t, delivery.DateTimeCostConditions, 2)
	assert.Equal(t, DeliveryDateTimeCostCondition{
		Days: []string{"1", "2", "3", "4", "5"}, TimeStart: "10:00", TimeEnd: "21:30", Value: 50, NetValue: 10,
	}, delivery.DateTimeCostConditions[0])
	assert.Equal(t, DeliveryDateTimeCostCondition{
		Days: []string{}, TimeStart: "10:00", TimeEnd: "22:00",
	}, delivery.DateTimeCostConditions[1])
	require.Len(t, delivery.RegionWeightCostConditions, 1)
	region := delivery.RegionWeightCostConditions[0]
	assert.Equal(t, "RU", region.Country)
	assert.Equal(t, "Регион", region.Region)
	assert.Equal(t, 1, region.RegionID)
	assert.Equal(t, "Город", region.City)
	assert.Equal(t, 2, region.CityID)
	require.NotNil(t, region.WeightStart)
	assert.Equal(t, float32(0), *region.WeightStart)
	require.NotNil(t, region.WeightEnd)
	assert.Equal(t, float32(5000), *region.WeightEnd)
	assert.Nil(t, region.SummStart)
	require.NotNil(t, region.SummEnd)
	assert.Equal(t, float32(5000), *region.SummEnd)
	assert.Equal(t, float32(500), region.Value)
	assert.Equal(t, float32(100), region.NetValue)
	assert.Equal(t, "subtract", region.NetValueType)
}

func TestDeliveryCostSettingsEmpty(t *testing.T) {
	for _, payload := range []string{
		`{}`,
		`{"dateTimeCostConditions":[],"regionWeightCostConditions":[]}`,
		`{"dateTimeCostConditions":null,"regionWeightCostConditions":null}`,
	} {
		var delivery DeliveryType
		require.NoError(t, json.Unmarshal([]byte(payload), &delivery))
		assert.Empty(t, delivery.DateTimeCostConditions)
		assert.Empty(t, delivery.RegionWeightCostConditions)
	}
}

func TestDeliveryTypesAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "true", request.URL.Query().Get("includeCostSettings"))
		assert.Empty(t, request.URL.Query()["filter[codes][]"])
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(writer, `{"success":false,"errorMsg":"Missing delivery type codes"}`)
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)
	client := New(server.URL, "test-key").WithHTTPClient(server.Client())
	_, status, err := client.DeliveryTypesWithParams(DeliveryTypesRequest{IncludeCostSettings: true})
	require.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, status)
	_, ok := AsAPIError(err)
	assert.True(t, ok)
}

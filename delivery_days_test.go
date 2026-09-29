package retailcrm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderDeliveryDays(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want StringOrNumber
	}{
		{name: "диапазон из ответа CRM", raw: `"1 - 2"`, want: "1 - 2"},
		{name: "строка с числом", raw: `"2"`, want: "2"},
		{name: "число", raw: `2`, want: "2"},
		{name: "ноль", raw: `0`, want: "0"},
		{name: "пустая строка", raw: `""`},
		{name: "null", raw: `null`},
		{name: "поле отсутствует"},
	}

	for _, scenario := range cases {
		t.Run(scenario.name, func(t *testing.T) {
			fields := `"status":"delivered","customField":"сохранено"`
			if scenario.raw != "" {
				fields += `,"days":` + scenario.raw
			}
			payload := `{"success":true,"order":{"id":123,"delivery":{"data":{` + fields + `}}}}`

			var response OrderResponse
			require.NoError(t, json.Unmarshal([]byte(payload), &response), "Ответ заказа должен читаться")
			require.NotNil(t, response.Order.Delivery, "Доставка должна сохраняться")
			require.NotNil(t, response.Order.Delivery.Data, "Данные доставки должны сохраняться")

			data := response.Order.Delivery.Data
			assert.Equal(t, scenario.want, data.Days, "Срок доставки должен сохраняться без потери формата")
			assert.Equal(t, "delivered", data.Status, "Статус доставки должен сохраняться")
			assert.Equal(t, "сохранено", data.AdditionalFields["customField"], "Дополнительные поля должны сохраняться")
			assert.NotContains(t, data.AdditionalFields, "days", "Срок не должен дублироваться в дополнительных полях")

			encoded, err := json.Marshal(data)
			require.NoError(t, err, "Данные доставки должны сериализоваться")
			var fieldsOut map[string]any
			require.NoError(t, json.Unmarshal(encoded, &fieldsOut))
			if scenario.want == "" {
				assert.NotContains(t, fieldsOut, "days", "Пустой срок должен пропускаться")
			} else {
				assert.Equal(t, string(scenario.want), fieldsOut["days"], "Срок должен передаваться строкой по контракту CRM")
			}
			assert.Equal(t, "сохранено", fieldsOut["customField"], "Дополнительные поля должны переживать повторную сериализацию")
		})
	}
}

func TestOrderDeliveryDaysInvalid(t *testing.T) {
	for _, raw := range []string{`true`, `[]`, `{}`} {
		t.Run(raw, func(t *testing.T) {
			payload := `{"success":true,"order":{"delivery":{"data":{"days":` + raw + `}}}}`
			var response OrderResponse
			require.Error(t, json.Unmarshal([]byte(payload), &response), "Неподдерживаемый тип срока не должен приниматься")
		})
	}
}

package retailcrm

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestAPIErrorsList_UnmarshalJSON(t *testing.T) {
	var list APIErrorsList

	require.NoError(t, json.Unmarshal([]byte(`["first", "second"]`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["0"], "first")
	assert.Equal(t, list["1"], "second")

	require.NoError(t, json.Unmarshal([]byte(`{"a": "first", "b": "second"}`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["a"], "first")
	assert.Equal(t, list["b"], "second")
}

func TestCustomFieldsList_UnmarshalJSON(t *testing.T) {
	var list StringMap

	require.NoError(t, json.Unmarshal([]byte(`["first", "second"]`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["0"], "first")
	assert.Equal(t, list["1"], "second")

	require.NoError(t, json.Unmarshal([]byte(`{"a": "first", "b": "second"}`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["a"], "first")
	assert.Equal(t, list["b"], "second")
}

func TestOrderPayments_UnmarshalJSON(t *testing.T) {
	var list OrderPayments

	require.NoError(t, json.Unmarshal([]byte(`[{"id": 1}, {"id": 2}]`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["0"], OrderPayment{ID: 1})
	assert.Equal(t, list["1"], OrderPayment{ID: 2})

	require.NoError(t, json.Unmarshal([]byte(`{"a": {"id": 1}, "b": {"id": 2}}`), &list))
	assert.Len(t, list, 2)
	assert.Equal(t, list["a"], OrderPayment{ID: 1})
	assert.Equal(t, list["b"], OrderPayment{ID: 2})
}

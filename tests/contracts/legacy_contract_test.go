package contracts

import (
	"encoding/json"
	"testing"
)

type responseShape map[string]any

func assertFields(t *testing.T, body string, fields []string) {
	t.Helper()
	var m responseShape
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	for _, f := range fields {
		if _, ok := m[f]; !ok {
			t.Fatalf("missing required field %s", f)
		}
	}
}

func TestLegacyCatalogContractShape(t *testing.T) {
	assertFields(t, `{"id":"p1","vendor_id":"v1","name":"Item"}`, []string{"id", "vendor_id", "name"})
}

func TestLegacyOrderContractShape(t *testing.T) {
	assertFields(t, `{"id":"o1","status":"CONFIRMED","lines":[]}`, []string{"id", "status", "lines"})
}

func TestLegacyInventoryReserveContractShape(t *testing.T) {
	assertFields(t, `{"order_id":"o1","sku":"s1","reserved":1}`, []string{"order_id", "sku", "reserved"})
}

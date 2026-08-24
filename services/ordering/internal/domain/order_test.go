package domain

import (
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"testing"
)

func TestOrderLifecycle(t *testing.T) {
	vendor := sharedkernel.MustVendorID("v-1")
	o := NewOrder("o-1", vendor)
	if err := o.AddLine("sku-1", 1, 100); err != nil {
		t.Fatal(err)
	}
	o.Submit()
	o.Confirm()
	if o.Status() != OrderStatusConfirmed {
		t.Fatalf("expected confirmed")
	}
}

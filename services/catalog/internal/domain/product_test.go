package domain

import (
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"testing"
)

func TestProductVendorIsolation(t *testing.T) {
	vendorA := sharedkernel.MustVendorID("v-a")
	vendorB := sharedkernel.MustVendorID("v-b")
	p, err := NewProduct("p-1", vendorA, "Item")
	if err != nil {
		t.Fatal(err)
	}
	if err := p.AddVariant(vendorB, "sku-1", 100); err == nil {
		t.Fatal("expected vendor mismatch error")
	}
}

package domain

import (
	"testing"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

func TestReserveAndRelease(t *testing.T) {
	vendor := sharedkernel.MustVendorID("v-1")
	s := NewStockItem("s-1", vendor, "sku-1", 5)
	_, err := s.Reserve("o-1", 3, 0)
	if err != nil {
		t.Fatal(err)
	}
	if s.Available() != 2 || s.Reserved() != 3 {
		t.Fatalf("unexpected stock state")
	}
	s.Release("o-1", 3)
	if s.Available() != 5 || s.Reserved() != 0 {
		t.Fatalf("unexpected release state")
	}
}

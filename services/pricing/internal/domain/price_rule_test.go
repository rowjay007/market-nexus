package domain

import (
	"testing"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

func TestQuote(t *testing.T) {
	rule := NewPriceRule("r1", sharedkernel.MustVendorID("v1"), 1000)
	quoted, err := rule.Quote("o1", 10000)
	if err != nil {
		t.Fatal(err)
	}
	if quoted.Total <= 0 {
		t.Fatal("expected positive total")
	}
}

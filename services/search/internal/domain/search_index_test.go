package domain

import (
	"testing"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

func TestSearchDocumentMatchAndRank(t *testing.T) {
	doc := NewSearchDocument("d1", "p1", sharedkernel.MustVendorID("v1"), "Gaming Laptop", "high performance", 2)
	if !doc.Matches("laptop") {
		t.Fatal("expected query to match")
	}
	if doc.RankScore() <= 0 {
		t.Fatal("expected positive rank")
	}
}

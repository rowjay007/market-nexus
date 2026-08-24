package domain

import (
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"testing"
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

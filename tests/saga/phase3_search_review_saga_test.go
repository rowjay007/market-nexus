package saga

import (
	"testing"

	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/reviewtrust"
	"github.com/rowjay007/market-nexus/services/search"
)

func TestPhase3_SearchVendorIsolationAndRanking(t *testing.T) {
	bus := events.NewInMemoryBus()
	searchSvc := search.NewInMemoryService(bus)

	vendorA := sharedkernel.MustVendorID("v-a")
	vendorB := sharedkernel.MustVendorID("v-b")

	_ = searchSvc.Index("d1", "p1", vendorA, "Laptop Pro", "fast laptop", 2)
	_ = searchSvc.Index("d2", "p2", vendorA, "Laptop Air", "portable laptop", 1)
	_ = searchSvc.Index("d3", "p3", vendorB, "Laptop Other", "other vendor", 3)

	resultsA := searchSvc.Query(vendorA, "laptop")
	if len(resultsA) != 2 {
		t.Fatalf("expected 2 vendor A results, got %d", len(resultsA))
	}
	if resultsA[0] != "p1" {
		t.Fatalf("expected tier-ranked first result p1, got %s", resultsA[0])
	}

	resultsB := searchSvc.Query(vendorB, "laptop")
	if len(resultsB) != 1 || resultsB[0] != "p3" {
		t.Fatalf("unexpected vendor B results: %+v", resultsB)
	}
}

func TestPhase3_ReviewTrustVendorMetrics(t *testing.T) {
	bus := events.NewInMemoryBus()
	rt := reviewtrust.NewInMemoryService(bus)
	vendor := sharedkernel.MustVendorID("v-rate")

	if err := rt.SubmitReview("r1", "p1", vendor, "b1", 5, "great"); err != nil {
		t.Fatal(err)
	}
	if err := rt.SubmitReview("r2", "p2", vendor, "b2", 3, "ok"); err != nil {
		t.Fatal(err)
	}

	avg := rt.VendorAverageRating(vendor)
	if avg != 4.0 {
		t.Fatalf("expected average 4.0, got %.2f", avg)
	}
}

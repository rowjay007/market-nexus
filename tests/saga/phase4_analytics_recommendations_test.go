package saga

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/analytics"
	"testing"
)

func TestPhase4_RecommendationsPipeline(t *testing.T) {
	bus := events.NewInMemoryBus()
	a := analytics.NewInMemoryService(bus)
	vendor := sharedkernel.MustVendorID("v-ana")

	if err := a.RecordBehavior("b1", "u1", "p1", vendor, "VIEW"); err != nil {
		t.Fatal(err)
	}
	if err := a.RecordBehavior("b2", "u1", "p2", vendor, "CLICK"); err != nil {
		t.Fatal(err)
	}
	if err := a.RecordBehavior("b3", "u1", "p3", vendor, "PURCHASE"); err != nil {
		t.Fatal(err)
	}

	recs := a.Recommendations("u1", 3)
	if len(recs) != 3 {
		t.Fatalf("expected 3 recommendations, got %d", len(recs))
	}
	if recs[0] != "p3" {
		t.Fatalf("expected purchase-ranked item first, got %s", recs[0])
	}
}

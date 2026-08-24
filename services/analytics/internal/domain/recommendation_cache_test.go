package domain

import (
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"testing"
)

func TestRecommendationWeights(t *testing.T) {
	cache := NewRecommendationCache()
	vendor := sharedkernel.MustVendorID("v1")

	view, _ := NewBehaviorEvent("e1", "u1", "p1", vendor, BehaviorView)
	click, _ := NewBehaviorEvent("e2", "u1", "p2", vendor, BehaviorClick)
	purchase, _ := NewBehaviorEvent("e3", "u1", "p3", vendor, BehaviorPurchase)
	cache.Apply(view)
	cache.Apply(click)
	cache.Apply(purchase)

	top := cache.TopProducts("u1", 3)
	if len(top) != 3 {
		t.Fatalf("expected 3 recommendations")
	}
	if top[0] != "p3" {
		t.Fatalf("expected purchase-weighted item first, got %s", top[0])
	}
}

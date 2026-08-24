package domain

import (
	"testing"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

func TestReviewValidationAndEvent(t *testing.T) {
	r, err := NewReview("r1", "p1", sharedkernel.MustVendorID("v1"), "b1", 5, "great")
	if err != nil {
		t.Fatal(err)
	}
	e := r.SubmittedEvent()
	if e.ReviewID != "r1" {
		t.Fatalf("unexpected review id")
	}
}

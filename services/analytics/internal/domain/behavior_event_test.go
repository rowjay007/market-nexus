package domain

import (
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"testing"
)

func TestBehaviorEventValidation(t *testing.T) {
	_, err := NewBehaviorEvent("e1", "u1", "p1", sharedkernel.MustVendorID("v1"), BehaviorClick)
	if err != nil {
		t.Fatal(err)
	}
}

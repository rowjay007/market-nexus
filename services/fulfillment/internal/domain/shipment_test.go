package domain

import (
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"testing"
)

func TestScheduleAndCancel(t *testing.T) {
	sh := NewShipment("s1", "o1", sharedkernel.MustVendorID("v1"))
	_, err := sh.Schedule("Street 1")
	if err != nil {
		t.Fatal(err)
	}
	cancelled := sh.Cancel("test")
	if cancelled.OrderID != "o1" {
		t.Fatalf("expected order id")
	}
}

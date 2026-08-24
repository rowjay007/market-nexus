package domain

import (
	"testing"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

func TestCaptureAndRefund(t *testing.T) {
	intent := NewPaymentIntent("pi1", "o1", sharedkernel.MustVendorID("v1"))
	_, err := intent.Capture(1000)
	if err != nil {
		t.Fatal(err)
	}
	ref := intent.Refund("test")
	if ref.Amount != 1000 {
		t.Fatalf("expected amount to refund")
	}
}

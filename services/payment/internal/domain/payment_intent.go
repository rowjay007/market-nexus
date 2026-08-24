package domain

import (
	"errors"
	"time"

	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

var ErrInvalidCaptureAmount = errors.New("invalid capture amount")

type PaymentIntentID string

type PaymentIntent struct {
	id       PaymentIntentID
	orderID  string
	vendorID sharedkernel.VendorID
	amount   int64
	status   string
}

func NewPaymentIntent(id PaymentIntentID, orderID string, vendorID sharedkernel.VendorID) *PaymentIntent {
	return &PaymentIntent{id: id, orderID: orderID, vendorID: vendorID, status: "CREATED"}
}

func (p *PaymentIntent) Capture(amount int64) (PaymentCaptured, error) {
	if amount <= 0 {
		return PaymentCaptured{}, ErrInvalidCaptureAmount
	}
	p.amount = amount
	p.status = "CAPTURED"
	return PaymentCaptured{
		BaseEvent: events.BaseEvent{Type: "PaymentCaptured", At: time.Now().UTC()},
		OrderID:   p.orderID,
		VendorID:  p.vendorID.String(),
		Amount:    amount,
	}, nil
}

func (p *PaymentIntent) Refund(reason string) PaymentRefunded {
	p.status = "REFUNDED"
	return PaymentRefunded{
		BaseEvent: events.BaseEvent{Type: "PaymentRefunded", At: time.Now().UTC()},
		OrderID:   p.orderID,
		VendorID:  p.vendorID.String(),
		Amount:    p.amount,
		Reason:    reason,
	}
}

func (p *PaymentIntent) OrderID() string {
	return p.orderID
}

type PaymentCaptured struct {
	events.BaseEvent
	OrderID  string
	VendorID string
	Amount   int64
}

type PaymentRefunded struct {
	events.BaseEvent
	OrderID  string
	VendorID string
	Amount   int64
	Reason   string
}

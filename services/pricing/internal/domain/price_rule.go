package domain

import (
	"errors"
	"time"

	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

var ErrInvalidSubtotal = errors.New("invalid subtotal")

type PriceRuleID string

type PriceRule struct {
	id          PriceRuleID
	vendorID    sharedkernel.VendorID
	discountBps int64
}

func NewPriceRule(id PriceRuleID, vendorID sharedkernel.VendorID, discountBps int64) *PriceRule {
	if discountBps < 0 {
		discountBps = 0
	}
	if discountBps > 10000 {
		discountBps = 10000
	}
	return &PriceRule{id: id, vendorID: vendorID, discountBps: discountBps}
}

func (r *PriceRule) Quote(orderID string, subtotal int64) (PriceQuoted, error) {
	if subtotal < 0 {
		return PriceQuoted{}, ErrInvalidSubtotal
	}
	discount := (subtotal * r.discountBps) / 10000
	tax := (subtotal - discount) / 10
	total := subtotal - discount + tax
	return PriceQuoted{
		BaseEvent: events.BaseEvent{Type: "PriceQuoted", At: time.Now().UTC()},
		OrderID:   orderID,
		VendorID:  r.vendorID.String(),
		Subtotal:  subtotal,
		Discount:  discount,
		Tax:       tax,
		Total:     total,
	}, nil
}

func (r *PriceRule) VendorID() sharedkernel.VendorID {
	return r.vendorID
}

type PriceQuoted struct {
	events.BaseEvent
	OrderID   string
	VendorID  string
	Subtotal  int64
	Discount  int64
	Tax       int64
	Total     int64
}

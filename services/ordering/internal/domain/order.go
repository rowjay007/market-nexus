package domain

import (
	"errors"
	"time"

	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

var (
	ErrInvalidOrderLine = errors.New("invalid order line")
)

type OrderID string
type OrderStatus string

const (
	OrderStatusDraft     OrderStatus = "DRAFT"
	OrderStatusSubmitted OrderStatus = "SUBMITTED"
	OrderStatusConfirmed OrderStatus = "CONFIRMED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type OrderLine struct {
	SKU      string
	Quantity int
	Price    int64
}

type Order struct {
	id       OrderID
	vendorID sharedkernel.VendorID
	status   OrderStatus
	lines    []OrderLine
}

func NewOrder(id OrderID, vendorID sharedkernel.VendorID) *Order {
	return &Order{id: id, vendorID: vendorID, status: OrderStatusDraft, lines: []OrderLine{}}
}

func (o *Order) AddLine(sku string, qty int, price int64) error {
	if sku == "" || qty <= 0 || price < 0 {
		return ErrInvalidOrderLine
	}
	o.lines = append(o.lines, OrderLine{SKU: sku, Quantity: qty, Price: price})
	return nil
}

func (o *Order) Submit() OrderPlaced {
	o.status = OrderStatusSubmitted
	return OrderPlaced{
		BaseEvent: events.BaseEvent{Type: "OrderPlaced", At: time.Now().UTC()},
		OrderID:   string(o.id),
		VendorID:  o.vendorID.String(),
		Lines:     o.lines,
	}
}

func (o *Order) Confirm() OrderConfirmed {
	o.status = OrderStatusConfirmed
	return OrderConfirmed{
		BaseEvent: events.BaseEvent{Type: "OrderConfirmed", At: time.Now().UTC()},
		OrderID:   string(o.id),
		VendorID:  o.vendorID.String(),
	}
}

func (o *Order) Cancel(reason string) OrderCancelled {
	o.status = OrderStatusCancelled
	return OrderCancelled{
		BaseEvent: events.BaseEvent{Type: "OrderCancelled", At: time.Now().UTC()},
		OrderID:   string(o.id),
		VendorID:  o.vendorID.String(),
		Reason:    reason,
	}
}

func (o *Order) ID() OrderID { return o.id }
func (o *Order) VendorID() sharedkernel.VendorID { return o.vendorID }
func (o *Order) Status() OrderStatus { return o.status }
func (o *Order) Lines() []OrderLine {
	out := make([]OrderLine, len(o.lines))
	copy(out, o.lines)
	return out
}

type OrderPlaced struct {
	events.BaseEvent
	OrderID  string
	VendorID string
	Lines    []OrderLine
}

type OrderConfirmed struct {
	events.BaseEvent
	OrderID  string
	VendorID string
}

type OrderCancelled struct {
	events.BaseEvent
	OrderID  string
	VendorID string
	Reason   string
}

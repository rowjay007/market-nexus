package domain

import (
	"errors"
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"time"
)

var ErrInvalidAddress = errors.New("invalid shipping address")

type ShipmentID string

type Shipment struct {
	id       ShipmentID
	orderID  string
	vendorID sharedkernel.VendorID
	status   string
	address  string
}

func NewShipment(id ShipmentID, orderID string, vendorID sharedkernel.VendorID) *Shipment {
	return &Shipment{id: id, orderID: orderID, vendorID: vendorID, status: "CREATED"}
}

func (s *Shipment) Schedule(address string) (FulfillmentScheduled, error) {
	if address == "" {
		return FulfillmentScheduled{}, ErrInvalidAddress
	}
	s.address = address
	s.status = "SCHEDULED"
	return FulfillmentScheduled{
		BaseEvent: events.BaseEvent{Type: "FulfillmentScheduled", At: time.Now().UTC()},
		OrderID:   s.orderID,
		VendorID:  s.vendorID.String(),
		Address:   address,
	}, nil
}

func (s *Shipment) Cancel(reason string) FulfillmentCancelled {
	s.status = "CANCELLED"
	return FulfillmentCancelled{
		BaseEvent: events.BaseEvent{Type: "FulfillmentCancelled", At: time.Now().UTC()},
		OrderID:   s.orderID,
		VendorID:  s.vendorID.String(),
		Reason:    reason,
	}
}

func (s *Shipment) OrderID() string {
	return s.orderID
}

type FulfillmentScheduled struct {
	events.BaseEvent
	OrderID  string
	VendorID string
	Address  string
}

type FulfillmentCancelled struct {
	events.BaseEvent
	OrderID  string
	VendorID string
	Reason   string
}

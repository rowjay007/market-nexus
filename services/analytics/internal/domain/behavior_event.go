package domain

import (
	"errors"
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"time"
)

var ErrInvalidBehaviorType = errors.New("invalid behavior type")

type BehaviorType string

const (
	BehaviorView     BehaviorType = "VIEW"
	BehaviorClick    BehaviorType = "CLICK"
	BehaviorPurchase BehaviorType = "PURCHASE"
)

type BehaviorEvent struct {
	id         string
	userID     string
	productID  string
	vendorID   sharedkernel.VendorID
	behavior   BehaviorType
	occurredAt time.Time
}

func NewBehaviorEvent(id string, userID string, productID string, vendorID sharedkernel.VendorID, behavior BehaviorType) (*BehaviorEvent, error) {
	if behavior != BehaviorView && behavior != BehaviorClick && behavior != BehaviorPurchase {
		return nil, ErrInvalidBehaviorType
	}
	return &BehaviorEvent{id: id, userID: userID, productID: productID, vendorID: vendorID, behavior: behavior, occurredAt: time.Now().UTC()}, nil
}

func (e *BehaviorEvent) UserID() string                  { return e.userID }
func (e *BehaviorEvent) ProductID() string               { return e.productID }
func (e *BehaviorEvent) VendorID() sharedkernel.VendorID { return e.vendorID }
func (e *BehaviorEvent) Behavior() BehaviorType          { return e.behavior }

func (e *BehaviorEvent) RecordedEvent() UserBehaviorRecorded {
	return UserBehaviorRecorded{
		BaseEvent: events.BaseEvent{Type: "UserBehaviorRecorded", At: time.Now().UTC()},
		EventID:   e.id,
		UserID:    e.userID,
		ProductID: e.productID,
		VendorID:  e.vendorID.String(),
		Behavior:  string(e.behavior),
	}
}

type UserBehaviorRecorded struct {
	events.BaseEvent
	EventID   string
	UserID    string
	ProductID string
	VendorID  string
	Behavior  string
}

package domain

import (
	"errors"
	"time"

	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
)

type StockItemID string

type StockItem struct {
	id        StockItemID
	vendorID  sharedkernel.VendorID
	sku       string
	available int
	reserved  int
	version   int
}

func NewStockItem(id StockItemID, vendorID sharedkernel.VendorID, sku string, available int) *StockItem {
	return &StockItem{id: id, vendorID: vendorID, sku: sku, available: available}
}

func (s *StockItem) Reserve(orderID string, qty int, expectedVersion int) (InventoryReserved, error) {
	if expectedVersion != s.version {
		return InventoryReserved{}, ErrInsufficientStock
	}
	if qty <= 0 || s.available < qty {
		return InventoryReserved{}, ErrInsufficientStock
	}
	s.available -= qty
	s.reserved += qty
	s.version++
	return InventoryReserved{
		BaseEvent: events.BaseEvent{Type: "InventoryReserved", At: time.Now().UTC()},
		OrderID:   orderID,
		SKU:       s.sku,
		VendorID:  s.vendorID.String(),
		Quantity:  qty,
	}, nil
}

func (s *StockItem) Release(orderID string, qty int) InventoryReleased {
	if qty > s.reserved {
		qty = s.reserved
	}
	s.reserved -= qty
	s.available += qty
	s.version++
	return InventoryReleased{
		BaseEvent: events.BaseEvent{Type: "InventoryReleased", At: time.Now().UTC()},
		OrderID:   orderID,
		SKU:       s.sku,
		VendorID:  s.vendorID.String(),
		Quantity:  qty,
	}
}

func (s *StockItem) SKU() string { return s.sku }
func (s *StockItem) VendorID() sharedkernel.VendorID { return s.vendorID }
func (s *StockItem) Available() int { return s.available }
func (s *StockItem) Reserved() int { return s.reserved }
func (s *StockItem) Version() int { return s.version }

type InventoryReserved struct {
	events.BaseEvent
	OrderID  string
	SKU      string
	VendorID string
	Quantity int
}

type InventoryReleased struct {
	events.BaseEvent
	OrderID  string
	SKU      string
	VendorID string
	Quantity int
}

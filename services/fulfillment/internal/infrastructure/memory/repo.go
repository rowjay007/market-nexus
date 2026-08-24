package memory

import (
	"sync"

	"github.com/rowjay007/market-nexus/services/fulfillment/internal/domain"
)

type ShipmentRepo struct {
	mu    sync.RWMutex
	items map[string]*domain.Shipment
}

func NewShipmentRepo() *ShipmentRepo {
	return &ShipmentRepo{items: map[string]*domain.Shipment{}}
}

func (r *ShipmentRepo) Save(shipment *domain.Shipment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[shipment.OrderID()] = shipment
	return nil
}

func (r *ShipmentRepo) GetByOrderID(orderID string) (*domain.Shipment, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[orderID]
	return item, ok
}

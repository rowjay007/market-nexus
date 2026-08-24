package application

import (
	"errors"

	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/inventory/internal/domain"
)

var ErrStockNotFound = errors.New("stock item not found")

type StockRepository interface {
	Save(item *domain.StockItem) error
	GetByVendorAndSKU(vendorID sharedkernel.VendorID, sku string) (*domain.StockItem, bool)
}

type InventoryService struct {
	repo StockRepository
	bus  *events.InMemoryBus
}

func NewInventoryService(repo StockRepository, bus *events.InMemoryBus) *InventoryService {
	return &InventoryService{repo: repo, bus: bus}
}

func (s *InventoryService) Reserve(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error {
	item, ok := s.repo.GetByVendorAndSKU(vendorID, sku)
	if !ok {
		return ErrStockNotFound
	}
	evt, err := item.Reserve(orderID, qty, item.Version())
	if err != nil {
		return err
	}
	if err := s.repo.Save(item); err != nil {
		return err
	}
	s.bus.Publish(evt)
	return nil
}

func (s *InventoryService) Release(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error {
	item, ok := s.repo.GetByVendorAndSKU(vendorID, sku)
	if !ok {
		return ErrStockNotFound
	}
	evt := item.Release(orderID, qty)
	if err := s.repo.Save(item); err != nil {
		return err
	}
	s.bus.Publish(evt)
	return nil
}

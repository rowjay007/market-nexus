package inventory

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/inventory/internal/application"
	"github.com/rowjay007/market-nexus/services/inventory/internal/domain"
	"github.com/rowjay007/market-nexus/services/inventory/internal/infrastructure/memory"
)

type Service struct {
	app  *application.InventoryService
	repo *memory.StockRepo
}

func NewInMemoryService(bus *events.InMemoryBus) *Service {
	repo := memory.NewStockRepo()
	return &Service{app: application.NewInventoryService(repo, bus), repo: repo}
}

func (s *Service) SeedStock(id string, vendorID sharedkernel.VendorID, sku string, available int) error {
	return s.repo.Save(domain.NewStockItem(domain.StockItemID(id), vendorID, sku, available))
}

func (s *Service) Reserve(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error {
	return s.app.Reserve(orderID, vendorID, sku, qty)
}

func (s *Service) Release(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error {
	return s.app.Release(orderID, vendorID, sku, qty)
}

func (s *Service) Stock(vendorID sharedkernel.VendorID, sku string) (available int, reserved int, ok bool) {
	item, found := s.repo.GetByVendorAndSKU(vendorID, sku)
	if !found {
		return 0, 0, false
	}
	return item.Available(), item.Reserved(), true
}

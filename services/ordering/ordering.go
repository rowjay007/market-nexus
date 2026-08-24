package ordering

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/ordering/internal/application"
	"github.com/rowjay007/market-nexus/services/ordering/internal/infrastructure/memory"
)

type LineInput struct {
	SKU      string
	Quantity int
}

type CatalogACL interface {
	ResolvePrice(vendorID string, sku string) (int64, bool)
}

type InventoryACL interface {
	Reserve(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error
	Release(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error
}

type Service struct {
	app *application.OrderingService
}

func NewInMemoryService(bus *events.InMemoryBus, catalogACL CatalogACL, inventoryACL InventoryACL) *Service {
	repo := memory.NewOrderRepo()
	return &Service{app: application.NewOrderingService(repo, bus, catalogACL, inventoryACL)}
}

func (s *Service) PlaceOrder(id string, vendorID sharedkernel.VendorID, lines []LineInput) (string, error) {
	in := make([]struct {
		SKU      string
		Quantity int
	}, 0, len(lines))
	for _, line := range lines {
		in = append(in, struct {
			SKU      string
			Quantity int
		}{SKU: line.SKU, Quantity: line.Quantity})
	}
	order, err := s.app.PlaceOrder(id, vendorID, in)
	if err != nil {
		return "", err
	}
	return string(order.Status()), nil
}

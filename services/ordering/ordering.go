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

type PricingACL interface {
	Quote(orderID string, vendorID sharedkernel.VendorID, subtotal int64) (int64, error)
}

type PaymentACL interface {
	Capture(orderID string, vendorID sharedkernel.VendorID, amount int64) error
	Refund(orderID string, reason string) error
}

type FulfillmentACL interface {
	Schedule(orderID string, vendorID sharedkernel.VendorID, address string) error
	Cancel(orderID string, reason string) error
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

func (s *Service) CheckoutOrder(
	id string,
	vendorID sharedkernel.VendorID,
	lines []LineInput,
	pricingACL PricingACL,
	paymentACL PaymentACL,
	fulfillmentACL FulfillmentACL,
	shippingAddress string,
) (string, error) {
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
	order, err := s.app.CheckoutOrder(id, vendorID, in, pricingACL, paymentACL, fulfillmentACL, shippingAddress)
	if err != nil {
		return "", err
	}
	return string(order.Status()), nil
}

package application

import (
	"errors"

	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	orderingdomain "github.com/rowjay007/market-nexus/services/ordering/internal/domain"
)

var (
	ErrCatalogSKUUnknown = errors.New("catalog sku unknown in ordering acl")
)

type OrderRepository interface {
	Save(order *orderingdomain.Order) error
	GetByID(id orderingdomain.OrderID) (*orderingdomain.Order, bool)
}

type CatalogACL interface {
	ResolvePrice(vendorID string, sku string) (int64, bool)
}

type InventoryACL interface {
	Reserve(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error
	Release(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error
}

type OrderingService struct {
	repo         OrderRepository
	bus          *events.InMemoryBus
	catalogACL   CatalogACL
	inventoryACL InventoryACL
}

func NewOrderingService(repo OrderRepository, bus *events.InMemoryBus, catalogACL CatalogACL, inventoryACL InventoryACL) *OrderingService {
	return &OrderingService{repo: repo, bus: bus, catalogACL: catalogACL, inventoryACL: inventoryACL}
}

func (s *OrderingService) PlaceOrder(id string, vendorID sharedkernel.VendorID, lines []struct {
	SKU      string
	Quantity int
}) (*orderingdomain.Order, error) {
	order := orderingdomain.NewOrder(orderingdomain.OrderID(id), vendorID)
	for _, line := range lines {
		price, ok := s.catalogACL.ResolvePrice(vendorID.String(), line.SKU)
		if !ok {
			return nil, ErrCatalogSKUUnknown
		}
		if err := order.AddLine(line.SKU, line.Quantity, price); err != nil {
			return nil, err
		}
	}
	if err := s.repo.Save(order); err != nil {
		return nil, err
	}
	placed := order.Submit()
	s.bus.Publish(placed)

	for _, line := range order.Lines() {
		if err := s.inventoryACL.Reserve(id, vendorID, line.SKU, line.Quantity); err != nil {
			cancel := order.Cancel(err.Error())
			_ = s.repo.Save(order)
			s.bus.Publish(cancel)
			for _, reserved := range order.Lines() {
				if reserved.SKU == line.SKU {
					break
				}
				_ = s.inventoryACL.Release(id, vendorID, reserved.SKU, reserved.Quantity)
			}
			return nil, err
		}
	}

	confirmed := order.Confirm()
	_ = s.repo.Save(order)
	s.bus.Publish(confirmed)
	return order, nil
}

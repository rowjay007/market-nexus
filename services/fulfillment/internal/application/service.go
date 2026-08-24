package application

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/fulfillment/internal/domain"
)

type ShipmentRepository interface {
	Save(shipment *domain.Shipment) error
	GetByOrderID(orderID string) (*domain.Shipment, bool)
}

type FulfillmentService struct {
	repo ShipmentRepository
	bus  *events.InMemoryBus
}

func NewFulfillmentService(repo ShipmentRepository, bus *events.InMemoryBus) *FulfillmentService {
	return &FulfillmentService{repo: repo, bus: bus}
}

func (s *FulfillmentService) Schedule(orderID string, vendorID sharedkernel.VendorID, address string) error {
	shipment, ok := s.repo.GetByOrderID(orderID)
	if !ok {
		shipment = domain.NewShipment(domain.ShipmentID("sh-"+orderID), orderID, vendorID)
	}
	evt, err := shipment.Schedule(address)
	if err != nil {
		return err
	}
	if err := s.repo.Save(shipment); err != nil {
		return err
	}
	s.bus.Publish(evt)
	return nil
}

func (s *FulfillmentService) Cancel(orderID string, reason string) error {
	shipment, ok := s.repo.GetByOrderID(orderID)
	if !ok {
		return nil
	}
	evt := shipment.Cancel(reason)
	if err := s.repo.Save(shipment); err != nil {
		return err
	}
	s.bus.Publish(evt)
	return nil
}

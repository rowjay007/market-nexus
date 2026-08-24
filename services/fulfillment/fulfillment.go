package fulfillment

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/fulfillment/internal/application"
	"github.com/rowjay007/market-nexus/services/fulfillment/internal/infrastructure/memory"
)

type Service struct {
	app *application.FulfillmentService
}

func NewInMemoryService(bus *events.InMemoryBus) *Service {
	repo := memory.NewShipmentRepo()
	return &Service{app: application.NewFulfillmentService(repo, bus)}
}

func (s *Service) Schedule(orderID string, vendorID sharedkernel.VendorID, address string) error {
	return s.app.Schedule(orderID, vendorID, address)
}

func (s *Service) Cancel(orderID string, reason string) error {
	return s.app.Cancel(orderID, reason)
}

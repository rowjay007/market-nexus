package payment

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/payment/internal/application"
	"github.com/rowjay007/market-nexus/services/payment/internal/infrastructure/memory"
)

type Service struct {
	app *application.PaymentService
}

func NewInMemoryService(bus *events.InMemoryBus) *Service {
	repo := memory.NewPaymentRepo()
	return &Service{app: application.NewPaymentService(repo, bus)}
}

func (s *Service) Capture(orderID string, vendorID sharedkernel.VendorID, amount int64) error {
	return s.app.Capture(orderID, vendorID, amount)
}

func (s *Service) Refund(orderID string, reason string) error {
	return s.app.Refund(orderID, reason)
}

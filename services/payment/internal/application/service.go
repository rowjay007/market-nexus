package application

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/payment/internal/domain"
)

type PaymentRepository interface {
	Save(intent *domain.PaymentIntent) error
	GetByOrderID(orderID string) (*domain.PaymentIntent, bool)
}

type PaymentService struct {
	repo PaymentRepository
	bus  *events.InMemoryBus
}

func NewPaymentService(repo PaymentRepository, bus *events.InMemoryBus) *PaymentService {
	return &PaymentService{repo: repo, bus: bus}
}

func (s *PaymentService) Capture(orderID string, vendorID sharedkernel.VendorID, amount int64) error {
	intent, ok := s.repo.GetByOrderID(orderID)
	if !ok {
		intent = domain.NewPaymentIntent(domain.PaymentIntentID("pi-"+orderID), orderID, vendorID)
	}
	evt, err := intent.Capture(amount)
	if err != nil {
		return err
	}
	if err := s.repo.Save(intent); err != nil {
		return err
	}
	s.bus.Publish(evt)
	return nil
}

func (s *PaymentService) Refund(orderID string, reason string) error {
	intent, ok := s.repo.GetByOrderID(orderID)
	if !ok {
		return nil
	}
	evt := intent.Refund(reason)
	if err := s.repo.Save(intent); err != nil {
		return err
	}
	s.bus.Publish(evt)
	return nil
}

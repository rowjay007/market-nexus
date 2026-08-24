package application

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/pricing/internal/domain"
)

type RuleRepository interface {
	Save(rule *domain.PriceRule) error
	GetByVendor(vendorID sharedkernel.VendorID) (*domain.PriceRule, bool)
}

type PricingService struct {
	repo RuleRepository
	bus  *events.InMemoryBus
}

func NewPricingService(repo RuleRepository, bus *events.InMemoryBus) *PricingService {
	return &PricingService{repo: repo, bus: bus}
}

func (s *PricingService) SeedRule(rule *domain.PriceRule) error {
	return s.repo.Save(rule)
}

func (s *PricingService) Quote(orderID string, vendorID sharedkernel.VendorID, subtotal int64) (int64, error) {
	rule, ok := s.repo.GetByVendor(vendorID)
	if !ok {
		rule = domain.NewPriceRule("default", vendorID, 0)
	}
	quoted, err := rule.Quote(orderID, subtotal)
	if err != nil {
		return 0, err
	}
	s.bus.Publish(quoted)
	return quoted.Total, nil
}

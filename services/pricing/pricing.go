package pricing

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/pricing/internal/application"
	"github.com/rowjay007/market-nexus/services/pricing/internal/domain"
	"github.com/rowjay007/market-nexus/services/pricing/internal/infrastructure/memory"
)

type Service struct {
	app *application.PricingService
}

func NewInMemoryService(bus *events.InMemoryBus) *Service {
	repo := memory.NewRuleRepo()
	return &Service{app: application.NewPricingService(repo, bus)}
}

func (s *Service) SeedRule(id string, vendorID sharedkernel.VendorID, discountBps int64) error {
	return s.app.SeedRule(domain.NewPriceRule(domain.PriceRuleID(id), vendorID, discountBps))
}

func (s *Service) Quote(orderID string, vendorID sharedkernel.VendorID, subtotal int64) (int64, error) {
	return s.app.Quote(orderID, vendorID, subtotal)
}

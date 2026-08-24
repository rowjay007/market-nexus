package analytics

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/analytics/internal/application"
	"github.com/rowjay007/market-nexus/services/analytics/internal/domain"
	"github.com/rowjay007/market-nexus/services/analytics/internal/infrastructure/memory"
)

type Service struct {
	app *application.AnalyticsService
}

func NewInMemoryService(bus *events.InMemoryBus) *Service {
	repo := memory.NewEventRepo()
	cache := domain.NewRecommendationCache()
	return &Service{app: application.NewAnalyticsService(repo, cache, bus)}
}

func (s *Service) RecordBehavior(id string, userID string, productID string, vendorID sharedkernel.VendorID, behavior string) error {
	return s.app.RecordBehavior(id, userID, productID, vendorID, domain.BehaviorType(behavior))
}

func (s *Service) Recommendations(userID string, limit int) []string {
	return s.app.Recommendations(userID, limit)
}

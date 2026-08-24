package application

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/analytics/internal/domain"
)

type EventRepository interface {
	Append(event *domain.BehaviorEvent) error
	ListByVendor(vendorID sharedkernel.VendorID) []*domain.BehaviorEvent
}

type AnalyticsService struct {
	repo  EventRepository
	cache *domain.RecommendationCache
	bus   *events.InMemoryBus
}

func NewAnalyticsService(repo EventRepository, cache *domain.RecommendationCache, bus *events.InMemoryBus) *AnalyticsService {
	return &AnalyticsService{repo: repo, cache: cache, bus: bus}
}

func (s *AnalyticsService) RecordBehavior(id string, userID string, productID string, vendorID sharedkernel.VendorID, behavior domain.BehaviorType) error {
	e, err := domain.NewBehaviorEvent(id, userID, productID, vendorID, behavior)
	if err != nil {
		return err
	}
	if err := s.repo.Append(e); err != nil {
		return err
	}
	s.cache.Apply(e)
	s.bus.Publish(e.RecordedEvent())
	s.bus.Publish(domain.RecommendationComputed(userID, s.cache.TopProducts(userID, 10)))
	return nil
}

func (s *AnalyticsService) Recommendations(userID string, limit int) []string {
	return s.cache.TopProducts(userID, limit)
}

package application

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/reviewtrust/internal/domain"
)

type ReviewRepository interface {
	Save(review *domain.Review) error
	ListByVendor(vendorID sharedkernel.VendorID) []*domain.Review
}

type ReviewTrustService struct {
	repo ReviewRepository
	bus  *events.InMemoryBus
}

func NewReviewTrustService(repo ReviewRepository, bus *events.InMemoryBus) *ReviewTrustService {
	return &ReviewTrustService{repo: repo, bus: bus}
}

func (s *ReviewTrustService) SubmitReview(review *domain.Review) error {
	if err := s.repo.Save(review); err != nil {
		return err
	}
	s.bus.Publish(review.SubmittedEvent())
	return nil
}

func (s *ReviewTrustService) OpenDispute(dispute *domain.Dispute) {
	s.bus.Publish(dispute.OpenedEvent())
}

func (s *ReviewTrustService) VendorAverageRating(vendorID sharedkernel.VendorID) float64 {
	items := s.repo.ListByVendor(vendorID)
	if len(items) == 0 {
		return 0
	}
	var total int
	for _, r := range items {
		total += r.Rating()
	}
	return float64(total) / float64(len(items))
}

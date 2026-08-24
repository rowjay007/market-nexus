package reviewtrust

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/reviewtrust/internal/application"
	"github.com/rowjay007/market-nexus/services/reviewtrust/internal/domain"
	"github.com/rowjay007/market-nexus/services/reviewtrust/internal/infrastructure/memory"
)

type Service struct {
	app *application.ReviewTrustService
}

func NewInMemoryService(bus *events.InMemoryBus) *Service {
	repo := memory.NewReviewRepo()
	return &Service{app: application.NewReviewTrustService(repo, bus)}
}

func (s *Service) SubmitReview(id string, productID string, vendorID sharedkernel.VendorID, buyerID string, rating int, comment string) error {
	review, err := domain.NewReview(domain.ReviewID(id), productID, vendorID, buyerID, rating, comment)
	if err != nil {
		return err
	}
	return s.app.SubmitReview(review)
}

func (s *Service) OpenDispute(id string, reviewID string, vendorID sharedkernel.VendorID, reason string, isFraud bool) {
	s.app.OpenDispute(domain.NewDispute(id, reviewID, vendorID, reason, isFraud))
}

func (s *Service) VendorAverageRating(vendorID sharedkernel.VendorID) float64 {
	return s.app.VendorAverageRating(vendorID)
}

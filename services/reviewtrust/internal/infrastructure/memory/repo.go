package memory

import (
	"sync"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/reviewtrust/internal/domain"
)

type ReviewRepo struct {
	mu       sync.RWMutex
	byVendor map[sharedkernel.VendorID][]*domain.Review
}

func NewReviewRepo() *ReviewRepo {
	return &ReviewRepo{byVendor: map[sharedkernel.VendorID][]*domain.Review{}}
}

func (r *ReviewRepo) Save(review *domain.Review) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byVendor[review.VendorID()] = append(r.byVendor[review.VendorID()], review)
	return nil
}

func (r *ReviewRepo) ListByVendor(vendorID sharedkernel.VendorID) []*domain.Review {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.byVendor[vendorID]
	out := make([]*domain.Review, len(items))
	copy(out, items)
	return out
}

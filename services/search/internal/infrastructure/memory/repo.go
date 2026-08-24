package memory

import (
	"sync"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/search/internal/domain"
)

type IndexRepo struct {
	mu      sync.RWMutex
	byVendor map[sharedkernel.VendorID][]*domain.SearchDocument
}

func NewIndexRepo() *IndexRepo {
	return &IndexRepo{byVendor: map[sharedkernel.VendorID][]*domain.SearchDocument{}}
}

func (r *IndexRepo) Save(doc *domain.SearchDocument) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byVendor[doc.VendorID()] = append(r.byVendor[doc.VendorID()], doc)
	return nil
}

func (r *IndexRepo) ListByVendor(vendorID sharedkernel.VendorID) []*domain.SearchDocument {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.byVendor[vendorID]
	out := make([]*domain.SearchDocument, len(items))
	copy(out, items)
	return out
}

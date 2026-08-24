package memory

import (
	"sync"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/catalog/internal/domain"
)

type ProductRepo struct {
	mu      sync.RWMutex
	byID    map[domain.ProductID]*domain.Product
	byVendor map[sharedkernel.VendorID][]*domain.Product
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{byID: map[domain.ProductID]*domain.Product{}, byVendor: map[sharedkernel.VendorID][]*domain.Product{}}
}

func (r *ProductRepo) Save(product *domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byID[product.ID()] = product
	r.byVendor[product.VendorID()] = append(r.byVendor[product.VendorID()], product)
	return nil
}

func (r *ProductRepo) GetByID(id domain.ProductID) (*domain.Product, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.byID[id]
	return p, ok
}

func (r *ProductRepo) ListByVendor(vendorID sharedkernel.VendorID) []*domain.Product {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.byVendor[vendorID]
	out := make([]*domain.Product, len(list))
	copy(out, list)
	return out
}

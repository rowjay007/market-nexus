package catalogreadmodel

import (
	"github.com/rowjay007/market-nexus/services/catalog"
	"sync"
)

type ProductSnapshot struct {
	ProductID string
	VendorID  string
	Name      string
	Variants  map[string]int64
}

type Projection struct {
	mu       sync.RWMutex
	products map[string]ProductSnapshot
}

func NewProjection() *Projection {
	return &Projection{products: map[string]ProductSnapshot{}}
}

func (p *Projection) HandleCatalogItemPublished(evt catalog.ItemPublished) {
	variants := map[string]int64{}
	for _, v := range evt.Variants {
		variants[v.SKU] = v.Price
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.products[evt.ProductID] = ProductSnapshot{
		ProductID: evt.ProductID,
		VendorID:  evt.VendorID,
		Name:      evt.Name,
		Variants:  variants,
	}
}

func (p *Projection) ResolvePrice(vendorID string, sku string) (int64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, product := range p.products {
		if product.VendorID != vendorID {
			continue
		}
		price, ok := product.Variants[sku]
		if ok {
			return price, true
		}
	}
	return 0, false
}

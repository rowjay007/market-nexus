package memory

import (
	"sync"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/inventory/internal/domain"
)

type key struct {
	vendorID sharedkernel.VendorID
	sku      string
}

type StockRepo struct {
	mu    sync.RWMutex
	items map[key]*domain.StockItem
}

func NewStockRepo() *StockRepo {
	return &StockRepo{items: map[key]*domain.StockItem{}}
}

func (r *StockRepo) Save(item *domain.StockItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[key{vendorID: item.VendorID(), sku: item.SKU()}] = item
	return nil
}

func (r *StockRepo) GetByVendorAndSKU(vendorID sharedkernel.VendorID, sku string) (*domain.StockItem, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[key{vendorID: vendorID, sku: sku}]
	return item, ok
}

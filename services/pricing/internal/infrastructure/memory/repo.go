package memory

import (
	"sync"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/pricing/internal/domain"
)

type RuleRepo struct {
	mu    sync.RWMutex
	items map[sharedkernel.VendorID]*domain.PriceRule
}

func NewRuleRepo() *RuleRepo {
	return &RuleRepo{items: map[sharedkernel.VendorID]*domain.PriceRule{}}
}

func (r *RuleRepo) Save(rule *domain.PriceRule) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[rule.VendorID()] = rule
	return nil
}

func (r *RuleRepo) GetByVendor(vendorID sharedkernel.VendorID) (*domain.PriceRule, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rule, ok := r.items[vendorID]
	return rule, ok
}

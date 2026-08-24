package memory

import (
	"sync"

	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/analytics/internal/domain"
)

type EventRepo struct {
	mu       sync.RWMutex
	byVendor map[sharedkernel.VendorID][]*domain.BehaviorEvent
}

func NewEventRepo() *EventRepo {
	return &EventRepo{byVendor: map[sharedkernel.VendorID][]*domain.BehaviorEvent{}}
}

func (r *EventRepo) Append(event *domain.BehaviorEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byVendor[event.VendorID()] = append(r.byVendor[event.VendorID()], event)
	return nil
}

func (r *EventRepo) ListByVendor(vendorID sharedkernel.VendorID) []*domain.BehaviorEvent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := r.byVendor[vendorID]
	out := make([]*domain.BehaviorEvent, len(items))
	copy(out, items)
	return out
}

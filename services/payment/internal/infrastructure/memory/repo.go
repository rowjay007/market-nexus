package memory

import (
	"sync"

	"github.com/rowjay007/market-nexus/services/payment/internal/domain"
)

type PaymentRepo struct {
	mu    sync.RWMutex
	items map[string]*domain.PaymentIntent
}

func NewPaymentRepo() *PaymentRepo {
	return &PaymentRepo{items: map[string]*domain.PaymentIntent{}}
}

func (r *PaymentRepo) Save(intent *domain.PaymentIntent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[intent.OrderID()] = intent
	return nil
}

func (r *PaymentRepo) GetByOrderID(orderID string) (*domain.PaymentIntent, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[orderID]
	return item, ok
}

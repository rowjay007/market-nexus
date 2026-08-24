package memory

import (
	orderingdomain "github.com/rowjay007/market-nexus/services/ordering/internal/domain"
	"sync"
)

type OrderRepo struct {
	mu    sync.RWMutex
	items map[orderingdomain.OrderID]*orderingdomain.Order
}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{items: map[orderingdomain.OrderID]*orderingdomain.Order{}}
}

func (r *OrderRepo) Save(order *orderingdomain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[order.ID()] = order
	return nil
}

func (r *OrderRepo) GetByID(id orderingdomain.OrderID) (*orderingdomain.Order, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[id]
	return item, ok
}

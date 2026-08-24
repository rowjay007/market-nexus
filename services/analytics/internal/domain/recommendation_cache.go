package domain

import (
	"sort"
	"time"

	"github.com/rowjay007/market-nexus/pkg/events"
)

type RecommendationCache struct {
	items map[string]map[string]int
}

func NewRecommendationCache() *RecommendationCache {
	return &RecommendationCache{items: map[string]map[string]int{}}
}

func (c *RecommendationCache) Apply(e *BehaviorEvent) {
	if _, ok := c.items[e.UserID()]; !ok {
		c.items[e.UserID()] = map[string]int{}
	}
	weight := 1
	if e.Behavior() == BehaviorClick {
		weight = 3
	}
	if e.Behavior() == BehaviorPurchase {
		weight = 8
	}
	c.items[e.UserID()][e.ProductID()] += weight
}

func (c *RecommendationCache) TopProducts(userID string, limit int) []string {
	if limit <= 0 {
		limit = 10
	}
	scores := c.items[userID]
	type pair struct {
		productID string
		score     int
	}
	arr := make([]pair, 0, len(scores))
	for pid, score := range scores {
		arr = append(arr, pair{productID: pid, score: score})
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].score == arr[j].score {
			return arr[i].productID < arr[j].productID
		}
		return arr[i].score > arr[j].score
	})
	if len(arr) > limit {
		arr = arr[:limit]
	}
	out := make([]string, 0, len(arr))
	for _, p := range arr {
		out = append(out, p.productID)
	}
	return out
}

func RecommendationComputed(userID string, products []string) RecommendationsComputed {
	return RecommendationsComputed{
		BaseEvent: events.BaseEvent{Type: "RecommendationsComputed", At: time.Now().UTC()},
		UserID:    userID,
		Products:  products,
	}
}

type RecommendationsComputed struct {
	events.BaseEvent
	UserID   string
	Products []string
}

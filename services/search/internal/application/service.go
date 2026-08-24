package application

import (
	"sort"

	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/search/internal/domain"
)

type IndexRepository interface {
	Save(doc *domain.SearchDocument) error
	ListByVendor(vendorID sharedkernel.VendorID) []*domain.SearchDocument
}

type SearchService struct {
	repo IndexRepository
	bus  *events.InMemoryBus
}

func NewSearchService(repo IndexRepository, bus *events.InMemoryBus) *SearchService {
	return &SearchService{repo: repo, bus: bus}
}

func (s *SearchService) Index(doc *domain.SearchDocument) error {
	if err := s.repo.Save(doc); err != nil {
		return err
	}
	s.bus.Publish(doc.PublishedEvent())
	return nil
}

func (s *SearchService) Query(vendorID sharedkernel.VendorID, q string) []string {
	docs := s.repo.ListByVendor(vendorID)
	filtered := make([]*domain.SearchDocument, 0, len(docs))
	for _, doc := range docs {
		if doc.Matches(q) {
			filtered = append(filtered, doc)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].RankScore() > filtered[j].RankScore()
	})
	out := make([]string, 0, len(filtered))
	for _, doc := range filtered {
		out = append(out, doc.ProductID())
	}
	return out
}

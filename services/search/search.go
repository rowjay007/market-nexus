package search

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/search/internal/application"
	"github.com/rowjay007/market-nexus/services/search/internal/domain"
	"github.com/rowjay007/market-nexus/services/search/internal/infrastructure/memory"
)

type Service struct {
	app *application.SearchService
}

func NewInMemoryService(bus *events.InMemoryBus) *Service {
	repo := memory.NewIndexRepo()
	return &Service{app: application.NewSearchService(repo, bus)}
}

func (s *Service) Index(id string, productID string, vendorID sharedkernel.VendorID, title string, body string, tier int) error {
	return s.app.Index(domain.NewSearchDocument(domain.DocumentID(id), productID, vendorID, title, body, tier))
}

func (s *Service) Query(vendorID sharedkernel.VendorID, q string) []string {
	return s.app.Query(vendorID, q)
}

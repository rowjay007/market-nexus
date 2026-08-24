package application

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/catalog/internal/domain"
)

type ProductRepository interface {
	Save(product *domain.Product) error
	GetByID(id domain.ProductID) (*domain.Product, bool)
	ListByVendor(vendorID sharedkernel.VendorID) []*domain.Product
}

type CatalogService struct {
	repo ProductRepository
	bus  *events.InMemoryBus
}

func NewCatalogService(repo ProductRepository, bus *events.InMemoryBus) *CatalogService {
	return &CatalogService{repo: repo, bus: bus}
}

func (s *CatalogService) PublishProduct(product *domain.Product) error {
	if err := s.repo.Save(product); err != nil {
		return err
	}
	s.bus.Publish(product.Publish())
	return nil
}

func (s *CatalogService) ListProducts(vendorID sharedkernel.VendorID) []*domain.Product {
	return s.repo.ListByVendor(vendorID)
}

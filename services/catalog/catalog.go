package catalog

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/catalog/internal/application"
	"github.com/rowjay007/market-nexus/services/catalog/internal/domain"
	"github.com/rowjay007/market-nexus/services/catalog/internal/infrastructure/memory"
)

type VariantInput struct {
	SKU   string
	Price int64
}

type ItemPublished struct {
	ProductID string
	VendorID  string
	Name      string
	Variants  []VariantInput
}

type Service struct {
	app *application.CatalogService
}

func NewInMemoryService(bus *events.InMemoryBus) *Service {
	repo := memory.NewProductRepo()
	return &Service{app: application.NewCatalogService(repo, bus)}
}

func (s *Service) PublishProduct(productID string, vendorID sharedkernel.VendorID, name string, variants []VariantInput) error {
	p, err := domain.NewProduct(domain.ProductID(productID), vendorID, name)
	if err != nil {
		return err
	}
	for _, v := range variants {
		if err := p.AddVariant(vendorID, domain.SKU(v.SKU), v.Price); err != nil {
			return err
		}
	}
	return s.app.PublishProduct(p)
}

func ToPublishedEvent(e events.Event) (ItemPublished, bool) {
	evt, ok := e.(domain.CatalogItemPublished)
	if !ok {
		return ItemPublished{}, false
	}
	variants := make([]VariantInput, 0, len(evt.Variants))
	for _, v := range evt.Variants {
		variants = append(variants, VariantInput{SKU: string(v.SKU), Price: v.Price})
	}
	return ItemPublished{ProductID: evt.ProductID, VendorID: evt.VendorID, Name: evt.Name, Variants: variants}, true
}

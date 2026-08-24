package domain

import (
	"errors"
	"time"

	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
)

var (
	ErrInvalidProductID = errors.New("invalid product id")
	ErrCrossVendorSKU  = errors.New("sku vendor mismatch")
)

type ProductID string
type SKU string

type Variant struct {
	SKU   SKU
	Price int64
}

type Product struct {
	id       ProductID
	vendorID sharedkernel.VendorID
	name     string
	variants []Variant
}

func NewProduct(id ProductID, vendorID sharedkernel.VendorID, name string) (*Product, error) {
	if id == "" {
		return nil, ErrInvalidProductID
	}
	return &Product{id: id, vendorID: vendorID, name: name, variants: []Variant{}}, nil
}

func (p *Product) AddVariant(vendorID sharedkernel.VendorID, sku SKU, price int64) error {
	if vendorID != p.vendorID {
		return ErrCrossVendorSKU
	}
	p.variants = append(p.variants, Variant{SKU: sku, Price: price})
	return nil
}

func (p *Product) Publish() CatalogItemPublished {
	return CatalogItemPublished{
		BaseEvent: events.BaseEvent{Type: "CatalogItemPublished", At: time.Now().UTC()},
		ProductID: string(p.id),
		VendorID:  p.vendorID.String(),
		Name:      p.name,
		Variants:  p.variants,
	}
}

func (p *Product) ID() ProductID {
	return p.id
}

func (p *Product) VendorID() sharedkernel.VendorID {
	return p.vendorID
}

func (p *Product) Name() string {
	return p.name
}

func (p *Product) Variants() []Variant {
	out := make([]Variant, len(p.variants))
	copy(out, p.variants)
	return out
}

type CatalogItemPublished struct {
	events.BaseEvent
	ProductID string
	VendorID  string
	Name      string
	Variants  []Variant
}

package saga

import (
	"errors"
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/catalog"
	"github.com/rowjay007/market-nexus/services/inventory"
	"github.com/rowjay007/market-nexus/services/ordering"
	"github.com/rowjay007/market-nexus/services/ordering/acl/catalogreadmodel"
	"testing"
)

type inventoryACLAdapter struct {
	service *inventory.Service
}

func (a inventoryACLAdapter) Reserve(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error {
	return a.service.Reserve(orderID, vendorID, sku, qty)
}

func (a inventoryACLAdapter) Release(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error {
	return a.service.Release(orderID, vendorID, sku, qty)
}

type failingInventoryACL struct {
	delegate ordering.InventoryACL
	failsOn  string
}

func (f failingInventoryACL) Reserve(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error {
	if sku == f.failsOn {
		return errors.New("forced reservation failure")
	}
	return f.delegate.Reserve(orderID, vendorID, sku, qty)
}

func (f failingInventoryACL) Release(orderID string, vendorID sharedkernel.VendorID, sku string, qty int) error {
	return f.delegate.Release(orderID, vendorID, sku, qty)
}

func TestOrderToInventorySaga_HappyPath(t *testing.T) {
	bus := events.NewInMemoryBus()
	vendor := sharedkernel.MustVendorID("v-22")

	catalogSvc := catalog.NewInMemoryService(bus)
	projection := catalogreadmodel.NewProjection()
	bus.Subscribe("CatalogItemPublished", func(e events.Event) {
		evt, ok := catalog.ToPublishedEvent(e)
		if !ok {
			return
		}
		projection.HandleCatalogItemPublished(evt)
	})

	err := catalogSvc.PublishProduct("p-22", vendor, "Headphones", []catalog.VariantInput{{SKU: "sku-22", Price: 4999}})
	if err != nil {
		t.Fatalf("catalog publish failed: %v", err)
	}

	inventorySvc := inventory.NewInMemoryService(bus)
	_ = inventorySvc.SeedStock("st-22", vendor, "sku-22", 5)

	orderingSvc := ordering.NewInMemoryService(bus, projection, inventoryACLAdapter{service: inventorySvc})
	status, err := orderingSvc.PlaceOrder("o-22", vendor, []ordering.LineInput{{SKU: "sku-22", Quantity: 3}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != "CONFIRMED" {
		t.Fatalf("expected CONFIRMED, got %s", status)
	}

	available, reserved, ok := inventorySvc.Stock(vendor, "sku-22")
	if !ok {
		t.Fatal("stock missing")
	}
	if available != 2 || reserved != 3 {
		t.Fatalf("unexpected stock: available=%d reserved=%d", available, reserved)
	}
}

func TestOrderToInventorySaga_Compensation(t *testing.T) {
	bus := events.NewInMemoryBus()
	vendor := sharedkernel.MustVendorID("v-33")

	catalogSvc := catalog.NewInMemoryService(bus)
	projection := catalogreadmodel.NewProjection()
	bus.Subscribe("CatalogItemPublished", func(e events.Event) {
		evt, ok := catalog.ToPublishedEvent(e)
		if !ok {
			return
		}
		projection.HandleCatalogItemPublished(evt)
	})

	err := catalogSvc.PublishProduct("p-33", vendor, "Bundle", []catalog.VariantInput{
		{SKU: "sku-ok", Price: 1000},
		{SKU: "sku-fail", Price: 2000},
	})
	if err != nil {
		t.Fatalf("catalog publish failed: %v", err)
	}

	inventorySvc := inventory.NewInMemoryService(bus)
	_ = inventorySvc.SeedStock("st-ok", vendor, "sku-ok", 10)
	_ = inventorySvc.SeedStock("st-fail", vendor, "sku-fail", 10)

	baseACL := inventoryACLAdapter{service: inventorySvc}
	orderingSvc := ordering.NewInMemoryService(bus, projection, failingInventoryACL{delegate: baseACL, failsOn: "sku-fail"})

	_, err = orderingSvc.PlaceOrder("o-33", vendor, []ordering.LineInput{
		{SKU: "sku-ok", Quantity: 3},
		{SKU: "sku-fail", Quantity: 1},
	})
	if err == nil {
		t.Fatal("expected reservation failure")
	}

	available, reserved, ok := inventorySvc.Stock(vendor, "sku-ok")
	if !ok {
		t.Fatal("sku-ok missing")
	}
	if available != 10 || reserved != 0 {
		t.Fatalf("expected compensation release for sku-ok, got available=%d reserved=%d", available, reserved)
	}
}

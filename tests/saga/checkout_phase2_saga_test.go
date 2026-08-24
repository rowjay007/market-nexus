package saga

import (
	"errors"
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"github.com/rowjay007/market-nexus/services/catalog"
	"github.com/rowjay007/market-nexus/services/fulfillment"
	"github.com/rowjay007/market-nexus/services/inventory"
	"github.com/rowjay007/market-nexus/services/ordering"
	"github.com/rowjay007/market-nexus/services/ordering/acl/catalogreadmodel"
	"github.com/rowjay007/market-nexus/services/payment"
	"github.com/rowjay007/market-nexus/services/pricing"
	"testing"
)

type pricingACLAdapter struct {
	service *pricing.Service
}

func (a pricingACLAdapter) Quote(orderID string, vendorID sharedkernel.VendorID, subtotal int64) (int64, error) {
	return a.service.Quote(orderID, vendorID, subtotal)
}

type paymentACLAdapter struct {
	service *payment.Service
}

func (a paymentACLAdapter) Capture(orderID string, vendorID sharedkernel.VendorID, amount int64) error {
	return a.service.Capture(orderID, vendorID, amount)
}

func (a paymentACLAdapter) Refund(orderID string, reason string) error {
	return a.service.Refund(orderID, reason)
}

type fulfillmentACLAdapter struct {
	service *fulfillment.Service
}

func (a fulfillmentACLAdapter) Schedule(orderID string, vendorID sharedkernel.VendorID, address string) error {
	return a.service.Schedule(orderID, vendorID, address)
}

func (a fulfillmentACLAdapter) Cancel(orderID string, reason string) error {
	return a.service.Cancel(orderID, reason)
}

type failingFulfillmentACL struct {
	delegate ordering.FulfillmentACL
}

func (f failingFulfillmentACL) Schedule(orderID string, vendorID sharedkernel.VendorID, address string) error {
	return errors.New("forced fulfillment failure")
}

func (f failingFulfillmentACL) Cancel(orderID string, reason string) error {
	return f.delegate.Cancel(orderID, reason)
}

func setupPhase2(t *testing.T) (*ordering.Service, *inventory.Service, ordering.PricingACL, ordering.PaymentACL, ordering.FulfillmentACL, sharedkernel.VendorID) {
	t.Helper()
	bus := events.NewInMemoryBus()
	vendor := sharedkernel.MustVendorID("v-phase2")

	catalogSvc := catalog.NewInMemoryService(bus)
	projection := catalogreadmodel.NewProjection()
	bus.Subscribe("CatalogItemPublished", func(e events.Event) {
		evt, ok := catalog.ToPublishedEvent(e)
		if !ok {
			return
		}
		projection.HandleCatalogItemPublished(evt)
	})

	err := catalogSvc.PublishProduct("p-phase2", vendor, "Console", []catalog.VariantInput{{SKU: "sku-phase2", Price: 30000}})
	if err != nil {
		t.Fatalf("catalog publish failed: %v", err)
	}

	inventorySvc := inventory.NewInMemoryService(bus)
	_ = inventorySvc.SeedStock("st-phase2", vendor, "sku-phase2", 10)

	pricingSvc := pricing.NewInMemoryService(bus)
	_ = pricingSvc.SeedRule("rule-phase2", vendor, 500)

	paymentSvc := payment.NewInMemoryService(bus)
	fulfillmentSvc := fulfillment.NewInMemoryService(bus)
	orderingSvc := ordering.NewInMemoryService(bus, projection, inventoryACLAdapter{service: inventorySvc})

	return orderingSvc, inventorySvc, pricingACLAdapter{service: pricingSvc}, paymentACLAdapter{service: paymentSvc}, fulfillmentACLAdapter{service: fulfillmentSvc}, vendor
}

func TestCheckoutPhase2_HappyPath(t *testing.T) {
	orderingSvc, inventorySvc, pricingACL, paymentACL, fulfillmentACL, vendor := setupPhase2(t)

	status, err := orderingSvc.CheckoutOrder(
		"o-phase2-1",
		vendor,
		[]ordering.LineInput{{SKU: "sku-phase2", Quantity: 2}},
		pricingACL,
		paymentACL,
		fulfillmentACL,
		"1 Main Street",
	)
	if err != nil {
		t.Fatalf("checkout failed: %v", err)
	}
	if status != "CONFIRMED" {
		t.Fatalf("expected CONFIRMED, got %s", status)
	}

	available, reserved, ok := inventorySvc.Stock(vendor, "sku-phase2")
	if !ok {
		t.Fatal("stock missing")
	}
	if available != 8 || reserved != 2 {
		t.Fatalf("unexpected inventory after checkout: available=%d reserved=%d", available, reserved)
	}
}

func TestCheckoutPhase2_CompensationOnFulfillmentFailure(t *testing.T) {
	orderingSvc, inventorySvc, pricingACL, paymentACL, baseFulfillmentACL, vendor := setupPhase2(t)
	failing := failingFulfillmentACL{delegate: baseFulfillmentACL}

	_, err := orderingSvc.CheckoutOrder(
		"o-phase2-2",
		vendor,
		[]ordering.LineInput{{SKU: "sku-phase2", Quantity: 2}},
		pricingACL,
		paymentACL,
		failing,
		"1 Main Street",
	)
	if err == nil {
		t.Fatal("expected fulfillment failure")
	}

	available, reserved, ok := inventorySvc.Stock(vendor, "sku-phase2")
	if !ok {
		t.Fatal("stock missing")
	}
	if available != 10 || reserved != 0 {
		t.Fatalf("expected compensation release, got available=%d reserved=%d", available, reserved)
	}
}

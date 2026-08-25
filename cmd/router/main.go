package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	catalogURL, _ := url.Parse("http://catalog-service:8080/graphql")
	inventoryURL, _ := url.Parse("http://inventory-service:8080/graphql")
	orderingURL, _ := url.Parse("http://ordering-service:8080/graphql")
	pricingURL, _ := url.Parse("http://pricing-service:8080/graphql")
	paymentURL, _ := url.Parse("http://payment-service:8080/graphql")
	fulfillmentURL, _ := url.Parse("http://fulfillment-service:8080/graphql")
	searchURL, _ := url.Parse("http://search-service:8080/graphql")
	reviewURL, _ := url.Parse("http://reviewtrust-service:8080/graphql")
	analyticsURL, _ := url.Parse("http://analytics-service:8080/graphql")

	mux := http.NewServeMux()
	mux.Handle("/subgraph/catalog", httputil.NewSingleHostReverseProxy(catalogURL))
	mux.Handle("/subgraph/inventory", httputil.NewSingleHostReverseProxy(inventoryURL))
	mux.Handle("/subgraph/ordering", httputil.NewSingleHostReverseProxy(orderingURL))
	mux.Handle("/subgraph/pricing", httputil.NewSingleHostReverseProxy(pricingURL))
	mux.Handle("/subgraph/payment", httputil.NewSingleHostReverseProxy(paymentURL))
	mux.Handle("/subgraph/fulfillment", httputil.NewSingleHostReverseProxy(fulfillmentURL))
	mux.Handle("/subgraph/search", httputil.NewSingleHostReverseProxy(searchURL))
	mux.Handle("/subgraph/reviewtrust", httputil.NewSingleHostReverseProxy(reviewURL))
	mux.Handle("/subgraph/analytics", httputil.NewSingleHostReverseProxy(analyticsURL))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	log.Println("router listening on :4000")
	if err := http.ListenAndServe(":4000", mux); err != nil {
		log.Fatal(err)
	}
}

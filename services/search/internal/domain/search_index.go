package domain

import (
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"strings"
	"time"
)

type DocumentID string

type SearchDocument struct {
	id        DocumentID
	productID string
	vendorID  sharedkernel.VendorID
	title     string
	body      string
	tier      int
}

func NewSearchDocument(id DocumentID, productID string, vendorID sharedkernel.VendorID, title string, body string, tier int) *SearchDocument {
	if tier < 0 {
		tier = 0
	}
	return &SearchDocument{id: id, productID: productID, vendorID: vendorID, title: title, body: body, tier: tier}
}

func (d *SearchDocument) Matches(q string) bool {
	query := strings.ToLower(strings.TrimSpace(q))
	if query == "" {
		return true
	}
	text := strings.ToLower(d.title + " " + d.body)
	return strings.Contains(text, query)
}

func (d *SearchDocument) RankScore() int {
	return len(d.title) + (d.tier * 100)
}

func (d *SearchDocument) ProductID() string               { return d.productID }
func (d *SearchDocument) VendorID() sharedkernel.VendorID { return d.vendorID }
func (d *SearchDocument) Title() string                   { return d.title }

func (d *SearchDocument) PublishedEvent() SearchIndexed {
	return SearchIndexed{
		BaseEvent: events.BaseEvent{Type: "SearchIndexed", At: time.Now().UTC()},
		ProductID: d.productID,
		VendorID:  d.vendorID.String(),
		Title:     d.title,
	}
}

type SearchIndexed struct {
	events.BaseEvent
	ProductID string
	VendorID  string
	Title     string
}

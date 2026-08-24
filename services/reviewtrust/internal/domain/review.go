package domain

import (
	"errors"
	"github.com/rowjay007/market-nexus/pkg/events"
	"github.com/rowjay007/market-nexus/pkg/sharedkernel"
	"time"
)

var ErrInvalidRating = errors.New("invalid rating")

type ReviewID string

type Review struct {
	id        ReviewID
	productID string
	vendorID  sharedkernel.VendorID
	buyerID   string
	rating    int
	comment   string
}

func NewReview(id ReviewID, productID string, vendorID sharedkernel.VendorID, buyerID string, rating int, comment string) (*Review, error) {
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}
	return &Review{id: id, productID: productID, vendorID: vendorID, buyerID: buyerID, rating: rating, comment: comment}, nil
}

func (r *Review) SubmittedEvent() ReviewSubmitted {
	return ReviewSubmitted{
		BaseEvent: events.BaseEvent{Type: "ReviewSubmitted", At: time.Now().UTC()},
		ReviewID:  string(r.id),
		ProductID: r.productID,
		VendorID:  r.vendorID.String(),
		BuyerID:   r.buyerID,
		Rating:    r.rating,
	}
}

type Dispute struct {
	id       string
	reviewID string
	vendorID sharedkernel.VendorID
	reason   string
	isFraud  bool
}

func NewDispute(id string, reviewID string, vendorID sharedkernel.VendorID, reason string, isFraud bool) *Dispute {
	return &Dispute{id: id, reviewID: reviewID, vendorID: vendorID, reason: reason, isFraud: isFraud}
}

func (d *Dispute) OpenedEvent() DisputeOpened {
	return DisputeOpened{
		BaseEvent: events.BaseEvent{Type: "DisputeOpened", At: time.Now().UTC()},
		DisputeID: d.id,
		ReviewID:  d.reviewID,
		VendorID:  d.vendorID.String(),
		Reason:    d.reason,
		IsFraud:   d.isFraud,
	}
}

func (r *Review) ProductID() string               { return r.productID }
func (r *Review) VendorID() sharedkernel.VendorID { return r.vendorID }
func (r *Review) Rating() int                     { return r.rating }

type ReviewSubmitted struct {
	events.BaseEvent
	ReviewID  string
	ProductID string
	VendorID  string
	BuyerID   string
	Rating    int
}

type DisputeOpened struct {
	events.BaseEvent
	DisputeID string
	ReviewID  string
	VendorID  string
	Reason    string
	IsFraud   bool
}

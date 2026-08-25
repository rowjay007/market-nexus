package stripe

import (
	"errors"
	stripe "github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"
)

type Client struct{}

func New(secretKey string) (*Client, error) {
	if secretKey == "" {
		return nil, errors.New("stripe secret key required")
	}
	stripe.Key = secretKey
	return &Client{}, nil
}

func (c *Client) CapturePaymentIntent(id string) (*stripe.PaymentIntent, error) {
	if id == "" {
		return nil, errors.New("payment intent id required")
	}
	return paymentintent.Capture(id, &stripe.PaymentIntentCaptureParams{})
}

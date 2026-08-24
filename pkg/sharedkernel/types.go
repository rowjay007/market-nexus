package sharedkernel

import (
	"errors"
	"strings"
)

var (
	ErrInvalidVendorID = errors.New("invalid vendor id")
)

type VendorID string

func NewVendorID(value string) (VendorID, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", ErrInvalidVendorID
	}
	return VendorID(trimmed), nil
}

func MustVendorID(value string) VendorID {
	v, err := NewVendorID(value)
	if err != nil {
		panic(err)
	}
	return v
}

func (v VendorID) String() string {
	return string(v)
}

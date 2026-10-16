package paymentcryptography

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type PaymentCryptography struct {
	Keys []Key `tree:"keys"`
}

type Key struct {
	resource.Resource `tree:"-"`
	Enabled           value.Bool `tree:"enabled"`
}

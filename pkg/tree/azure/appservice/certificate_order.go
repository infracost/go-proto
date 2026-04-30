package appservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CertificateOrder struct {
	resource.Resource `tree:"-"`
	ProductType       value.Value[CertificateProductType] `tree:"product_type"`
}

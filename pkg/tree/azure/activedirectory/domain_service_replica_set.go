package activedirectory

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type DomainServiceReplicaSet struct {
	resource.Resource `tree:"-"`
	SKU               value.String `tree:"sku"`
	DomainServiceID   value.String `tree:"domain_service_id"`
}

package appsync

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type APICache struct {
	resource.Resource        `tree:"-"`
	AtRestEncryptionEnabled  value.Bool `tree:"at_rest_encryption_enabled"`
	TransitEncryptionEnabled value.Bool `tree:"transit_encryption_enabled"`
}

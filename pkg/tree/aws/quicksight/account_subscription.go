package quicksight

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type AccountSubscription struct {
	resource.Resource `tree:"-"`
	Edition           value.Value[Edition] `tree:"edition"`
}

type Edition uint32

const (
	EditionUnknown Edition = iota
	EditionStandard
	EditionEnterprise
	EditionEnterpriseAndQ
)

package opensearchserverless

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

// Collection costs are usage-based: OCUs are shared by collections with the same KMS key.
type Collection struct {
	resource.Resource `tree:"-"`
}

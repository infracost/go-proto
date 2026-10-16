package backup

import "github.com/infracost/go-proto/pkg/tree/resource"

// Framework is an AWS Backup Audit Manager framework.
type Framework struct {
	resource.Resource `tree:"-"`
}

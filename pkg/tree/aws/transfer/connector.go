package transfer

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Connector struct {
	resource.Resource `tree:"-"`
	SFTP              value.Bool `tree:"sftp"`
}

type WebApp struct {
	resource.Resource `tree:"-"`
	Units             value.Int `tree:"units"`
}

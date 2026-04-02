package dms

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Endpoint struct {
	resource.Resource `tree:"-"`
	SslMode           value.Value[SslMode] `tree:"ssl_mode"`
}

type SslMode uint32

const (
	SslModeNone      SslMode = 0
	SslModeRequire   SslMode = 1
	SslModeVerifyCA  SslMode = 2
	SslModeVerifyFull SslMode = 3
)

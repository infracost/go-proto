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
	SslModeUnknown SslMode = iota
	SslModeNone
	SslModeRequire
	SslModeVerifyCA
	SslModeVerifyFull
)

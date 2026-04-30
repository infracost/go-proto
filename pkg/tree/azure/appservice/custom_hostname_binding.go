package appservice

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type CustomHostnameBinding struct {
	resource.Resource `tree:"-"`
	SSLState          value.Value[SSLState] `tree:"ssl_state"`
}

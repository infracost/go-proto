package sagemaker

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type MLflowTrackingServer struct {
	resource.Resource  `tree:"-"`
	TrackingServerSize value.String `tree:"tracking_server_size"`
}

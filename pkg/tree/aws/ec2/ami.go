package ec2

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

// AMI covers aws_ami_copy and aws_ami_from_instance, which create new EBS
// snapshots. Snapshot size is not known from config.
type AMI struct {
	resource.Resource `tree:"-"`
}

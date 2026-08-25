package karpenter

import (
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// EC2NodeClass is a karpenter.k8s.aws EC2NodeClass — the AWS-specific half of a
// Karpenter node definition, holding what a NodePool's constraints do not: the
// AMI, the disks, the networking and the instance profile.
//
// The disks are the cost that is easy to miss. A NodePool bounds instance
// types, and instance cost is what a reader thinks of, but every node this
// class launches also gets the EBS volumes in BlockDeviceMappings — charged per
// provisioned GB-month whether or not they are used, plus provisioned IOPS and
// throughput on gp3 and io2. On a large, churning Karpenter fleet that is real
// money attached to a field nobody edits after the first week.
//
// EC2NodeClass is cluster-scoped, so it has no namespace of its own; see the
// note on NodePool.
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource; the class's own name on the embedded
// meta.ObjectMeta; and its Kubernetes labels are stored as the base resource's
// Tags.
type EC2NodeClass struct {
	resource.Resource `tree:"-"`
	meta.ObjectMeta   `tree:"-"`

	// AMIFamily is spec.amiFamily — AL2023, Bottlerocket, Windows2022 and so
	// on. It decides the node's OS, which decides its licensing: a Windows node
	// carries a per-core charge on top of the instance rate that no Linux
	// family does.
	AMIFamily value.String `tree:"ami_family"`

	// BlockDeviceMappings is spec.blockDeviceMappings — the EBS volumes
	// attached to every node this class launches. Empty means the AMI family's
	// own defaults apply, which are not stated here, so an empty list is not
	// the same as no disks.
	BlockDeviceMappings []BlockDeviceMapping `tree:"block_device_mappings"`

	// SubnetSelectorTerms and SecurityGroupSelectorTerms are
	// spec.subnetSelectorTerms and spec.securityGroupSelectorTerms — how the
	// class discovers where to place nodes.
	//
	// The subnets matter for cost beyond placement: which availability zones
	// they span decides whether traffic between pods crosses a zone boundary
	// and picks up cross-AZ data transfer charges.
	SubnetSelectorTerms        []SelectorTerm `tree:"subnet_selector_terms"`
	SecurityGroupSelectorTerms []SelectorTerm `tree:"security_group_selector_terms"`

	// Annotations are the class's Kubernetes annotations, surfaced verbatim.
	Annotations []resource.Tag `tree:"annotations"`
}

// SelectorTerm is one entry of an EC2NodeClass's subnet or security-group
// selector terms — a discovery rule rather than a fixed reference.
//
// Terms resolve against live AWS state, so what a term actually selects is not
// knowable from the repository. Both forms are kept: Tags for the tag-matching
// form, ID for a direct reference.
type SelectorTerm struct {
	// Tags is the tag key/value set a subnet or security group must carry to
	// match. The common form, usually keyed on a cluster-discovery tag.
	Tags []resource.Tag `tree:"tags"`

	// ID is a direct subnet or security-group id, used instead of tags when a
	// term names one outright.
	ID value.String `tree:"id"`
}

// BlockDeviceMapping is one entry of an EC2NodeClass's spec.blockDeviceMappings
// — one EBS volume attached to every node the class launches.
type BlockDeviceMapping struct {
	// DeviceName is the device the volume attaches at, e.g. "/dev/xvda".
	DeviceName value.String `tree:"device_name"`

	// VolumeSizeBytes is ebs.volumeSize, in bytes.
	//
	// The manifest states it as a Kubernetes-style quantity ("100Gi", "500G"),
	// which is reduced to bytes here for the same reason container memory is:
	// an exact integer in a single base unit, with no GiB-versus-GB ambiguity
	// left for a consumer to guess at. The provider plugin converts to the GB
	// its rate is quoted in.
	VolumeSizeBytes value.Int `tree:"volume_size_bytes"`

	// VolumeType is ebs.volumeType — gp3, gp2, io2 and so on. The primary
	// determinant of the per-GB rate.
	VolumeType value.String `tree:"volume_type"`

	// IOPS and Throughput are ebs.iops and ebs.throughput. Both are billable
	// above the baseline included with the volume type, and both are unset far
	// more often than not — unset means the type's baseline rather than zero.
	IOPS       value.Int `tree:"iops"`
	Throughput value.Int `tree:"throughput"`

	// Encrypted is ebs.encrypted. Not itself billable, but it is a tagging and
	// compliance signal of the same kind the annotations carry.
	Encrypted value.Bool `tree:"encrypted"`

	// DeleteOnTermination is ebs.deleteOnTermination. When false, the volume
	// outlives the node that created it and keeps being charged — under
	// Karpenter, whose nodes are replaced continuously, that accumulates
	// orphaned volumes rather than leaving one behind.
	DeleteOnTermination value.Bool `tree:"delete_on_termination"`
}

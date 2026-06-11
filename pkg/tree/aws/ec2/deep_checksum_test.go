package ec2

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/proto/gen/go/infracost/parser"
	"github.com/stretchr/testify/assert"
)

// minimalDefinition returns a Definition with a non-nil CallStack, since
// ToProviderResource dereferences base.Definition.CallStack.Frames.
func minimalDefinition() resource.Definition {
	return resource.Definition{CallStack: &parser.CallStack{}}
}

// Fixtures captured from the canonical implementation in
// github.com/infracost/providers/pkg/processor.CalculateDeepResourceChecksum
// for shapes with no linked resources. The walker in this repo must produce
// the same values byte-for-byte so existing DeepChecksums in production
// don't change for bare resources.
const (
	bareInstanceFixture = "ed89ed84468e3ffd"
	bareSubnetFixture   = "a51a65aac4f58be1"
)

func TestCalculateDeepChecksum_BareInstance_MatchesProviders(t *testing.T) {
	i := &Instance{Resource: resource.Resource{ID: "i-1", FullChecksum: "i1-fc"}}
	assert.Equal(t, bareInstanceFixture, resource.CalculateDeepChecksum(i))
}

func TestCalculateDeepChecksum_BareSubnet_MatchesProviders(t *testing.T) {
	s := &Subnet{Resource: resource.Resource{ID: "sn-1", FullChecksum: "sn1-fc"}}
	assert.Equal(t, bareSubnetFixture, resource.CalculateDeepChecksum(s))
}

// Instance.Relationships.LaunchTemplate is a pointer-to-Impl nested inside a
// plain wrapper struct. The old providers walker couldn't reach it (skipped
// the wrapper struct AND had a pointer-receiver bug). The new walker should
// reach it, so different LaunchTemplates must produce different parent
// checksums.
func TestCalculateDeepChecksum_WrapperStructPointer_IsTraversed(t *testing.T) {
	bare := &Instance{Resource: resource.Resource{ID: "i-1", FullChecksum: "i1-fc"}}
	withLT1 := &Instance{
		Resource:      resource.Resource{ID: "i-1", FullChecksum: "i1-fc"},
		Relationships: InstanceRelationships{LaunchTemplate: &LaunchTemplate{Resource: resource.Resource{ID: "lt-1", FullChecksum: "lt1-fc"}}},
	}
	withLT2 := &Instance{
		Resource:      resource.Resource{ID: "i-1", FullChecksum: "i1-fc"},
		Relationships: InstanceRelationships{LaunchTemplate: &LaunchTemplate{Resource: resource.Resource{ID: "lt-2", FullChecksum: "lt2-fc"}}},
	}

	bareCS := resource.CalculateDeepChecksum(bare)
	lt1CS := resource.CalculateDeepChecksum(withLT1)
	lt2CS := resource.CalculateDeepChecksum(withLT2)

	assert.NotEqual(t, bareCS, lt1CS, "linked LaunchTemplate should change the parent checksum")
	assert.NotEqual(t, lt1CS, lt2CS, "different LaunchTemplates should produce different parent checksums")
}

// Subnet.Relationships.NATGateways is a slice-of-pointer-to-Impl nested
// inside a plain wrapper struct. The old providers walker missed it (didn't
// descend through the wrapper). The new walker should reach it, and order
// of elements should not affect the result (sorted internally).
func TestCalculateDeepChecksum_WrapperStructSlice_IsTraversed(t *testing.T) {
	bare := &Subnet{Resource: resource.Resource{ID: "sn-1", FullChecksum: "sn1-fc"}}
	ngw1 := &NATGateway{Resource: resource.Resource{ID: "ngw-1", FullChecksum: "ngw1-fc"}}
	ngw2 := &NATGateway{Resource: resource.Resource{ID: "ngw-2", FullChecksum: "ngw2-fc"}}

	withGateways := &Subnet{
		Resource:      resource.Resource{ID: "sn-1", FullChecksum: "sn1-fc"},
		Relationships: SubnetRelationships{NATGateways: []*NATGateway{ngw1, ngw2}},
	}
	withGatewaysReversed := &Subnet{
		Resource:      resource.Resource{ID: "sn-1", FullChecksum: "sn1-fc"},
		Relationships: SubnetRelationships{NATGateways: []*NATGateway{ngw2, ngw1}},
	}

	bareCS := resource.CalculateDeepChecksum(bare)
	withCS := resource.CalculateDeepChecksum(withGateways)
	reversedCS := resource.CalculateDeepChecksum(withGatewaysReversed)

	assert.NotEqual(t, bareCS, withCS, "linked NATGateways should change the parent checksum")
	assert.Equal(t, withCS, reversedCS, "slice order should not affect the checksum (sorted internally)")
}

// ToProviderResource was previously a method on *Resource and called the
// walker on the base struct, which lost all wrapper context and produced
// hash(FullChecksum) regardless of linked resources. Now it takes the
// wrapping Implementation, so the DeepChecksum on the produced provider
// resource actually reflects the linked LaunchTemplate.
func TestToProviderResource_DeepChecksum_IncludesWrapperStructLinks(t *testing.T) {
	bare := &Instance{Resource: resource.Resource{ID: "i-1", FullChecksum: "i1-fc", Definition: minimalDefinition()}}
	withLT := &Instance{
		Resource:      resource.Resource{ID: "i-1", FullChecksum: "i1-fc", Definition: minimalDefinition()},
		Relationships: InstanceRelationships{LaunchTemplate: &LaunchTemplate{Resource: resource.Resource{ID: "lt-1", FullChecksum: "lt1-fc"}}},
	}

	bareCS := resource.ToProviderResource(bare).Metadata.DeepChecksum
	withLTCS := resource.ToProviderResource(withLT).Metadata.DeepChecksum

	assert.NotEqual(t, bareCS, withLTCS, "DeepChecksum on the produced provider.Resource must change when a linked resource is added through a wrapper struct")
}

func TestCalculateDeepChecksum_Deterministic(t *testing.T) {
	i := &Instance{
		Resource: resource.Resource{ID: "i-1", FullChecksum: "i1-fc"},
		Relationships: InstanceRelationships{
			LaunchTemplate: &LaunchTemplate{Resource: resource.Resource{ID: "lt-1", FullChecksum: "lt1-fc"}},
		},
	}
	assert.Equal(t, resource.CalculateDeepChecksum(i), resource.CalculateDeepChecksum(i))
}

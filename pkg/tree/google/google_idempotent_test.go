package google

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infracost/go-proto/pkg/tree/google/compute"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

func TestGooglePostProcess_IsIdempotent(t *testing.T) {
	g := &Google{
		Compute: compute.Compute{
			Addresses: []compute.Address{
				{
					Resource: resource.Resource{ID: "addr-1"},
					Address:  value.New("10.0.0.1", 0, "", nil),
				},
			},
			Instances: []compute.Instance{
				{
					Resource: resource.Resource{ID: "inst-1"},
					NATIP:    value.New("10.0.0.1", 0, "", nil),
				},
			},
		},
	}

	g.PostProcess()
	inst := g.Compute.Addresses[0].Relationships.Instance

	g.PostProcess()
	assert.Equal(t, inst, g.Compute.Addresses[0].Relationships.Instance)
	assert.NotNil(t, inst)
}

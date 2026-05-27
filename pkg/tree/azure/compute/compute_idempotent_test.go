package compute

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	c := &Compute{
		ManagedDisks: []ManagedDisk{
			{Resource: resource.Resource{ID: "md-1"}},
		},
		VirtualMachines: []VirtualMachine{
			{Resource: resource.Resource{ID: "vm-1"}},
		},
		DiskAttachments: []DiskAttachment{
			{
				Resource:         resource.Resource{ID: "da-1"},
				ManagedDiskID:    value.New("md-1", 0, "", nil),
				VirtualMachineID: value.New("vm-1", 0, "", nil),
			},
		},
	}

	c.PostProcess()
	vm := c.ManagedDisks[0].Relationships.VirtualMachine

	c.PostProcess()
	assert.Equal(t, vm, c.ManagedDisks[0].Relationships.VirtualMachine)
	assert.NotNil(t, vm)
}

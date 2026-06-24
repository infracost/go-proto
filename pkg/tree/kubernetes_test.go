package tree

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/kubernetes"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/apps"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/batch"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/workload"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKubernetesRoundTrip exercises the embedded-struct flattening: a Deployment
// embeds workload.Workload, and both its base fields (containers, labels,
// annotations) and its own field (replicas) must survive ToProto -> FromProto at
// the same level.
func TestKubernetesRoundTrip(t *testing.T) {
	original := &Tree{
		Kubernetes: kubernetes.Kubernetes{
			Apps: apps.Apps{
				Deployments: []apps.Deployment{
					{
						Workload: workload.Workload{
							Resource: resource.Resource{
								ID:     "api",
								Region: "us-east-1",
								// Kubernetes labels are stored as the base resource's Tags.
								Tags: resource.Tags{
									{Key: value.New("app", 0, "", nil), Value: value.New("api", 0, "", nil)},
								},
							},
							Annotations: []resource.Tag{
								{
									Key:   value.New("eks.amazonaws.com/role-arn", 0, "", nil),
									Value: value.New("arn:aws:iam::123456789012:role/api", 0, "", nil),
								},
							},
							Containers: []workload.Container{
								{
									Name:                 value.New("api", 0, "", nil),
									CPURequestMillicores: value.New[int64](500, 0, "", nil),
									MemoryRequestBytes:   value.New[int64](536870912, 0, "", nil),
								},
							},
						},
						Replicas: value.New[int64](3, 0, "", nil),
					},
				},
				DaemonSets: []workload.Workload{
					{
						Resource: resource.Resource{ID: "node-agent"},
						Containers: []workload.Container{
							{Name: value.New("agent", 0, "", nil)},
						},
					},
				},
			},
			Batch: batch.Batch{
				CronJobs: []batch.CronJob{
					{
						Job: batch.Job{
							Workload:    workload.Workload{Resource: resource.Resource{ID: "report"}},
							Completions: value.New[int64](5, 0, "", nil),
							Parallelism: value.New[int64](2, 0, "", nil),
						},
						Schedule: value.New("0 * * * *", 0, "", nil),
					},
				},
			},
		},
	}

	origDep := original.Kubernetes.Apps.Deployments[0]
	origDaemon := original.Kubernetes.Apps.DaemonSets[0]
	origCron := original.Kubernetes.Batch.CronJobs[0]

	proto, err := original.ToProto()
	require.NoError(t, err)

	// The Deployment's base fields and its own replicas are flattened into a
	// single attribute object — not nested under an embedded-struct key.
	depAttrs := proto.Providers["kubernetes"].Services["apps"].Resources[0].Attributes.Entries
	assert.Equal(t, origDep.Replicas.Value(), depAttrs["replicas"].GetIntValue())
	require.NotNil(t, depAttrs["containers"], "embedded base fields must be flattened")
	require.NotNil(t, depAttrs["annotations"], "embedded base fields must be flattened")

	result, err := FromProto(proto)
	require.NoError(t, err)

	require.Len(t, result.Kubernetes.Apps.Deployments, 1)
	dep := result.Kubernetes.Apps.Deployments[0]
	assert.Equal(t, origDep.ID, dep.ID)
	assert.Equal(t, origDep.Region, dep.Region)
	assert.Equal(t, origDep.Replicas.Value(), dep.Replicas.Value())
	require.Len(t, dep.Tags, len(origDep.Tags))
	assert.Equal(t, origDep.Tags[0].Key.Value(), dep.Tags[0].Key.Value())
	assert.Equal(t, origDep.Tags[0].Value.Value(), dep.Tags[0].Value.Value())
	require.Len(t, dep.Annotations, len(origDep.Annotations))
	assert.Equal(t, origDep.Annotations[0].Key.Value(), dep.Annotations[0].Key.Value())
	assert.Equal(t, origDep.Annotations[0].Value.Value(), dep.Annotations[0].Value.Value())
	require.Len(t, dep.Containers, len(origDep.Containers))
	assert.Equal(t, origDep.Containers[0].CPURequestMillicores.Value(), dep.Containers[0].CPURequestMillicores.Value())
	assert.Equal(t, origDep.Containers[0].MemoryRequestBytes.Value(), dep.Containers[0].MemoryRequestBytes.Value())

	require.Len(t, result.Kubernetes.Apps.DaemonSets, 1)
	daemon := result.Kubernetes.Apps.DaemonSets[0]
	assert.Equal(t, origDaemon.ID, daemon.ID)
	require.Len(t, daemon.Containers, len(origDaemon.Containers))
	assert.Equal(t, origDaemon.Containers[0].Name.Value(), daemon.Containers[0].Name.Value())

	// Multi-level embedding: CronJob -> Job -> Workload all flatten together.
	require.Len(t, result.Kubernetes.Batch.CronJobs, 1)
	cj := result.Kubernetes.Batch.CronJobs[0]
	assert.Equal(t, origCron.ID, cj.ID)
	assert.Equal(t, origCron.Completions.Value(), cj.Completions.Value())
	assert.Equal(t, origCron.Parallelism.Value(), cj.Parallelism.Value())
	assert.Equal(t, origCron.Schedule.Value(), cj.Schedule.Value())
}

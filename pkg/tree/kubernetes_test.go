package tree

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/kubernetes"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/apps"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/batch"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/core"
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
				StatefulSets: []apps.StatefulSet{
					{
						Workload: workload.Workload{Resource: resource.Resource{ID: "db"}},
						Replicas: value.New[int64](3, 0, "", nil),
						VolumeClaimTemplates: []core.StorageRequest{
							{
								StorageClassName: value.New("gp3", 0, "", nil),
								RequestBytes:     value.New[int64](10737418240, 0, "", nil),
							},
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
			Core: core.Core{
				PersistentVolumeClaims: []core.PersistentVolumeClaim{
					{
						Resource: resource.Resource{
							ID: "data",
							Tags: resource.Tags{
								{Key: value.New("app", 0, "", nil), Value: value.New("api", 0, "", nil)},
							},
						},
						StorageRequest: core.StorageRequest{
							StorageClassName: value.New("io2", 0, "", nil),
							RequestBytes:     value.New[int64](21474836480, 0, "", nil),
						},
						Annotations: []resource.Tag{
							{
								Key:   value.New("volume.beta.kubernetes.io/storage-provisioner", 0, "", nil),
								Value: value.New("ebs.csi.aws.com", 0, "", nil),
							},
						},
					},
				},
				Services: []core.Service{
					{
						Resource: resource.Resource{ID: "api-lb"},
						Type:     value.New("LoadBalancer", 0, "", nil),
						Annotations: []resource.Tag{
							{
								Key:   value.New("service.beta.kubernetes.io/aws-load-balancer-type", 0, "", nil),
								Value: value.New("nlb", 0, "", nil),
							},
						},
						Ports: []core.ServicePort{
							{Port: value.New[int64](443, 0, "", nil), Protocol: value.New("TCP", 0, "", nil)},
						},
					},
				},
				Namespaces: []core.Namespace{
					{
						Resource: resource.Resource{
							ID:           "prod",
							SupportsTags: true,
							// A Namespace carries no cost, so its labels are the
							// whole point of surfacing it — they must survive the
							// round trip for tag policies to see them.
							Tags: resource.Tags{
								{Key: value.New("team", 0, "", nil), Value: value.New("platform", 0, "", nil)},
							},
						},
						Annotations: []resource.Tag{
							{
								Key:   value.New("scheduler.alpha.kubernetes.io/node-selector", 0, "", nil),
								Value: value.New("env=prod", 0, "", nil),
							},
						},
					},
				},
			},
		},
	}

	origDep := original.Kubernetes.Apps.Deployments[0]
	origDaemon := original.Kubernetes.Apps.DaemonSets[0]
	origSts := original.Kubernetes.Apps.StatefulSets[0]
	origCron := original.Kubernetes.Batch.CronJobs[0]
	origPVC := original.Kubernetes.Core.PersistentVolumeClaims[0]
	origSvc := original.Kubernetes.Core.Services[0]
	origNS := original.Kubernetes.Core.Namespaces[0]

	proto, err := original.ToProto()
	require.NoError(t, err)

	// The Deployment's base fields and its own replicas are flattened into a
	// single attribute object — not nested under an embedded-struct key.
	depAttrs := proto.Providers["kubernetes"].Services["apps"].Resources[0].Attributes.Entries
	assert.Equal(t, origDep.Replicas.Value(), depAttrs["replicas"].GetIntValue())
	require.NotNil(t, depAttrs["containers"], "embedded base fields must be flattened")
	require.NotNil(t, depAttrs["annotations"], "embedded base fields must be flattened")

	// The PersistentVolumeClaim embeds StorageRequest, so its storage fields are
	// flattened to the claim's top level rather than nested under an embed key.
	pvcAttrs := proto.Providers["kubernetes"].Services["core"].Resources[0].Attributes.Entries
	assert.Equal(t, origPVC.StorageClassName.Value(), pvcAttrs["storage_class_name"].GetStringValue())
	assert.Equal(t, origPVC.RequestBytes.Value(), pvcAttrs["request_bytes"].GetIntValue())

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

	// StatefulSet volumeClaimTemplates: a slice of nested StorageRequest structs.
	require.Len(t, result.Kubernetes.Apps.StatefulSets, 1)
	sts := result.Kubernetes.Apps.StatefulSets[0]
	assert.Equal(t, origSts.ID, sts.ID)
	assert.Equal(t, origSts.Replicas.Value(), sts.Replicas.Value())
	require.Len(t, sts.VolumeClaimTemplates, len(origSts.VolumeClaimTemplates))
	assert.Equal(t, origSts.VolumeClaimTemplates[0].StorageClassName.Value(), sts.VolumeClaimTemplates[0].StorageClassName.Value())
	assert.Equal(t, origSts.VolumeClaimTemplates[0].RequestBytes.Value(), sts.VolumeClaimTemplates[0].RequestBytes.Value())

	// core group: PersistentVolumeClaim (embedded StorageRequest) and Service.
	require.Len(t, result.Kubernetes.Core.PersistentVolumeClaims, 1)
	pvc := result.Kubernetes.Core.PersistentVolumeClaims[0]
	assert.Equal(t, origPVC.ID, pvc.ID)
	assert.Equal(t, origPVC.StorageClassName.Value(), pvc.StorageClassName.Value())
	assert.Equal(t, origPVC.RequestBytes.Value(), pvc.RequestBytes.Value())
	require.Len(t, pvc.Tags, len(origPVC.Tags))
	assert.Equal(t, origPVC.Tags[0].Key.Value(), pvc.Tags[0].Key.Value())
	require.Len(t, pvc.Annotations, len(origPVC.Annotations))
	assert.Equal(t, origPVC.Annotations[0].Key.Value(), pvc.Annotations[0].Key.Value())
	assert.Equal(t, origPVC.Annotations[0].Value.Value(), pvc.Annotations[0].Value.Value())

	require.Len(t, result.Kubernetes.Core.Services, 1)
	svc := result.Kubernetes.Core.Services[0]
	assert.Equal(t, origSvc.ID, svc.ID)
	assert.Equal(t, origSvc.Type.Value(), svc.Type.Value())
	require.Len(t, svc.Annotations, len(origSvc.Annotations))
	assert.Equal(t, origSvc.Annotations[0].Key.Value(), svc.Annotations[0].Key.Value())
	require.Len(t, svc.Ports, len(origSvc.Ports))
	assert.Equal(t, origSvc.Ports[0].Port.Value(), svc.Ports[0].Port.Value())
	assert.Equal(t, origSvc.Ports[0].Protocol.Value(), svc.Ports[0].Protocol.Value())

	// A Namespace has no cost-relevant fields of its own — the labels it carries
	// as the base resource's Tags are what tag policies act on, so they and the
	// tag-support flag are what must survive.
	require.Len(t, result.Kubernetes.Core.Namespaces, 1)
	ns := result.Kubernetes.Core.Namespaces[0]
	assert.Equal(t, origNS.ID, ns.ID)
	assert.True(t, ns.SupportsTags, "a namespace must stay tag-supporting so tag policies apply to it")
	require.Len(t, ns.Tags, len(origNS.Tags))
	assert.Equal(t, origNS.Tags[0].Key.Value(), ns.Tags[0].Key.Value())
	assert.Equal(t, origNS.Tags[0].Value.Value(), ns.Tags[0].Value.Value())
	require.Len(t, ns.Annotations, len(origNS.Annotations))
	assert.Equal(t, origNS.Annotations[0].Key.Value(), ns.Annotations[0].Key.Value())
	assert.Equal(t, origNS.Annotations[0].Value.Value(), ns.Annotations[0].Value.Value())

	// Multi-level embedding: CronJob -> Job -> Workload all flatten together.
	require.Len(t, result.Kubernetes.Batch.CronJobs, 1)
	cj := result.Kubernetes.Batch.CronJobs[0]
	assert.Equal(t, origCron.ID, cj.ID)
	assert.Equal(t, origCron.Completions.Value(), cj.Completions.Value())
	assert.Equal(t, origCron.Parallelism.Value(), cj.Parallelism.Value())
	assert.Equal(t, origCron.Schedule.Value(), cj.Schedule.Value())
}

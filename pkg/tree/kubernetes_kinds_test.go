package tree

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/kubernetes"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/autoscaling"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/core"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/karpenter"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/meta"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/networking"
	"github.com/infracost/go-proto/pkg/tree/kubernetes/policy"
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func str(s string) value.String { return value.New(s, 0, "", nil) }
func i64(n int64) value.Int     { return value.New(n, 0, "", nil) }

// TestKubernetesGovernanceKindsRoundTrip covers the kinds that provision nothing
// but decide whether a change to a workload takes effect, plus the two that do
// provision. Together they exercise the shapes the older kinds do not:
// pointer-backed string lists, structs nested two deep, and booleans.
//
// The case worth stating is the VerticalPodAutoscaler's TargetRef.Name. Every
// kind embeds meta.ObjectMeta, whose Name flattens into the kind's own attribute
// object, and TargetRef carries a Name of its own — a different thing, the
// workload being targeted rather than the autoscaler. They must both survive and
// stay distinct. TestNoFlattenedTagCollisions proves they cannot collide
// structurally; this proves the values do not cross on the wire.
func TestKubernetesGovernanceKindsRoundTrip(t *testing.T) {
	original := &Tree{
		Kubernetes: kubernetes.Kubernetes{
			Autoscaling: autoscaling.Autoscaling{
				VerticalPodAutoscalers: []autoscaling.VerticalPodAutoscaler{
					{
						Resource:   resource.Resource{ID: "api-vpa.1a2b3c4d"},
						ObjectMeta: meta.ObjectMeta{Name: str("api-vpa"), Namespace: str("shop")},
						TargetRef: autoscaling.TargetRef{
							APIVersion: str("apps/v1"),
							Kind:       str("Deployment"),
							Name:       str("api"),
						},
						UpdateMode: str(autoscaling.UpdateModeAuto),
						ContainerPolicies: []autoscaling.ContainerPolicy{
							{
								ContainerName: str("sidecar"),
								Mode:          str(autoscaling.ContainerPolicyModeOff),
								MinAllowed:    autoscaling.ResourceAmounts{CPUMillicores: i64(100)},
								MaxAllowed:    autoscaling.ResourceAmounts{MemoryBytes: i64(536870912)},
							},
							{
								ContainerName:       str("api"),
								ControlledResources: value.NewList([]value.String{str("memory")}, 0, "", nil),
								ControlledValues:    str(autoscaling.ControlledValuesRequestsOnly),
							},
						},
					},
				},
				HorizontalPodAutoscalers: []autoscaling.HorizontalPodAutoscaler{
					{
						Resource:       resource.Resource{ID: "api-hpa.5e6f7a8b"},
						ObjectMeta:     meta.ObjectMeta{Name: str("api-hpa"), Namespace: str("shop")},
						ScaleTargetRef: autoscaling.TargetRef{Kind: str("Deployment"), Name: str("api")},
						MinReplicas:    i64(2),
						MaxReplicas:    i64(20),
						Metrics: []autoscaling.Metric{
							{
								Type:              str(autoscaling.MetricSourceTypeResource),
								Name:              str("cpu"),
								TargetType:        str(autoscaling.MetricTargetTypeUtilization),
								TargetUtilization: i64(50),
							},
							{
								Type:              str(autoscaling.MetricSourceTypeContainerResource),
								Name:              str("memory"),
								ContainerName:     str("api"),
								TargetType:        str(autoscaling.MetricTargetTypeUtilization),
								TargetUtilization: i64(80),
							},
							{
								Type:        str(autoscaling.MetricSourceTypeExternal),
								Name:        str("sqs_queue_depth"),
								TargetType:  str(autoscaling.MetricTargetTypeAverageValue),
								TargetValue: str("30"),
							},
						},
					},
				},
			},
			Policy: policy.Policy{
				PodDisruptionBudgets: []policy.PodDisruptionBudget{
					{
						Resource:     resource.Resource{ID: "api-pdb.9c8d7e6f"},
						ObjectMeta:   meta.ObjectMeta{Name: str("api-pdb"), Namespace: str("shop")},
						MinAvailable: str("50%"),
						// Both halves of the selector, because the point of the
						// type is that they are not interchangeable: the
						// expressions below say three things no key/value pair
						// can, and the old flat shape dropped all three.
						Selector: meta.LabelSelector{
							MatchLabels: []resource.Tag{
								{Key: str("app"), Value: str("api")},
							},
							MatchExpressions: []meta.LabelSelectorRequirement{
								{Key: str("tier"), Operator: str("In"), Values: value.NewList([]value.String{str("web"), str("edge")}, 0, "", nil)},
								{Key: str("batch"), Operator: str("NotIn"), Values: value.NewList([]value.String{str("true")}, 0, "", nil)},
								{Key: str("managed"), Operator: str("Exists")},
							},
						},
					},
				},
			},
			Networking: networking.Networking{
				Ingresses: []networking.Ingress{
					{
						Resource:         resource.Resource{ID: "shop-ingress.11223344"},
						ObjectMeta:       meta.ObjectMeta{Name: str("shop-ingress"), Namespace: str("shop")},
						IngressClassName: str("alb"),
						Annotations: []resource.Tag{
							{Key: str("alb.ingress.kubernetes.io/group.name"), Value: str("shared")},
						},
						Rules: []networking.IngressRule{
							{
								Host: str("shop.example.com"),
								Paths: []networking.IngressPath{
									{
										Path:        str("/api"),
										PathType:    str("Prefix"),
										ServiceName: str("api"),
										ServicePort: str("http"),
									},
								},
							},
						},
						DefaultBackendServiceName: str("storefront"),
						DefaultBackendServicePort: str("80"),
					},
				},
			},
			Karpenter: karpenter.Karpenter{
				NodePools: []karpenter.NodePool{
					{
						Resource:   resource.Resource{ID: "default.aabbccdd"},
						ObjectMeta: meta.ObjectMeta{Name: str("default")},
						Requirements: []karpenter.Requirement{
							{
								Key:      str("karpenter.sh/capacity-type"),
								Operator: str("In"),
								Values:   value.NewList([]value.String{str("spot"), str("on-demand")}, 0, "", nil),
							},
						},
						Taints: []karpenter.Taint{
							{Key: str("dedicated"), Value: str("gpu"), Effect: str("NoSchedule")},
						},
						NodeClassRef: karpenter.NodeClassRef{
							Group: str("karpenter.k8s.aws"),
							Kind:  str("EC2NodeClass"),
							Name:  str("default"),
						},
						Limits: karpenter.NodePoolLimits{CPUMillicores: i64(64000)},
						Disruption: karpenter.Disruption{
							ConsolidationPolicy: str(karpenter.ConsolidationPolicyWhenEmptyOrUnderutilized),
							ConsolidateAfter:    str("30s"),
							Budgets: []karpenter.DisruptionBudget{
								{
									Nodes:   str("20%"),
									Reasons: value.NewList([]value.String{str("Underutilized")}, 0, "", nil),
								},
							},
						},
					},
				},
				EC2NodeClasses: []karpenter.EC2NodeClass{
					{
						Resource:   resource.Resource{ID: "default.eeff0011"},
						ObjectMeta: meta.ObjectMeta{Name: str("default")},
						AMIFamily:  str("AL2023"),
						BlockDeviceMappings: []karpenter.BlockDeviceMapping{
							{
								DeviceName:          str("/dev/xvda"),
								VolumeSizeBytes:     i64(107374182400),
								VolumeType:          str("gp3"),
								IOPS:                i64(3000),
								Encrypted:           value.New(true, 0, "", nil),
								DeleteOnTermination: value.New(false, 0, "", nil),
							},
						},
						SubnetSelectorTerms: []karpenter.SelectorTerm{
							{Tags: []resource.Tag{{Key: str("karpenter.sh/discovery"), Value: str("prod")}}},
						},
					},
				},
			},
			Core: core.Core{
				LimitRanges: []core.LimitRange{
					{
						Resource:   resource.Resource{ID: "defaults.55667788"},
						ObjectMeta: meta.ObjectMeta{Name: str("defaults"), Namespace: str("shop")},
						Limits: []core.LimitRangeItem{
							{
								Type:           str(core.LimitRangeTypeContainer),
								DefaultRequest: core.LimitRangeAmounts{CPUMillicores: i64(100), MemoryBytes: i64(134217728)},
								Max:            core.LimitRangeAmounts{CPUMillicores: i64(4000)},
								MaxLimitRequestRatio: core.LimitRangeRatios{
									CPU: value.New(4.0, 0, "", nil),
								},
							},
						},
					},
				},
			},
		},
	}

	proto, err := original.ToProto()
	require.NoError(t, err)

	result, err := FromProto(proto)
	require.NoError(t, err)

	// VerticalPodAutoscaler: the two Names must not have crossed. The VPA's own
	// name comes off the flattened ObjectMeta; the target's stays inside the
	// nested TargetRef object.
	require.Len(t, result.Kubernetes.Autoscaling.VerticalPodAutoscalers, 1)
	vpa := result.Kubernetes.Autoscaling.VerticalPodAutoscalers[0]
	assert.Equal(t, "api-vpa", vpa.Name.Value(), "the VPA's own name, promoted from ObjectMeta")
	assert.Equal(t, "api", vpa.TargetRef.Name.Value(), "the workload it targets")
	assert.Equal(t, "shop", vpa.Namespace.Value())
	assert.Equal(t, autoscaling.UpdateModeAuto, vpa.UpdateMode.Value())
	require.Len(t, vpa.ContainerPolicies, 2)
	assert.Equal(t, autoscaling.ContainerPolicyModeOff, vpa.ContainerPolicies[0].Mode.Value())
	assert.Equal(t, int64(100), vpa.ContainerPolicies[0].MinAllowed.CPUMillicores.Value())
	assert.Equal(t, int64(536870912), vpa.ContainerPolicies[0].MaxAllowed.MemoryBytes.Value())
	assert.Nil(t, vpa.ContainerPolicies[0].ControlledResources,
		"an omitted controlledResources must stay nil: nil defaults to both resources, an empty list would read as neither")

	// The second policy is the partly-governed container: memory only, requests
	// only. A CPU recommendation on it is actionable and a memory one is not, so
	// the narrowing has to survive as the list it was written as — a policy that
	// arrived back with both resources, or with none, reads as the opposite
	// case.
	memOnly := vpa.ContainerPolicies[1]
	assert.Equal(t, "api", memOnly.ContainerName.Value())
	require.NotNil(t, memOnly.ControlledResources)
	assert.Equal(t, []string{"memory"}, listValues(memOnly.ControlledResources))
	assert.Equal(t, autoscaling.ControlledValuesRequestsOnly, memOnly.ControlledValues.Value())

	require.Len(t, result.Kubernetes.Autoscaling.HorizontalPodAutoscalers, 1)
	hpa := result.Kubernetes.Autoscaling.HorizontalPodAutoscalers[0]
	assert.Equal(t, "api", hpa.ScaleTargetRef.Name.Value())
	assert.Equal(t, int64(2), hpa.MinReplicas.Value())
	assert.Equal(t, int64(20), hpa.MaxReplicas.Value())

	// Metrics: the utilization targets are the setpoints a rightsizing pass has
	// to read observed usage against, so the percentages must survive as
	// numbers, and the container-scoped one must stay bound to its container.
	// The external metric carries no utilization at all, which is how a
	// consumer tells "held at 50% of its request" from "driven by a queue".
	require.Len(t, hpa.Metrics, 3)
	assert.Equal(t, autoscaling.MetricSourceTypeResource, hpa.Metrics[0].Type.Value())
	assert.Equal(t, "cpu", hpa.Metrics[0].Name.Value())
	assert.Equal(t, int64(50), hpa.Metrics[0].TargetUtilization.Value())
	assert.Empty(t, hpa.Metrics[0].ContainerName.Value(), "a pod-total metric names no container")

	assert.Equal(t, autoscaling.MetricSourceTypeContainerResource, hpa.Metrics[1].Type.Value())
	assert.Equal(t, "memory", hpa.Metrics[1].Name.Value())
	assert.Equal(t, "api", hpa.Metrics[1].ContainerName.Value())
	assert.Equal(t, int64(80), hpa.Metrics[1].TargetUtilization.Value())

	assert.Equal(t, autoscaling.MetricSourceTypeExternal, hpa.Metrics[2].Type.Value())
	assert.Equal(t, "sqs_queue_depth", hpa.Metrics[2].Name.Value())
	assert.Equal(t, autoscaling.MetricTargetTypeAverageValue, hpa.Metrics[2].TargetType.Value())
	assert.Equal(t, "30", hpa.Metrics[2].TargetValue.Value(), "an arbitrary-unit quantity stays a string")
	assert.Zero(t, hpa.Metrics[2].TargetUtilization.Value(), "no utilization target on an external metric")

	// PodDisruptionBudget: minAvailable stays a string so "50%" survives as a
	// percentage rather than being read as an absolute count.
	require.Len(t, result.Kubernetes.Policy.PodDisruptionBudgets, 1)
	pdb := result.Kubernetes.Policy.PodDisruptionBudgets[0]
	assert.Equal(t, "50%", pdb.MinAvailable.Value())
	require.Len(t, pdb.Selector.MatchLabels, 1)
	assert.Equal(t, "app", pdb.Selector.MatchLabels[0].Key.Value())
	assert.Equal(t, "api", pdb.Selector.MatchLabels[0].Value.Value())
	// The expressions are the half the flat shape could not hold. Each of these
	// three states something a key/value pair cannot, so a round trip that
	// silently dropped them would leave the selector wider than the manifest
	// wrote it.
	require.Len(t, pdb.Selector.MatchExpressions, 3)
	assert.Equal(t, "tier", pdb.Selector.MatchExpressions[0].Key.Value())
	assert.Equal(t, "In", pdb.Selector.MatchExpressions[0].Operator.Value())
	assert.Equal(t, []string{"web", "edge"}, listValues(pdb.Selector.MatchExpressions[0].Values))
	assert.Equal(t, "NotIn", pdb.Selector.MatchExpressions[1].Operator.Value())
	assert.Equal(t, []string{"true"}, listValues(pdb.Selector.MatchExpressions[1].Values))
	assert.Equal(t, "Exists", pdb.Selector.MatchExpressions[2].Operator.Value())
	assert.Nil(t, pdb.Selector.MatchExpressions[2].Values,
		"Exists takes no values, so the list is absent rather than empty — an empty list under In matches nothing")

	// Ingress: rules nest paths, so this is a struct two levels inside a slice.
	require.Len(t, result.Kubernetes.Networking.Ingresses, 1)
	ing := result.Kubernetes.Networking.Ingresses[0]
	assert.Equal(t, "alb", ing.IngressClassName.Value())
	require.Len(t, ing.Rules, 1)
	assert.Equal(t, "shop.example.com", ing.Rules[0].Host.Value())
	require.Len(t, ing.Rules[0].Paths, 1)
	assert.Equal(t, "/api", ing.Rules[0].Paths[0].Path.Value())
	assert.Equal(t, "http", ing.Rules[0].Paths[0].ServicePort.Value(), "a named port must stay a string")
	assert.Equal(t, "storefront", ing.DefaultBackendServiceName.Value(),
		"the default backend is a peer of the rules, and the only Service join on a rule-less Ingress")
	assert.Equal(t, "80", ing.DefaultBackendServicePort.Value())

	// NodePool: the pointer-backed string lists are the least-exercised type
	// here, and a nil-versus-empty slip in either direction would be silent.
	require.Len(t, result.Kubernetes.Karpenter.NodePools, 1)
	pool := result.Kubernetes.Karpenter.NodePools[0]
	require.Len(t, pool.Requirements, 1)
	require.NotNil(t, pool.Requirements[0].Values)
	assert.Equal(t, []string{"spot", "on-demand"}, listValues(pool.Requirements[0].Values))
	require.Len(t, pool.Taints, 1)
	assert.Equal(t, "NoSchedule", pool.Taints[0].Effect.Value())
	assert.Equal(t, "EC2NodeClass", pool.NodeClassRef.Kind.Value())
	assert.Equal(t, int64(64000), pool.Limits.CPUMillicores.Value())
	assert.Equal(t, karpenter.ConsolidationPolicyWhenEmptyOrUnderutilized, pool.Disruption.ConsolidationPolicy.Value())
	require.Len(t, pool.Disruption.Budgets, 1)
	assert.Equal(t, "20%", pool.Disruption.Budgets[0].Nodes.Value())
	require.NotNil(t, pool.Disruption.Budgets[0].Reasons)
	assert.Equal(t, []string{"Underutilized"}, listValues(pool.Disruption.Budgets[0].Reasons))

	// EC2NodeClass: booleans, and a false that must stay false rather than
	// reading as unset — deleteOnTermination false is the case that leaves
	// volumes behind, so losing it loses the finding.
	require.Len(t, result.Kubernetes.Karpenter.EC2NodeClasses, 1)
	nc := result.Kubernetes.Karpenter.EC2NodeClasses[0]
	assert.Equal(t, "AL2023", nc.AMIFamily.Value())
	require.Len(t, nc.BlockDeviceMappings, 1)
	bdm := nc.BlockDeviceMappings[0]
	assert.Equal(t, int64(107374182400), bdm.VolumeSizeBytes.Value())
	assert.Equal(t, "gp3", bdm.VolumeType.Value())
	assert.Equal(t, int64(3000), bdm.IOPS.Value())
	assert.True(t, bdm.Encrypted.Value())
	assert.False(t, bdm.DeleteOnTermination.Value())
	require.Len(t, nc.SubnetSelectorTerms, 1)
	require.Len(t, nc.SubnetSelectorTerms[0].Tags, 1)
	assert.Equal(t, "karpenter.sh/discovery", nc.SubnetSelectorTerms[0].Tags[0].Key.Value())

	// LimitRange: amounts nested inside items inside the kind, and the ratio
	// which is a bare multiple rather than a quantity.
	require.Len(t, result.Kubernetes.Core.LimitRanges, 1)
	lr := result.Kubernetes.Core.LimitRanges[0]
	assert.Equal(t, "shop", lr.Namespace.Value())
	require.Len(t, lr.Limits, 1)
	item := lr.Limits[0]
	assert.Equal(t, core.LimitRangeTypeContainer, item.Type.Value())
	assert.Equal(t, int64(100), item.DefaultRequest.CPUMillicores.Value())
	assert.Equal(t, int64(134217728), item.DefaultRequest.MemoryBytes.Value())
	assert.Equal(t, int64(4000), item.Max.CPUMillicores.Value())
	assert.InDelta(t, 4.0, item.MaxLimitRequestRatio.CPU.Value(), 0.0001)
}

// listValues unwraps a value.List into the bare strings it holds, so the
// assertions above read as the list the manifest stated.
func listValues(l *value.List[string]) []string {
	if l == nil {
		return nil
	}
	items := l.Items()
	out := make([]string, 0, len(items))
	for _, v := range items {
		out = append(out, v.Value())
	}
	return out
}

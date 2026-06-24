// Package workload holds the shared Kubernetes workload resource type. Every
// workload kind (Deployment, StatefulSet, Job, ...) has the same cost-relevant
// shape — replica count plus per-container sizing — so they share a single
// Workload type. The kind itself is carried on the embedded resource (its
// Definition.ResourceType and address), and the API group it belongs to is the
// service that holds it in the tree (see the apps and batch packages).
package workload

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// Workload is the shape common to every Kubernetes workload kind: the
// per-container sizing and the annotations. It is used directly for kinds that
// have no extra cost-relevant fields (e.g. DaemonSet, whose replica count is
// dynamic — one pod per node) and embedded by the kinds that do
// (Deployment/StatefulSet/ReplicaSet add a replica count; Job adds
// completions/parallelism; see the apps and batch packages).
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource (Definition.ResourceType / Definition.Address /
// Definition.SourceRange). The workload's Kubernetes labels are stored as the
// base resource's Tags (each carries its own source range), so the
// label/tag-enforcement use case reuses the existing tag machinery. The sizing
// read out of the manifest — or, for Helm, out of values.yaml — is first-class
// and typed below.
type Workload struct {
	resource.Resource `tree:"-"`

	// Annotations are the workload's Kubernetes annotations, surfaced verbatim.
	// The parser stays provider-agnostic: cloud-provider signals (IRSA role
	// ARNs, GKE/Azure workload-identity annotations, etc.) live here, and it is
	// up to downstream consumers to interpret them and decide a provider.
	Annotations []resource.Tag `tree:"annotations"`

	// Containers holds the per-container resource requests and limits.
	Containers []Container `tree:"containers"`
}

// Container is the cost-relevant sizing for a single container in a workload's
// pod template.
//
// CPU is stored in millicores and memory in bytes — the canonical base units a
// Kubernetes quantity reduces to ("500m" -> 500 millicores, "512Mi" ->
// 536870912 bytes). These are exact integers regardless of the suffix the
// author used, which avoids any GiB-vs-GB ambiguity; the quantity parsing lives
// in the parser and the downstream provider plugin converts to vCores/GB as
// part of its rate.
type Container struct {
	Name                 value.String `tree:"name"`
	CPURequestMillicores value.Int    `tree:"cpu_request_millicores"`
	CPULimitMillicores   value.Int    `tree:"cpu_limit_millicores"`
	MemoryRequestBytes   value.Int    `tree:"memory_request_bytes"`
	MemoryLimitBytes     value.Int    `tree:"memory_limit_bytes"`
}

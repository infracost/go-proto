package core

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

// StorageRequest is the cost-relevant sizing of a persistent volume: the
// requested capacity and the storage class that selects the backing disk SKU.
// The storage class is the price driver (e.g. "gp3" vs "io2" on AWS,
// "premium-rwo" vs "standard-rwo" on GKE differ several-fold per GB), and the
// request is the exact provisioned size — so this pair is enough to price the
// volume without cluster or API access.
//
// It is shared by a standalone PersistentVolumeClaim (embedded, so its fields
// sit at the claim's top level) and by the volumeClaimTemplates of a StatefulSet
// (see the apps package), where one such volume is provisioned per replica.
type StorageRequest struct {
	// StorageClassName is spec.storageClassName. Empty means the cluster's
	// default StorageClass — the parser leaves it empty rather than guessing, so
	// the downstream pricer can fall back to a documented default.
	StorageClassName value.String `tree:"storage_class_name"`

	// RequestBytes is spec.resources.requests.storage reduced to bytes — the
	// canonical base unit a Kubernetes quantity reduces to ("10Gi" ->
	// 10737418240), exact regardless of the suffix the author used. The
	// downstream pricer converts to GB as part of its rate.
	RequestBytes value.Int `tree:"request_bytes"`
}

// PersistentVolumeClaim is a core/v1 PersistentVolumeClaim: a request that a
// dynamic provisioner satisfies with a real cloud volume (EBS / Persistent Disk
// / Azure disk). That volume is billed independently of the nodes, so — unlike
// per-workload compute — pricing it does not double-count anything the node-pool
// IaC already covers.
//
// The kind, address ([namespace, kind, name]) and source range live on the
// embedded resource.Resource; the claim's Kubernetes labels are stored as the
// base resource's Tags (reusing the tag machinery, as workloads do).
type PersistentVolumeClaim struct {
	resource.Resource `tree:"-"`
	StorageRequest    `tree:"-"`

	// Annotations are the claim's Kubernetes annotations, surfaced verbatim so
	// downstream consumers can sniff cloud-provider signals (IRSA role ARNs,
	// GKE/Azure workload-identity annotations, etc.) — matching how the workload
	// type carries them.
	Annotations []resource.Tag `tree:"annotations"`
}
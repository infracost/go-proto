// Package meta holds the Kubernetes object metadata shared by every kind in the
// tree. It mirrors the metadata.name / metadata.namespace pair that every
// Kubernetes object carries, as typed values rather than as segments of the
// resource address.
//
// The address ([namespace, kind, name]) still carries the same identity, but it
// is a display and uniqueness construct: the name segment is joined with a hash
// on the resource ID to disambiguate reuse across namespaces, and the namespace
// segment cannot express where a value came from. Consumers that match a tree
// resource against a real cluster object — code search, cost attribution — need
// the bare values and their provenance, which is what these fields are for. It
// is the same shape every other provider already exposes (an RDS instance's
// Identifier, an S3 bucket's name), so those consumers can read a name off a
// typed field instead of parsing an address.
package meta

import (
	"github.com/infracost/go-proto/pkg/tree/value"
)

// ObjectMeta is the identity every Kubernetes object declares under its
// metadata block. It is embedded by each kind in the tree, so its fields
// serialize inline alongside the kind's own.
type ObjectMeta struct {
	// Name is metadata.name verbatim — the name the object is created with in
	// the cluster, and the name a metrics pipeline reports it under. Unlike the
	// resource ID it carries no disambiguating hash, so it is directly
	// comparable to an observed workload name.
	Name value.String `tree:"name"`

	// Namespace is the namespace the object is deployed to, which is not always
	// something the manifest states. Kubernetes resolves it at apply time: an
	// explicit metadata.namespace wins, otherwise it comes from the apply
	// context (kubectl's -n, a Helm release namespace, a Kustomize overlay),
	// and only in the absence of all of those does it fall back to "default".
	//
	// A repository holds the manifest but not the apply context, so the value
	// here carries its origin in its flags and consumers must read them:
	//
	//   - no flag       read from metadata.namespace, with a source range
	//   - flag.Config   supplied by the caller's configuration, standing in for
	//                   the -n the deploy would have used
	//   - flag.Synthetic  nothing declared it; "default" is our assumption, and
	//                     is as likely to be wrong as the deploy is to name a
	//                     namespace
	//
	// A synthetic namespace must not be used to exclude a candidate when
	// matching against a real cluster — it is an assumption, not evidence.
	Namespace value.String `tree:"namespace"`
}

// GetObjectMeta returns the metadata itself, so that every kind embedding
// ObjectMeta promotes an accessor for it. Embedding promotes fields too, but
// only a method survives type erasure: a consumer walking a tree holds each
// resource as an any, and Go generics cannot constrain on a field. This is the
// same reason resource.Resource carries GetBase.
func (m *ObjectMeta) GetObjectMeta() *ObjectMeta {
	return m
}

// Object is satisfied by every Kubernetes kind in the tree, through the
// promoted GetObjectMeta above — at whatever embedding depth it sits (a CronJob
// reaches it through Job through Workload). Consumers that match tree resources
// against real cluster objects assert on this to reach the typed name and
// namespace, as they assert on resource.Implementation to reach the base.
type Object interface {
	GetObjectMeta() *ObjectMeta
}

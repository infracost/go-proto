package core

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
)

// Namespace is a core/v1 Namespace. Unlike the other kinds in this package it
// provisions nothing and carries no cost of its own — it is surfaced so that
// label/tag-enforcement policies can be applied to it. Namespaces are where
// teams commonly record ownership (team, cost-centre, environment), so a
// namespace missing those labels is exactly the kind of gap a tagging policy
// exists to catch, and it cannot be caught if the namespace never reaches the
// tree.
//
// The kind, address and source range live on the embedded resource.Resource;
// the Namespace's Kubernetes labels are stored as the base resource's Tags,
// reusing the tag machinery as the other kinds do.
//
// A Namespace is cluster-scoped, so it has no metadata.namespace of its own.
// The parser addresses it as [name, kind, name] — scoped to itself — which
// keeps every Kubernetes kind on the same three-segment address shape and
// groups the namespace with the resources it contains.
type Namespace struct {
	resource.Resource `tree:"-"`

	// Annotations are the Namespace's Kubernetes annotations, surfaced verbatim
	// so downstream consumers can read the cloud-provider signals the other
	// kinds also carry.
	Annotations []resource.Tag `tree:"annotations"`
}

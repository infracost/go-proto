package kms

import "github.com/infracost/go-proto/pkg/tree/resource"

type ReplicaKey struct {
	resource.Resource `tree:"-"`
}

type ReplicaExternalKey struct {
	resource.Resource `tree:"-"`
}

package storage

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type Bucket struct {
	resource.Resource `tree:"-"`
	Location          value.String                `tree:"location"`
	StorageClass      value.Value[StorageClass]   `tree:"storage_class"`
	Versioning        value.Bool                  `tree:"versioning"`
	LifecycleRules    []LifecycleRule             `tree:"lifecycle_rules"`
}

type LifecycleRule struct {
	ActionType              value.Value[LifecycleActionType] `tree:"action_type"`
	ActionStorageClass      value.Value[StorageClass]  `tree:"action_storage_class"`
	Age                     value.Int                  `tree:"age"`
	NumNewerVersions        value.Int                  `tree:"num_newer_versions"`
	DaysSinceNoncurrentTime value.Int                  `tree:"days_since_noncurrent_time"`
}

type LifecycleActionType uint32

const (
	LifecycleActionTypeUnknown                        LifecycleActionType = iota
	LifecycleActionTypeDelete
	LifecycleActionTypeSetStorageClass
	LifecycleActionTypeAbortIncompleteMultipartUpload
)

type StorageClass uint32

const (
	StorageClassUnknown       StorageClass = iota
	StorageClassStandard
	StorageClassNearline
	StorageClassColdline
	StorageClassArchive
	StorageClassMultiRegional
	StorageClassRegional
)

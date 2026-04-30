package s3

import (
	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
)

type LifecycleConfiguration struct {
	resource.Resource `tree:"-"`
	BucketName        value.String    `tree:"bucket"`
	Rules             []LifecycleRule `tree:"rules"`
}

type LifecycleRule struct {
	resource.Resource                    `tree:"-"`
	Enabled                              value.Bool                    `tree:"enabled"`
	AbortIncompleteMultipartUploadDays   value.Int                     `tree:"abort_incomplete_multipart_upload_days"`
	NoncurrentVersionExpirationDays      value.Int                     `tree:"noncurrent_version_expiration_days"`
	NoncurrentVersionTransitions         []NoncurrentVersionTransition `tree:"noncurrent_version_transitions"`
	Transitions                          []LifecycleTransition         `tree:"transitions"`
	ExpirationDays                       value.Int                     `tree:"expiration_days"`
	ExpirationDate                       value.String                  `tree:"expiration_date"`
}

type NoncurrentVersionTransition struct {
	Days         value.Int                    `tree:"days"`
	StorageClass value.Value[S3StorageClass] `tree:"storage_class"`
}

type LifecycleTransition struct {
	Days         value.Int                    `tree:"days"`
	StorageClass value.Value[S3StorageClass] `tree:"storage_class"`
}

type S3StorageClass uint32

const (
	S3StorageClassUnknown             S3StorageClass = iota
	S3StorageClassStandard
	S3StorageClassStandardIA
	S3StorageClassOneZoneIA
	S3StorageClassIntelligentTiering
	S3StorageClassGlacier
	S3StorageClassGlacierIR
	S3StorageClassDeepArchive
)

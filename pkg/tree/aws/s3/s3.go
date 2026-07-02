package s3

import "strings"

type S3 struct {
	Buckets                          []Bucket                          `tree:"buckets"`
	LifecycleConfigurations          []LifecycleConfiguration          `tree:"lifecycle_configurations"`
	IntelligentTieringConfigurations []IntelligentTieringConfiguration `tree:"intelligent_tiering_configurations"`
	BucketInventories                []BucketInventory                 `tree:"bucket_inventories"`
	BucketAnalyticsConfigurations    []BucketAnalyticsConfiguration    `tree:"bucket_analytics_configurations"`
	BucketVersioningConfigurations   []BucketVersioningConfiguration   `tree:"bucket_versioning_configurations"`
	BucketPolicies                   []BucketPolicy                    `tree:"bucket_policies"`
	AccessPoints                     []AccessPoint                     `tree:"access_points"`
}

func (s *S3) PostProcess() {
	// Reset relationships this method writes to so PostProcess is idempotent.
	for i := range s.Buckets {
		s.Buckets[i].Relationships = BucketRelationships{}
	}

	bucketMap := make(map[string]*Bucket)
	for i := range s.Buckets {
		if name := s.Buckets[i].Name.Value(); name != "" {
			bucketMap[name] = &s.Buckets[i]
		}
		if id := s.Buckets[i].ID; id != "" {
			bucketMap[id] = &s.Buckets[i]
		}
	}

	// findBucket locates the owning bucket for a sub-resource. It matches on the
	// referenced bucket name first, then falls back to the sub-resource's own ID.
	// CloudFormation/CDK buckets often have no explicit BucketName (the physical
	// name is generated at deploy time), so inline configurations (lifecycle,
	// intelligent-tiering, versioning) carry an empty bucket name. For those, the
	// sub-resource shares the owning bucket's logical ID, so the ID fallback links
	// them correctly. Standalone resources (e.g. a separate lifecycle configuration
	// or bucket policy) have their own distinct ID and so never false-match.
	findBucket := func(bucketName, id string) (*Bucket, bool) {
		if b, ok := bucketMap[bucketName]; ok {
			return b, true
		}
		if id != "" {
			if b, ok := bucketMap[id]; ok {
				return b, true
			}
		}
		return nil, false
	}

	for i, lc := range s.LifecycleConfigurations {
		if b, ok := findBucket(lc.BucketName.Value(), lc.ID); ok {
			b.Relationships.LifecycleConfigurations = append(b.Relationships.LifecycleConfigurations, &s.LifecycleConfigurations[i])
		}
	}

	for i, itc := range s.IntelligentTieringConfigurations {
		if b, ok := findBucket(itc.BucketName.Value(), itc.ID); ok {
			b.Relationships.IntelligentTieringConfigurations = append(b.Relationships.IntelligentTieringConfigurations, &s.IntelligentTieringConfigurations[i])
		}
	}

	for i, vc := range s.BucketVersioningConfigurations {
		if b, ok := findBucket(vc.BucketName.Value(), vc.ID); ok {
			b.Relationships.BucketVersioningConfigurations = append(b.Relationships.BucketVersioningConfigurations, &s.BucketVersioningConfigurations[i])
		}
	}

	for i, bp := range s.BucketPolicies {
		b, ok := findBucket(bp.BucketName.Value(), bp.ID)
		if ok {
			b.Relationships.BucketPolicies = append(b.Relationships.BucketPolicies, &s.BucketPolicies[i])
		}
		finalizeDeniesInsecureTransport(&s.BucketPolicies[i], b, ok)
	}
}

// finalizeDeniesInsecureTransport resolves whether a bucket policy denies non-SSL
// transport for the bucket it is attached to. The SSL-deny statement itself is
// extracted at parse time (BucketPolicy.SSLDeny); the coverage decision — whether
// that statement's Resource ARNs actually refer to THIS bucket and cover both the
// bucket and its objects — is made here because it needs the owning bucket's
// identities, which are only known after linking. Matching against both the
// bucket's resolved name and its synthetic instance token makes the result
// correct no matter whether either side resolved to a literal ARN or to an
// unresolved reference.
func finalizeDeniesInsecureTransport(bp *BucketPolicy, bucket *Bucket, linked bool) {
	d := bp.SSLDeny
	if d == nil {
		// A converter that does not compute SSL-deny facts (e.g. one that has not
		// adopted this path). Leave the value it already set untouched.
		return
	}

	var denies bool
	switch {
	case d.Opaque:
		// The policy document was not parseable but was synthetic: its real
		// contents are unknowable, so give it the benefit of the doubt.
		denies = true
	case !d.Present:
		denies = false
	case !linked:
		// A valid non-SSL deny statement exists but the policy could not be linked
		// to a bucket in this project, so there is no bucket to verify coverage
		// against. Accept it rather than emit a false positive; unlinked policies
		// are not surfaced by the S3.5 finding anyway.
		denies = true
	default:
		denies = denyResourcesCoverBucket(d.Resources, bucket)
	}

	bp.DeniesInsecureTransport = bp.DeniesInsecureTransport.WithValue(denies)
}

// denyResourcesCoverBucket reports whether the non-SSL deny statements' combined
// Resource list protects both the bucket and its objects for the given bucket. A
// policy that only covers "<arn>/*" (objects) or that targets a different bucket
// leaves this bucket unprotected and is not compliant.
func denyResourcesCoverBucket(resources []string, bucket *Bucket) bool {
	hasBucket := false
	hasObjects := false
	for _, r := range resources {
		if r == "*" {
			// A wildcard resource covers the bucket and all of its objects.
			return true
		}
		bare := strings.TrimSuffix(r, "/*")
		isObjects := bare != r
		if !arnRefersToBucket(bare, bucket) {
			continue
		}
		if isObjects {
			hasObjects = true
		} else {
			hasBucket = true
		}
	}
	return hasBucket && hasObjects
}

// arnRefersToBucket reports whether an S3 ARN (with any "/*" object suffix already
// removed) refers to the given bucket. It matches either of the bucket's two
// identities:
//
//   - its synthetic instance token: the placeholder the parser emits for the
//     bucket's computed attributes (.arn/.id) when they are unresolved. This
//     equals the bucket's own ID, so a deny statement written against
//     `aws_s3_bucket.x.arn` matches even when the policy's `bucket` argument
//     resolved to a literal name; and
//   - its resolved name: matching an `arn:aws:s3:::<name>` ARN or a bare "<name>".
//
// Checking both makes coverage detection independent of which side happened to
// resolve, while still rejecting ARNs that refer to a different bucket.
func arnRefersToBucket(arn string, bucket *Bucket) bool {
	if bucket == nil {
		return false
	}
	if id := bucket.ID; id != "" && arn == id {
		return true
	}
	if name := bucket.Name.Value(); name != "" {
		if arn == name || strings.HasSuffix(arn, ":"+name) {
			return true
		}
	}
	return false
}

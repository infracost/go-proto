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
	// sub-resource either shares the owning bucket's logical ID, or is a child
	// synthesized by the parser whose ID is "<bucketID>:<path>" — the ID fallback
	// links both correctly. Standalone resources (e.g. a separate lifecycle
	// configuration or bucket policy) have their own distinct ID and so never
	// false-match.
	findBucket := func(bucketName, id string) (*Bucket, bool) {
		if b, ok := bucketMap[bucketName]; ok {
			return b, true
		}
		if id != "" {
			if b, ok := bucketMap[id]; ok {
				return b, true
			}
			if parentID, _, isChild := strings.Cut(id, ":"); isChild {
				if b, ok := bucketMap[parentID]; ok {
					return b, true
				}
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
		if b, ok := findBucket(bp.BucketName.Value(), bp.ID); ok {
			b.Relationships.BucketPolicies = append(b.Relationships.BucketPolicies, &s.BucketPolicies[i])
		}
	}
}

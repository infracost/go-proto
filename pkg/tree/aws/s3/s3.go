package s3

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

	for i, lc := range s.LifecycleConfigurations {
		if b, ok := bucketMap[lc.BucketName.Value()]; ok {
			b.Relationships.LifecycleConfigurations = append(b.Relationships.LifecycleConfigurations, &s.LifecycleConfigurations[i])
		}
	}

	for i, itc := range s.IntelligentTieringConfigurations {
		if b, ok := bucketMap[itc.BucketName.Value()]; ok {
			b.Relationships.IntelligentTieringConfigurations = append(b.Relationships.IntelligentTieringConfigurations, &s.IntelligentTieringConfigurations[i])
		}
	}

	for i, vc := range s.BucketVersioningConfigurations {
		if b, ok := bucketMap[vc.BucketName.Value()]; ok {
			b.Relationships.BucketVersioningConfigurations = append(b.Relationships.BucketVersioningConfigurations, &s.BucketVersioningConfigurations[i])
		}
	}

	for i, bp := range s.BucketPolicies {
		if b, ok := bucketMap[bp.BucketName.Value()]; ok {
			b.Relationships.BucketPolicies = append(b.Relationships.BucketPolicies, &s.BucketPolicies[i])
		}
	}
}

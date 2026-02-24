package event

import (
	"testing"

	eventpb "github.com/infracost/proto/gen/go/infracost/parser/event"
	"github.com/infracost/proto/gen/go/infracost/provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func defaultProjectInfo() *provider.ProjectInfo {
	return &provider.ProjectInfo{
		Name:       "my-project",
		BranchName: "main",
	}
}

func mkResource(name, resourceType string, tags []*provider.Tag) *provider.Resource {
	return &provider.Resource{
		Name: name,
		Type: resourceType,
		Tagging: &provider.Tagging{
			SupportsTags: true,
			Tags:         tags,
		},
	}
}

func mkTag(key, value string) *provider.Tag {
	return &provider.Tag{Key: key, Value: value}
}

func mkDefaultTag(key, value string) *provider.Tag {
	return &provider.Tag{Key: key, Value: value, IsDefault: true}
}

func TestEvaluateAgainstResources_EnvTagListPolicy(t *testing.T) {
	allowedEnvs := []string{"production", "development", "staging", "test"}
	policy := &eventpb.TagPolicy{
		Id:   "env-policy",
		Name: "Environment Tag Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:           "Env",
				Type:          eventpb.TagPolicyRequirement_LIST,
				AllowedValues: allowedEnvs,
				Mandatory:     true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	tests := []struct {
		name           string
		resources      []*provider.Resource
		wantPassing    int
		wantFailing    int
		wantSuggestion string
	}{
		{
			name: "valid env tag passes",
			resources: []*provider.Resource{
				mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
					mkTag("Env", "production"),
				}),
			},
			wantPassing: 1,
			wantFailing: 0,
		},
		{
			name: "all valid env values pass",
			resources: []*provider.Resource{
				mkResource("aws_instance.a", "aws_instance", []*provider.Tag{mkTag("Env", "production")}),
				mkResource("aws_instance.b", "aws_instance", []*provider.Tag{mkTag("Env", "development")}),
				mkResource("aws_instance.c", "aws_instance", []*provider.Tag{mkTag("Env", "staging")}),
				mkResource("aws_instance.d", "aws_instance", []*provider.Tag{mkTag("Env", "test")}),
			},
			wantPassing: 4,
			wantFailing: 0,
		},
		{
			name: "invalid env value fails",
			resources: []*provider.Resource{
				mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
					mkTag("Env", "prod"),
				}),
			},
			wantPassing:    0,
			wantFailing:    1,
			wantSuggestion: "production",
		},
		{
			name: "missing mandatory env tag fails",
			resources: []*provider.Resource{
				mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
					mkTag("Name", "my-server"),
				}),
			},
			wantPassing: 0,
			wantFailing: 1,
		},
		{
			name: "mix of passing and failing",
			resources: []*provider.Resource{
				mkResource("aws_instance.good", "aws_instance", []*provider.Tag{mkTag("Env", "production")}),
				mkResource("aws_instance.bad", "aws_instance", []*provider.Tag{mkTag("Env", "prd")}),
			},
			wantPassing: 1,
			wantFailing: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := policies.EvaluateAgainstResources(tt.resources, project)
			require.Len(t, results, 1)
			result := results[0]
			assert.Equal(t, tt.wantPassing, len(result.PassingResources), "passing resource count")
			assert.Equal(t, tt.wantFailing, len(result.FailingResources), "failing resource count")
			if tt.wantSuggestion != "" && tt.wantFailing > 0 {
				require.NotEmpty(t, result.FailingResources[0].InvalidTags)
				assert.Equal(t, tt.wantSuggestion, result.FailingResources[0].InvalidTags[0].Suggestion)
			}
		})
	}
}

func TestEvaluateAgainstResources_RegexRequirement(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "owner-policy",
		Name: "Owner Tag Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:        "Owner",
				Type:       eventpb.TagPolicyRequirement_REGEX,
				ValueRegex: `^team-[a-z]+$`,
				Mandatory:  true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	tests := []struct {
		name        string
		resources   []*provider.Resource
		wantPassing int
		wantFailing int
	}{
		{
			name: "valid owner passes regex",
			resources: []*provider.Resource{
				mkResource("aws_s3_bucket.data", "aws_s3_bucket", []*provider.Tag{
					mkTag("Owner", "team-platform"),
				}),
			},
			wantPassing: 1,
			wantFailing: 0,
		},
		{
			name: "invalid owner fails regex",
			resources: []*provider.Resource{
				mkResource("aws_s3_bucket.data", "aws_s3_bucket", []*provider.Tag{
					mkTag("Owner", "john"),
				}),
			},
			wantPassing: 0,
			wantFailing: 1,
		},
		{
			name: "owner with uppercase fails regex",
			resources: []*provider.Resource{
				mkResource("aws_s3_bucket.data", "aws_s3_bucket", []*provider.Tag{
					mkTag("Owner", "team-Platform"),
				}),
			},
			wantPassing: 0,
			wantFailing: 1,
		},
		{
			name: "missing mandatory owner tag fails",
			resources: []*provider.Resource{
				mkResource("aws_s3_bucket.data", "aws_s3_bucket", []*provider.Tag{
					mkTag("Name", "data-bucket"),
				}),
			},
			wantPassing: 0,
			wantFailing: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := policies.EvaluateAgainstResources(tt.resources, project)
			require.Len(t, results, 1)
			assert.Equal(t, tt.wantPassing, len(results[0].PassingResources))
			assert.Equal(t, tt.wantFailing, len(results[0].FailingResources))
		})
	}
}

func TestEvaluateAgainstResources_JSRegexWithFlags(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "regex-flags",
		Name: "Case-Insensitive Regex",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:        "Env",
				Type:       eventpb.TagPolicyRequirement_REGEX,
				ValueRegex: `/^(production|staging|test)$/i`,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	tests := []struct {
		name        string
		value       string
		wantPassing int
		wantFailing int
	}{
		{"lowercase matches", "production", 1, 0},
		{"uppercase matches with i flag", "PRODUCTION", 1, 0},
		{"mixed case matches", "Staging", 1, 0},
		{"invalid value fails", "dev", 0, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources := []*provider.Resource{
				mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
					mkTag("Env", tt.value),
				}),
			}
			results := policies.EvaluateAgainstResources(resources, project)
			require.Len(t, results, 1)
			assert.Equal(t, tt.wantPassing, len(results[0].PassingResources))
			assert.Equal(t, tt.wantFailing, len(results[0].FailingResources))
		})
	}
}

func TestEvaluateAgainstResources_AnyRequirement(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "any-policy",
		Name: "Any Value Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:       "CostCenter",
				Type:      eventpb.TagPolicyRequirement_ANY,
				Mandatory: true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	t.Run("any value passes", func(t *testing.T) {
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
				mkTag("CostCenter", "anything-goes"),
			}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Len(t, results[0].PassingResources, 1)
		assert.Empty(t, results[0].FailingResources)
	})

	t.Run("missing mandatory any tag fails", func(t *testing.T) {
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
				mkTag("Env", "prod"),
			}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Empty(t, results[0].PassingResources)
		assert.Len(t, results[0].FailingResources, 1)
		assert.Contains(t, results[0].FailingResources[0].MissingMandatoryTags, "CostCenter")
	})
}

func TestEvaluateAgainstResources_NonTaggableResourcesSkipped(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "env-policy",
		Name: "Env Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:       "Env",
				Type:      eventpb.TagPolicyRequirement_ANY,
				Mandatory: true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	resources := []*provider.Resource{
		{
			Name: "aws_iam_policy.readonly",
			Type: "aws_iam_policy",
			Tagging: &provider.Tagging{
				SupportsTags: false,
			},
		},
		{
			Name:    "data.aws_ami.latest",
			Type:    "data.aws_ami",
			Tagging: nil,
		},
	}

	results := policies.EvaluateAgainstResources(resources, project)
	assert.Empty(t, results, "non-taggable resources should produce no results")
}

func TestEvaluateAgainstResources_ResourceFilter(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "ec2-only",
		Name: "EC2 Only Policy",
		ResourceFilter: &eventpb.StringFilter{
			Include: []string{"aws_instance"},
		},
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:       "Env",
				Type:      eventpb.TagPolicyRequirement_ANY,
				Mandatory: true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	resources := []*provider.Resource{
		mkResource("aws_instance.web", "aws_instance", []*provider.Tag{mkTag("Env", "prod")}),
		mkResource("aws_s3_bucket.data", "aws_s3_bucket", []*provider.Tag{mkTag("Env", "prod")}),
	}

	results := policies.EvaluateAgainstResources(resources, project)
	require.Len(t, results, 1)
	assert.Len(t, results[0].PassingResources, 1)
	assert.Equal(t, "aws_instance.web", results[0].PassingResources[0].Address)
}

func TestEvaluateAgainstResources_BranchFilter(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "main-only",
		Name: "Main Branch Only",
		BranchFilter: &eventpb.StringFilter{
			Include: []string{"main"},
		},
		Requirements: []*eventpb.TagPolicyRequirement{
			{Key: "Env", Type: eventpb.TagPolicyRequirement_ANY, Mandatory: true},
		},
	}
	policies := TagPolicies{policy}

	t.Run("matching branch applies policy", func(t *testing.T) {
		project := &provider.ProjectInfo{Name: "proj", BranchName: "main"}
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{mkTag("Env", "prod")}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		assert.Len(t, results, 1)
	})

	t.Run("non-matching branch skips policy", func(t *testing.T) {
		project := &provider.ProjectInfo{Name: "proj", BranchName: "feature-branch"}
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{mkTag("Env", "prod")}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		assert.Empty(t, results)
	})
}

func TestEvaluateAgainstResources_ProjectFilter(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "project-filter",
		Name: "Project Filter Policy",
		ProjectFilter: &eventpb.StringFilter{
			Include: []string{"infra-*"},
		},
		Requirements: []*eventpb.TagPolicyRequirement{
			{Key: "Env", Type: eventpb.TagPolicyRequirement_ANY, Mandatory: true},
		},
	}
	policies := TagPolicies{policy}
	resources := []*provider.Resource{
		mkResource("aws_instance.web", "aws_instance", []*provider.Tag{mkTag("Env", "prod")}),
	}

	t.Run("matching project applies policy", func(t *testing.T) {
		project := &provider.ProjectInfo{Name: "infra-core", BranchName: "main"}
		results := policies.EvaluateAgainstResources(resources, project)
		assert.Len(t, results, 1)
	})

	t.Run("non-matching project skips policy", func(t *testing.T) {
		project := &provider.ProjectInfo{Name: "app-backend", BranchName: "main"}
		results := policies.EvaluateAgainstResources(resources, project)
		assert.Empty(t, results)
	})
}

func TestEvaluateAgainstResources_MultipleRequirements(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "multi-req",
		Name: "Multiple Requirements",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:           "Env",
				Type:          eventpb.TagPolicyRequirement_LIST,
				AllowedValues: []string{"production", "staging"},
				Mandatory:     true,
			},
			{
				Key:        "Owner",
				Type:       eventpb.TagPolicyRequirement_REGEX,
				ValueRegex: `^team-[a-z]+$`,
				Mandatory:  true,
			},
			{
				Key:       "CostCenter",
				Type:      eventpb.TagPolicyRequirement_ANY,
				Mandatory: true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	t.Run("all requirements met passes", func(t *testing.T) {
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
				mkTag("Env", "production"),
				mkTag("Owner", "team-platform"),
				mkTag("CostCenter", "12345"),
			}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Len(t, results[0].PassingResources, 1)
		assert.Empty(t, results[0].FailingResources)
	})

	t.Run("one requirement fails means resource fails", func(t *testing.T) {
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
				mkTag("Env", "production"),
				mkTag("Owner", "not-a-team"),
				mkTag("CostCenter", "12345"),
			}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Empty(t, results[0].PassingResources)
		assert.Len(t, results[0].FailingResources, 1)
	})

	t.Run("all requirements fail", func(t *testing.T) {
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
				mkTag("Name", "my-server"),
			}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Empty(t, results[0].PassingResources)
		assert.Len(t, results[0].FailingResources, 1)
		failing := results[0].FailingResources[0]
		assert.Contains(t, failing.MissingMandatoryTags, "Owner")
		assert.Contains(t, failing.MissingMandatoryTags, "CostCenter")
		// Env is LIST+mandatory so it appears as an InvalidTag with MissingMandatory=true
		require.NotEmpty(t, failing.InvalidTags)
		var found bool
		for _, it := range failing.InvalidTags {
			if it.Key == "Env" && it.MissingMandatory {
				found = true
			}
		}
		assert.True(t, found, "expected Env to appear as a missing mandatory invalid tag")
	})
}

func TestEvaluateAgainstResources_NonMandatoryTagMissing(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "optional-env",
		Name: "Optional Env",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:           "Env",
				Type:          eventpb.TagPolicyRequirement_LIST,
				AllowedValues: []string{"production", "staging"},
				Mandatory:     false,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	resources := []*provider.Resource{
		mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
			mkTag("Name", "web"),
		}),
	}

	results := policies.EvaluateAgainstResources(resources, project)
	require.Len(t, results, 1)
	assert.Len(t, results[0].PassingResources, 1, "resource with missing non-mandatory tag should pass")
	assert.Empty(t, results[0].FailingResources)
}

func TestEvaluateAgainstResources_DefaultTags(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "env-policy",
		Name: "Env Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:           "Env",
				Type:          eventpb.TagPolicyRequirement_LIST,
				AllowedValues: []string{"production", "staging"},
				Mandatory:     true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	t.Run("default tag with valid value passes", func(t *testing.T) {
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
				mkDefaultTag("Env", "production"),
			}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Len(t, results[0].PassingResources, 1)
	})

	t.Run("default tag with invalid value reports fromDefaultTags", func(t *testing.T) {
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
				mkDefaultTag("Env", "prod"),
			}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		require.Len(t, results[0].FailingResources, 1)
		require.NotEmpty(t, results[0].FailingResources[0].InvalidTags)
		assert.True(t, results[0].FailingResources[0].InvalidTags[0].FromDefaultTags)
	})
}

func TestEvaluateAgainstResources_ChildResources(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "env-policy",
		Name: "Env Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:           "Env",
				Type:          eventpb.TagPolicyRequirement_LIST,
				AllowedValues: []string{"production", "staging"},
				Mandatory:     true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	t.Run("child resource is also evaluated", func(t *testing.T) {
		resources := []*provider.Resource{
			{
				Name: "aws_autoscaling_group.web",
				Type: "aws_autoscaling_group",
				Tagging: &provider.Tagging{
					SupportsTags: true,
					Tags:         []*provider.Tag{mkTag("Env", "production")},
				},
				ChildResources: []*provider.Resource{
					{
						Name: "tag",
						Tagging: &provider.Tagging{
							SupportsTags: true,
							Tags:         []*provider.Tag{mkTag("Env", "invalid")},
						},
					},
				},
			},
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Len(t, results[0].PassingResources, 1, "parent should pass")
		assert.Len(t, results[0].FailingResources, 1, "child should fail")
		assert.Equal(t, "aws_autoscaling_group.web.tag", results[0].FailingResources[0].Address)
	})

	t.Run("non-taggable child is skipped", func(t *testing.T) {
		resources := []*provider.Resource{
			{
				Name: "aws_autoscaling_group.web",
				Type: "aws_autoscaling_group",
				Tagging: &provider.Tagging{
					SupportsTags: true,
					Tags:         []*provider.Tag{mkTag("Env", "production")},
				},
				ChildResources: []*provider.Resource{
					{
						Name:    "launch_config",
						Tagging: nil,
					},
				},
			},
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Len(t, results[0].PassingResources, 1)
		assert.Empty(t, results[0].FailingResources)
	})
}

func TestEvaluateAgainstResources_SyntheticTags(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "env-policy",
		Name: "Env Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:       "Env",
				Type:      eventpb.TagPolicyRequirement_ANY,
				Mandatory: true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	t.Run("synthetic key skips mandatory check", func(t *testing.T) {
		resources := []*provider.Resource{
			{
				Name: "aws_instance.web",
				Type: "aws_instance",
				Tagging: &provider.Tagging{
					SupportsTags: true,
					Tags: []*provider.Tag{
						{Key: "Name", Value: "web", IsKeySynthetic: true},
					},
				},
			},
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		// When there are synthetic keys, missing mandatory tags are not reported
		assert.Len(t, results[0].PassingResources, 1)
		assert.Empty(t, results[0].FailingResources)
	})

	t.Run("synthetic value is ignored for validation", func(t *testing.T) {
		resources := []*provider.Resource{
			{
				Name: "aws_instance.web",
				Type: "aws_instance",
				Tagging: &provider.Tagging{
					SupportsTags: true,
					Tags: []*provider.Tag{
						{Key: "Env", Value: "unknown", IsValueSynthetic: true},
					},
				},
			},
		}
		listPolicy := &eventpb.TagPolicy{
			Id:   "env-list",
			Name: "Env List",
			Requirements: []*eventpb.TagPolicyRequirement{
				{
					Key:           "Env",
					Type:          eventpb.TagPolicyRequirement_LIST,
					AllowedValues: []string{"production"},
				},
			},
		}
		results := TagPolicies{listPolicy}.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Len(t, results[0].PassingResources, 1, "synthetic value should not be validated")
	})
}

func TestEvaluateAgainstResources_PropagationProblems(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "env-policy",
		Name: "Env Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{Key: "Env", Type: eventpb.TagPolicyRequirement_ANY},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	resources := []*provider.Resource{
		{
			Name: "aws_autoscaling_group.web",
			Type: "aws_autoscaling_group",
			Tagging: &provider.Tagging{
				SupportsTags: true,
				Tags:         []*provider.Tag{mkTag("Env", "production")},
				PropagationProblems: []*provider.TagPropagationProblem{
					{
						Attribute:    "propagate_at_launch",
						ActualValue:  "false",
						TagRecipient: "EC2 instances",
						ValidValues:  []string{"true"},
						AffectedTags: []string{"Env"},
					},
				},
			},
		},
	}

	results := policies.EvaluateAgainstResources(resources, project)
	require.Len(t, results, 1)
	// The resource passes the tag value check but has propagation problems
	assert.Len(t, results[0].PassingResources, 1)
}

func TestEvaluateAgainstResources_TagFilter(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "prod-only",
		Name: "Only check prod resources",
		TagFilter: &eventpb.MapFilter{
			Include: map[string]string{
				"Env": "production",
			},
		},
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:        "Owner",
				Type:       eventpb.TagPolicyRequirement_REGEX,
				ValueRegex: `^team-`,
				Mandatory:  true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	t.Run("resource matching tag filter is evaluated", func(t *testing.T) {
		resources := []*provider.Resource{
			mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
				mkTag("Env", "production"),
				mkTag("Owner", "team-core"),
			}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		require.Len(t, results, 1)
		assert.Len(t, results[0].PassingResources, 1)
	})

	t.Run("resource not matching tag filter is skipped", func(t *testing.T) {
		resources := []*provider.Resource{
			mkResource("aws_instance.dev", "aws_instance", []*provider.Tag{
				mkTag("Env", "development"),
			}),
		}
		results := policies.EvaluateAgainstResources(resources, project)
		assert.Empty(t, results, "resource not matching tag filter should be skipped entirely")
	})
}

func TestEvaluateAgainstResources_BlockPRAndPRComment(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:        "blocking-policy",
		Name:      "Blocking Policy",
		BlockPr:   true,
		PrComment: true,
		Requirements: []*eventpb.TagPolicyRequirement{
			{Key: "Env", Type: eventpb.TagPolicyRequirement_ANY, Mandatory: true},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	resources := []*provider.Resource{
		mkResource("aws_instance.web", "aws_instance", []*provider.Tag{mkTag("Env", "prod")}),
	}
	results := policies.EvaluateAgainstResources(resources, project)
	require.Len(t, results, 1)
	assert.True(t, results[0].BlockPR)
	assert.True(t, results[0].PRComment)
}

func TestEvaluateAgainstResources_MultiplePolicies(t *testing.T) {
	policies := TagPolicies{
		{
			Id:   "b-policy",
			Name: "Policy B",
			Requirements: []*eventpb.TagPolicyRequirement{
				{Key: "Env", Type: eventpb.TagPolicyRequirement_ANY, Mandatory: true},
			},
		},
		{
			Id:   "a-policy",
			Name: "Policy A",
			Requirements: []*eventpb.TagPolicyRequirement{
				{Key: "Owner", Type: eventpb.TagPolicyRequirement_ANY, Mandatory: true},
			},
		},
	}
	project := defaultProjectInfo()
	resources := []*provider.Resource{
		mkResource("aws_instance.web", "aws_instance", []*provider.Tag{
			mkTag("Env", "prod"),
			mkTag("Owner", "team-core"),
		}),
	}

	results := policies.EvaluateAgainstResources(resources, project)
	require.Len(t, results, 2)
	// Results should be sorted by TagPolicyID
	assert.Equal(t, "a-policy", results[0].TagPolicyID)
	assert.Equal(t, "b-policy", results[1].TagPolicyID)
}

func TestEvaluateAgainstResources_EmptyResources(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "env-policy",
		Name: "Env Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{Key: "Env", Type: eventpb.TagPolicyRequirement_ANY, Mandatory: true},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	results := policies.EvaluateAgainstResources([]*provider.Resource{}, project)
	assert.Empty(t, results, "no resources means no results")
}

func TestEvaluateAgainstResources_Metadata(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "env-policy",
		Name: "Env Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{Key: "Env", Type: eventpb.TagPolicyRequirement_ANY, Mandatory: true},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	resources := []*provider.Resource{
		{
			Name:         "aws_instance.web",
			Type:         "aws_instance",
			ProviderLink: "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance",
			Tagging: &provider.Tagging{
				SupportsTags:        true,
				SupportsDefaultTags: true,
				Tags:                []*provider.Tag{mkTag("Name", "web")},
			},
			Metadata: &provider.ResourceMetadata{
				Filename:  "main.tf",
				StartLine: 42,
			},
		},
	}

	results := policies.EvaluateAgainstResources(resources, project)
	require.Len(t, results, 1)
	require.Len(t, results[0].FailingResources, 1)
	failing := results[0].FailingResources[0]
	assert.Equal(t, "aws_instance", failing.ResourceType)
	assert.Equal(t, "main.tf", failing.Path)
	assert.Equal(t, 42, failing.Line)
	assert.Equal(t, "https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/instance", failing.ProviderLink)
	assert.True(t, failing.SupportsDefaultTags)
	assert.Equal(t, []string{"my-project"}, failing.ProjectNames)
}

func TestEvaluateAgainstResources_TotalCounts(t *testing.T) {
	policy := &eventpb.TagPolicy{
		Id:   "env-policy",
		Name: "Env Policy",
		Requirements: []*eventpb.TagPolicyRequirement{
			{
				Key:           "Env",
				Type:          eventpb.TagPolicyRequirement_LIST,
				AllowedValues: []string{"production", "staging"},
				Mandatory:     true,
			},
		},
	}
	policies := TagPolicies{policy}
	project := defaultProjectInfo()

	resources := []*provider.Resource{
		mkResource("aws_instance.a", "aws_instance", []*provider.Tag{mkTag("Env", "production")}),
		mkResource("aws_instance.b", "aws_instance", []*provider.Tag{mkTag("Env", "bad")}),
		mkResource("aws_instance.c", "aws_instance", []*provider.Tag{mkTag("Env", "staging")}),
	}

	results := policies.EvaluateAgainstResources(resources, project)
	require.Len(t, results, 1)
	assert.Equal(t, 3, results[0].TotalDetectedResources)
	assert.Equal(t, 3, results[0].TotalTaggableResources)
}

func TestCreateSimplifiedInvalidTagForList(t *testing.T) {
	validValues := []string{"production", "development", "staging", "test"}

	t.Run("close typo suggests correct value", func(t *testing.T) {
		result := createSimplifiedInvalidTagForList("Env", "producton", validValues, false, "")
		assert.Equal(t, "production", result.Suggestion)
	})

	t.Run("prefix match suggests correct value", func(t *testing.T) {
		result := createSimplifiedInvalidTagForList("Env", "prod", validValues, false, "")
		assert.Equal(t, "production", result.Suggestion)
	})

	t.Run("empty value returns truncated valid values", func(t *testing.T) {
		result := createSimplifiedInvalidTagForList("Env", "", validValues, false, "")
		assert.Empty(t, result.Suggestion)
		assert.Equal(t, 4, result.ValidValueCount)
	})

	t.Run("no close match returns sorted by distance", func(t *testing.T) {
		result := createSimplifiedInvalidTagForList("Env", "zzzzz", validValues, false, "")
		assert.Empty(t, result.Suggestion)
		assert.NotEmpty(t, result.ValidValues)
	})

	t.Run("from default tags is propagated", func(t *testing.T) {
		result := createSimplifiedInvalidTagForList("Env", "prod", validValues, true, "")
		assert.True(t, result.FromDefaultTags)
	})

	t.Run("custom message is propagated", func(t *testing.T) {
		result := createSimplifiedInvalidTagForList("Env", "bad", validValues, false, "Use a valid environment")
		assert.Equal(t, "Use a valid environment", result.Message)
	})

	t.Run("large valid values list is truncated", func(t *testing.T) {
		large := make([]string, 300)
		for i := range large {
			large[i] = "value-" + string(rune('a'+i%26))
		}
		result := createSimplifiedInvalidTagForList("Env", "bad", large, false, "")
		assert.True(t, result.ValidValuesTruncated)
		assert.LessOrEqual(t, len(result.ValidValues), maxStoredValidTagValues)
	})
}

func TestConvertJSRegexToGo(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		match   string
		noMatch string
	}{
		{
			name:    "plain regex without delimiters",
			input:   `^prod$`,
			match:   "prod",
			noMatch: "production",
		},
		{
			name:    "JS regex with i flag",
			input:   `/^prod$/i`,
			match:   "PROD",
			noMatch: "production",
		},
		{
			name:    "JS regex without flags",
			input:   `/^(a|b)$/`,
			match:   "a",
			noMatch: "c",
		},
		{
			name:    "JS regex with multiple flags",
			input:   `/^hello$/im`,
			match:   "hello",
			noMatch: "world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiled, err := convertJSRegexToGo(tt.input)
			require.NoError(t, err)
			assert.True(t, compiled.MatchString(tt.match), "expected %q to match %q", tt.input, tt.match)
			assert.False(t, compiled.MatchString(tt.noMatch), "expected %q not to match %q", tt.input, tt.noMatch)
		})
	}
}

func TestHasDefaultTags(t *testing.T) {
	assert.False(t, hasDefaultTags(nil))
	assert.False(t, hasDefaultTags([]*provider.Tag{}))
	assert.False(t, hasDefaultTags([]*provider.Tag{mkTag("Env", "prod")}))
	assert.True(t, hasDefaultTags([]*provider.Tag{mkDefaultTag("Env", "prod")}))
	assert.True(t, hasDefaultTags([]*provider.Tag{
		mkTag("Name", "web"),
		mkDefaultTag("Env", "prod"),
	}))
}

package event

import (
	"testing"

	eventpb "github.com/infracost/proto/gen/go/infracost/parser/event"
	"github.com/infracost/proto/gen/go/infracost/provider"
	"github.com/stretchr/testify/assert"
)

func TestMapFilter_Matches(t1 *testing.T) {
	mkTags := func(m map[string]string) []*provider.Tag {
		var tags []*provider.Tag
		for k, v := range m {
			tags = append(tags, &provider.Tag{
				Key:       k,
				Value:     v,
				IsDefault: false,
			})
		}
		return tags
	}
	type fields struct {
		MapFilter *eventpb.MapFilter
	}
	type args struct {
		tags []*provider.Tag
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "empty filter returns true",
			fields: fields{
				MapFilter: &eventpb.MapFilter{},
			},
			args: args{tags: mkTags(map[string]string{"env": "prod"})},
			want: true,
		},
		{
			name: "exclude only - all excludes match -> false",
			fields: fields{
				MapFilter: &eventpb.MapFilter{
					Exclude: map[string]string{
						"env":   "prod",
						"owner": "team-*",
					},
				},
			},
			args: args{tags: mkTags(map[string]string{"env": "prod", "owner": "team-core"})},
			want: false,
		},
		{
			name: "exclude only - one exclude does not match -> true",
			fields: fields{
				MapFilter: &eventpb.MapFilter{
					Exclude: map[string]string{
						"env":   "prod",
						"owner": "team-*",
					},
				},
			},
			args: args{tags: mkTags(map[string]string{"env": "prod", "owner": "other"})},
			want: true,
		},
		{
			name: "include only - all includes match -> true",
			fields: fields{
				MapFilter: &eventpb.MapFilter{
					Include: map[string]string{
						"env":   "prod",
						"owner": "team-core",
					},
				},
			},
			args: args{tags: mkTags(map[string]string{"env": "prod", "owner": "team-core", "region": "us-east-1"})},
			want: true,
		},
		{
			name: "include only - missing key -> false",
			fields: fields{
				MapFilter: &eventpb.MapFilter{
					Include: map[string]string{
						"env":   "prod",
						"owner": "team-core",
					},
				},
			},
			args: args{tags: mkTags(map[string]string{"env": "prod"})},
			want: false,
		},
		{
			name: "include only - wildcard pattern matches -> true",
			fields: fields{
				MapFilter: &eventpb.MapFilter{
					Include: map[string]string{
						"service": "api-*",
					},
				},
			},
			args: args{tags: mkTags(map[string]string{"service": "api-billing"})},
			want: true,
		},
		{
			name: "include only - wildcard pattern does not match -> false",
			fields: fields{
				MapFilter: &eventpb.MapFilter{
					Include: map[string]string{
						"service": "api-*",
					},
				},
			},
			args: args{tags: mkTags(map[string]string{"service": "web"})},
			want: false,
		},
		{
			name: "include+exclude - excludes fully match -> false (exclude takes precedence)",
			fields: fields{
				MapFilter: &eventpb.MapFilter{
					Include: map[string]string{
						"env": "prod",
					},
					Exclude: map[string]string{
						"env":   "prod",
						"owner": "team-core",
					},
				},
			},
			args: args{tags: mkTags(map[string]string{"env": "prod", "owner": "team-core"})},
			want: false,
		},
		{
			name: "include+exclude - excludes not fully match, includes all match -> true",
			fields: fields{
				MapFilter: &eventpb.MapFilter{
					Include: map[string]string{
						"env":   "prod",
						"owner": "team-core",
					},
					Exclude: map[string]string{
						"env":   "dev",
						"owner": "team-core",
					},
				},
			},
			args: args{tags: mkTags(map[string]string{"env": "prod", "owner": "team-core"})},
			want: true,
		},
		{
			name: "include+exclude - excludes not fully match, includes partially match -> false",
			fields: fields{
				MapFilter: &eventpb.MapFilter{
					Include: map[string]string{
						"env":   "prod",
						"owner": "team-core",
					},
					Exclude: map[string]string{
						"env": "dev",
					},
				},
			},
			args: args{tags: mkTags(map[string]string{"env": "prod"})},
			want: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &MapFilter{
				MapFilter: tt.fields.MapFilter,
			}
			assert.Equalf(t1, tt.want, t.Matches(tt.args.tags), "Matches(%v)", tt.args.tags)
		})
	}
}

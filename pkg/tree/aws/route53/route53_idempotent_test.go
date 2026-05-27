package route53

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	r := &Route53{
		Records: []Record{
			{
				Resource:  resource.Resource{ID: "rec-1"},
				AliasName: value.New("target", 0, "", nil),
			},
			{
				Resource: resource.Resource{ID: "rec-2"},
				Name:     value.New("target", 0, "", nil),
			},
		},
	}

	r.PostProcess()
	alias := r.Records[0].Relationships.AliasRecord

	r.PostProcess()
	assert.Equal(t, alias, r.Records[0].Relationships.AliasRecord)
	assert.NotNil(t, alias)
}

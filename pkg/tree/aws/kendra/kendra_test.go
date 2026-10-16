package kendra

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_LinksDataSourcesToIndex(t *testing.T) {
	k := &Kendra{
		Indexes: []Index{{Resource: resource.Resource{ID: "idx-1"}}},
		DataSources: []DataSource{
			{IndexID: value.New("idx-1", 0, "", nil)},
			{IndexID: value.New("idx-1", 0, "", nil)},
			{IndexID: value.New("other", 0, "", nil)},
		},
	}

	k.PostProcess()
	k.PostProcess()
	assert.Equal(t, &k.Indexes[0], k.DataSources[0].Relationships.Index)
	assert.True(t, k.DataSources[0].Relationships.FirstForIndex)
	assert.Equal(t, &k.Indexes[0], k.DataSources[1].Relationships.Index)
	assert.False(t, k.DataSources[1].Relationships.FirstForIndex)
	assert.Nil(t, k.DataSources[2].Relationships.Index)
}

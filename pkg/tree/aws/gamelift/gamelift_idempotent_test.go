package gamelift

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree/resource"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
)

func TestPostProcess_IsIdempotent(t *testing.T) {
	g := &GameLift{
		Builds: []Build{{Resource: resource.Resource{ID: "build-1"}}},
		Fleets: []Fleet{{BuildID: value.New("build-1", 0, "", nil)}},
	}

	g.PostProcess()
	build := g.Fleets[0].Relationships.Build

	g.PostProcess()
	assert.NotNil(t, build)
	assert.Equal(t, build, g.Fleets[0].Relationships.Build)
}

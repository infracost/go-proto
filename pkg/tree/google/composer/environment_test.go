package composer_test

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree"
	"github.com/infracost/go-proto/pkg/tree/google"
	"github.com/infracost/go-proto/pkg/tree/google/composer"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentProtoRoundTrip(t *testing.T) {
	original := &tree.Tree{
		Google: google.Google{
			Composer: composer.Composer{
				Environments: []composer.Environment{{
					ImageVersion:       value.New("composer-2.9.7-airflow-2.9.3", 0, "", nil),
					EnvironmentSize:    value.New("ENVIRONMENT_SIZE_SMALL", 0, "", nil),
					SchedulerCPU:       value.New(0.5, 0, "", nil),
					SchedulerMemoryGB:  value.New(2.0, 0, "", nil),
					SchedulerStorageGB: value.New(1.0, 0, "", nil),
					SchedulerCount:     value.New(int64(1), 0, "", nil),
					WebServerCPU:       value.New(0.5, 0, "", nil),
					WebServerMemoryGB:  value.New(2.0, 0, "", nil),
					WebServerStorageGB: value.New(1.0, 0, "", nil),
					WorkerCPU:          value.New(0.5, 0, "", nil),
					WorkerMemoryGB:     value.New(2.0, 0, "", nil),
					WorkerStorageGB:    value.New(1.0, 0, "", nil),
					WorkerMinCount:     value.New(int64(2), 0, "", nil),
				}},
			},
		},
	}

	encoded, err := original.ToProto()
	require.NoError(t, err)

	decoded, err := tree.FromProto(encoded)
	require.NoError(t, err)
	require.Len(t, decoded.Google.Composer.Environments, 1)

	env := decoded.Google.Composer.Environments[0]
	assert.Equal(t, "composer-2.9.7-airflow-2.9.3", env.ImageVersion.Value())
	assert.Equal(t, "ENVIRONMENT_SIZE_SMALL", env.EnvironmentSize.Value())
	assert.InDelta(t, 0.5, env.SchedulerCPU.Value(), 1e-9)
	assert.InDelta(t, 2.0, env.WebServerMemoryGB.Value(), 1e-9)
	assert.InDelta(t, 1.0, env.WorkerStorageGB.Value(), 1e-9)
	assert.Equal(t, int64(1), env.SchedulerCount.Value())
	assert.Equal(t, int64(2), env.WorkerMinCount.Value())
}

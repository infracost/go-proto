package redis_test

import (
	"testing"

	"github.com/infracost/go-proto/pkg/tree"
	"github.com/infracost/go-proto/pkg/tree/google"
	"github.com/infracost/go-proto/pkg/tree/google/redis"
	"github.com/infracost/go-proto/pkg/tree/value"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemorystoreInstanceProtoRoundTrip(t *testing.T) {
	original := &tree.Tree{
		Google: google.Google{
			Redis: redis.Redis{
				MemorystoreInstances: []redis.MemorystoreInstance{{
					NodeType:               value.New(redis.MemorystoreNodeTypeHighmemMedium, 0, "node_type", nil),
					ShardCount:             value.New(int64(3), 0, "shard_count", nil),
					ReplicaCount:           value.New(int64(1), 0, "replica_count", nil),
					PersistenceMode:        value.New("AOF", 0, "persistence_config.0.mode", nil),
					AutomatedBackupEnabled: value.New(true, 0, "automated_backup_config", nil),
				}},
			},
		},
	}

	encoded, err := original.ToProto()
	require.NoError(t, err)

	decoded, err := tree.FromProto(encoded)
	require.NoError(t, err)
	require.Len(t, decoded.Google.Redis.MemorystoreInstances, 1)

	instance := decoded.Google.Redis.MemorystoreInstances[0]
	assert.Equal(t, redis.MemorystoreNodeTypeHighmemMedium, instance.NodeType.Value())
	assert.Equal(t, int64(3), instance.ShardCount.Value())
	assert.Equal(t, int64(1), instance.ReplicaCount.Value())
	assert.Equal(t, "AOF", instance.PersistenceMode.Value())
	assert.True(t, instance.AutomatedBackupEnabled.Value())
}

//go:build unit

package repository

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCapacityKeysShareRedisClusterHashTag(t *testing.T) {
	groupKey := capacityGroupSlotKey(42)
	memberKey := capacityMemberSlotKey(42, 99)
	metricsKey := capacityMetricsKey(42)

	require.Contains(t, groupKey, "{set:42}")
	require.Contains(t, memberKey, "{set:42}")
	require.Contains(t, metricsKey, "{set:42}")
	require.Equal(t, redisHashTagForTest(groupKey), redisHashTagForTest(memberKey))
	require.Equal(t, redisHashTagForTest(groupKey), redisHashTagForTest(metricsKey))
	require.NotEqual(t, redisHashTagForTest(groupKey), redisHashTagForTest(capacityGroupSlotKey(43)))
}

func redisHashTagForTest(key string) string {
	start := strings.IndexByte(key, '{')
	if start < 0 {
		return key
	}
	end := strings.IndexByte(key[start+1:], '}')
	if end < 0 || end == 0 {
		return key
	}
	return key[start+1 : start+1+end]
}

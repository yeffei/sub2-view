//go:build integration

package repository

import (
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUpstreamCapacitySnapshotRepository_StoresAndAggregates(t *testing.T) {
	ctx := t.Context()
	var setID int64
	require.NoError(t, integrationDB.QueryRowContext(ctx, `INSERT INTO upstream_account_sets (name, code, platform, enabled, shared_concurrency_limit) VALUES ('snapshot-test', 'snapshot-test', 'openai', true, 3000) RETURNING id`).Scan(&setID))
	t.Cleanup(func() {
		_, _ = integrationDB.ExecContext(ctx, `DELETE FROM upstream_account_sets WHERE id = $1`, setID)
	})
	repo := NewAccountMonitorRepository(nil, integrationDB).(*accountMonitorRepository)
	now := time.Now().UTC()
	rows := []*service.UpstreamCapacitySnapshotRow{
		{SetID: setID, CapacityLimit: 3000, CurrentConcurrency: 1200, LoadRate: 40, CheckedAt: now.Add(-4 * time.Minute)},
		{SetID: setID, CapacityLimit: 3000, CurrentConcurrency: 2100, LoadRate: 70, CheckedAt: now.Add(-2 * time.Minute)},
		{SetID: setID, CapacityLimit: 3000, CurrentConcurrency: 1800, LoadRate: 90, CheckedAt: now.Add(-time.Minute)},
	}
	require.NoError(t, repo.InsertUpstreamCapacitySnapshots(ctx, rows))
	stats, err := repo.ListUpstreamCapacitySnapshotStats(ctx, []int64{setID}, now.Add(-5*time.Minute))
	require.NoError(t, err)
	require.Equal(t, 2100, stats[setID].PeakConcurrency5m)
	require.Equal(t, 88, stats[setID].P95LoadRate5m)
}

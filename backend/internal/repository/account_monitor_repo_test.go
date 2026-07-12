package repository

import (
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUpsertPoolRuntimeWeightStatesPersistsFactorAndStreaks(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })
	now := time.Date(2026, 7, 12, 10, 0, 0, 0, time.UTC)
	mock.ExpectExec("INSERT INTO upstream_pool_runtime_weights").
		WithArgs(int64(1), int64(2), 0.75, 0.5, 0, 3, "probe_failed", now, now).
		WillReturnResult(sqlmock.NewResult(1, 1))
	repo := &accountMonitorRepository{db: db}
	err = repo.UpsertPoolRuntimeWeightStates(t.Context(), []*service.PoolRuntimeWeightState{{
		PoolID: 1, AccountID: 2, Factor: 0.75, TargetFactor: 0.5,
		UnhealthyStreak: 3, Reason: "probe_failed", LastObservedAt: now, UpdatedAt: now,
	}})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePoolAvailabilityBeforeUsesBoundedBatches(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	cutoff := time.Date(2026, 6, 12, 0, 0, 0, 0, time.UTC)
	query := regexp.QuoteMeta(accountMonitorPrunePoolAvailabilitySQL)
	mock.ExpectExec(query).
		WithArgs(cutoff, channelMonitorPruneBatchSize).
		WillReturnResult(sqlmock.NewResult(0, channelMonitorPruneBatchSize))
	mock.ExpectExec(query).
		WithArgs(cutoff, channelMonitorPruneBatchSize).
		WillReturnResult(sqlmock.NewResult(0, 7))

	repo := &accountMonitorRepository{db: db}
	deleted, err := repo.DeletePoolAvailabilityBefore(t.Context(), cutoff)
	require.NoError(t, err)
	require.Equal(t, int64(channelMonitorPruneBatchSize+7), deleted)
	require.NoError(t, mock.ExpectationsWereMet())
}

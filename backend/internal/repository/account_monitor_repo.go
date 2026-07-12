package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

type accountMonitorRepository struct {
	client *dbent.Client
	db     *sql.DB
	sql    sqlExecutor
}

func NewAccountMonitorRepository(client *dbent.Client, db *sql.DB) service.AccountMonitorRepository {
	return &accountMonitorRepository{client: client, db: db, sql: db}
}

func (r *accountMonitorRepository) InsertHistoryBatch(ctx context.Context, rows []*service.AccountMonitorHistoryRow) error {
	if len(rows) == 0 {
		return nil
	}
	exec := txAwareSQLExecutor(ctx, r.sql, r.client)
	if exec == nil {
		return fmt.Errorf("sql executor is not configured")
	}
	valueParts := make([]string, 0, len(rows))
	args := make([]any, 0, len(rows)*10)
	n := 0
	for _, row := range rows {
		if row == nil {
			continue
		}
		base := n*10 + 1
		n++
		valueParts = append(valueParts, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", base, base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8, base+9))
		args = append(args,
			row.AccountID,
			nullableInt64(row.PoolID),
			nullableInt64(row.GroupID),
			row.Provider,
			row.Model,
			row.Status,
			nullableInt(row.LatencyMs),
			nullableInt(row.PingLatencyMs),
			row.Message,
			row.CheckedAt,
		)
	}
	if len(valueParts) == 0 {
		return nil
	}
	q := `
		INSERT INTO account_monitor_histories (
		    account_id, pool_id, group_id, provider, model, status,
		    latency_ms, ping_latency_ms, message, checked_at
		)
		VALUES ` + strings.Join(valueParts, ",")
	if _, err := exec.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("insert account monitor histories: %w", err)
	}
	return nil
}

func (r *accountMonitorRepository) InsertPoolAvailabilitySnapshots(ctx context.Context, rows []*service.PoolAvailabilitySnapshotRow) error {
	if len(rows) == 0 {
		return nil
	}
	valueParts := make([]string, 0, len(rows))
	args := make([]any, 0, len(rows)*5)
	for _, row := range rows {
		if row == nil || row.PoolID <= 0 {
			continue
		}
		base := len(args) + 1
		valueParts = append(valueParts, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)", base, base+1, base+2, base+3, base+4))
		args = append(args, row.PoolID, row.Status, row.TotalMembers, row.AvailableMembers, row.CheckedAt)
	}
	if len(valueParts) == 0 {
		return nil
	}
	q := `INSERT INTO pool_availability_snapshots
		(pool_id, status, total_members, available_members, checked_at) VALUES ` + strings.Join(valueParts, ",")
	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("insert pool availability snapshots: %w", err)
	}
	return nil
}

func (r *accountMonitorRepository) ListPoolAvailabilitySince(ctx context.Context, poolIDs []int64, since time.Time) (map[int64][]*service.PoolAvailabilitySnapshotEntry, error) {
	out := make(map[int64][]*service.PoolAvailabilitySnapshotEntry, len(poolIDs))
	if len(poolIDs) == 0 {
		return out, nil
	}
	const q = `SELECT id, pool_id, status, total_members, available_members, checked_at
		FROM pool_availability_snapshots
		WHERE pool_id = ANY($1) AND checked_at >= $2
		ORDER BY pool_id, checked_at DESC`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(poolIDs), since)
	if err != nil {
		return nil, fmt.Errorf("query pool availability snapshots: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		entry := &service.PoolAvailabilitySnapshotEntry{}
		if err := rows.Scan(&entry.ID, &entry.PoolID, &entry.Status, &entry.TotalMembers, &entry.AvailableMembers, &entry.CheckedAt); err != nil {
			return nil, fmt.Errorf("scan pool availability snapshot: %w", err)
		}
		out[entry.PoolID] = append(out[entry.PoolID], entry)
	}
	return out, rows.Err()
}

func (r *accountMonitorRepository) ListPoolRuntimeWeightStates(ctx context.Context, poolIDs []int64) (map[int64]map[int64]*service.PoolRuntimeWeightState, error) {
	out := make(map[int64]map[int64]*service.PoolRuntimeWeightState, len(poolIDs))
	if len(poolIDs) == 0 {
		return out, nil
	}
	const q = `SELECT pool_id, account_id, factor, target_factor, healthy_streak,
		unhealthy_streak, reason, last_observed_at, updated_at
		FROM upstream_pool_runtime_weights
		WHERE pool_id = ANY($1)`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(poolIDs))
	if err != nil {
		return nil, fmt.Errorf("query pool runtime weight states: %w", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		state := &service.PoolRuntimeWeightState{}
		if err := rows.Scan(
			&state.PoolID, &state.AccountID, &state.Factor, &state.TargetFactor,
			&state.HealthyStreak, &state.UnhealthyStreak, &state.Reason,
			&state.LastObservedAt, &state.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan pool runtime weight state: %w", err)
		}
		if out[state.PoolID] == nil {
			out[state.PoolID] = map[int64]*service.PoolRuntimeWeightState{}
		}
		out[state.PoolID][state.AccountID] = state
	}
	return out, rows.Err()
}

func (r *accountMonitorRepository) UpsertPoolRuntimeWeightStates(ctx context.Context, states []*service.PoolRuntimeWeightState) error {
	valueParts := make([]string, 0, len(states))
	args := make([]any, 0, len(states)*9)
	for _, state := range states {
		if state == nil || state.PoolID <= 0 || state.AccountID <= 0 {
			continue
		}
		base := len(args) + 1
		valueParts = append(valueParts, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", base, base+1, base+2, base+3, base+4, base+5, base+6, base+7, base+8))
		args = append(args, state.PoolID, state.AccountID, state.Factor, state.TargetFactor,
			state.HealthyStreak, state.UnhealthyStreak, state.Reason, state.LastObservedAt, state.UpdatedAt)
	}
	if len(valueParts) == 0 {
		return nil
	}
	q := `INSERT INTO upstream_pool_runtime_weights (
		pool_id, account_id, factor, target_factor, healthy_streak, unhealthy_streak,
		reason, last_observed_at, updated_at
	) VALUES ` + strings.Join(valueParts, ",") + `
	ON CONFLICT (pool_id, account_id) DO UPDATE SET
		factor = EXCLUDED.factor,
		target_factor = EXCLUDED.target_factor,
		healthy_streak = EXCLUDED.healthy_streak,
		unhealthy_streak = EXCLUDED.unhealthy_streak,
		reason = EXCLUDED.reason,
		last_observed_at = EXCLUDED.last_observed_at,
		updated_at = EXCLUDED.updated_at`
	if _, err := r.db.ExecContext(ctx, q, args...); err != nil {
		return fmt.Errorf("upsert pool runtime weight states: %w", err)
	}
	return nil
}

func (r *accountMonitorRepository) ListLatestForAccountIDs(ctx context.Context, ids []int64) (map[int64][]*service.AccountMonitorLatest, error) {
	out := make(map[int64][]*service.AccountMonitorLatest, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	const q = `
		SELECT DISTINCT ON (account_id, model)
		    account_id, provider, model, status, latency_ms, ping_latency_ms, checked_at
		FROM account_monitor_histories
		WHERE account_id = ANY($1)
		ORDER BY account_id, model, checked_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(ids))
	if err != nil {
		return nil, fmt.Errorf("query account monitor latest: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var accountID int64
		latest := &service.AccountMonitorLatest{}
		var latency, ping sql.NullInt64
		if err := rows.Scan(&accountID, &latest.Provider, &latest.Model, &latest.Status, &latency, &ping, &latest.CheckedAt); err != nil {
			return nil, fmt.Errorf("scan account monitor latest: %w", err)
		}
		latest.AccountID = accountID
		assignNullInt(&latest.LatencyMs, latency)
		assignNullInt(&latest.PingLatencyMs, ping)
		out[accountID] = append(out[accountID], latest)
	}
	return out, rows.Err()
}

func (r *accountMonitorRepository) ComputeAvailabilityForAccounts(ctx context.Context, ids []int64, windowDays int) (map[int64][]*service.AccountMonitorAvailability, error) {
	out := make(map[int64][]*service.AccountMonitorAvailability, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	if windowDays <= 0 {
		windowDays = 7
	}
	const q = `
		SELECT account_id,
		       model,
		       COUNT(*)                                                     AS total,
		       COUNT(*) FILTER (WHERE status IN ('operational','degraded')) AS ok,
		       CASE WHEN COUNT(latency_ms) > 0
		            THEN SUM(latency_ms) FILTER (WHERE latency_ms IS NOT NULL)::float8 / COUNT(latency_ms)
		            ELSE NULL END                                           AS avg_latency_ms
		FROM account_monitor_histories
		WHERE account_id = ANY($1)
		  AND checked_at >= NOW() - ($2::int || ' days')::interval
		GROUP BY account_id, model
	`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(ids), windowDays)
	if err != nil {
		return nil, fmt.Errorf("query account monitor availability: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var accountID int64
		row := &service.AccountMonitorAvailability{WindowDays: windowDays}
		var avgLatency sql.NullFloat64
		if err := rows.Scan(&accountID, &row.Model, &row.TotalChecks, &row.OperationalChecks, &avgLatency); err != nil {
			return nil, fmt.Errorf("scan account monitor availability: %w", err)
		}
		row.AccountID = accountID
		finalizeAccountAvailabilityRow(row, avgLatency)
		out[accountID] = append(out[accountID], row)
	}
	return out, rows.Err()
}

func (r *accountMonitorRepository) ListHistorySinceForAccounts(
	ctx context.Context,
	ids []int64,
	primaryModels map[int64]string,
	since time.Time,
) (map[int64][]*service.AccountMonitorHistoryEntry, error) {
	out := make(map[int64][]*service.AccountMonitorHistoryEntry, len(ids))
	pairIDs, pairModels := buildAccountModelPairs(ids, primaryModels)
	if len(pairIDs) == 0 {
		return r.listHistorySinceForAccountsAllModels(ctx, ids, since)
	}
	const q = `
		WITH targets AS (
		    SELECT unnest($1::bigint[]) AS account_id,
		           unnest($2::text[])   AS model
		)
		SELECT h.id, h.account_id, h.pool_id, h.group_id, h.provider, h.model,
		       h.status, h.latency_ms, h.ping_latency_ms, h.message, h.checked_at
		FROM account_monitor_histories h
		JOIN targets t
		  ON t.account_id = h.account_id AND t.model = h.model
		WHERE h.checked_at >= $3
		ORDER BY h.account_id, h.checked_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(pairIDs), pq.Array(pairModels), since)
	if err != nil {
		return nil, fmt.Errorf("query account monitor history window: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		entry := &service.AccountMonitorHistoryEntry{}
		var poolID, groupID, latency, ping sql.NullInt64
		if err := rows.Scan(&entry.ID, &entry.AccountID, &poolID, &groupID, &entry.Provider, &entry.Model, &entry.Status, &latency, &ping, &entry.Message, &entry.CheckedAt); err != nil {
			return nil, fmt.Errorf("scan account monitor history window: %w", err)
		}
		assignNullInt64(&entry.PoolID, poolID)
		assignNullInt64(&entry.GroupID, groupID)
		assignNullInt(&entry.LatencyMs, latency)
		assignNullInt(&entry.PingLatencyMs, ping)
		out[entry.AccountID] = append(out[entry.AccountID], entry)
	}
	return out, rows.Err()
}

func (r *accountMonitorRepository) listHistorySinceForAccountsAllModels(
	ctx context.Context,
	ids []int64,
	since time.Time,
) (map[int64][]*service.AccountMonitorHistoryEntry, error) {
	out := make(map[int64][]*service.AccountMonitorHistoryEntry, len(ids))
	if len(ids) == 0 {
		return out, nil
	}
	const q = `
		SELECT id, account_id, pool_id, group_id, provider, model,
		       status, latency_ms, ping_latency_ms, message, checked_at
		FROM account_monitor_histories
		WHERE account_id = ANY($1)
		  AND checked_at >= $2
		ORDER BY account_id, checked_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, pq.Array(ids), since)
	if err != nil {
		return nil, fmt.Errorf("query account monitor history all models: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		entry := &service.AccountMonitorHistoryEntry{}
		var poolID, groupID, latency, ping sql.NullInt64
		if err := rows.Scan(&entry.ID, &entry.AccountID, &poolID, &groupID, &entry.Provider, &entry.Model, &entry.Status, &latency, &ping, &entry.Message, &entry.CheckedAt); err != nil {
			return nil, fmt.Errorf("scan account monitor history all models: %w", err)
		}
		assignNullInt64(&entry.PoolID, poolID)
		assignNullInt64(&entry.GroupID, groupID)
		assignNullInt(&entry.LatencyMs, latency)
		assignNullInt(&entry.PingLatencyMs, ping)
		out[entry.AccountID] = append(out[entry.AccountID], entry)
	}
	return out, rows.Err()
}

func (r *accountMonitorRepository) DeleteHistoryBefore(ctx context.Context, before time.Time) (int64, error) {
	var total int64
	for {
		res, err := r.db.ExecContext(ctx, accountMonitorPruneHistorySQL, before, channelMonitorPruneBatchSize)
		if err != nil {
			return total, fmt.Errorf("account monitor prune batch: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return total, fmt.Errorf("account monitor prune rows affected: %w", err)
		}
		total += affected
		if affected == 0 {
			break
		}
	}
	return total, nil
}

// DeletePoolAvailabilityBefore removes at most four small batches per daily
// maintenance pass, keeping cleanup bounded even after a long outage.
func (r *accountMonitorRepository) DeletePoolAvailabilityBefore(ctx context.Context, before time.Time) (int64, error) {
	const maxBatches = 4
	var total int64
	for batch := 0; batch < maxBatches; batch++ {
		res, err := r.db.ExecContext(ctx, accountMonitorPrunePoolAvailabilitySQL, before, channelMonitorPruneBatchSize)
		if err != nil {
			return total, fmt.Errorf("pool availability snapshot prune batch: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return total, fmt.Errorf("pool availability snapshot prune rows affected: %w", err)
		}
		total += affected
		if affected < channelMonitorPruneBatchSize {
			break
		}
	}
	return total, nil
}

func finalizeAccountAvailabilityRow(row *service.AccountMonitorAvailability, avgLatency sql.NullFloat64) {
	if row.TotalChecks > 0 {
		row.AvailabilityPct = float64(row.OperationalChecks) * 100.0 / float64(row.TotalChecks)
	}
	if avgLatency.Valid {
		v := int(avgLatency.Float64)
		row.AvgLatencyMs = &v
	}
}

func buildAccountModelPairs(ids []int64, primaryModels map[int64]string) ([]int64, []string) {
	if len(ids) == 0 || len(primaryModels) == 0 {
		return nil, nil
	}
	pairIDs := make([]int64, 0, len(ids))
	pairModels := make([]string, 0, len(ids))
	for _, id := range ids {
		model := strings.TrimSpace(primaryModels[id])
		if id <= 0 || model == "" {
			continue
		}
		pairIDs = append(pairIDs, id)
		pairModels = append(pairModels, model)
	}
	return pairIDs, pairModels
}

func assignNullInt64(dst **int64, n sql.NullInt64) {
	if !n.Valid {
		return
	}
	v := n.Int64
	*dst = &v
}

func nullableInt64(v *int64) any {
	if v == nil {
		return nil
	}
	return *v
}

func nullableInt(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}

const accountMonitorPruneHistorySQL = `
WITH batch AS (
    SELECT id FROM account_monitor_histories
    WHERE checked_at < $1
    ORDER BY id
    LIMIT $2
)
DELETE FROM account_monitor_histories
WHERE id IN (SELECT id FROM batch)
`

const accountMonitorPrunePoolAvailabilitySQL = `
WITH batch AS (
    SELECT id FROM pool_availability_snapshots
    WHERE checked_at < $1
    ORDER BY checked_at, id
    LIMIT $2
)
DELETE FROM pool_availability_snapshots
WHERE id IN (SELECT id FROM batch)
`

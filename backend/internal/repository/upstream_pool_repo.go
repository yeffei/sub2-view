package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

type upstreamPoolRepository struct {
	db *sql.DB
}

func NewUpstreamPoolRepository(_ *dbent.Client, sqlDB *sql.DB) service.UpstreamPoolRepository {
	return &upstreamPoolRepository{db: sqlDB}
}

func (r *upstreamPoolRepository) ListUpstreamPools(ctx context.Context) ([]service.UpstreamPool, error) {
	if r == nil || r.db == nil {
		return []service.UpstreamPool{}, nil
	}

	const query = `
SELECT
  id, name, code, platform, description, enabled, scheduler_mode,
  default_required_capability, default_required_transport, sticky_enabled,
  sticky_ttl_seconds, sticky_escape_enabled, sticky_escape_error_rate_threshold,
  sticky_escape_ttft_ms_threshold, load_balance_enabled, failover_enabled,
  top_k, max_failover_hops, wait_timeout_ms, max_waiting, policy_json,
  created_at, updated_at
FROM upstream_pools
ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query upstream pools: %w", err)
	}
	defer rows.Close()

	out := make([]service.UpstreamPool, 0)
	for rows.Next() {
		pool, err := scanUpstreamPool(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, pool)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate upstream pools: %w", err)
	}
	return out, nil
}

func (r *upstreamPoolRepository) GetUpstreamPoolByID(ctx context.Context, id int64) (*service.UpstreamPool, error) {
	if r == nil || r.db == nil || id <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
SELECT
  id, name, code, platform, description, enabled, scheduler_mode,
  default_required_capability, default_required_transport, sticky_enabled,
  sticky_ttl_seconds, sticky_escape_enabled, sticky_escape_error_rate_threshold,
  sticky_escape_ttft_ms_threshold, load_balance_enabled, failover_enabled,
  top_k, max_failover_hops, wait_timeout_ms, max_waiting, policy_json,
  created_at, updated_at
FROM upstream_pools
WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)
	pool, err := scanUpstreamPoolRow(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, err
	}
	return &pool, nil
}

func (r *upstreamPoolRepository) CreateUpstreamPool(ctx context.Context, input *service.UpstreamPool) (*service.UpstreamPool, error) {
	if r == nil || r.db == nil || input == nil {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
INSERT INTO upstream_pools (
  name, code, platform, description, enabled, scheduler_mode,
  default_required_capability, default_required_transport, sticky_enabled,
  sticky_ttl_seconds, sticky_escape_enabled, sticky_escape_error_rate_threshold,
  sticky_escape_ttft_ms_threshold, load_balance_enabled, failover_enabled,
  top_k, max_failover_hops, wait_timeout_ms, max_waiting, policy_json
) VALUES (
  $1,$2,$3,$4,$5,$6,
  $7,$8,$9,
  $10,$11,$12,
  $13,$14,$15,
  $16,$17,$18,$19,$20
) RETURNING id, created_at, updated_at`

	policyJSON, err := json.Marshal(normalizeAnyMapJSON(input.PolicyJSON))
	if err != nil {
		return nil, fmt.Errorf("marshal upstream pool policy: %w", err)
	}
	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.Name,
		input.Code,
		input.Platform,
		input.Description,
		input.Enabled,
		input.SchedulerMode,
		input.DefaultRequiredCapability,
		input.DefaultRequiredTransport,
		input.StickyEnabled,
		input.StickyTTLSeconds,
		input.StickyEscapeEnabled,
		input.StickyEscapeErrorRateThreshold,
		input.StickyEscapeTTFTMSThreshold,
		input.LoadBalanceEnabled,
		input.FailoverEnabled,
		input.TopK,
		input.MaxFailoverHops,
		input.WaitTimeoutMS,
		input.MaxWaiting,
		policyJSON,
	).Scan(&input.ID, &input.CreatedAt, &input.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create upstream pool: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) UpdateUpstreamPool(ctx context.Context, input *service.UpstreamPool) (*service.UpstreamPool, error) {
	if r == nil || r.db == nil || input == nil || input.ID <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
UPDATE upstream_pools SET
  name = $2,
  code = $3,
  platform = $4,
  description = $5,
  enabled = $6,
  scheduler_mode = $7,
  default_required_capability = $8,
  default_required_transport = $9,
  sticky_enabled = $10,
  sticky_ttl_seconds = $11,
  sticky_escape_enabled = $12,
  sticky_escape_error_rate_threshold = $13,
  sticky_escape_ttft_ms_threshold = $14,
  load_balance_enabled = $15,
  failover_enabled = $16,
  top_k = $17,
  max_failover_hops = $18,
  wait_timeout_ms = $19,
  max_waiting = $20,
  policy_json = $21,
  updated_at = NOW()
WHERE id = $1
RETURNING created_at, updated_at`

	policyJSON, err := json.Marshal(normalizeAnyMapJSON(input.PolicyJSON))
	if err != nil {
		return nil, fmt.Errorf("marshal upstream pool policy: %w", err)
	}
	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.ID,
		input.Name,
		input.Code,
		input.Platform,
		input.Description,
		input.Enabled,
		input.SchedulerMode,
		input.DefaultRequiredCapability,
		input.DefaultRequiredTransport,
		input.StickyEnabled,
		input.StickyTTLSeconds,
		input.StickyEscapeEnabled,
		input.StickyEscapeErrorRateThreshold,
		input.StickyEscapeTTFTMSThreshold,
		input.LoadBalanceEnabled,
		input.FailoverEnabled,
		input.TopK,
		input.MaxFailoverHops,
		input.WaitTimeoutMS,
		input.MaxWaiting,
		policyJSON,
	).Scan(&input.CreatedAt, &input.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, fmt.Errorf("update upstream pool: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) DeleteUpstreamPool(ctx context.Context, id int64) error {
	if r == nil || r.db == nil || id <= 0 {
		return service.ErrUpstreamPoolNotFound
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM upstream_pools WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete upstream pool: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete upstream pool rows affected: %w", err)
	}
	if affected == 0 {
		return service.ErrUpstreamPoolNotFound
	}
	return nil
}

func (r *upstreamPoolRepository) GetUpstreamPoolPlatformUsage(ctx context.Context, poolID int64) (service.UpstreamPoolPlatformUsage, error) {
	if r == nil || r.db == nil || poolID <= 0 {
		return service.UpstreamPoolPlatformUsage{}, nil
	}
	var usage service.UpstreamPoolPlatformUsage
	const query = `
SELECT
  (SELECT COUNT(*) FROM upstream_pool_members WHERE pool_id = $1),
  (SELECT COUNT(*) FROM upstream_pool_member_sets WHERE pool_id = $1),
  (SELECT COUNT(*) FROM upstream_pool_bindings WHERE pool_id = $1)`
	if err := r.db.QueryRowContext(ctx, query, poolID).Scan(&usage.DirectMemberCount, &usage.MemberSetCount, &usage.BindingCount); err != nil {
		return service.UpstreamPoolPlatformUsage{}, fmt.Errorf("query upstream pool platform usage: %w", err)
	}
	return usage, nil
}

func (r *upstreamPoolRepository) ListUpstreamPoolMembers(ctx context.Context, poolID int64) ([]service.UpstreamPoolMember, error) {
	if r == nil || r.db == nil || poolID <= 0 {
		return []service.UpstreamPoolMember{}, nil
	}

	const query = `
WITH direct_members AS (
  SELECT
    m.id,
    m.pool_id,
    m.account_id,
    COALESCE(a.name, '') AS account_name,
    COALESCE(a.platform, '') AS account_platform,
    COALESCE(a.type, '') AS account_type,
    m.enabled,
    m.schedulable_override,
    m.manual_drained,
    m.weight,
    m.priority_override,
    m.max_concurrency_override,
    m.notes,
    m.joined_at,
    m.updated_at,
    'direct'::text AS source_type,
    NULL::BIGINT AS source_set_id,
    ''::text AS source_set_name,
    TRUE AS editable,
    0 AS source_priority
  FROM upstream_pool_members m
  LEFT JOIN accounts a ON a.id = m.account_id AND a.deleted_at IS NULL
  WHERE m.pool_id = $1
), set_members AS (
  SELECT
    0::BIGINT AS id,
    pms.pool_id,
    a.id AS account_id,
    COALESCE(a.name, '') AS account_name,
    COALESCE(a.platform, '') AS account_platform,
    COALESCE(a.type, '') AS account_type,
    (pms.enabled AND uas.enabled) AS enabled,
    NULL::BOOLEAN AS schedulable_override,
    FALSE AS manual_drained,
    100 AS weight,
    NULL::INTEGER AS priority_override,
    NULL::INTEGER AS max_concurrency_override,
    pms.notes,
    pms.joined_at,
    pms.updated_at,
    'account_set'::text AS source_type,
    uas.id AS source_set_id,
    uas.name AS source_set_name,
    FALSE AS editable,
    1 AS source_priority
  FROM upstream_pool_member_sets pms
  JOIN upstream_account_sets uas ON uas.id = pms.set_id
  JOIN upstream_account_set_members uasm ON uasm.set_id = uas.id
  JOIN accounts a ON a.id = uasm.account_id
  WHERE pms.pool_id = $1
    AND a.deleted_at IS NULL
), ranked AS (
  SELECT DISTINCT ON (account_id)
    id, pool_id, account_id, account_name, account_platform, account_type, enabled,
    schedulable_override, manual_drained, weight, priority_override,
    max_concurrency_override, notes, joined_at, updated_at,
    source_type, source_set_id, source_set_name, editable
  FROM (
    SELECT * FROM direct_members
    UNION ALL
    SELECT * FROM set_members
  ) merged
  ORDER BY account_id, source_priority ASC, source_set_id ASC NULLS FIRST
)
SELECT
  id, pool_id, account_id, account_name, account_platform, account_type,
  enabled, schedulable_override, manual_drained,
  weight, priority_override, max_concurrency_override, notes, joined_at, updated_at,
  source_type, source_set_id, source_set_name, editable
FROM ranked
ORDER BY account_name ASC, account_id ASC`

	rows, err := r.db.QueryContext(ctx, query, poolID)
	if err != nil {
		return nil, fmt.Errorf("query upstream pool members: %w", err)
	}
	defer rows.Close()

	out := make([]service.UpstreamPoolMember, 0)
	for rows.Next() {
		var member service.UpstreamPoolMember
		var schedulableOverride sql.NullBool
		var priorityOverride sql.NullInt64
		var maxConcurrencyOverride sql.NullInt64
		var sourceSetID sql.NullInt64
		if err := rows.Scan(
			&member.ID,
			&member.PoolID,
			&member.AccountID,
			&member.AccountName,
			&member.AccountPlatform,
			&member.AccountType,
			&member.Enabled,
			&schedulableOverride,
			&member.ManualDrained,
			&member.Weight,
			&priorityOverride,
			&maxConcurrencyOverride,
			&member.Notes,
			&member.JoinedAt,
			&member.UpdatedAt,
			&member.SourceType,
			&sourceSetID,
			&member.SourceSetName,
			&member.Editable,
		); err != nil {
			return nil, fmt.Errorf("scan upstream pool member: %w", err)
		}
		if schedulableOverride.Valid {
			member.SchedulableOverride = &schedulableOverride.Bool
		}
		if priorityOverride.Valid {
			v := int(priorityOverride.Int64)
			member.PriorityOverride = &v
		}
		if maxConcurrencyOverride.Valid {
			v := int(maxConcurrencyOverride.Int64)
			member.MaxConcurrencyOverride = &v
		}
		if sourceSetID.Valid {
			value := sourceSetID.Int64
			member.SourceSetID = &value
		}
		out = append(out, member)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate upstream pool members: %w", err)
	}
	return out, nil
}

func (r *upstreamPoolRepository) CreateUpstreamPoolMember(ctx context.Context, input *service.UpstreamPoolMember) (*service.UpstreamPoolMember, error) {
	if r == nil || r.db == nil || input == nil {
		return nil, service.ErrUpstreamPoolNotFound
	}
	const query = `
INSERT INTO upstream_pool_members (
  pool_id, account_id, enabled, schedulable_override, manual_drained,
  weight, priority_override, max_concurrency_override, notes
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
RETURNING id, joined_at, updated_at`
	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.PoolID,
		input.AccountID,
		input.Enabled,
		input.SchedulableOverride,
		input.ManualDrained,
		input.Weight,
		input.PriorityOverride,
		input.MaxConcurrencyOverride,
		input.Notes,
	).Scan(&input.ID, &input.JoinedAt, &input.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create upstream pool member: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) UpdateUpstreamPoolMember(ctx context.Context, input *service.UpstreamPoolMember) (*service.UpstreamPoolMember, error) {
	if r == nil || r.db == nil || input == nil || input.ID <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}
	const query = `
UPDATE upstream_pool_members SET
  enabled = $2,
  schedulable_override = $3,
  manual_drained = $4,
  weight = $5,
  priority_override = $6,
  max_concurrency_override = $7,
  notes = $8,
  updated_at = NOW()
WHERE id = $1
RETURNING joined_at, updated_at`
	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.ID,
		input.Enabled,
		input.SchedulableOverride,
		input.ManualDrained,
		input.Weight,
		input.PriorityOverride,
		input.MaxConcurrencyOverride,
		input.Notes,
	).Scan(&input.JoinedAt, &input.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, fmt.Errorf("update upstream pool member: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) GetUpstreamPoolMemberByID(ctx context.Context, id int64) (*service.UpstreamPoolMember, error) {
	if r == nil || r.db == nil || id <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
SELECT
  m.id, m.pool_id, m.account_id, COALESCE(a.name, '') AS account_name, COALESCE(a.platform, '') AS account_platform, COALESCE(a.type, '') AS account_type,
  m.enabled, m.schedulable_override, m.manual_drained,
  m.weight, m.priority_override, m.max_concurrency_override, m.notes, m.joined_at, m.updated_at
FROM upstream_pool_members m
LEFT JOIN accounts a ON a.id = m.account_id
WHERE m.id = $1`

	var member service.UpstreamPoolMember
	var schedulableOverride sql.NullBool
	var priorityOverride sql.NullInt64
	var maxConcurrencyOverride sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&member.ID,
		&member.PoolID,
		&member.AccountID,
		&member.AccountName,
		&member.AccountPlatform,
		&member.AccountType,
		&member.Enabled,
		&schedulableOverride,
		&member.ManualDrained,
		&member.Weight,
		&priorityOverride,
		&maxConcurrencyOverride,
		&member.Notes,
		&member.JoinedAt,
		&member.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, fmt.Errorf("query upstream pool member by id: %w", err)
	}
	if schedulableOverride.Valid {
		member.SchedulableOverride = &schedulableOverride.Bool
	}
	if priorityOverride.Valid {
		v := int(priorityOverride.Int64)
		member.PriorityOverride = &v
	}
	if maxConcurrencyOverride.Valid {
		v := int(maxConcurrencyOverride.Int64)
		member.MaxConcurrencyOverride = &v
	}
	member.SourceType = "direct"
	member.Editable = true
	return &member, nil
}

func (r *upstreamPoolRepository) DeleteUpstreamPoolMember(ctx context.Context, id int64) error {
	if r == nil || r.db == nil || id <= 0 {
		return service.ErrUpstreamPoolNotFound
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM upstream_pool_members WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete upstream pool member: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete upstream pool member rows affected: %w", err)
	}
	if affected == 0 {
		return service.ErrUpstreamPoolNotFound
	}
	return nil
}

func (r *upstreamPoolRepository) SyncUpstreamPoolDirectMembers(ctx context.Context, poolID int64, targets []service.UpstreamPoolMember, mode service.UpstreamPoolMemberSyncMode) (*service.UpstreamPoolMemberSyncResult, error) {
	if r == nil || r.db == nil || poolID <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}
	pool, err := r.GetUpstreamPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}
	current, err := r.ListUpstreamPoolMembers(ctx, poolID)
	if err != nil {
		return nil, err
	}
	result := service.BuildUpstreamPoolMemberSyncResult(pool, current, targets, mode)

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin upstream pool member sync: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	for _, change := range result.Creates {
		target := findUpstreamPoolMemberSyncTarget(targets, change.AccountID)
		if target.AccountID <= 0 {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
INSERT INTO upstream_pool_members (
  pool_id, account_id, enabled, schedulable_override, manual_drained,
  weight, priority_override, max_concurrency_override, notes
) VALUES ($1,$2,$3,NULL,$4,$5,NULL,NULL,$6)
ON CONFLICT (pool_id, account_id) DO NOTHING`,
			poolID,
			target.AccountID,
			target.Enabled,
			target.ManualDrained,
			target.Weight,
			target.Notes,
		); err != nil {
			return nil, fmt.Errorf("sync create upstream pool member: %w", err)
		}
	}

	if result.Mode == service.UpstreamPoolMemberSyncModeOverwriteSchedulerFields {
		for _, change := range result.Updates {
			target := findUpstreamPoolMemberSyncTarget(targets, change.AccountID)
			if target.AccountID <= 0 {
				continue
			}
			if _, err := tx.ExecContext(ctx, `
UPDATE upstream_pool_members SET
  enabled = $3,
  schedulable_override = NULL,
  manual_drained = $4,
  weight = $5,
  priority_override = NULL,
  max_concurrency_override = NULL,
  notes = $6,
  updated_at = NOW()
WHERE pool_id = $1 AND account_id = $2`,
				poolID,
				target.AccountID,
				target.Enabled,
				target.ManualDrained,
				target.Weight,
				target.Notes,
			); err != nil {
				return nil, fmt.Errorf("sync update upstream pool member: %w", err)
			}
		}
	}

	for _, change := range result.Deletes {
		if _, err := tx.ExecContext(ctx, `DELETE FROM upstream_pool_members WHERE pool_id = $1 AND account_id = $2`, poolID, change.AccountID); err != nil {
			return nil, fmt.Errorf("sync delete upstream pool member: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit upstream pool member sync: %w", err)
	}
	return result, nil
}

func (r *upstreamPoolRepository) ListUpstreamAccountSets(ctx context.Context) ([]service.UpstreamAccountSet, error) {
	if r == nil || r.db == nil {
		return []service.UpstreamAccountSet{}, nil
	}

	const query = `
SELECT
  s.id, s.name, s.code, s.platform, s.description, s.enabled,
  COUNT(m.account_id) AS account_count,
  s.created_at, s.updated_at
FROM upstream_account_sets s
LEFT JOIN upstream_account_set_members m ON m.set_id = s.id
GROUP BY s.id
ORDER BY s.id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query upstream account sets: %w", err)
	}
	defer rows.Close()

	out := make([]service.UpstreamAccountSet, 0)
	for rows.Next() {
		var item service.UpstreamAccountSet
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Code,
			&item.Platform,
			&item.Description,
			&item.Enabled,
			&item.AccountCount,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan upstream account set: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate upstream account sets: %w", err)
	}
	return out, nil
}

func (r *upstreamPoolRepository) GetUpstreamAccountSetByID(ctx context.Context, id int64) (*service.UpstreamAccountSet, error) {
	if r == nil || r.db == nil || id <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
SELECT
  s.id, s.name, s.code, s.platform, s.description, s.enabled,
  COUNT(m.account_id) AS account_count,
  s.created_at, s.updated_at
FROM upstream_account_sets s
LEFT JOIN upstream_account_set_members m ON m.set_id = s.id
WHERE s.id = $1
GROUP BY s.id`

	var item service.UpstreamAccountSet
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Code,
		&item.Platform,
		&item.Description,
		&item.Enabled,
		&item.AccountCount,
		&item.CreatedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, fmt.Errorf("query upstream account set by id: %w", err)
	}
	return &item, nil
}

func (r *upstreamPoolRepository) CreateUpstreamAccountSet(ctx context.Context, input *service.UpstreamAccountSet) (*service.UpstreamAccountSet, error) {
	if r == nil || r.db == nil || input == nil {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
INSERT INTO upstream_account_sets (
  name, code, platform, description, enabled
) VALUES ($1,$2,$3,$4,$5)
RETURNING id, created_at, updated_at`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.Name,
		input.Code,
		input.Platform,
		input.Description,
		input.Enabled,
	).Scan(&input.ID, &input.CreatedAt, &input.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create upstream account set: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) UpdateUpstreamAccountSet(ctx context.Context, input *service.UpstreamAccountSet) (*service.UpstreamAccountSet, error) {
	if r == nil || r.db == nil || input == nil || input.ID <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
UPDATE upstream_account_sets SET
  name = $2,
  code = $3,
  platform = $4,
  description = $5,
  enabled = $6,
  updated_at = NOW()
WHERE id = $1
RETURNING created_at, updated_at`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.ID,
		input.Name,
		input.Code,
		input.Platform,
		input.Description,
		input.Enabled,
	).Scan(&input.CreatedAt, &input.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, fmt.Errorf("update upstream account set: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) DeleteUpstreamAccountSet(ctx context.Context, id int64) error {
	if r == nil || r.db == nil || id <= 0 {
		return service.ErrUpstreamPoolNotFound
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM upstream_account_sets WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete upstream account set: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete upstream account set rows affected: %w", err)
	}
	if affected == 0 {
		return service.ErrUpstreamPoolNotFound
	}
	return nil
}

func (r *upstreamPoolRepository) GetUpstreamAccountSetPlatformUsage(ctx context.Context, setID int64) (service.UpstreamAccountSetPlatformUsage, error) {
	if r == nil || r.db == nil || setID <= 0 {
		return service.UpstreamAccountSetPlatformUsage{}, nil
	}
	var usage service.UpstreamAccountSetPlatformUsage
	const query = `
SELECT
  (SELECT COUNT(*) FROM upstream_account_set_members WHERE set_id = $1),
  (SELECT COUNT(*) FROM upstream_pool_member_sets WHERE set_id = $1)`
	if err := r.db.QueryRowContext(ctx, query, setID).Scan(&usage.MemberCount, &usage.PoolBindingCount); err != nil {
		return service.UpstreamAccountSetPlatformUsage{}, fmt.Errorf("query upstream account set platform usage: %w", err)
	}
	return usage, nil
}

func (r *upstreamPoolRepository) ListUpstreamAccountSetMembers(ctx context.Context, setID int64) ([]service.UpstreamAccountSetMember, error) {
	if r == nil || r.db == nil || setID <= 0 {
		return []service.UpstreamAccountSetMember{}, nil
	}

	const query = `
SELECT
  m.set_id,
  m.account_id,
  COALESCE(a.name, '') AS account_name,
  COALESCE(a.platform, '') AS account_platform,
  COALESCE(a.type, '') AS account_type,
  COALESCE(a.status, '') AS account_status,
  m.added_at
FROM upstream_account_set_members m
LEFT JOIN accounts a ON a.id = m.account_id
WHERE m.set_id = $1
ORDER BY account_name ASC, m.account_id ASC`

	rows, err := r.db.QueryContext(ctx, query, setID)
	if err != nil {
		return nil, fmt.Errorf("query upstream account set members: %w", err)
	}
	defer rows.Close()

	out := make([]service.UpstreamAccountSetMember, 0)
	for rows.Next() {
		var item service.UpstreamAccountSetMember
		if err := rows.Scan(
			&item.SetID,
			&item.AccountID,
			&item.AccountName,
			&item.AccountPlatform,
			&item.AccountType,
			&item.AccountStatus,
			&item.AddedAt,
		); err != nil {
			return nil, fmt.Errorf("scan upstream account set member: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate upstream account set members: %w", err)
	}
	return out, nil
}

func (r *upstreamPoolRepository) AddUpstreamAccountSetMembers(ctx context.Context, setID int64, accountIDs []int64) error {
	if r == nil || r.db == nil || setID <= 0 || len(accountIDs) == 0 {
		return nil
	}
	const query = `
INSERT INTO upstream_account_set_members (set_id, account_id)
SELECT $1, UNNEST($2::BIGINT[])
ON CONFLICT (set_id, account_id) DO NOTHING`
	if _, err := r.db.ExecContext(ctx, query, setID, pq.Array(accountIDs)); err != nil {
		return fmt.Errorf("add upstream account set members: %w", err)
	}
	return nil
}

func (r *upstreamPoolRepository) DeleteUpstreamAccountSetMember(ctx context.Context, setID, accountID int64) error {
	if r == nil || r.db == nil || setID <= 0 || accountID <= 0 {
		return service.ErrUpstreamPoolNotFound
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM upstream_account_set_members WHERE set_id = $1 AND account_id = $2`, setID, accountID)
	if err != nil {
		return fmt.Errorf("delete upstream account set member: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete upstream account set member rows affected: %w", err)
	}
	if affected == 0 {
		return service.ErrUpstreamPoolNotFound
	}
	return nil
}

func (r *upstreamPoolRepository) ListUpstreamPoolMemberSets(ctx context.Context, poolID int64) ([]service.UpstreamPoolMemberSet, error) {
	if r == nil || r.db == nil || poolID <= 0 {
		return []service.UpstreamPoolMemberSet{}, nil
	}

	const query = `
SELECT
  pms.id, pms.pool_id, pms.set_id,
  COALESCE(s.name, '') AS set_name,
  COALESCE(s.code, '') AS set_code,
  COALESCE(s.platform, '') AS set_platform,
  pms.enabled, pms.notes, pms.joined_at, pms.updated_at
FROM upstream_pool_member_sets pms
LEFT JOIN upstream_account_sets s ON s.id = pms.set_id
WHERE pms.pool_id = $1
ORDER BY pms.id ASC`

	rows, err := r.db.QueryContext(ctx, query, poolID)
	if err != nil {
		return nil, fmt.Errorf("query upstream pool member sets: %w", err)
	}
	defer rows.Close()

	out := make([]service.UpstreamPoolMemberSet, 0)
	for rows.Next() {
		var item service.UpstreamPoolMemberSet
		if err := rows.Scan(
			&item.ID,
			&item.PoolID,
			&item.SetID,
			&item.SetName,
			&item.SetCode,
			&item.SetPlatform,
			&item.Enabled,
			&item.Notes,
			&item.JoinedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan upstream pool member set: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate upstream pool member sets: %w", err)
	}
	return out, nil
}

func (r *upstreamPoolRepository) GetUpstreamPoolMemberSetByID(ctx context.Context, id int64) (*service.UpstreamPoolMemberSet, error) {
	if r == nil || r.db == nil || id <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
SELECT
  pms.id, pms.pool_id, pms.set_id,
  COALESCE(s.name, '') AS set_name,
  COALESCE(s.code, '') AS set_code,
  COALESCE(s.platform, '') AS set_platform,
  pms.enabled, pms.notes, pms.joined_at, pms.updated_at
FROM upstream_pool_member_sets pms
LEFT JOIN upstream_account_sets s ON s.id = pms.set_id
WHERE pms.id = $1`

	var item service.UpstreamPoolMemberSet
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.PoolID,
		&item.SetID,
		&item.SetName,
		&item.SetCode,
		&item.SetPlatform,
		&item.Enabled,
		&item.Notes,
		&item.JoinedAt,
		&item.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, fmt.Errorf("query upstream pool member set by id: %w", err)
	}
	return &item, nil
}

func (r *upstreamPoolRepository) CreateUpstreamPoolMemberSet(ctx context.Context, input *service.UpstreamPoolMemberSet) (*service.UpstreamPoolMemberSet, error) {
	if r == nil || r.db == nil || input == nil {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
INSERT INTO upstream_pool_member_sets (
  pool_id, set_id, enabled, notes
) VALUES ($1,$2,$3,$4)
RETURNING id, joined_at, updated_at`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.PoolID,
		input.SetID,
		input.Enabled,
		input.Notes,
	).Scan(&input.ID, &input.JoinedAt, &input.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create upstream pool member set: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) UpdateUpstreamPoolMemberSet(ctx context.Context, input *service.UpstreamPoolMemberSet) (*service.UpstreamPoolMemberSet, error) {
	if r == nil || r.db == nil || input == nil || input.ID <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}

	const query = `
UPDATE upstream_pool_member_sets SET
  enabled = $2,
  notes = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING joined_at, updated_at`

	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.ID,
		input.Enabled,
		input.Notes,
	).Scan(&input.JoinedAt, &input.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, fmt.Errorf("update upstream pool member set: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) DeleteUpstreamPoolMemberSet(ctx context.Context, id int64) error {
	if r == nil || r.db == nil || id <= 0 {
		return service.ErrUpstreamPoolNotFound
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM upstream_pool_member_sets WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete upstream pool member set: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete upstream pool member set rows affected: %w", err)
	}
	if affected == 0 {
		return service.ErrUpstreamPoolNotFound
	}
	return nil
}

func (r *upstreamPoolRepository) ListUpstreamPoolBindings(ctx context.Context) ([]service.UpstreamPoolBinding, error) {
	if r == nil || r.db == nil {
		return []service.UpstreamPoolBinding{}, nil
	}

	const query = `
SELECT
  b.id, b.group_id, COALESCE(g.name, '') AS group_name, COALESCE(g.platform, '') AS group_platform,
  b.pool_id, b.platform, b.models, b.request_path_scope, b.priority, b.enabled, b.created_at, b.updated_at
FROM upstream_pool_bindings b
LEFT JOIN groups g ON g.id = b.group_id
ORDER BY b.priority ASC, b.id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query upstream pool bindings: %w", err)
	}
	defer rows.Close()

	out := make([]service.UpstreamPoolBinding, 0)
	for rows.Next() {
		var binding service.UpstreamPoolBinding
		var modelsJSON []byte
		var requestPathScopeJSON []byte
		if err := rows.Scan(
			&binding.ID,
			&binding.GroupID,
			&binding.GroupName,
			&binding.GroupPlatform,
			&binding.PoolID,
			&binding.Platform,
			&modelsJSON,
			&requestPathScopeJSON,
			&binding.Priority,
			&binding.Enabled,
			&binding.CreatedAt,
			&binding.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan upstream pool binding: %w", err)
		}
		binding.Models = cloneStringSliceJSON(modelsJSON)
		binding.RequestPathScope = cloneStringSliceJSON(requestPathScopeJSON)
		out = append(out, binding)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate upstream pool bindings: %w", err)
	}
	return out, nil
}

func (r *upstreamPoolRepository) CreateUpstreamPoolBinding(ctx context.Context, input *service.UpstreamPoolBinding) (*service.UpstreamPoolBinding, error) {
	if r == nil || r.db == nil || input == nil {
		return nil, service.ErrUpstreamPoolNotFound
	}
	const query = `
INSERT INTO upstream_pool_bindings (
  group_id, pool_id, platform, models, request_path_scope, priority, enabled
) VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING id, created_at, updated_at`
	modelsJSON, err := json.Marshal(normalizeStringSliceJSON(input.Models))
	if err != nil {
		return nil, fmt.Errorf("marshal upstream pool binding models: %w", err)
	}
	pathJSON, err := json.Marshal(normalizeStringSliceJSON(input.RequestPathScope))
	if err != nil {
		return nil, fmt.Errorf("marshal upstream pool binding request_path_scope: %w", err)
	}
	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.GroupID,
		input.PoolID,
		input.Platform,
		modelsJSON,
		pathJSON,
		input.Priority,
		input.Enabled,
	).Scan(&input.ID, &input.CreatedAt, &input.UpdatedAt); err != nil {
		return nil, fmt.Errorf("create upstream pool binding: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) UpdateUpstreamPoolBinding(ctx context.Context, input *service.UpstreamPoolBinding) (*service.UpstreamPoolBinding, error) {
	if r == nil || r.db == nil || input == nil || input.ID <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}
	const query = `
UPDATE upstream_pool_bindings SET
  group_id = $2,
  pool_id = $3,
  platform = $4,
  models = $5,
  request_path_scope = $6,
  priority = $7,
  enabled = $8,
  updated_at = NOW()
WHERE id = $1
RETURNING created_at, updated_at`
	modelsJSON, err := json.Marshal(normalizeStringSliceJSON(input.Models))
	if err != nil {
		return nil, fmt.Errorf("marshal upstream pool binding models: %w", err)
	}
	pathJSON, err := json.Marshal(normalizeStringSliceJSON(input.RequestPathScope))
	if err != nil {
		return nil, fmt.Errorf("marshal upstream pool binding request_path_scope: %w", err)
	}
	if err := r.db.QueryRowContext(
		ctx,
		query,
		input.ID,
		input.GroupID,
		input.PoolID,
		input.Platform,
		modelsJSON,
		pathJSON,
		input.Priority,
		input.Enabled,
	).Scan(&input.CreatedAt, &input.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, fmt.Errorf("update upstream pool binding: %w", err)
	}
	return input, nil
}

func (r *upstreamPoolRepository) GetUpstreamPoolBindingByID(ctx context.Context, id int64) (*service.UpstreamPoolBinding, error) {
	if r == nil || r.db == nil || id <= 0 {
		return nil, service.ErrUpstreamPoolNotFound
	}
	const query = `
SELECT
  b.id, b.group_id, COALESCE(g.name, '') AS group_name, COALESCE(g.platform, '') AS group_platform,
  b.pool_id, b.platform, b.models, b.request_path_scope, b.priority, b.enabled, b.created_at, b.updated_at
FROM upstream_pool_bindings b
LEFT JOIN groups g ON g.id = b.group_id
WHERE b.id = $1`
	var binding service.UpstreamPoolBinding
	var modelsJSON []byte
	var requestPathScopeJSON []byte
	if err := r.db.QueryRowContext(ctx, query, id).Scan(
		&binding.ID,
		&binding.GroupID,
		&binding.GroupName,
		&binding.GroupPlatform,
		&binding.PoolID,
		&binding.Platform,
		&modelsJSON,
		&requestPathScopeJSON,
		&binding.Priority,
		&binding.Enabled,
		&binding.CreatedAt,
		&binding.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrUpstreamPoolNotFound
		}
		return nil, fmt.Errorf("query upstream pool binding by id: %w", err)
	}
	binding.Models = cloneStringSliceJSON(modelsJSON)
	binding.RequestPathScope = cloneStringSliceJSON(requestPathScopeJSON)
	return &binding, nil
}

func (r *upstreamPoolRepository) DeleteUpstreamPoolBinding(ctx context.Context, id int64) error {
	if r == nil || r.db == nil || id <= 0 {
		return service.ErrUpstreamPoolNotFound
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM upstream_pool_bindings WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete upstream pool binding: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete upstream pool binding rows affected: %w", err)
	}
	if affected == 0 {
		return service.ErrUpstreamPoolNotFound
	}
	return nil
}

func (r *upstreamPoolRepository) ListEnabledMemberAccountIDsByGroupAndPlatform(ctx context.Context, groupID int64, platform string) (map[int64]struct{}, error) {
	resolved, err := r.GetResolvedBindingByGroupAndPlatform(ctx, groupID, platform)
	if err != nil {
		return nil, err
	}
	if resolved == nil {
		return nil, nil
	}
	return resolved.MemberIDs, nil
}

func (r *upstreamPoolRepository) GetResolvedBindingByGroupAndPlatform(ctx context.Context, groupID int64, platform string) (*service.UpstreamPoolResolvedBinding, error) {
	if r == nil || r.db == nil || groupID <= 0 || platform == "" {
		return nil, nil
	}

	const bindingQuery = `
SELECT
  b.id,
  b.group_id,
  b.pool_id,
  b.platform,
  b.models,
  b.request_path_scope,
  b.priority,
  b.enabled,
  b.created_at,
  b.updated_at,
  p.id,
  p.name,
  p.code,
  p.platform,
  p.description,
  p.enabled,
  p.scheduler_mode,
  p.default_required_capability,
  p.default_required_transport,
  p.sticky_enabled,
  p.sticky_ttl_seconds,
  p.sticky_escape_enabled,
  p.sticky_escape_error_rate_threshold,
  p.sticky_escape_ttft_ms_threshold,
  p.load_balance_enabled,
  p.failover_enabled,
  p.top_k,
  p.max_failover_hops,
  p.wait_timeout_ms,
  p.max_waiting,
  p.policy_json,
  p.created_at,
  p.updated_at
FROM upstream_pool_bindings b
JOIN upstream_pools p ON p.id = b.pool_id
WHERE b.group_id = $1
  AND b.platform = $2
  AND b.enabled = TRUE
ORDER BY CASE WHEN p.enabled THEN 0 ELSE 1 END ASC, b.priority ASC, b.id ASC
LIMIT 1`

	var (
		binding service.UpstreamPoolBinding
		pool    service.UpstreamPool

		modelsJSON           []byte
		requestPathScopeJSON []byte
		policyJSON           []byte
	)

	err := r.db.QueryRowContext(ctx, bindingQuery, groupID, platform).Scan(
		&binding.ID,
		&binding.GroupID,
		&binding.PoolID,
		&binding.Platform,
		&modelsJSON,
		&requestPathScopeJSON,
		&binding.Priority,
		&binding.Enabled,
		&binding.CreatedAt,
		&binding.UpdatedAt,
		&pool.ID,
		&pool.Name,
		&pool.Code,
		&pool.Platform,
		&pool.Description,
		&pool.Enabled,
		&pool.SchedulerMode,
		&pool.DefaultRequiredCapability,
		&pool.DefaultRequiredTransport,
		&pool.StickyEnabled,
		&pool.StickyTTLSeconds,
		&pool.StickyEscapeEnabled,
		&pool.StickyEscapeErrorRateThreshold,
		&pool.StickyEscapeTTFTMSThreshold,
		&pool.LoadBalanceEnabled,
		&pool.FailoverEnabled,
		&pool.TopK,
		&pool.MaxFailoverHops,
		&pool.WaitTimeoutMS,
		&pool.MaxWaiting,
		&policyJSON,
		&pool.CreatedAt,
		&pool.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query upstream pool binding: %w", err)
	}

	binding.Models = cloneStringSliceJSON(modelsJSON)
	binding.RequestPathScope = cloneStringSliceJSON(requestPathScopeJSON)
	pool.PolicyJSON = cloneAnyMapJSON(policyJSON)
	pool.AccountTypeStrategy = service.UpstreamPoolAccountTypeStrategyFromPolicyJSON(pool.PolicyJSON)

	if !pool.Enabled {
		return &service.UpstreamPoolResolvedBinding{
			Binding:       &binding,
			Pool:          &pool,
			MemberIDs:     map[int64]struct{}{},
			MemberConfigs: map[int64]service.UpstreamPoolResolvedMemberConfig{},
		}, nil
	}

	const membersQuery = `
SELECT DISTINCT ON (account_id)
  account_id,
  schedulable_override,
  weight,
  priority_override,
  max_concurrency_override
FROM (
  SELECT
    m.account_id,
    m.schedulable_override,
    m.weight,
    m.priority_override,
    m.max_concurrency_override,
    0 AS source_priority,
    NULL::BIGINT AS source_set_id
  FROM upstream_pool_members m
  JOIN accounts a ON a.id = m.account_id
  WHERE m.pool_id = $1
    AND m.enabled = TRUE
    AND m.manual_drained = FALSE
    AND a.deleted_at IS NULL
  UNION ALL
  SELECT
    uasm.account_id,
    NULL::BOOLEAN AS schedulable_override,
    100 AS weight,
    NULL::INTEGER AS priority_override,
    NULL::INTEGER AS max_concurrency_override,
    1 AS source_priority,
    pms.set_id AS source_set_id
  FROM upstream_pool_member_sets pms
  JOIN upstream_account_sets uas ON uas.id = pms.set_id
  JOIN upstream_account_set_members uasm ON uasm.set_id = uas.id
  JOIN accounts a ON a.id = uasm.account_id
  WHERE pms.pool_id = $1
    AND pms.enabled = TRUE
    AND uas.enabled = TRUE
    AND a.deleted_at IS NULL
) members
ORDER BY account_id, source_priority ASC, source_set_id ASC NULLS FIRST`

	rows, err := r.db.QueryContext(ctx, membersQuery, pool.ID)
	if err != nil {
		return nil, fmt.Errorf("query upstream pool members: %w", err)
	}
	defer rows.Close()

	memberIDs := make(map[int64]struct{})
	memberConfigs := make(map[int64]service.UpstreamPoolResolvedMemberConfig)
	for rows.Next() {
		var (
			accountID               int64
			schedulableOverrideNull sql.NullBool
			weight                  int
			priorityOverrideNull    sql.NullInt64
			maxConcurrencyNull      sql.NullInt64
		)
		if err := rows.Scan(
			&accountID,
			&schedulableOverrideNull,
			&weight,
			&priorityOverrideNull,
			&maxConcurrencyNull,
		); err != nil {
			return nil, fmt.Errorf("scan upstream pool member: %w", err)
		}
		if accountID > 0 {
			memberIDs[accountID] = struct{}{}
			cfg := service.UpstreamPoolResolvedMemberConfig{
				AccountID: accountID,
				Weight:    weight,
			}
			if schedulableOverrideNull.Valid {
				value := schedulableOverrideNull.Bool
				cfg.SchedulableOverride = &value
			}
			if priorityOverrideNull.Valid {
				value := int(priorityOverrideNull.Int64)
				cfg.PriorityOverride = &value
			}
			if maxConcurrencyNull.Valid {
				value := int(maxConcurrencyNull.Int64)
				cfg.MaxConcurrencyOverride = &value
			}
			memberConfigs[accountID] = cfg
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate upstream pool members: %w", err)
	}

	return &service.UpstreamPoolResolvedBinding{
		Binding:       &binding,
		Pool:          &pool,
		MemberIDs:     memberIDs,
		MemberConfigs: memberConfigs,
	}, nil
}

func (r *upstreamPoolRepository) GetOpenAIRoutingPolicy(ctx context.Context, groupID int64) (*service.OpenAIRoutingPolicy, error) {
	resolved, err := r.GetResolvedBindingByGroupAndPlatform(ctx, groupID, service.PlatformOpenAI)
	if err != nil || resolved == nil || resolved.Binding == nil || resolved.Pool == nil {
		return &service.OpenAIRoutingPolicy{}, err
	}

	pool := resolved.Pool
	stickyEscapeEnabled := pool.StickyEscapeEnabled
	policy := &service.OpenAIRoutingPolicy{
		HasBinding:                     true,
		PoolID:                         pool.ID,
		PoolCode:                       pool.Code,
		PoolName:                       pool.Name,
		SchedulerMode:                  pool.SchedulerMode,
		AccountTypeStrategy:            pool.AccountTypeStrategy,
		StickyEnabled:                  pool.StickyEnabled,
		StickyEscapeEnabled:            &stickyEscapeEnabled,
		StickyEscapeErrorRateThreshold: pool.StickyEscapeErrorRateThreshold,
		StickyEscapeTTFTMSThreshold:    pool.StickyEscapeTTFTMSThreshold,
		LoadBalanceEnabled:             pool.LoadBalanceEnabled,
		FailoverEnabled:                pool.FailoverEnabled,
		TopK:                           pool.TopK,
		MaxFailoverHops:                pool.MaxFailoverHops,
		WaitTimeout:                    time.Duration(pool.WaitTimeoutMS) * time.Millisecond,
		MaxWaiting:                     pool.MaxWaiting,
	}
	service.ApplyOpenAIRoutingPolicyJSON(policy, pool.PolicyJSON)
	return policy, nil
}

func scanUpstreamPoolRow(scanner interface {
	Scan(dest ...any) error
}) (service.UpstreamPool, error) {
	var pool service.UpstreamPool
	var policyJSON []byte
	if err := scanner.Scan(
		&pool.ID,
		&pool.Name,
		&pool.Code,
		&pool.Platform,
		&pool.Description,
		&pool.Enabled,
		&pool.SchedulerMode,
		&pool.DefaultRequiredCapability,
		&pool.DefaultRequiredTransport,
		&pool.StickyEnabled,
		&pool.StickyTTLSeconds,
		&pool.StickyEscapeEnabled,
		&pool.StickyEscapeErrorRateThreshold,
		&pool.StickyEscapeTTFTMSThreshold,
		&pool.LoadBalanceEnabled,
		&pool.FailoverEnabled,
		&pool.TopK,
		&pool.MaxFailoverHops,
		&pool.WaitTimeoutMS,
		&pool.MaxWaiting,
		&policyJSON,
		&pool.CreatedAt,
		&pool.UpdatedAt,
	); err != nil {
		return service.UpstreamPool{}, err
	}
	pool.PolicyJSON = cloneAnyMapJSON(policyJSON)
	pool.AccountTypeStrategy = service.UpstreamPoolAccountTypeStrategyFromPolicyJSON(pool.PolicyJSON)
	return pool, nil
}

func scanUpstreamPool(rows *sql.Rows) (service.UpstreamPool, error) {
	pool, err := scanUpstreamPoolRow(rows)
	if err != nil {
		return service.UpstreamPool{}, fmt.Errorf("scan upstream pool: %w", err)
	}
	return pool, nil
}

func cloneStringSliceJSON(raw []byte) []string {
	if len(raw) == 0 {
		return nil
	}
	var out []string
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil
	}
	return out
}

func cloneAnyMapJSON(raw []byte) map[string]any {
	if len(raw) == 0 {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil
	}
	return out
}

func normalizeAnyMapJSON(in map[string]any) map[string]any {
	if in == nil {
		return map[string]any{}
	}
	return in
}

func normalizeStringSliceJSON(in []string) []string {
	if in == nil {
		return []string{}
	}
	return in
}

func findUpstreamPoolMemberSyncTarget(targets []service.UpstreamPoolMember, accountID int64) service.UpstreamPoolMember {
	for _, target := range targets {
		if target.AccountID == accountID {
			return target
		}
	}
	return service.UpstreamPoolMember{}
}

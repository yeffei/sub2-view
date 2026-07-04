package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
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

func (r *upstreamPoolRepository) ListUpstreamPoolMembers(ctx context.Context, poolID int64) ([]service.UpstreamPoolMember, error) {
	if r == nil || r.db == nil || poolID <= 0 {
		return []service.UpstreamPoolMember{}, nil
	}

	const query = `
SELECT
  m.id, m.pool_id, m.account_id, COALESCE(a.name, '') AS account_name, COALESCE(a.platform, '') AS account_platform,
  m.enabled, m.schedulable_override, m.manual_drained,
  m.weight, m.priority_override, m.max_concurrency_override, m.notes, m.joined_at, m.updated_at
FROM upstream_pool_members m
LEFT JOIN accounts a ON a.id = m.account_id
WHERE m.pool_id = $1
ORDER BY m.id ASC`

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
		if err := rows.Scan(
			&member.ID,
			&member.PoolID,
			&member.AccountID,
			&member.AccountName,
			&member.AccountPlatform,
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
  m.id, m.pool_id, m.account_id, COALESCE(a.name, '') AS account_name, COALESCE(a.platform, '') AS account_platform,
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
  AND p.enabled = TRUE
ORDER BY b.priority ASC, b.id ASC
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

	const membersQuery = `
SELECT
  m.account_id,
  m.schedulable_override,
  m.weight,
  m.priority_override,
  m.max_concurrency_override
FROM upstream_pool_members m
WHERE m.pool_id = $1
  AND m.enabled = TRUE
  AND m.manual_drained = FALSE`

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

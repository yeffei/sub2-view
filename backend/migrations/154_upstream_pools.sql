CREATE TABLE IF NOT EXISTS upstream_pools (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    platform VARCHAR(50) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    scheduler_mode VARCHAR(20) NOT NULL DEFAULT 'advanced',
    default_required_capability VARCHAR(64) NOT NULL DEFAULT '',
    default_required_transport VARCHAR(64) NOT NULL DEFAULT '',
    sticky_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sticky_ttl_seconds INTEGER NOT NULL DEFAULT 1800,
    sticky_escape_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sticky_escape_error_rate_threshold DECIMAL(10,4) NOT NULL DEFAULT 0.3000,
    sticky_escape_ttft_ms_threshold INTEGER NOT NULL DEFAULT 6000,
    load_balance_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    failover_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    top_k INTEGER NOT NULL DEFAULT 2,
    max_failover_hops INTEGER NOT NULL DEFAULT 3,
    wait_timeout_ms INTEGER NOT NULL DEFAULT 30000,
    max_waiting INTEGER NOT NULL DEFAULT 100,
    policy_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_upstream_pools_code ON upstream_pools(code);
CREATE INDEX IF NOT EXISTS idx_upstream_pools_platform_enabled ON upstream_pools(platform, enabled);

CREATE TABLE IF NOT EXISTS upstream_pool_members (
    id BIGSERIAL PRIMARY KEY,
    pool_id BIGINT NOT NULL REFERENCES upstream_pools(id) ON DELETE CASCADE,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    schedulable_override BOOLEAN NULL,
    manual_drained BOOLEAN NOT NULL DEFAULT FALSE,
    weight INTEGER NOT NULL DEFAULT 100,
    priority_override INTEGER NULL,
    max_concurrency_override INTEGER NULL,
    notes TEXT NOT NULL DEFAULT '',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_upstream_pool_members_pool_account ON upstream_pool_members(pool_id, account_id);
CREATE INDEX IF NOT EXISTS idx_upstream_pool_members_pool_enabled ON upstream_pool_members(pool_id, enabled, manual_drained);
CREATE INDEX IF NOT EXISTS idx_upstream_pool_members_account_id ON upstream_pool_members(account_id);

CREATE TABLE IF NOT EXISTS upstream_pool_bindings (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    pool_id BIGINT NOT NULL REFERENCES upstream_pools(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    models JSONB NOT NULL DEFAULT '[]'::jsonb,
    request_path_scope JSONB NOT NULL DEFAULT '[]'::jsonb,
    priority INTEGER NOT NULL DEFAULT 100,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_upstream_pool_bindings_group_platform_enabled
    ON upstream_pool_bindings(group_id, platform, enabled, priority);
CREATE UNIQUE INDEX IF NOT EXISTS uq_upstream_pool_bindings_group_pool_platform
    ON upstream_pool_bindings(group_id, pool_id, platform);

CREATE TABLE IF NOT EXISTS upstream_account_sets (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(100) NOT NULL,
    platform VARCHAR(50) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_upstream_account_sets_code ON upstream_account_sets(code);
CREATE INDEX IF NOT EXISTS idx_upstream_account_sets_platform_enabled
    ON upstream_account_sets(platform, enabled);

CREATE TABLE IF NOT EXISTS upstream_account_set_members (
    set_id BIGINT NOT NULL REFERENCES upstream_account_sets(id) ON DELETE CASCADE,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (set_id, account_id)
);

CREATE INDEX IF NOT EXISTS idx_upstream_account_set_members_account_id
    ON upstream_account_set_members(account_id);

CREATE TABLE IF NOT EXISTS upstream_pool_member_sets (
    id BIGSERIAL PRIMARY KEY,
    pool_id BIGINT NOT NULL REFERENCES upstream_pools(id) ON DELETE CASCADE,
    set_id BIGINT NOT NULL REFERENCES upstream_account_sets(id) ON DELETE CASCADE,
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT NOT NULL DEFAULT '',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_upstream_pool_member_sets_pool_set
    ON upstream_pool_member_sets(pool_id, set_id);
CREATE INDEX IF NOT EXISTS idx_upstream_pool_member_sets_pool_enabled
    ON upstream_pool_member_sets(pool_id, enabled);
CREATE INDEX IF NOT EXISTS idx_upstream_pool_member_sets_set_id
    ON upstream_pool_member_sets(set_id);

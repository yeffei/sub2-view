-- Ephemeral health-based weight factors. Manual member weights remain the
-- baseline; stale runtime observations are ignored by scheduler queries.

CREATE TABLE IF NOT EXISTS upstream_pool_runtime_weights (
    pool_id           BIGINT       NOT NULL REFERENCES upstream_pools(id) ON DELETE CASCADE,
    account_id        BIGINT       NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    factor            NUMERIC(4,2) NOT NULL DEFAULT 1.00,
    target_factor     NUMERIC(4,2) NOT NULL DEFAULT 1.00,
    healthy_streak    INTEGER      NOT NULL DEFAULT 0,
    unhealthy_streak  INTEGER      NOT NULL DEFAULT 0,
    reason            TEXT         NOT NULL DEFAULT '',
    last_observed_at  TIMESTAMPTZ  NOT NULL,
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pool_id, account_id),
    CONSTRAINT upstream_pool_runtime_weights_factor_check
        CHECK (factor >= 0.25 AND factor <= 1.25),
    CONSTRAINT upstream_pool_runtime_weights_target_factor_check
        CHECK (target_factor >= 0.25 AND target_factor <= 1.25),
    CONSTRAINT upstream_pool_runtime_weights_streak_check
        CHECK (healthy_streak >= 0 AND unhealthy_streak >= 0)
);

CREATE INDEX IF NOT EXISTS idx_upstream_pool_runtime_weights_observed
    ON upstream_pool_runtime_weights (last_observed_at);

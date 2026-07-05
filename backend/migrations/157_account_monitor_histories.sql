-- Migration: 157_account_monitor_histories
-- Account-level active probe history. User-facing pool health aggregates these
-- rows by upstream pool membership instead of reading manually maintained
-- channel_monitors.

CREATE TABLE IF NOT EXISTS account_monitor_histories (
    id              BIGSERIAL PRIMARY KEY,
    account_id      BIGINT      NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    pool_id         BIGINT      NULL REFERENCES upstream_pools(id) ON DELETE SET NULL,
    group_id        BIGINT      NULL REFERENCES groups(id) ON DELETE SET NULL,
    provider        TEXT        NOT NULL,
    model           TEXT        NOT NULL,
    status          TEXT        NOT NULL,
    latency_ms      INTEGER     NULL,
    ping_latency_ms INTEGER     NULL,
    message         TEXT        NOT NULL DEFAULT '',
    checked_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT account_monitor_histories_provider_check
        CHECK (provider IN ('openai', 'anthropic', 'gemini')),
    CONSTRAINT account_monitor_histories_status_check
        CHECK (status IN ('operational', 'degraded', 'failed', 'error'))
);

CREATE INDEX IF NOT EXISTS idx_account_monitor_histories_account_model_checked
    ON account_monitor_histories (account_id, model, checked_at DESC);

CREATE INDEX IF NOT EXISTS idx_account_monitor_histories_checked_at
    ON account_monitor_histories (checked_at);

CREATE INDEX IF NOT EXISTS idx_account_monitor_histories_pool_checked
    ON account_monitor_histories (pool_id, checked_at DESC)
    WHERE pool_id IS NOT NULL;


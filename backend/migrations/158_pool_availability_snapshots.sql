-- Pool-level runtime snapshots. A pool is available whenever at least one
-- current member is schedulable; individual account probe failures do not
-- reduce pool availability while failover capacity remains.

CREATE TABLE IF NOT EXISTS pool_availability_snapshots (
    id                BIGSERIAL PRIMARY KEY,
    pool_id           BIGINT      NOT NULL REFERENCES upstream_pools(id) ON DELETE CASCADE,
    status            TEXT        NOT NULL,
    total_members     INTEGER     NOT NULL DEFAULT 0,
    available_members INTEGER     NOT NULL DEFAULT 0,
    checked_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pool_availability_snapshots_status_check
        CHECK (status IN ('operational', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_pool_availability_snapshots_pool_checked
    ON pool_availability_snapshots (pool_id, checked_at DESC);

CREATE INDEX IF NOT EXISTS idx_pool_availability_snapshots_checked_at
    ON pool_availability_snapshots (checked_at);

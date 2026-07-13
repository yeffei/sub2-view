CREATE TABLE IF NOT EXISTS upstream_capacity_snapshots (
  id BIGSERIAL PRIMARY KEY,
  set_id BIGINT NOT NULL REFERENCES upstream_account_sets(id) ON DELETE CASCADE,
  capacity_limit INTEGER NOT NULL,
  current_concurrency INTEGER NOT NULL,
  waiting_count INTEGER NOT NULL,
  load_rate INTEGER NOT NULL,
  checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_upstream_capacity_snapshots_set_checked
  ON upstream_capacity_snapshots (set_id, checked_at DESC);

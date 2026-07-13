-- Shared upstream concurrency capacity attached to account sets.
-- Normal account-set membership remains independent; capacity membership is
-- opt-in and an account may belong to at most one shared capacity group.

ALTER TABLE upstream_account_sets
    ADD COLUMN IF NOT EXISTS shared_concurrency_limit INTEGER NULL;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'upstream_account_sets_shared_concurrency_limit_check'
    ) THEN
        ALTER TABLE upstream_account_sets
            ADD CONSTRAINT upstream_account_sets_shared_concurrency_limit_check
            CHECK (shared_concurrency_limit IS NULL OR shared_concurrency_limit > 0);
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS upstream_account_set_capacity_members (
    account_id BIGINT PRIMARY KEY,
    set_id BIGINT NOT NULL,
    hard_concurrency_limit INTEGER NULL,
    soft_concurrency_share INTEGER NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT upstream_account_set_capacity_members_membership_fk
        FOREIGN KEY (set_id, account_id)
        REFERENCES upstream_account_set_members(set_id, account_id)
        ON DELETE CASCADE,
    CONSTRAINT upstream_account_set_capacity_members_hard_limit_check
        CHECK (hard_concurrency_limit IS NULL OR hard_concurrency_limit > 0),
    CONSTRAINT upstream_account_set_capacity_members_soft_share_check
        CHECK (soft_concurrency_share IS NULL OR soft_concurrency_share > 0)
);

CREATE INDEX IF NOT EXISTS idx_upstream_account_set_capacity_members_set_id
    ON upstream_account_set_capacity_members(set_id);

CREATE OR REPLACE FUNCTION validate_upstream_account_set_capacity_member()
RETURNS TRIGGER AS $$
DECLARE
    group_limit INTEGER;
    set_platform TEXT;
    account_platform TEXT;
BEGIN
    SELECT shared_concurrency_limit, platform
      INTO group_limit, set_platform
      FROM upstream_account_sets
     WHERE id = NEW.set_id;

    IF group_limit IS NULL THEN
        RAISE EXCEPTION 'shared_concurrency_limit is required for capacity members';
    END IF;

    IF NEW.hard_concurrency_limit IS NOT NULL
       AND NEW.hard_concurrency_limit > group_limit THEN
        RAISE EXCEPTION 'hard_concurrency_limit must not exceed shared_concurrency_limit';
    END IF;

    SELECT platform
      INTO account_platform
      FROM accounts
     WHERE id = NEW.account_id;

    IF LOWER(COALESCE(account_platform, '')) <> LOWER(COALESCE(set_platform, '')) THEN
        RAISE EXCEPTION 'capacity member platform must match account set platform';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_validate_upstream_account_set_capacity_member
    ON upstream_account_set_capacity_members;
CREATE TRIGGER trg_validate_upstream_account_set_capacity_member
BEFORE INSERT OR UPDATE
ON upstream_account_set_capacity_members
FOR EACH ROW
EXECUTE FUNCTION validate_upstream_account_set_capacity_member();

CREATE OR REPLACE FUNCTION validate_upstream_account_set_shared_limit()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.shared_concurrency_limit IS NULL
       AND EXISTS (
           SELECT 1
             FROM upstream_account_set_capacity_members
            WHERE set_id = NEW.id
       ) THEN
        RAISE EXCEPTION 'shared_concurrency_limit cannot be cleared while capacity members exist';
    END IF;

    IF NEW.shared_concurrency_limit IS NOT NULL
       AND EXISTS (
           SELECT 1
             FROM upstream_account_set_capacity_members
            WHERE set_id = NEW.id
              AND hard_concurrency_limit > NEW.shared_concurrency_limit
       ) THEN
        RAISE EXCEPTION 'shared_concurrency_limit must not be lower than a member hard_concurrency_limit';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_validate_upstream_account_set_shared_limit
    ON upstream_account_sets;
CREATE TRIGGER trg_validate_upstream_account_set_shared_limit
BEFORE UPDATE OF shared_concurrency_limit
ON upstream_account_sets
FOR EACH ROW
EXECUTE FUNCTION validate_upstream_account_set_shared_limit();

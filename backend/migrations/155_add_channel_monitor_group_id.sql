ALTER TABLE channel_monitors
  ADD COLUMN IF NOT EXISTS group_id BIGINT NULL REFERENCES groups(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS channel_monitors_group_id_idx
  ON channel_monitors(group_id);

UPDATE channel_monitors AS cm
SET group_id = g.id
FROM groups AS g
WHERE cm.group_id IS NULL
  AND COALESCE(BTRIM(cm.group_name), '') <> ''
  AND LOWER(g.name) = LOWER(BTRIM(cm.group_name))
  AND LOWER(g.platform) = LOWER(cm.provider);

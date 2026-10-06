ALTER TABLE workflows
DROP COLUMN IF EXISTS last_run_at,
DROP COLUMN IF EXISTS next_run_at,
DROP COLUMN IF EXISTS schedule_value,
DROP COLUMN IF EXISTS schedule_type,
DROP COLUMN IF EXISTS schedule_enabled;
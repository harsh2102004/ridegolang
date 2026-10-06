ALTER TABLE workflows
ADD COLUMN schedule_enabled BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN schedule_type TEXT,
ADD COLUMN schedule_value TEXT,
ADD COLUMN next_run_at TIMESTAMP,
ADD COLUMN last_run_at TIMESTAMP;
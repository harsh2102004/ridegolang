ALTER TABLE workflows
ADD COLUMN is_schedule_running BOOLEAN NOT NULL DEFAULT false;
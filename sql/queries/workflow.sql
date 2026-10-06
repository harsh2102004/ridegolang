-- name: CreateWorkflow :one
INSERT INTO workflows (
    id, user_id, name, description, trigger_type, is_active
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetWorkflowByID :one
SELECT * FROM workflows
WHERE id = $1 AND user_id = $2;

-- name: ListWorkflowByUser :many
SELECT * FROM workflows
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteWorkflow :exec
DELETE FROM workflows
WHERE id = $1 AND user_id = $2;

-- name: CreateWorkflowStep :one
INSERT INTO workflow_steps (
    id,
    workflow_id,
    frontend_node_id,
    step_order,
    step_type,
    config,
    position_x,
    position_y
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListWorkflowSteps :many
SELECT * FROM workflow_steps
WHERE workflow_id = $1
ORDER BY step_order ASC;

-- name: CreateWorkflowRun :one
INSERT INTO workflow_runs (
    id, workflow_id, user_id, status, started_at
) VALUES (
    $1, $2, $3, $4, NOW()
)
RETURNING *;

-- name: UpdateWorkflowRunStatus :exec
UPDATE workflow_runs
SET status = $2,
    error_message = $3,
    finished_at = NOW()
WHERE id = $1;

-- name: ListWorkflowRuns :many
SELECT * FROM workflow_runs
WHERE workflow_id = $1 AND user_id = $2
ORDER BY created_at DESC;

-- name: CreateWorkflowStepRun :one
INSERT INTO workflow_step_runs (
    id,
    workflow_run_id,
    workflow_step_id,
    status,
    input,
    output,
    error_message,
    started_at,
    finished_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
)
RETURNING *;


-- name: ListWorkflowStepRuns :many
SELECT *
FROM workflow_step_runs
WHERE workflow_run_id = $1
ORDER BY created_at ASC;


-- name: UpdateWorkflowStepRunStatus :exec
UPDATE workflow_step_runs
SET status = $2,
    output = $3,
    error_message = $4,
    finished_at = NOW()
WHERE id = $1;


-- name: DeleteWorkflowSteps :exec
DELETE FROM workflow_steps
WHERE workflow_id = $1;


-- name: CreateWorkflowEdge :one
INSERT INTO workflow_edges (
    id,
    workflow_id,
    source_step_id,
    target_step_id,
    condition_branch
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;


-- name: ListWorkflowEdges :many
SELECT
    we.id,
    we.workflow_id,
    we.source_step_id,
    we.target_step_id,
    we.condition_branch,
    we.created_at,

    source.frontend_node_id AS source_frontend_node_id,
    target.frontend_node_id AS target_frontend_node_id
FROM workflow_edges we
JOIN workflow_steps source ON source.id = we.source_step_id
JOIN workflow_steps target ON target.id = we.target_step_id
WHERE we.workflow_id = $1
ORDER BY we.created_at ASC;

-- name: DeleteWorkflowEdges :exec
DELETE FROM workflow_edges
WHERE workflow_id = $1;


-- name: ListWorkflowEdgesForExecution :many
SELECT *
FROM workflow_edges
WHERE workflow_id = $1
ORDER BY created_at ASC;


-- name: CreateWebhookTrigger :one
INSERT INTO webhook_triggers (
    id,
    workflow_id,
    user_id,
    webhook_url_id,
    frontend_node_id
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetWebhookTriggerByURLID :one
SELECT * FROM webhook_triggers
WHERE webhook_url_id = $1;

-- name: DeleteWebhookTriggersByWorkflow :exec
DELETE FROM webhook_triggers
WHERE workflow_id = $1;

-- name: ListWebhookTriggersByWorkflow :many
SELECT * FROM webhook_triggers
WHERE workflow_id = $1;



-- name: UpdateWorkflowSchedule :exec
UPDATE workflows
SET
    schedule_enabled = $2,
    schedule_type = $3,
    schedule_value = $4,
    next_run_at = $5,
    updated_at = NOW()
WHERE id = $1 AND user_id = $6;

-- name: ListDueScheduledWorkflows :many
SELECT *
FROM workflows
WHERE schedule_enabled = true
  AND next_run_at IS NOT NULL
  AND next_run_at <= NOW()
  AND is_active = true
  AND is_schedule_running = false
  AND trigger_type = 'schedule';

-- name: MarkWorkflowScheduleRun :exec
UPDATE workflows
SET
    last_run_at = NOW(),
    next_run_at = $2,
    is_schedule_running = false,
    updated_at = NOW()
WHERE id = $1;

-- name: MarkScheduleRunning :exec
UPDATE workflows
SET
    is_schedule_running = $2,
    updated_at = NOW()
WHERE id = $1;


-- name: GetWorkflowSchedule :one
SELECT
    schedule_enabled,
    schedule_type,
    schedule_value,
    next_run_at,
    last_run_at
FROM workflows
WHERE id = $1 AND user_id = $2;


-- name: UpdateWorkflow :one
UPDATE workflows
SET
    name = $3,
    description = $4,
    updated_at = NOW()
WHERE id = $1 AND user_id = $2
RETURNING *;
-- name: GetMonthlyUsage :one
SELECT *
FROM workflow_usage
WHERE user_id = $1 AND month = $2;

-- name: IncrementWorkflowUsage :one
INSERT INTO workflow_usage (
    id,
    user_id,
    month,
    workflow_runs
)
VALUES (
    $1,
    $2,
    $3,
    1
)
ON CONFLICT (user_id, month)
DO UPDATE SET
    workflow_runs = workflow_usage.workflow_runs + 1,
    updated_at = NOW()
RETURNING *;
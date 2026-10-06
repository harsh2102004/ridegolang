package repository

import (
	"context"
	"org/api-core/internal/db"

	"github.com/google/uuid"
)

type SQLCWorkflowRepository struct {
	q *db.Queries
}

func NewSQLCWorkflowRepository(q *db.Queries) *SQLCWorkflowRepository {
	return &SQLCWorkflowRepository{q: q}
}

func (r *SQLCWorkflowRepository) CreateWorkflow(ctx context.Context, arg db.CreateWorkflowParams) (db.Workflow, error) {
	return r.q.CreateWorkflow(ctx, arg)
}

func (r *SQLCWorkflowRepository) ListWorkflowByUser(ctx context.Context, userID uuid.UUID) ([]db.Workflow, error) {
	return r.q.ListWorkflowByUser(ctx, userID)
}

func (r *SQLCWorkflowRepository) GetWorkflowByID(ctx context.Context, arg db.GetWorkflowByIDParams) (db.Workflow, error) {
	return r.q.GetWorkflowByID(ctx, arg)
}

func (r *SQLCWorkflowRepository) UpdateWorkflow(
	ctx context.Context,
	arg db.UpdateWorkflowParams,
) (db.Workflow, error) {
	return r.q.UpdateWorkflow(ctx, arg)
}

func (r *SQLCWorkflowRepository) DeleteWorkflow(ctx context.Context, arg db.DeleteWorkflowParams) error {
	return r.q.DeleteWorkflow(ctx, arg)
}

func (r *SQLCWorkflowRepository) CreateWorkflowStep(ctx context.Context, arg db.CreateWorkflowStepParams) (db.WorkflowStep, error) {
	return r.q.CreateWorkflowStep(ctx, arg)
}

func (r *SQLCWorkflowRepository) ListWorkflowSteps(ctx context.Context, workflowID uuid.UUID) ([]db.WorkflowStep, error) {
	return r.q.ListWorkflowSteps(ctx, workflowID)
}

func (r *SQLCWorkflowRepository) CreateWorkflowRun(ctx context.Context, arg db.CreateWorkflowRunParams) (db.WorkflowRun, error) {
	return r.q.CreateWorkflowRun(ctx, arg)
}

func (r *SQLCWorkflowRepository) UpdateWorkflowRunStatus(ctx context.Context, arg db.UpdateWorkflowRunStatusParams) error {
	return r.q.UpdateWorkflowRunStatus(ctx, arg)
}

func (r *SQLCWorkflowRepository) ListWorkflowRuns(ctx context.Context, arg db.ListWorkflowRunsParams) ([]db.WorkflowRun, error) {
	return r.q.ListWorkflowRuns(ctx, arg)
}

func (r *SQLCWorkflowRepository) CreateWorkflowStepRun(ctx context.Context, arg db.CreateWorkflowStepRunParams) (db.WorkflowStepRun, error) {
	return r.q.CreateWorkflowStepRun(ctx, arg)
}

func (r *SQLCWorkflowRepository) ListWorkflowStepRuns(ctx context.Context, workflowRunID uuid.UUID) ([]db.WorkflowStepRun, error) {
	return r.q.ListWorkflowStepRuns(ctx, workflowRunID)
}

func (r *SQLCWorkflowRepository) UpdateWorkflowStepRunStatus(ctx context.Context, arg db.UpdateWorkflowStepRunStatusParams) error {
	return r.q.UpdateWorkflowStepRunStatus(ctx, arg)
}

func (r *SQLCWorkflowRepository) DeleteWorkflowSteps(ctx context.Context, workflowID uuid.UUID) error {
	return r.q.DeleteWorkflowSteps(ctx, workflowID)
}

func (r *SQLCWorkflowRepository) CreateWorkflowEdge(ctx context.Context, arg db.CreateWorkflowEdgeParams) (db.WorkflowEdge, error) {
	return r.q.CreateWorkflowEdge(ctx, arg)
}

func (r *SQLCWorkflowRepository) ListWorkflowEdges(ctx context.Context, workflowID uuid.UUID) ([]db.ListWorkflowEdgesRow, error) {
	return r.q.ListWorkflowEdges(ctx, workflowID)
}

func (r *SQLCWorkflowRepository) DeleteWorkflowEdges(ctx context.Context, workflowID uuid.UUID) error {
	return r.q.DeleteWorkflowEdges(ctx, workflowID)
}

func (r *SQLCWorkflowRepository) ListWorkflowEdgesForExecution(ctx context.Context, workflowID uuid.UUID) ([]db.WorkflowEdge, error) {
	return r.q.ListWorkflowEdgesForExecution(ctx, workflowID)
}

func (r *SQLCWorkflowRepository) CreateWebhookTrigger(ctx context.Context, arg db.CreateWebhookTriggerParams) (db.WebhookTrigger, error) {
	return r.q.CreateWebhookTrigger(ctx, arg)
}

func (r *SQLCWorkflowRepository) GetWebhookTriggerByURLID(ctx context.Context, webhookURLID string) (db.WebhookTrigger, error) {
	return r.q.GetWebhookTriggerByURLID(ctx, webhookURLID)
}

func (r *SQLCWorkflowRepository) DeleteWebhookTriggersByWorkflow(ctx context.Context, workflowID uuid.UUID) error {
	return r.q.DeleteWebhookTriggersByWorkflow(ctx, workflowID)
}

func (r *SQLCWorkflowRepository) ListWebhookTriggersByWorkflow(ctx context.Context, workflowID uuid.UUID) ([]db.WebhookTrigger, error) {
	return r.q.ListWebhookTriggersByWorkflow(ctx, workflowID)
}

func (r *SQLCWorkflowRepository) UpdateWorkflowSchedule(ctx context.Context, params db.UpdateWorkflowScheduleParams) error {
	return r.q.UpdateWorkflowSchedule(ctx, params)
}

func (r *SQLCWorkflowRepository) ListDueScheduledWorkflows(ctx context.Context) ([]db.Workflow, error) {
	return r.q.ListDueScheduledWorkflows(ctx)
}

func (r *SQLCWorkflowRepository) MarkWorkflowScheduleRun(ctx context.Context, params db.MarkWorkflowScheduleRunParams) error {
	return r.q.MarkWorkflowScheduleRun(ctx, params)
}

func (r *SQLCWorkflowRepository) MarkScheduleRunning(ctx context.Context, params db.MarkScheduleRunningParams) error {
	return r.q.MarkScheduleRunning(ctx, params)
}

func (r *SQLCWorkflowRepository) GetWorkflowSchedule(ctx context.Context, params db.GetWorkflowScheduleParams) (db.GetWorkflowScheduleRow, error) {
	return r.q.GetWorkflowSchedule(ctx, params)
}

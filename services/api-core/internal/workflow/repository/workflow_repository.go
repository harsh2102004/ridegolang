package repository

import (
	"context"

	"github.com/google/uuid"
	"org/api-core/internal/db"
)

type WorkflowRepository interface {
	// Workflow
	CreateWorkflow(ctx context.Context, arg db.CreateWorkflowParams) (db.Workflow, error)
	ListWorkflowByUser(ctx context.Context, userID uuid.UUID) ([]db.Workflow, error)
    UpdateWorkflow(ctx context.Context, arg db.UpdateWorkflowParams) (db.Workflow, error)
	GetWorkflowByID(ctx context.Context, arg db.GetWorkflowByIDParams) (db.Workflow, error)
	DeleteWorkflow(ctx context.Context, arg db.DeleteWorkflowParams) error

	// Steps
	CreateWorkflowStep(ctx context.Context, arg db.CreateWorkflowStepParams) (db.WorkflowStep, error)
	ListWorkflowSteps(ctx context.Context, workflowID uuid.UUID) ([]db.WorkflowStep, error)
	DeleteWorkflowSteps(ctx context.Context, workflowID uuid.UUID) error // ✅ ADDED

	// Runs
	CreateWorkflowRun(ctx context.Context, arg db.CreateWorkflowRunParams) (db.WorkflowRun, error)
	UpdateWorkflowRunStatus(ctx context.Context, arg db.UpdateWorkflowRunStatusParams) error
	ListWorkflowRuns(ctx context.Context, arg db.ListWorkflowRunsParams) ([]db.WorkflowRun, error)

	CreateWorkflowStepRun(ctx context.Context, arg db.CreateWorkflowStepRunParams) (db.WorkflowStepRun, error)
	ListWorkflowStepRuns(ctx context.Context, workflowRunID uuid.UUID) ([]db.WorkflowStepRun, error)
	UpdateWorkflowStepRunStatus(ctx context.Context, arg db.UpdateWorkflowStepRunStatusParams) error

	//Edges
	CreateWorkflowEdge(ctx context.Context, arg db.CreateWorkflowEdgeParams) (db.WorkflowEdge, error)
	ListWorkflowEdges(ctx context.Context, workflowID uuid.UUID) ([]db.ListWorkflowEdgesRow, error)
	DeleteWorkflowEdges(ctx context.Context, workflowID uuid.UUID) error

	ListWorkflowEdgesForExecution(ctx context.Context, workflowID uuid.UUID) ([]db.WorkflowEdge, error)

	CreateWebhookTrigger(ctx context.Context, arg db.CreateWebhookTriggerParams) (db.WebhookTrigger, error)
	GetWebhookTriggerByURLID(ctx context.Context, webhookURLID string) (db.WebhookTrigger, error)
	DeleteWebhookTriggersByWorkflow(ctx context.Context, workflowID uuid.UUID) error
	ListWebhookTriggersByWorkflow(ctx context.Context, workflowID uuid.UUID) ([]db.WebhookTrigger, error)

	// Schedule
	UpdateWorkflowSchedule(ctx context.Context, params db.UpdateWorkflowScheduleParams) error
	ListDueScheduledWorkflows(ctx context.Context) ([]db.Workflow, error)
	MarkWorkflowScheduleRun(ctx context.Context, params db.MarkWorkflowScheduleRunParams) error
    
    
    MarkScheduleRunning(ctx context.Context, params db.MarkScheduleRunningParams) error
    GetWorkflowSchedule(ctx context.Context, params db.GetWorkflowScheduleParams) (db.GetWorkflowScheduleRow, error)
}

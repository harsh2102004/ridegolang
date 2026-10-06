package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"

	"org/api-core/internal/db"
	"org/api-core/internal/queue"
	"org/api-core/internal/workflow/executor"
	"org/api-core/internal/workflow/repository"

	"github.com/robfig/cron/v3"
)

type WorkflowService struct {
	repo          repository.WorkflowRepository
	queue         *queue.RedisQueue
	logger        *slog.Logger
	geminiService executor.GeminiService
    
	subscribers   map[string][]chan any
	mu            sync.Mutex
}

func NewWorkflowService(
	repo repository.WorkflowRepository,
	workflowQueue *queue.RedisQueue,
    geminiService executor.GeminiService,
	logger *slog.Logger,
) *WorkflowService {
	return &WorkflowService{
		repo:        repo,
		queue:       workflowQueue,
		logger:      logger,
		subscribers: make(map[string][]chan any),
        geminiService: geminiService,
	}
}

type SaveWorkflowStepRequest struct {
	FrontendNodeID string          `json:"frontend_node_id"`
	StepOrder      int             `json:"step_order"`
	StepType       string          `json:"step_type"`
	Config         json.RawMessage `json:"config"`
	PositionX      float64         `json:"position_x"`
	PositionY      float64         `json:"position_y"`
}

type SaveWorkflowEdgeRequest struct {
	Source          string `json:"source"`
	Target          string `json:"target"`
	ConditionBranch string `json:"condition_branch"`
}

type SaveWorkflowStepsRequest struct {
	Steps []SaveWorkflowStepRequest `json:"steps"`
	Edges []SaveWorkflowEdgeRequest `json:"edges"`
}

type UpdateWorkflowScheduleRequest struct {
	Enabled       bool   `json:"enabled"`
	ScheduleType  string `json:"schedule_type"`
	ScheduleValue string `json:"schedule_value"`
}

//
// 🔹 WORKFLOW METHODS
//

func (s *WorkflowService) CreateWorkflow(ctx context.Context, userID uuid.UUID, name string, description string) (db.Workflow, error) {
	s.logger.Info("creating workflow",
		"user_id", userID.String(),
		"name", name,
	)

	workflow, err := s.repo.CreateWorkflow(ctx, db.CreateWorkflowParams{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        name,
		Description: description,
		TriggerType: "manual",
		IsActive:    true,
	})
	if err != nil {
		s.logger.Error("failed to create workflow",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return workflow, err
	}

	s.logger.Info("workflow created",
		"user_id", userID.String(),
		"workflow_id", workflow.ID.String(),
	)
	return workflow, nil
}

func (s *WorkflowService) ListWorkflow(ctx context.Context, userID uuid.UUID) ([]db.Workflow, error) {
	s.logger.Info("listing workflows",
		"user_id", userID.String(),
	)

	workflows, err := s.repo.ListWorkflowByUser(ctx, userID)
	if err != nil {
		s.logger.Error("failed to list workflows",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	s.logger.Info("workflows listed",
		"user_id", userID.String(),
		"count", len(workflows),
	)
	return workflows, nil
}

func (s *WorkflowService) DeleteWorkflow(ctx context.Context, userID uuid.UUID, workflowID uuid.UUID) error {
	s.logger.Info("deleting workflow",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
	)

	err := s.repo.DeleteWorkflow(ctx, db.DeleteWorkflowParams{
		ID:     workflowID,
		UserID: userID,
	})
	if err != nil {
		s.logger.Error("failed to delete workflow",
			"user_id", userID.String(),
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return err
	}

	s.logger.Info("workflow deleted",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
	)
	return nil
}

//
// 🔹 STEP METHODS
//

func (s *WorkflowService) CreateStep(ctx context.Context, workflowID uuid.UUID, stepOrder int32, stepType string, config []byte) (db.WorkflowStep, error) {
	s.logger.Info("creating step",
		"workflow_id", workflowID.String(),
		"step_type", stepType,
	)

	step, err := s.repo.CreateWorkflowStep(ctx, db.CreateWorkflowStepParams{
		ID:         uuid.New(),
		WorkflowID: workflowID,
		StepOrder:  stepOrder,
		StepType:   stepType,
		Config:     config,
	})
	if err != nil {
		s.logger.Error("failed to create step",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return step, err
	}

	s.logger.Info("step created",
		"workflow_id", workflowID.String(),
		"step_id", step.ID.String(),
	)
	return step, nil
}

func (s *WorkflowService) ListSteps(ctx context.Context, workflowID uuid.UUID) ([]db.WorkflowStep, error) {
	s.logger.Info("listing steps",
		"workflow_id", workflowID.String(),
	)

	steps, err := s.repo.ListWorkflowSteps(ctx, workflowID)
	if err != nil {
		s.logger.Error("failed to list steps",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	s.logger.Info("steps listed",
		"workflow_id", workflowID.String(),
		"count", len(steps),
	)
	return steps, nil
}

func (s *WorkflowService) ListWorkflowRuns(ctx context.Context, workflowID, userID uuid.UUID) ([]db.WorkflowRun, error) {
	s.logger.Info("listing workflow runs",
		"workflow_id", workflowID.String(),
		"user_id", userID.String(),
	)

	runs, err := s.repo.ListWorkflowRuns(ctx, db.ListWorkflowRunsParams{
		WorkflowID: workflowID,
		UserID:     userID,
	})
	if err != nil {
		s.logger.Error("failed to list workflow runs",
			"workflow_id", workflowID.String(),
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	s.logger.Info("workflow runs listed",
		"workflow_id", workflowID.String(),
		"count", len(runs),
	)
	return runs, nil
}

func (s *WorkflowService) ListWorkflowStepRuns(ctx context.Context, workflowRunID uuid.UUID) ([]db.WorkflowStepRun, error) {
	s.logger.Info("listing workflow step runs",
		"workflow_run_id", workflowRunID.String(),
	)

	steps, err := s.repo.ListWorkflowStepRuns(ctx, workflowRunID)
	if err != nil {
		s.logger.Error("failed to list workflow step runs",
			"workflow_run_id", workflowRunID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	s.logger.Info("workflow step runs listed",
		"workflow_run_id", workflowRunID.String(),
		"count", len(steps),
	)
	return steps, nil
}

func (s *WorkflowService) SaveWorkflowSteps(
	ctx context.Context,
	workflowID uuid.UUID,
	userID uuid.UUID,
	steps []SaveWorkflowStepRequest,
	edges []SaveWorkflowEdgeRequest,
) error {
	s.logger.Info("saving workflow steps",
		"workflow_id", workflowID.String(),
		"user_id", userID.String(),
		"steps_count", len(steps),
		"edges_count", len(edges),
	)

	_, err := s.repo.GetWorkflowByID(ctx, db.GetWorkflowByIDParams{
		ID:     workflowID,
		UserID: userID,
	})
	if err != nil {
		s.logger.Error("workflow not found for saving steps",
			"workflow_id", workflowID.String(),
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return err
	}

	err = s.repo.DeleteWorkflowEdges(ctx, workflowID)
	if err != nil {
		s.logger.Error("failed to delete workflow edges",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return err
	}

	err = s.repo.DeleteWebhookTriggersByWorkflow(ctx, workflowID)
	if err != nil {
		s.logger.Error("failed to delete webhook triggers",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return err
	}

	err = s.repo.DeleteWorkflowSteps(ctx, workflowID)
	if err != nil {
		s.logger.Error("failed to delete workflow steps",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return err
	}

	nodeIDToStepID := make(map[string]uuid.UUID)

	for _, step := range steps {
		stepID := uuid.New()

		config := step.Config

		if step.StepType == "webhook_trigger" {
			webhookURLID := uuid.NewString()

			var configMap map[string]interface{}

			if len(config) > 0 {
				_ = json.Unmarshal(config, &configMap)
			}

			if configMap == nil {
				configMap = make(map[string]interface{})
			}

			configMap["webhook_url_id"] = webhookURLID
			apiBaseURL := os.Getenv("API_BASE_URL")
			if apiBaseURL == "" {
				apiBaseURL = "http://localhost:8080"
			}
			configMap["webhook_url"] = apiBaseURL + "/webhooks/" + webhookURLID

			updatedConfig, err := json.Marshal(configMap)
			if err != nil {
				s.logger.Error("failed to marshal webhook config",
					"workflow_id", workflowID.String(),
					"error", err.Error(),
				)
				return err
			}

			config = updatedConfig

			_, err = s.repo.CreateWebhookTrigger(ctx, db.CreateWebhookTriggerParams{
				ID:             uuid.New(),
				WorkflowID:     workflowID,
				UserID:         userID,
				WebhookUrlID:   webhookURLID,
				FrontendNodeID: step.FrontendNodeID,
			})
			if err != nil {
				s.logger.Error("failed to create webhook trigger",
					"workflow_id", workflowID.String(),
					"error", err.Error(),
				)
				return err
			}
		}

		_, err := s.repo.CreateWorkflowStep(ctx, db.CreateWorkflowStepParams{
			ID:             stepID,
			WorkflowID:     workflowID,
			FrontendNodeID: step.FrontendNodeID,
			StepOrder:      int32(step.StepOrder),
			StepType:       step.StepType,
			Config:         config,
			PositionX:      step.PositionX,
			PositionY:      step.PositionY,
		})
		if err != nil {
			s.logger.Error("failed to create workflow step",
				"workflow_id", workflowID.String(),
				"step_type", step.StepType,
				"error", err.Error(),
			)
			return err
		}

		nodeIDToStepID[step.FrontendNodeID] = stepID
	}

	for _, edge := range edges {
		sourceStepID, ok := nodeIDToStepID[edge.Source]
		if !ok {
			err := fmt.Errorf("source node not found: %s", edge.Source)
			s.logger.Error("failed to create edge: source node missing",
				"workflow_id", workflowID.String(),
				"source", edge.Source,
			)
			return err
		}

		targetStepID, ok := nodeIDToStepID[edge.Target]
		if !ok {
			err := fmt.Errorf("target node not found: %s", edge.Target)
			s.logger.Error("failed to create edge: target node missing",
				"workflow_id", workflowID.String(),
				"target", edge.Target,
			)
			return err
		}

		_, err := s.repo.CreateWorkflowEdge(ctx, db.CreateWorkflowEdgeParams{
			ID:           uuid.New(),
			WorkflowID:   workflowID,
			SourceStepID: sourceStepID,
			TargetStepID: targetStepID,
			ConditionBranch: sql.NullString{
				String: edge.ConditionBranch,
				Valid:  edge.ConditionBranch != "",
			},
		})
		if err != nil {
			s.logger.Error("failed to create workflow edge",
				"workflow_id", workflowID.String(),
				"error", err.Error(),
			)
			return err
		}
	}

	s.logger.Info("workflow steps saved successfully",
		"workflow_id", workflowID.String(),
	)
	return nil
}

func (s *WorkflowService) GetWorkflowSteps(
	ctx context.Context,
	workflowID uuid.UUID,
	userID uuid.UUID,
) ([]db.WorkflowStep, error) {
	s.logger.Info("getting workflow steps",
		"workflow_id", workflowID.String(),
		"user_id", userID.String(),
	)

	_, err := s.repo.GetWorkflowByID(ctx, db.GetWorkflowByIDParams{
		ID:     workflowID,
		UserID: userID,
	})
	if err != nil {
		s.logger.Error("workflow not found for getting steps",
			"workflow_id", workflowID.String(),
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	steps, err := s.repo.ListWorkflowSteps(ctx, workflowID)
	if err != nil {
		s.logger.Error("failed to list workflow steps",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	s.logger.Info("workflow steps retrieved",
		"workflow_id", workflowID.String(),
		"count", len(steps),
	)
	return steps, nil
}

func (s *WorkflowService) GetWorkflowEdges(
	ctx context.Context,
	workflowID uuid.UUID,
	userID uuid.UUID,
) ([]db.ListWorkflowEdgesRow, error) {
	s.logger.Info("getting workflow edges",
		"workflow_id", workflowID.String(),
		"user_id", userID.String(),
	)

	_, err := s.repo.GetWorkflowByID(ctx, db.GetWorkflowByIDParams{
		ID:     workflowID,
		UserID: userID,
	})
	if err != nil {
		s.logger.Error("workflow not found for getting edges",
			"workflow_id", workflowID.String(),
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	edges, err := s.repo.ListWorkflowEdges(ctx, workflowID)
	if err != nil {
		s.logger.Error("failed to list workflow edges",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return nil, err
	}

	s.logger.Info("workflow edges retrieved",
		"workflow_id", workflowID.String(),
		"count", len(edges),
	)
	return edges, nil
}

func (s *WorkflowService) RunWorkflowFromWebhook(
	ctx context.Context,
	webhookID string,
	payload map[string]interface{},
) (uuid.UUID, error) {
	s.logger.Info("running workflow from webhook",
		"webhook_id", webhookID,
	)

	trigger, err := s.repo.GetWebhookTriggerByURLID(ctx, webhookID)
	if err != nil {
		s.logger.Error("webhook trigger not found",
			"webhook_id", webhookID,
			"error", err.Error(),
		)
		return uuid.Nil, err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.logger.Error("failed to marshal webhook payload",
			"webhook_id", webhookID,
			"error", err.Error(),
		)
		return uuid.Nil, err
	}

	runID, err := s.RunWorkflowWithInput(
		ctx,
		trigger.WorkflowID,
		trigger.UserID,
		payloadBytes,
	)
	if err != nil {
		s.logger.Error("failed to run workflow from webhook",
			"webhook_id", webhookID,
			"workflow_id", trigger.WorkflowID.String(),
			"error", err.Error(),
		)
		return runID, err
	}

	s.logger.Info("workflow triggered from webhook",
		"webhook_id", webhookID,
		"workflow_id", trigger.WorkflowID.String(),
		"run_id", runID.String(),
	)
	return runID, nil
}

func (s *WorkflowService) RunWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
	userID uuid.UUID,
) (uuid.UUID, error) {
	return s.RunWorkflowWithInput(ctx, workflowID, userID, nil)
}

func (s *WorkflowService) RunWorkflowWithInput(
	ctx context.Context,
	workflowID uuid.UUID,
	userID uuid.UUID,
	initialInput []byte,
) (uuid.UUID, error) {
	s.logger.Info("queuing workflow run",
		"workflow_id", workflowID.String(),
		"user_id", userID.String(),
	)

	runID := uuid.New()

	_, err := s.repo.CreateWorkflowRun(ctx, db.CreateWorkflowRunParams{
		ID:         runID,
		WorkflowID: workflowID,
		UserID:     userID,
		Status:     "queued",
	})
	if err != nil {
		s.logger.Error("failed to create workflow run record",
			"workflow_id", workflowID.String(),
			"user_id", userID.String(),
			"error", err.Error(),
		)
		return runID, err
	}

	job := queue.WorkflowJob{
		WorkflowID: workflowID.String(),
		UserID:     userID.String(),
		RunID:      runID.String(),
		Input:      initialInput,
		Source:     "manual",
	}

	payload, err := json.Marshal(job)
	if err != nil {
		s.logger.Error("failed to marshal workflow job",
			"workflow_id", workflowID.String(),
			"run_id", runID.String(),
			"error", err.Error(),
		)
		return runID, err
	}

	err = s.queue.Push(ctx, payload)
	if err != nil {
		s.logger.Error("failed to push job to queue",
			"workflow_id", workflowID.String(),
			"run_id", runID.String(),
			"error", err.Error(),
		)
		return runID, err
	}

	s.logger.Info("workflow queued successfully",
		"workflow_id", workflowID.String(),
		"run_id", runID.String(),
	)
	return runID, nil
}

func (s *WorkflowService) ExecuteWorkflowRun(
	ctx context.Context,
	workflowID uuid.UUID,
	userID uuid.UUID,
	runID uuid.UUID,
	initialInput []byte,
) (uuid.UUID, error) {
	s.logger.Info("executing workflow run",
		"workflow_id", workflowID.String(),
		"user_id", userID.String(),
		"run_id", runID.String(),
	)

	_ = s.repo.UpdateWorkflowRunStatus(ctx, db.UpdateWorkflowRunStatusParams{
		ID:     runID,
		Status: "running",
		ErrorMessage: sql.NullString{
			Valid: false,
		},
	})

	failRun := func(execCtx context.Context, execErr error) error {
		s.logger.Error("workflow run failed",
			"workflow_id", workflowID.String(),
			"run_id", runID.String(),
			"error", execErr.Error(),
		)
		_ = s.repo.UpdateWorkflowRunStatus(execCtx, db.UpdateWorkflowRunStatusParams{
			ID:     runID,
			Status: "failed",
			ErrorMessage: sql.NullString{
				String: execErr.Error(),
				Valid:  true,
			},
		})
		return execErr
	}

	steps, err := s.repo.ListWorkflowSteps(ctx, workflowID)
	if err != nil {
		return runID, failRun(ctx, err)
	}

	edges, err := s.repo.ListWorkflowEdgesForExecution(ctx, workflowID)
	if err != nil {
		return runID, failRun(ctx, err)
	}

	if len(steps) == 0 {
		return runID, failRun(ctx, errors.New("workflow has no steps"))
	}

	stepMap := make(map[uuid.UUID]db.WorkflowStep)
	incomingCount := make(map[uuid.UUID]int)

	type NextEdge struct {
		TargetStepID    uuid.UUID
		ConditionBranch string
	}

	nextSteps := make(map[uuid.UUID][]NextEdge)

	for _, step := range steps {
		stepMap[step.ID] = step
		incomingCount[step.ID] = 0
	}

	for _, edge := range edges {
		branch := ""
		if edge.ConditionBranch.Valid {
			branch = edge.ConditionBranch.String
		}
		nextSteps[edge.SourceStepID] = append(
			nextSteps[edge.SourceStepID],
			NextEdge{
				TargetStepID:    edge.TargetStepID,
				ConditionBranch: branch,
			},
		)
		incomingCount[edge.TargetStepID]++
	}

	var startStepID uuid.UUID
	foundStart := false

	for _, step := range steps {
		if incomingCount[step.ID] == 0 {
			startStepID = step.ID
			foundStart = true
			break
		}
	}

	if !foundStart {
		return runID, failRun(ctx, errors.New("no start node found"))
	}

	exec := executor.NewExecutor(s.geminiService)
	visited := make(map[uuid.UUID]bool)
	nodeInputs := make(map[uuid.UUID][]byte)

	if len(initialInput) > 0 {
		nodeInputs[startStepID] = initialInput
	}

	var executeNode func(execCtx context.Context, stepID uuid.UUID) error

	executeNode = func(execCtx context.Context, stepID uuid.UUID) error {
		if visited[stepID] {
			return nil
		}

		step, ok := stepMap[stepID]
		if !ok {
			return errors.New("step not found in workflow")
		}

		visited[stepID] = true

		stepRunID := uuid.New()
		input := nodeInputs[stepID]

		if len(input) == 0 {
			input = step.Config
		}

		_, err := s.repo.CreateWorkflowStepRun(execCtx, db.CreateWorkflowStepRunParams{
			ID:             stepRunID,
			WorkflowRunID:  runID,
			WorkflowStepID: step.ID,
			Status:         "running",
			Input: pqtype.NullRawMessage{
				RawMessage: input,
				Valid:      len(input) > 0,
			},
			Output: pqtype.NullRawMessage{
				Valid: false,
			},
			ErrorMessage: sql.NullString{
				Valid: false,
			},
		})
		if err != nil {
			return err
		}

		s.Publish(workflowID.String(), map[string]any{
			"type":      "step_status",
			"step_id":   step.ID.String(),
			"step_type": step.StepType,
			"status":    "running",
			"message":   step.StepType + " started",
		})

		output, err := exec.ExecuteStep(userID, step, nodeInputs[stepID])
		fmt.Println("STEP TYPE:", step.StepType)
		fmt.Println("STEP INPUT:", string(nodeInputs[stepID]))
		fmt.Println("STEP OUTPUT:", string(output))
		if err != nil {
			s.Publish(workflowID.String(), map[string]any{
				"type":      "step_status",
				"step_id":   step.ID.String(),
				"step_type": step.StepType,
				"status":    "failed",
				"message":   step.StepType + " failed: " + err.Error(),
				"error":     err.Error(),
			})

			_ = s.repo.UpdateWorkflowStepRunStatus(execCtx, db.UpdateWorkflowStepRunStatusParams{
				ID:     stepRunID,
				Status: "failed",
				Output: pqtype.NullRawMessage{
					RawMessage: output,
					Valid:      output != nil,
				},
				ErrorMessage: sql.NullString{
					String: err.Error(),
					Valid:  true,
				},
			})
			return err
		}

		err = s.repo.UpdateWorkflowStepRunStatus(execCtx, db.UpdateWorkflowStepRunStatusParams{
			ID:     stepRunID,
			Status: "success",
			Output: pqtype.NullRawMessage{
				RawMessage: output,
				Valid:      output != nil,
			},
			ErrorMessage: sql.NullString{
				Valid: false,
			},
		})
		if err != nil {
			return err
		}

		s.Publish(workflowID.String(), map[string]any{
			"type":      "step_status",
			"step_id":   step.ID.String(),
			"step_type": step.StepType,
			"status":    "success",
			"message":   step.StepType + " completed successfully",
		})

		// Debug logs (kept as fmt.Println for compatibility, but could be replaced)
		fmt.Println("===================================")
		log.Println("CURRENT STEP:", step.StepType)
		fmt.Println("CURRENT STEP ID:", stepID)
		fmt.Println("NEXT EDGES:", nextSteps[stepID])

		for _, edge := range nextSteps[stepID] {
			fmt.Println("EDGE TARGET:", edge.TargetStepID)
			fmt.Println("EDGE BRANCH:", edge.ConditionBranch)

			if step.StepType == "condition" {
				var conditionResult struct {
					Result bool `json:"result"`
				}
				if err := json.Unmarshal(output, &conditionResult); err != nil {
					return err
				}
				fmt.Println("CONDITION RESULT:", conditionResult.Result)

				if conditionResult.Result && edge.ConditionBranch != "true" {
					fmt.Println("SKIPPING FALSE EDGE")
					continue
				}
				if !conditionResult.Result && edge.ConditionBranch != "false" {
					fmt.Println("SKIPPING TRUE EDGE")
					continue
				}
			}

			fmt.Println("EXECUTING NEXT NODE:", edge.TargetStepID)

			if step.StepType == "condition" {
				nodeInputs[edge.TargetStepID] = nodeInputs[stepID]
			} else {
				nodeInputs[edge.TargetStepID] = output
			}

			if err := executeNode(execCtx, edge.TargetStepID); err != nil {
				return err
			}
		}
		return nil
	}

	if err := executeNode(ctx, startStepID); err != nil {
		_ = failRun(ctx, err)
		return runID, err
	}

	_ = s.repo.UpdateWorkflowRunStatus(ctx, db.UpdateWorkflowRunStatusParams{
		ID:     runID,
		Status: "success",
		ErrorMessage: sql.NullString{
			Valid: false,
		},
	})

	s.logger.Info("workflow run completed successfully",
		"workflow_id", workflowID.String(),
		"run_id", runID.String(),
	)
	return runID, nil
}

func (s *WorkflowService) UpdateWorkflowSchedule(
	ctx context.Context,
	workflowID uuid.UUID,
	userID uuid.UUID,
	req UpdateWorkflowScheduleRequest,
) error {
	s.logger.Info("updating workflow schedule",
		"workflow_id", workflowID.String(),
		"user_id", userID.String(),
		"enabled", req.Enabled,
	)

	var nextRunAt sql.NullTime

	if req.Enabled && req.ScheduleValue == "" {
		err := errors.New("schedule value is required")
		s.logger.Warn("schedule update failed: missing schedule value",
			"workflow_id", workflowID.String(),
		)
		return err
	}

	if req.Enabled {
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

		schedule, err := parser.Parse(req.ScheduleValue)
		if err != nil {
			s.logger.Error("invalid cron expression",
				"workflow_id", workflowID.String(),
				"schedule_value", req.ScheduleValue,
				"error", err.Error(),
			)
			return err
		}

		nextRunAt = sql.NullTime{
			Time:  schedule.Next(time.Now()),
			Valid: true,
		}
	}

	err := s.repo.UpdateWorkflowSchedule(ctx, db.UpdateWorkflowScheduleParams{
		ID:              workflowID,
		ScheduleEnabled: req.Enabled,
		ScheduleType: sql.NullString{
			String: req.ScheduleType,
			Valid:  req.ScheduleType != "",
		},
		ScheduleValue: sql.NullString{
			String: req.ScheduleValue,
			Valid:  req.ScheduleValue != "",
		},
		NextRunAt: nextRunAt,
		UserID:    userID,
	})
	if err != nil {
		s.logger.Error("failed to update workflow schedule",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return err
	}

	s.logger.Info("workflow schedule updated",
		"workflow_id", workflowID.String(),
	)
	return nil
}

func (s *WorkflowService) ListDueScheduledWorkflows(
	ctx context.Context,
) ([]db.Workflow, error) {
	s.logger.Info("listing due scheduled workflows")

	workflows, err := s.repo.ListDueScheduledWorkflows(ctx)
	if err != nil {
		s.logger.Error("failed to list due scheduled workflows",
			"error", err.Error(),
		)
		return nil, err
	}

	s.logger.Info("due scheduled workflows listed",
		"count", len(workflows),
	)
	return workflows, nil
}

func (s *WorkflowService) RunScheduledWorkflow(
	ctx context.Context,
	workflow db.Workflow,
) (err error) {
	s.logger.Info("executing scheduled workflow",
		"workflow_id", workflow.ID.String(),
	)

	err = s.repo.MarkScheduleRunning(
		ctx,
		db.MarkScheduleRunningParams{
			ID:                workflow.ID,
			IsScheduleRunning: true,
		},
	)
	if err != nil {
		s.logger.Error("failed to mark schedule as running",
			"workflow_id", workflow.ID.String(),
			"error", err.Error(),
		)
		return err
	}

	defer func() {
		if err != nil {
			s.unlockSchedule(ctx, workflow.ID)
		}
	}()

	_, err = s.RunWorkflow(
		ctx,
		workflow.ID,
		workflow.UserID,
	)
	if err != nil {
		s.logger.Error("scheduled workflow execution failed",
			"workflow_id", workflow.ID.String(),
			"error", err.Error(),
		)
		return err
	}

	if !workflow.ScheduleValue.Valid {
		err := errors.New("schedule value is empty")
		s.logger.Error("scheduled workflow missing schedule value",
			"workflow_id", workflow.ID.String(),
		)
		return err
	}

	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	schedule, err := parser.Parse(workflow.ScheduleValue.String)
	if err != nil {
		s.logger.Error("failed to parse schedule for next run",
			"workflow_id", workflow.ID.String(),
			"error", err.Error(),
		)
		return err
	}

	nextRunAt := schedule.Next(time.Now())

	err = s.repo.MarkWorkflowScheduleRun(
		ctx,
		db.MarkWorkflowScheduleRunParams{
			ID: workflow.ID,
			NextRunAt: sql.NullTime{
				Time:  nextRunAt,
				Valid: true,
			},
		},
	)
	if err != nil {
		s.logger.Error("failed to mark schedule run completion",
			"workflow_id", workflow.ID.String(),
			"error", err.Error(),
		)
		return err
	}

	s.logger.Info("scheduled workflow executed successfully",
		"workflow_id", workflow.ID.String(),
		"next_run_at", nextRunAt,
	)
	return nil
}

func (s *WorkflowService) unlockSchedule(
	ctx context.Context,
	workflowID uuid.UUID,
) {
	err := s.repo.MarkScheduleRunning(
		ctx,
		db.MarkScheduleRunningParams{
			ID:                workflowID,
			IsScheduleRunning: false,
		},
	)
	if err != nil {
		s.logger.Error("failed to unlock schedule",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
	}
}

func (s *WorkflowService) GetWorkflowSchedule(
	ctx context.Context,
	workflowID uuid.UUID,
	userID uuid.UUID,
) (db.GetWorkflowScheduleRow, error) {
	s.logger.Info("getting workflow schedule",
		"workflow_id", workflowID.String(),
		"user_id", userID.String(),
	)

	schedule, err := s.repo.GetWorkflowSchedule(
		ctx,
		db.GetWorkflowScheduleParams{
			ID:     workflowID,
			UserID: userID,
		},
	)
	if err != nil {
		s.logger.Error("failed to get workflow schedule",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		return schedule, err
	}

	s.logger.Info("workflow schedule retrieved",
		"workflow_id", workflowID.String(),
	)
	return schedule, nil
}

func (s *WorkflowService) Subscribe(workflowID string) chan any {
	s.logger.Info("client subscribed to workflow events",
		"workflow_id", workflowID,
	)
	ch := make(chan any, 10)

	s.mu.Lock()
	s.subscribers[workflowID] = append(s.subscribers[workflowID], ch)
	s.mu.Unlock()

	return ch
}

func (s *WorkflowService) Unsubscribe(workflowID string, ch chan any) {
	s.logger.Info("client unsubscribed from workflow events",
		"workflow_id", workflowID,
	)
	s.mu.Lock()
	defer s.mu.Unlock()

	channels := s.subscribers[workflowID]

	for i, subscriber := range channels {
		if subscriber == ch {
			s.subscribers[workflowID] = append(channels[:i], channels[i+1:]...)
			close(ch)
			break
		}
	}
}

func (s *WorkflowService) Publish(workflowID string, data any) {
	s.mu.Lock()
	channels := s.subscribers[workflowID]
	s.mu.Unlock()

	for _, ch := range channels {
		select {
		case ch <- data:
		default:
		}
	}
}

func (s *WorkflowService) UpdateWorkflow(
	ctx context.Context,
	workflowID uuid.UUID,
	userID uuid.UUID,
	name string,
	description string,
) (db.Workflow, error) {
	return s.repo.UpdateWorkflow(ctx, db.UpdateWorkflowParams{
		ID:          workflowID,
		UserID:      userID,
		Name:        name,
		Description: description,
	})
}

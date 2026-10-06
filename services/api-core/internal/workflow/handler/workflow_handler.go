package handler

import (
	"encoding/json"
	"io"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"org/api-core/internal/billing"
	usageService "org/api-core/internal/usage/service"
	"org/api-core/internal/utils/response"
	"org/api-core/internal/workflow/service"
)

type WorkflowHandler struct {
	workflowService *service.WorkflowService
	usageService    *usageService.Service
	logger          *slog.Logger
}

func NewWorkflowHandler(workflowService *service.WorkflowService, usageService *usageService.Service, logger *slog.Logger) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
		usageService:    usageService,
		logger:          logger,
	}
}

type SaveWorkflowStepsRequest struct {
	Steps []service.SaveWorkflowStepRequest `json:"steps"`
	Edges []service.SaveWorkflowEdgeRequest `json:"edges"`
}

type SaveWorkflowRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateWorkflowScheduleRequest struct {
	ScheduleEnabled bool   `json:"schedule_enabled"`
	ScheduleType    string `json:"schedule_type"`
	ScheduleValue   string `json:"schedule_value"`
}

// CREATE WORKFLOW
func (h *WorkflowHandler) CreateWorkflow(c *gin.Context) {
	h.logger.Info("create workflow request received")

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn("create workflow: invalid input",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid input", err.Error())
		return
	}

	if body.Name == "" {
		h.logger.Warn("create workflow: missing name")
		response.BadRequest(c, "workflow name is required", nil)
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("create workflow: user not authenticated")
		response.Unauthorized(c, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDValue.(string))
	if err != nil {
		h.logger.Warn("create workflow: invalid user id",
			"user_id", userIDValue,
		)
		response.BadRequest(c, "invalid user id", err.Error())
		return
	}

	workflow, err := h.workflowService.CreateWorkflow(
		c.Request.Context(),
		userID,
		body.Name,
		body.Description,
	)

	if err != nil {
		h.logger.Error("create workflow: failed",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to create workflow", err.Error())
		return
	}

	h.logger.Info("create workflow: success",
		"user_id", userID.String(),
		"workflow_id", workflow.ID.String(),
		"name", workflow.Name,
	)
	response.Created(c, "workflow created successfully", workflow)
}

// LIST WORKFLOWS
func (h *WorkflowHandler) ListWorkflows(c *gin.Context) {
	h.logger.Info("list workflows request received")

	userIDValue, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("list workflows: user not authenticated")
		response.Unauthorized(c, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDValue.(string))
	if err != nil {
		h.logger.Warn("list workflows: invalid user id",
			"user_id", userIDValue,
		)
		response.BadRequest(c, "invalid user id", err.Error())
		return
	}

	workflows, err := h.workflowService.ListWorkflow(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("list workflows: failed",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to fetch workflows", err.Error())
		return
	}

	h.logger.Info("list workflows: success",
		"user_id", userID.String(),
		"count", len(workflows),
	)
	response.OK(c, "workflows fetched successfully", workflows)
}

// DELETE WORKFLOW
func (h *WorkflowHandler) DeleteWorkflow(c *gin.Context) {
	h.logger.Info("delete workflow request received")

	userIDValue, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("delete workflow: user not authenticated")
		response.Unauthorized(c, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDValue.(string))
	if err != nil {
		h.logger.Warn("delete workflow: invalid user id",
			"user_id", userIDValue,
		)
		response.BadRequest(c, "invalid user id", err.Error())
		return
	}

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("delete workflow: invalid workflow id",
			"workflow_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow id", err.Error())
		return
	}

	err = h.workflowService.DeleteWorkflow(c.Request.Context(), userID, workflowID)
	if err != nil {
		h.logger.Error("delete workflow: failed",
			"user_id", userID.String(),
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to delete workflow", err.Error())
		return
	}

	h.logger.Info("delete workflow: success",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
	)
	response.OK(c, "workflow deleted successfully", nil)
}

// CREATE SINGLE STEP
func (h *WorkflowHandler) CreateStep(c *gin.Context) {
	h.logger.Info("create step request received")

	var body struct {
		StepOrder int32                  `json:"step_order"`
		StepType  string                 `json:"step_type"`
		Config    map[string]interface{} `json:"config"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn("create step: invalid input",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid input", err.Error())
		return
	}

	if body.StepType == "" {
		h.logger.Warn("create step: missing step type")
		response.BadRequest(c, "step type is required", nil)
		return
	}

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("create step: invalid workflow id",
			"workflow_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow id", err.Error())
		return
	}

	configBytes, err := json.Marshal(body.Config)
	if err != nil {
		h.logger.Warn("create step: invalid config",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid config", err.Error())
		return
	}

	step, err := h.workflowService.CreateStep(
		c.Request.Context(),
		workflowID,
		body.StepOrder,
		body.StepType,
		configBytes,
	)

	if err != nil {
		h.logger.Error("create step: failed",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to create step", err.Error())
		return
	}

	h.logger.Info("create step: success",
		"workflow_id", workflowID.String(),
		"step_id", step.ID.String(),
		"step_type", body.StepType,
	)
	response.Created(c, "workflow step created successfully", step)
}

// LIST STEPS
func (h *WorkflowHandler) ListSteps(c *gin.Context) {
	h.logger.Info("list steps request received")

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("list steps: invalid workflow id",
			"workflow_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow id", err.Error())
		return
	}

	steps, err := h.workflowService.ListSteps(c.Request.Context(), workflowID)
	if err != nil {
		h.logger.Error("list steps: failed",
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to fetch steps", err.Error())
		return
	}

	h.logger.Info("list steps: success",
		"workflow_id", workflowID.String(),
		"count", len(steps),
	)
	response.OK(c, "workflow steps fetched successfully", steps)
}

// SAVE CANVAS STEPS
func (h *WorkflowHandler) SaveWorkflowSteps(c *gin.Context) {
	h.logger.Info("save workflow steps request received")

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("save workflow steps: invalid workflow id",
			"workflow_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow id", err.Error())
		return
	}

	userIDString := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		h.logger.Warn("save workflow steps: invalid user id",
			"user_id", userIDString,
		)
		response.Unauthorized(c, "invalid user")
		return
	}

	var body SaveWorkflowStepsRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn("save workflow steps: invalid request body",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid request body", err.Error())
		return
	}

	err = h.workflowService.SaveWorkflowSteps(
		c.Request.Context(),
		workflowID,
		userID,
		body.Steps,
		body.Edges,
	)

	if err != nil {
		h.logger.Error("save workflow steps: failed",
			"user_id", userID.String(),
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to save workflow steps", err.Error())
		return
	}

	h.logger.Info("save workflow steps: success",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
		"steps_count", len(body.Steps),
		"edges_count", len(body.Edges),
	)
	response.OK(c, "workflow steps saved successfully", nil)
}

// GET CANVAS STEPS
func (h *WorkflowHandler) GetWorkflowSteps(c *gin.Context) {
	h.logger.Info("get workflow steps request received")

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("get workflow steps: invalid workflow id",
			"workflow_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow id", err.Error())
		return
	}

	userIDString := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		h.logger.Warn("get workflow steps: invalid user id",
			"user_id", userIDString,
		)
		response.Unauthorized(c, "invalid user")
		return
	}

	steps, err := h.workflowService.GetWorkflowSteps(
		c.Request.Context(),
		workflowID,
		userID,
	)

	if err != nil {
		h.logger.Error("get workflow steps: failed",
			"user_id", userID.String(),
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to fetch workflow steps", err.Error())
		return
	}

	h.logger.Info("get workflow steps: success",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
	)
	response.OK(c, "workflow steps fetched successfully", steps)
}

// RUN WORKFLOW
func (h *WorkflowHandler) RunWorkflow(c *gin.Context) {
	h.logger.Info("run workflow request received")

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("run workflow: invalid workflow id",
			"workflow_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow id", err.Error())
		return
	}

	userIDString := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		h.logger.Warn("run workflow: invalid user id",
			"user_id", userIDString,
		)
		response.Unauthorized(c, "invalid user")
		return
	}

	// CHECK SUBSCRIPTION
	plan, status, err := billing.GetUserSubscription(userIDString)
	if err != nil {
		h.logger.Error("run workflow: failed to check subscription",
			"user_id", userIDString,
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to check subscription", err.Error())
		return
	}

	if status != "active" {
		h.logger.Warn("run workflow: subscription inactive",
			"user_id", userIDString,
			"status", status,
		)
		response.Forbidden(c, "subscription inactive")
		return
	}

	if plan == "free" {
		h.logger.Warn("run workflow: free plan upgrade required",
			"user_id", userIDString,
		)
		response.Forbidden(c, "upgrade required to run workflows")
		return
	}

	// CHECK USAGE LIMIT BEFORE RUNNING
	allowed, message, err := h.usageService.CanRunWorkflow(
		c.Request.Context(),
		userID,
		plan,
	)
	if err != nil {
		h.logger.Error("run workflow: failed to check usage",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to check usage", err.Error())
		return
	}

	if !allowed {
		h.logger.Warn("run workflow: usage limit reached",
			"user_id", userID.String(),
			"message", message,
		)
		response.Forbidden(c, message)
		return
	}

	runID, err := h.workflowService.RunWorkflow(
		c.Request.Context(),
		workflowID,
		userID,
	)
	if err != nil {
		h.logger.Error("run workflow: execution failed",
			"user_id", userID.String(),
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "workflow execution failed", err.Error())
		return
	}

	// INCREMENT USAGE AFTER SUCCESSFUL RUN
	if err := h.usageService.IncrementWorkflowRun(c.Request.Context(), userID); err != nil {
		// do not block response, just log
		h.logger.Warn("run workflow: failed to increment usage (non-blocking)",
			"user_id", userID.String(),
			"error", err.Error(),
		)
	}

	h.logger.Info("run workflow: success",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
		"run_id", runID,
	)
	response.OK(c, "workflow executed successfully", gin.H{
		"run_id": runID,
	})
}

// LIST WORKFLOW RUNS
func (h *WorkflowHandler) ListWorkflowRuns(c *gin.Context) {
	h.logger.Info("list workflow runs request received")

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("list workflow runs: invalid workflow id",
			"workflow_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow id", err.Error())
		return
	}

	userIDString := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		h.logger.Warn("list workflow runs: invalid user id",
			"user_id", userIDString,
		)
		response.Unauthorized(c, "invalid user")
		return
	}

	runs, err := h.workflowService.ListWorkflowRuns(c.Request.Context(), workflowID, userID)
	if err != nil {
		h.logger.Error("list workflow runs: failed",
			"user_id", userID.String(),
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to fetch workflow runs", err.Error())
		return
	}

	h.logger.Info("list workflow runs: success",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
		"count", len(runs),
	)
	response.OK(c, "workflow runs fetched successfully", runs)
}

// LIST STEP RUN LOGS
func (h *WorkflowHandler) ListWorkflowStepRuns(c *gin.Context) {
	h.logger.Info("list workflow step runs request received")

	runID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("list workflow step runs: invalid run id",
			"run_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow run id", err.Error())
		return
	}

	steps, err := h.workflowService.ListWorkflowStepRuns(c.Request.Context(), runID)
	if err != nil {
		h.logger.Error("list workflow step runs: failed",
			"run_id", runID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to fetch step logs", err.Error())
		return
	}

	h.logger.Info("list workflow step runs: success",
		"run_id", runID.String(),
		"count", len(steps),
	)
	response.OK(c, "workflow step logs fetched successfully", steps)
}

func (h *WorkflowHandler) GetWorkflowEdges(c *gin.Context) {
	h.logger.Info("get workflow edges request received")

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("get workflow edges: invalid workflow id",
			"workflow_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow id", err.Error())
		return
	}

	userIDString := c.MustGet("user_id").(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		h.logger.Warn("get workflow edges: invalid user id",
			"user_id", userIDString,
		)
		response.Unauthorized(c, "invalid user")
		return
	}

	edges, err := h.workflowService.GetWorkflowEdges(
		c.Request.Context(),
		workflowID,
		userID,
	)
	if err != nil {
		h.logger.Error("get workflow edges: failed",
			"user_id", userID.String(),
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to fetch workflow edges", err.Error())
		return
	}

	h.logger.Info("get workflow edges: success",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
		"count", len(edges),
	)
	response.OK(c, "workflow edges fetched successfully", edges)
}

func (h *WorkflowHandler) HandleWebhookTrigger(c *gin.Context) {
	webhookId := c.Param("webhookID")
	h.logger.Info("webhook trigger received",
		"webhook_id", webhookId,
	)

	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		h.logger.Warn("webhook trigger: invalid payload",
			"webhook_id", webhookId,
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid webhook payload", err.Error())
		return
	}

	runID, err := h.workflowService.RunWorkflowFromWebhook(c.Request.Context(), webhookId, payload)
	if err != nil {
		h.logger.Error("webhook trigger: execution failed",
			"webhook_id", webhookId,
			"error", err.Error(),
		)
		response.InternalServerError(c, "webhook execution failed", err.Error())
		return
	}

	h.logger.Info("webhook trigger: success",
		"webhook_id", webhookId,
		"run_id", runID,
	)
	response.OK(c, "webhook received", gin.H{
		"run_id": runID,
	})
}

func (h *WorkflowHandler) UpdateWorkflowSchedule(c *gin.Context) {
	h.logger.Info("update workflow schedule request received")

	workflowIDParam := c.Param("id")
	workflowID, err := uuid.Parse(workflowIDParam)
	if err != nil {
		h.logger.Warn("update workflow schedule: invalid workflow id",
			"workflow_id", workflowIDParam,
		)
		response.BadRequest(c, "invalid workflow id", nil)
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("update workflow schedule: user not authenticated")
		response.Unauthorized(c, "unauthorized")
		return
	}

	userIDString, ok := userIDValue.(string)
	if !ok {
		h.logger.Warn("update workflow schedule: invalid user id type")
		response.Unauthorized(c, "invalid user")
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		h.logger.Warn("update workflow schedule: invalid user id",
			"user_id", userIDString,
		)
		response.Unauthorized(c, "invalid user id")
		return
	}

	var req service.UpdateWorkflowScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("update workflow schedule: invalid request body",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid request body", nil)
		return
	}

	err = h.workflowService.UpdateWorkflowSchedule(
		c.Request.Context(),
		workflowID,
		userID,
		req,
	)

	if err != nil {
		h.logger.Error("update workflow schedule: failed",
			"user_id", userID.String(),
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to update workflow schedule", err.Error())
		return
	}

	h.logger.Info("update workflow schedule: success",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
		// "schedule_enabled", req.ScheduleEnabled,
		"schedule_type", req.ScheduleType,
	)
	response.OK(c, "workflow schedule updated", nil)
}

func (h *WorkflowHandler) GetWorkflowSchedule(c *gin.Context) {
	h.logger.Info("get workflow schedule request received")

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		h.logger.Warn("get workflow schedule: invalid workflow id",
			"workflow_id", c.Param("id"),
		)
		response.BadRequest(c, "invalid workflow id", nil)
		return
	}

	userIDValue, _ := c.Get("user_id")
	userIDString := userIDValue.(string)
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		h.logger.Warn("get workflow schedule: invalid user id",
			"user_id", userIDString,
		)
		response.BadRequest(c, "invalid user id", nil)
		return
	}

	schedule, err := h.workflowService.GetWorkflowSchedule(
		c.Request.Context(),
		workflowID,
		userID,
	)

	if err != nil {
		h.logger.Error("get workflow schedule: failed",
			"user_id", userID.String(),
			"workflow_id", workflowID.String(),
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to get workflow schedule", err.Error())
		return
	}

	h.logger.Info("get workflow schedule: success",
		"user_id", userID.String(),
		"workflow_id", workflowID.String(),
	)
	response.OK(c, "workflow schedule fetched", schedule)
}

func (h *WorkflowHandler) StreamWorkflowEvents(c *gin.Context) {
	workflowID := c.Param("id")
	h.logger.Info("SSE stream started",
		"workflow_id", workflowID,
	)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "http://localhost:5173")

	eventChan := h.workflowService.Subscribe(workflowID)
	defer func() {
		h.workflowService.Unsubscribe(workflowID, eventChan)
		h.logger.Info("SSE stream closed",
			"workflow_id", workflowID,
		)
	}()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-eventChan; ok {
			c.SSEvent("workflow_update", msg)
			return true
		}
		return false
	})
}



func (h *WorkflowHandler) UpdateWorkflow(c *gin.Context) {
	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, "invalid input", err.Error())
		return
	}

	if body.Name == "" {
		response.BadRequest(c, "workflow name is required", nil)
		return
	}

	workflowID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid workflow id", err.Error())
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDValue.(string))
	if err != nil {
		response.BadRequest(c, "invalid user id", err.Error())
		return
	}

	workflow, err := h.workflowService.UpdateWorkflow(
		c.Request.Context(),
		workflowID,
		userID,
		body.Name,
		body.Description,
	)
	if err != nil {
		response.InternalServerError(c, "failed to update workflow", err.Error())
		return
	}

	response.OK(c, "workflow updated successfully", workflow)
}

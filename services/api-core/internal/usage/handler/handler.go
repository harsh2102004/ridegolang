package handler

import (
	"log/slog"
	"net/http"

	usageService "org/api-core/internal/usage/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UsageHandler struct {
	usageService *usageService.Service
	logger       *slog.Logger
}

func NewUsageHandler(usageService *usageService.Service, logger *slog.Logger) *UsageHandler {
	return &UsageHandler{
		usageService: usageService,
		logger:       logger,
	}
}

func (h *UsageHandler) GetUsage(c *gin.Context) {
	h.logger.Info("get usage request received")

	userIDString := c.GetString("user_id")
	if userIDString == "" {
		h.logger.Warn("get usage failed: user_id missing from context")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user",
		})
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		h.logger.Warn("get usage failed: invalid user_id format",
			"user_id", userIDString,
			"error", err.Error(),
		)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user",
		})
		return
	}

	usage, err := h.usageService.GetCurrentMonthUsage(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		// User may have no usage records yet; return default zero values
		h.logger.Warn("get usage: no usage data found or error, returning defaults",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		c.JSON(http.StatusOK, gin.H{
			"month":         h.usageService.CurrentMonth(),
			"workflow_runs": 0,
		})
		return
	}

	runs := int32(0)
	if usage.WorkflowRuns.Valid {
		runs = usage.WorkflowRuns.Int32
	}

	h.logger.Info("get usage successful",
		"user_id", userID.String(),
		"month", usage.Month,
		"workflow_runs", runs,
	)

	c.JSON(http.StatusOK, gin.H{
		"month":         usage.Month,
		"workflow_runs": runs,
	})
}
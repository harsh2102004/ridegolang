package handler

import (
	"net/http"

	"org/api-core/internal/gemini/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GeminiHandler struct {
	service service.GeminiService
}

func NewGeminiHandler(service service.GeminiService) *GeminiHandler {
	return &GeminiHandler{
		service: service,
	}
}

type SaveGeminiKeyRequest struct {
	ApiKey string `json:"api_key" binding:"required"`
}

func (h *GeminiHandler) GetKey(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	key, err := h.service.GetKey(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "gemini key not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_key": true,
		"id":      key.ID,
	})
}


func (h *GeminiHandler) SaveKey(c *gin.Context) {
	var req SaveGeminiKeyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	_, err = h.service.SaveKey(
		c.Request.Context(),
		userID,
		req.ApiKey,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "gemini api key saved successfully",
	})
}


func (h *GeminiHandler) UpdateKey(c *gin.Context) {
	var req SaveGeminiKeyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	_, err = h.service.UpdateKey(
		c.Request.Context(),
		userID,
		req.ApiKey,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "gemini api key updated successfully",
	})
}


func (h *GeminiHandler) DeleteKey(c *gin.Context) {
	userIDStr := c.GetString("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	err = h.service.DeleteKey(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "gemini api key deleted successfully",
	})
}
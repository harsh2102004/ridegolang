package handler

import (
	"log/slog"
	"net/http"

	"org/api-core/internal/oauth/provider"
	"org/api-core/internal/oauth/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type OAuthHandler struct {
	oauthService *service.OAuthService
	logger       *slog.Logger
}

func NewOAuthHandler(oauthService *service.OAuthService, logger *slog.Logger) *OAuthHandler {
	return &OAuthHandler{
		oauthService: oauthService,
		logger:       logger,
	}
}

func (h *OAuthHandler) GoogleStart(c *gin.Context) {
	h.logger.Info("google oauth start - connecting account")

	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("google oauth start failed: user not authenticated")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	state := userID.(string)

	config := provider.GoogleAuthConfig()
	url := config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.ApprovalForce,
	)

	h.logger.Info("google oauth start - redirecting to google",
		"user_id", state,
	)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	h.logger.Info("google oauth callback received")

	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		h.logger.Warn("google oauth callback missing code or state")
		c.JSON(http.StatusBadRequest, gin.H{"message": "code or state missing"})
		return
	}

	userID, err := uuid.Parse(state)
	if err != nil {
		h.logger.Warn("google oauth callback invalid state",
			"state", state,
		)
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid state"})
		return
	}

	config := provider.GoogleAuthConfig()
	token, err := config.Exchange(c.Request.Context(), code)
	if err != nil {
		h.logger.Error("google oauth callback token exchange failed",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to exchange token",
			"error":   err.Error(),
		})
		return
	}

	err = h.oauthService.SaveGoogleAccount(
		c.Request.Context(),
		userID,
		token,
		config.Scopes,
	)
	if err != nil {
		h.logger.Error("google oauth callback failed to save account",
			"user_id", userID.String(),
			"error", err.Error(),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to save google account",
			"error":   err.Error(),
		})
		return
	}

	h.logger.Info("google oauth callback successful - account connected",
		"user_id", userID.String(),
	)
	c.JSON(http.StatusOK, gin.H{
		"message": "google account connected successfully",
	})
}
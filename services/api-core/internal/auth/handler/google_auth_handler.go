package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"org/api-core/internal/auth/service"
	"org/api-core/internal/oauth/provider"
	"org/api-core/internal/utils/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type GoogleAuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

func NewGoogleAuthHandler(authService *service.AuthService, logger *slog.Logger) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *GoogleAuthHandler) GoogleAuthStart(c *gin.Context) {
	h.logger.Info("google oauth start - redirecting to google")

	config := provider.GoogleAuthConfig()
	url := config.AuthCodeURL(
		"oauth_signup",
		oauth2.AccessTypeOffline,
	)

	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *GoogleAuthHandler) GoogleAuthCallback(c *gin.Context) {
	h.logger.Info("google oauth callback received")

	code := c.Query("code")
	if code == "" {
		h.logger.Warn("google oauth callback: missing code")
		response.BadRequest(c, "code missing", nil)
		return
	}

	config := provider.GoogleAuthConfig()

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		h.logger.Error("google oauth callback: token exchange failed",
			"error", err.Error(),
		)
		response.InternalServerError(c, "token exchange failed", err.Error())
		return
	}

	client := config.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		h.logger.Error("google oauth callback: failed to fetch user info",
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to fetch google user", err.Error())
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		h.logger.Error("google oauth callback: failed to decode user info",
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to decode google user", err.Error())
		return
	}

	if userInfo.Email == "" {
		h.logger.Warn("google oauth callback: google account email missing")
		response.BadRequest(c, "google account email missing", nil)
		return
	}

	tokenPair, err := h.authService.GetOrCreateGoogleUser(
		c.Request.Context(),
		userInfo.Email,
	)
	if err != nil {
		h.logger.Error("google oauth callback: failed to get or create user",
			"email", userInfo.Email,
			"error", err.Error(),
		)
		response.InternalServerError(c, "google auth failed", err.Error())
		return
	}

	frontendURL := os.Getenv("APP_URL")
	redirectURL := fmt.Sprintf(
		"%s/oauth/callback?access_token=%s&refresh_token=%s",
		frontendURL,
		tokenPair.AccessToken,
		tokenPair.RefreshToken,
	)

	h.logger.Info("google oauth callback successful",
		"email", userInfo.Email,
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
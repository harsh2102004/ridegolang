package handler

import (
	"context"
	"log/slog"
	"org/api-core/internal/auth/service"
	"org/api-core/internal/billing"
	"org/api-core/internal/utils/response"
	"time"

	pb "org/api-core/proto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

type MeResponse struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Plan       string    `json:"plan"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  string    `json:"created_at"`
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn("signup invalid request body",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid input", err.Error())
		return
	}

	if body.Email == "" || body.Password == "" {
		h.logger.Warn("signup missing required fields")
		response.BadRequest(c, "email and password required", nil)
		return
	}

	_, err := h.authService.Signup(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		h.logger.Warn("signup failed",
			"email", body.Email,
			"error", err.Error(),
		)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	h.logger.Info("signup successful",
		"email", body.Email,
	)
	response.Created(c, "user created. please check your email to verify account", nil)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn("login invalid request body",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid input", err.Error())
		return
	}

	if body.Email == "" || body.Password == "" {
		h.logger.Warn("login missing required fields")
		response.BadRequest(c, "email and password required", nil)
		return
	}

	tokenPair, err := h.authService.Login(c.Request.Context(), body.Email, body.Password)
	if err != nil {
		h.logger.Warn("login failed",
			"email", body.Email,
			"error", err.Error(),
		)
		response.Unauthorized(c, err.Error())
		return
	}

	h.logger.Info("login successful",
		"email", body.Email,
	)
	response.OK(c, "login successful", gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	h.logger.Info("me profile fetch requested")

	userIDValue, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("me failed: user not authenticated")
		response.Unauthorized(c, "user not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDValue.(string))
	if err != nil {
		h.logger.Warn("me failed: invalid user id format",
			"user_id", userIDValue,
		)
		response.BadRequest(c, "invalid user id", nil)
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("me failed: failed to fetch user profile",
			"user_id", userID,
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to fetch user profile", err)
		return
	}

	res := MeResponse{
		ID:         user.ID,
		Email:      user.Email,
		Plan:       user.Plan.String,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt.Time.Format(time.RFC3339),
	}

	h.logger.Info("me profile fetched successfully",
		"user_id", user.ID.String(),
		"email", user.Email,
	)
	response.OK(c, "user profile fetched successfully", res)
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	h.logger.Info("verify email request received")

	token := c.Query("token")
	if token == "" {
		h.logger.Warn("verify email missing token")
		response.BadRequest(c, "token is required", nil)
		return
	}

	err := h.authService.VerifyEmail(c.Request.Context(), token)
	if err != nil {
		h.logger.Warn("verify email failed",
			"error", err.Error(),
		)
		response.BadRequest(c, err.Error(), nil)
		return
	}

	h.logger.Info("email verified successfully")
	response.OK(c, "email verified successfully", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	h.logger.Info("forgot password request received")

	var body struct {
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn("forgot password invalid request body",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid input", err.Error())
		return
	}

	// We log only that the request was processed, not whether the email exists.
	err := h.authService.ForgotPassword(c.Request.Context(), body.Email)
	if err != nil {
		h.logger.Warn("forgot password service returned error",
			"error", err.Error(),
		)
		// Still return same success message to avoid email enumeration
		response.OK(c, "if this email exists, reset link has been sent", nil)
		return
	}

	h.logger.Info("forgot password request processed (email may or may not exist)",
		"email", body.Email,
	)
	response.OK(c, "if this email exists, reset link has been sent", nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	h.logger.Info("reset password request received")

	var body struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn("reset password invalid request body",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid input", err.Error())
		return
	}

	if body.Token == "" || body.NewPassword == "" {
		h.logger.Warn("reset password missing token or password")
		response.BadRequest(c, "token and new password required", nil)
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), body.Token, body.NewPassword)
	if err != nil {
		h.logger.Error("reset password failed",
			"error", err.Error(),
		)
		response.InternalServerError(c, "reset failed", err.Error())
		return
	}

	h.logger.Info("password reset successfully")
	response.OK(c, "password reset successfully", nil)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	h.logger.Info("refresh token request received")

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn("refresh token invalid request body",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid input", err.Error())
		return
	}

	tokenPair, err := h.authService.RefreshAccessToken(c.Request.Context(), body.RefreshToken)
	if err != nil {
		h.logger.Warn("refresh token failed: invalid token")
		response.Unauthorized(c, "invalid refresh token")
		return
	}

	h.logger.Info("token refreshed successfully")
	response.OK(c, "token refreshed", gin.H{
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.logger.Info("logout request received")

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		h.logger.Warn("logout invalid request body",
			"error", err.Error(),
		)
		response.BadRequest(c, "invalid input", err.Error())
		return
	}

	err := h.authService.Logout(c.Request.Context(), body.RefreshToken)
	if err != nil {
		h.logger.Error("logout failed",
			"error", err.Error(),
		)
		response.InternalServerError(c, "logout failed", err.Error())
		return
	}

	h.logger.Info("logout successful")
	response.OK(c, "logged out successfully", nil)
}

func TestBilling(c *gin.Context) {
	// Note: This is a standalone function, not part of AuthHandler.
	// We add basic logging using a background logger if needed.
	// Since it's a test endpoint, simple log is added.
	slog.Info("TestBilling called - fetching user subscription")

	res, err := billing.Client.GetUserSubscription(
		context.Background(),
		&pb.GetUserSubscriptionRequest{
			UserId: "11111111-1111-1111-1111-111111111111",
		},
	)

	if err != nil {
		slog.Error("TestBilling failed",
			"error", err.Error(),
		)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.Info("TestBilling successful",
		"plan", res.Plan,
		"status", res.Status,
	)
	c.JSON(200, gin.H{
		"plan":   res.Plan,
		"status": res.Status,
	})
}
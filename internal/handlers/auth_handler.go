package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/akatakan/nobetgo/internal/services"
	"github.com/akatakan/nobetgo/util"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := h.utilBindJSON(c, &req); err != nil {
		return
	}

	token, role, err := h.service.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		util.Unauthorized(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		Role:  role,
	})
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := h.utilBindJSON(c, &req); err != nil {
		return
	}

	token, err := h.service.GenerateResetToken(c.Request.Context(), req.Email)
	if err != nil {
		// We return OK even if user not found to avoid email enumeration
		// but log the error for internal debugging
		slog.Warn("Password reset requested for unknown/error email", "email", req.Email, "error", err)
		c.JSON(http.StatusOK, gin.H{"message": "Talimatlar e-posta adresinize gönderildi (eğer kayıtlıysa)."})
		return
	}

	// For dev/test only: log at Debug level (invisible in production info/warn/error levels)
	slog.Debug("PASSWORD RESET LINK (DEV)",
		"email", req.Email,
		"link", fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token),
	)

	c.JSON(http.StatusOK, gin.H{"message": "Talimatlar e-posta adresinize gönderildi."})
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := h.utilBindJSON(c, &req); err != nil {
		return
	}

	if err := h.service.ResetPassword(req.Token, req.Password); err != nil {
		util.BadRequest(c, err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Şifreniz başarıyla güncellendi."})
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := h.utilBindJSON(c, &req); err != nil {
		return
	}

	// userID is set by AuthMiddleware
	userID, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "Oturum geçersiz")
		return
	}

	if err := h.service.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		util.BadRequest(c, err.Error(), nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Şifreniz başarıyla değiştirildi."})
}

// utilBindJSON is a helper to DRY bind/error handling
func (h *AuthHandler) utilBindJSON(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		util.BadRequest(c, "Geçersiz veri", err)
		return err
	}
	return nil
}

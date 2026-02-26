package handlers

import (
	"net/http"

	"github.com/akatakan/nobetgo/config"
	"github.com/akatakan/nobetgo/internal/repositories"
	"github.com/akatakan/nobetgo/util"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repo   *repositories.EmployeeRepository
	config config.Config
}

func NewAuthHandler(repo *repositories.EmployeeRepository, config config.Config) *AuthHandler {
	return &AuthHandler{repo: repo, config: config}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "Geçersiz istek", err)
		return
	}

	employee, err := h.repo.GetByEmail(req.Email)
	if err != nil {
		util.Unauthorized(c, "Geçersiz e-posta veya şifre")
		return
	}

	if !util.CheckPasswordHash(req.Password, employee.PasswordHash) {
		util.Unauthorized(c, "Geçersiz e-posta veya şifre")
		return
	}

	token, err := util.GenerateToken(employee.ID, employee.Role, h.config.Server.JWTSecret)
	if err != nil {
		util.InternalError(c, "Token oluşturulamadı", err)
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		Role:  employee.Role,
	})
}

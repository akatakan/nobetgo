package handlers

import (
	"net/http"
	"strconv"

	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifService *services.NotificationService
}

func NewNotificationHandler(s *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: s}
}

func (h *NotificationHandler) GetUnread(c *gin.Context) {
	// For MVP, employee ID is passed via query param (later via auth token)
	empIDStr := c.Query("employee_id")
	empID, err := strconv.ParseUint(empIDStr, 10, 32)
	if err != nil || empID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçerli bir employee_id parametresi gereklidir"})
		return
	}

	notifs, err := h.notifService.GetUnread(uint(empID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bildirimler alınamadı"})
		return
	}

	c.JSON(http.StatusOK, notifs)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz bildirim ID"})
		return
	}

	if err := h.notifService.MarkAsRead(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bildirim okundu olarak işaretlenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bildirim okundu olarak işaretlendi"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	var input struct {
		EmployeeID uint `json:"employee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id parametresi gereklidir"})
		return
	}

	if err := h.notifService.MarkAllAsRead(input.EmployeeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bildirimler güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tüm bildirimler okundu olarak işaretlendi"})
}

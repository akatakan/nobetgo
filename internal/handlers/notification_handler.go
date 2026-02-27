package handlers

import (
	"net/http"
	"strconv"

	"github.com/akatakan/nobetgo/internal/services"
	"github.com/akatakan/nobetgo/util"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifService *services.NotificationService
}

func NewNotificationHandler(s *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifService: s}
}

func (h *NotificationHandler) GetUnread(c *gin.Context) {
	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	employeeID, ok := resolveEmployeeAccess(c, actor, "employee_id", false)
	if !ok {
		return
	}

	notifs, err := h.notifService.GetUnread(employeeID)
	if err != nil {
		util.InternalError(c, "Bildirimler alinamadi", err)
		return
	}

	c.JSON(http.StatusOK, notifs)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		util.BadRequest(c, "Gecersiz bildirim ID", err)
		return
	}

	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	updated, err := h.notifService.MarkAsRead(uint(id), actor.id)
	if err != nil {
		util.InternalError(c, "Bildirim guncellenemedi", err)
		return
	}
	if !updated {
		util.JSONError(c, http.StatusNotFound, "Bildirim bulunamadi", nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bildirim okundu olarak isaretlendi"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	if err := h.notifService.MarkAllAsRead(actor.id); err != nil {
		util.InternalError(c, "Bildirimler guncellenemedi", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tum bildirimler okundu olarak isaretlendi"})
}

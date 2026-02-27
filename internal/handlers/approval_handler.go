package handlers

import (
	"net/http"
	"strconv"

	"github.com/akatakan/nobetgo/internal/services"
	"github.com/akatakan/nobetgo/util"
	"github.com/gin-gonic/gin"
)

// ApprovalHandler handles HTTP requests for approval workflows and audit logs.
type ApprovalHandler struct {
	service *services.ApprovalService
}

// NewApprovalHandler creates a new ApprovalHandler.
func NewApprovalHandler(service *services.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{service: service}
}

// GetPendingApprovals handles GET /approvals/pending
func (h *ApprovalHandler) GetPendingApprovals(c *gin.Context) {
	pending, err := h.service.GetPendingApprovals()
	if err != nil {
		util.InternalError(c, "Bekleyen onaylar getirilemedi", err)
		return
	}

	c.JSON(http.StatusOK, pending)
}

// ApproveTimeEntry handles POST /approvals/time-entry/:id/approve
func (h *ApprovalHandler) ApproveTimeEntry(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	entry, err := h.service.ApproveTimeEntry(id, actor.id)
	if err != nil {
		util.InternalError(c, "Giris onaylanamadi", err)
		return
	}

	c.JSON(http.StatusOK, entry)
}

// RejectTimeEntry handles POST /approvals/time-entry/:id/reject
func (h *ApprovalHandler) RejectTimeEntry(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	entry, err := h.service.RejectTimeEntry(id, actor.id)
	if err != nil {
		util.InternalError(c, "Giris reddedilemedi", err)
		return
	}

	c.JSON(http.StatusOK, entry)
}

// GetAuditLogs handles GET /audit-logs?entity_type=&entity_id=
func (h *ApprovalHandler) GetAuditLogs(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityIDStr := c.Query("entity_id")

	if entityType == "" {
		util.BadRequest(c, "entity_type gerekli", nil)
		return
	}

	// If entity_id is provided, get logs for specific entity
	if entityIDStr != "" {
		entityID, err := strconv.ParseUint(entityIDStr, 10, 32)
		if err != nil {
			util.BadRequest(c, "Gecersiz entity_id", err)
			return
		}

		logs, err := h.service.GetAuditLogs(entityType, uint(entityID))
		if err != nil {
			util.InternalError(c, "Denetim gunlukleri getirilemedi", err)
			return
		}
		c.JSON(http.StatusOK, logs)
		return
	}

	// No entity_id, get logs by date range
	start, end, err := parseDateRange(c)
	if err != nil {
		util.BadRequest(c, "Tarih araligi gecersiz", err)
		return
	}

	logs, err := h.service.GetAuditLogsByDateRange(start, end)
	if err != nil {
		util.InternalError(c, "Denetim gunlukleri getirilemedi", err)
		return
	}

	c.JSON(http.StatusOK, logs)
}

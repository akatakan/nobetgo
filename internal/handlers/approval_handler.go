package handlers

import (
	"net/http"

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

	var body struct {
		ApproverID uint `json:"approver_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		util.BadRequest(c, "Geçersiz veri", err)
		return
	}

	entry, err := h.service.ApproveTimeEntry(id, body.ApproverID)
	if err != nil {
		util.InternalError(c, "Giriş onaylanamadı", err)
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

	var body struct {
		ApproverID uint `json:"approver_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		util.BadRequest(c, "Geçersiz veri", err)
		return
	}

	entry, err := h.service.RejectTimeEntry(id, body.ApproverID)
	if err != nil {
		util.InternalError(c, "Giriş reddedilemedi", err)
		return
	}

	c.JSON(http.StatusOK, entry)
}

// GetAuditLogs handles GET /audit-logs?entity_type=&entity_id=
func (h *ApprovalHandler) GetAuditLogs(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID, err := parseUintParam(c, "entity_id")
	if entityType == "" {
		util.BadRequest(c, "entity_type gerekli", nil)
		return
	}

	// If entity_id not provided, get logs by date range
	if err != nil {
		start, end, err := parseDateRange(c)
		if err != nil {
			util.BadRequest(c, "Tarih aralığı geçersiz", err)
			return
		}
		logs, err := h.service.GetAuditLogsByDateRange(start, end)
		if err != nil {
			util.InternalError(c, "Denetim günlükleri getirilemedi", err)
			return
		}
		c.JSON(http.StatusOK, logs)
		return
	}

	logs, err2 := h.service.GetAuditLogs(entityType, entityID)
	if err2 != nil {
		util.InternalError(c, "Denetim günlükleri getirilemedi", err2)
		return
	}

	c.JSON(http.StatusOK, logs)
}

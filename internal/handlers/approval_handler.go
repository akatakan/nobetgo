package handlers

import (
	"net/http"

	"github.com/akatakan/nobetgo/internal/services"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.service.ApproveTimeEntry(id, body.ApproverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.service.RejectTimeEntry(id, body.ApproverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// GetAuditLogs handles GET /audit-logs?entity_type=&entity_id=
func (h *ApprovalHandler) GetAuditLogs(c *gin.Context) {
	entityType := c.Query("entity_type")
	entityID, err := parseUintParam(c, "entity_id")
	if entityType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "entity_type gerekli"})
		return
	}

	// If entity_id not provided, get logs by date range
	if err != nil {
		start, end, err := parseDateRange(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		logs, err := h.service.GetAuditLogsByDateRange(start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, logs)
		return
	}

	logs, err2 := h.service.GetAuditLogs(entityType, entityID)
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err2.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

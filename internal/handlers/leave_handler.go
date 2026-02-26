package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
)

// LeaveHandler handles HTTP requests for leaves and leave types.
type LeaveHandler struct {
	service *services.LeaveService
}

// NewLeaveHandler creates a new LeaveHandler.
func NewLeaveHandler(service *services.LeaveService) *LeaveHandler {
	return &LeaveHandler{service: service}
}

// --- Leave Requests ---

// RequestLeave handles POST /leaves
func (h *LeaveHandler) RequestLeave(c *gin.Context) {
	var req core.LeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leave, err := h.service.RequestLeave(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, leave)
}

// GetLeave handles GET /leaves/:id
func (h *LeaveHandler) GetLeave(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	leave, err := h.service.GetLeave(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leave)
}

// ListLeaves handles GET /leaves with query params
func (h *LeaveHandler) ListLeaves(c *gin.Context) {
	startStr := c.DefaultQuery("start", "")
	endStr := c.DefaultQuery("end", "")

	var start, end time.Time
	if startStr != "" && endStr != "" {
		var err error
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz start tarihi"})
			return
		}
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz end tarihi"})
			return
		}
	} else {
		now := time.Now()
		start = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end = start.AddDate(0, 1, 0)
	}

	if empIDStr := c.Query("employee_id"); empIDStr != "" {
		empID, _ := strconv.ParseUint(empIDStr, 10, 32)
		leaves, err := h.service.GetEmployeeLeaves(uint(empID), start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, leaves)
		return
	}

	if deptIDStr := c.Query("department_id"); deptIDStr != "" {
		deptID, _ := strconv.ParseUint(deptIDStr, 10, 32)
		leaves, err := h.service.GetDepartmentLeaves(uint(deptID), start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, leaves)
		return
	}

	// Default: pending leaves
	leaves, err := h.service.GetPendingLeaves()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaves)
}

// ApproveLeave handles POST /leaves/:id/approve
func (h *LeaveHandler) ApproveLeave(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	// In a real system, approverID comes from auth context
	var body struct {
		ApproverID uint `json:"approver_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	leave, err := h.service.ApproveLeave(id, body.ApproverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leave)
}

// RejectLeave handles POST /leaves/:id/reject
func (h *LeaveHandler) RejectLeave(c *gin.Context) {
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

	leave, err := h.service.RejectLeave(id, body.ApproverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, leave)
}

// GetLeaveBalance handles GET /leaves/balance?employee_id=&year=
func (h *LeaveHandler) GetLeaveBalance(c *gin.Context) {
	empIDStr := c.Query("employee_id")
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	if empIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id gerekli"})
		return
	}

	empID, _ := strconv.ParseUint(empIDStr, 10, 32)
	year, _ := strconv.Atoi(yearStr)

	balances, err := h.service.GetLeaveBalance(uint(empID), year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, balances)
}

// --- LeaveType CRUD ---

// CreateLeaveType handles POST /leave-types
func (h *LeaveHandler) CreateLeaveType(c *gin.Context) {
	var lt core.LeaveType
	if err := c.ShouldBindJSON(&lt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateLeaveType(&lt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, lt)
}

// GetAllLeaveTypes handles GET /leave-types
func (h *LeaveHandler) GetAllLeaveTypes(c *gin.Context) {
	types, err := h.service.GetAllLeaveTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, types)
}

// GetLeaveType handles GET /leave-types/:id
func (h *LeaveHandler) GetLeaveType(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	lt, err := h.service.GetLeaveType(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lt)
}

// UpdateLeaveType handles PUT /leave-types/:id
func (h *LeaveHandler) UpdateLeaveType(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var lt core.LeaveType
	if err := c.ShouldBindJSON(&lt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	lt.ID = id

	if err := h.service.UpdateLeaveType(&lt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, lt)
}

// DeleteLeaveType handles DELETE /leave-types/:id
func (h *LeaveHandler) DeleteLeaveType(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.service.DeleteLeaveType(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "İzin türü silindi"})
}

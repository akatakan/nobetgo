package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/akatakan/nobetgo/util"
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
		util.BadRequest(c, "Geçersiz veri", err)
		return
	}

	leave, err := h.service.RequestLeave(req)
	if err != nil {
		util.JSONError(c, http.StatusConflict, "İzin talebi çakışıyor", err)
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
		util.JSONError(c, http.StatusNotFound, "İzin bulunamadı", err)
		return
	}

	c.JSON(http.StatusOK, leave)
}

// ListLeaves handles GET /leaves with query params
func (h *LeaveHandler) ListLeaves(c *gin.Context) {
	var params core.PaginationParams
	params.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	params.Search = c.Query("search")

	startStr := c.DefaultQuery("start", "")
	endStr := c.DefaultQuery("end", "")

	var start, end time.Time
	if startStr != "" && endStr != "" {
		start, _ = time.Parse("2006-01-02", startStr)
		end, _ = time.Parse("2006-01-02", endStr)
	}

	employeeID, _ := strconv.ParseUint(c.Query("employee_id"), 10, 32)
	departmentID, _ := strconv.ParseUint(c.Query("department_id"), 10, 32)

	result, err := h.service.GetPaginatedLeaves(params, uint(employeeID), uint(departmentID), start, end)
	if err != nil {
		util.InternalError(c, "İzinler getirilemedi", err)
		return
	}

	c.JSON(http.StatusOK, result)
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
		util.BadRequest(c, "Geçersiz veri", err)
		return
	}

	leave, err := h.service.ApproveLeave(id, body.ApproverID)
	if err != nil {
		util.InternalError(c, "İzin onaylanamadı", err)
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
		util.BadRequest(c, "Geçersiz veri", err)
		return
	}

	leave, err := h.service.RejectLeave(id, body.ApproverID)
	if err != nil {
		util.InternalError(c, "İzin reddedilemedi", err)
		return
	}

	c.JSON(http.StatusOK, leave)
}

// GetLeaveBalance handles GET /leaves/balance?employee_id=&year=
func (h *LeaveHandler) GetLeaveBalance(c *gin.Context) {
	empIDStr := c.Query("employee_id")
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	if empIDStr == "" {
		util.BadRequest(c, "employee_id gerekli", nil)
		return
	}

	empID, _ := strconv.ParseUint(empIDStr, 10, 32)
	year, _ := strconv.Atoi(yearStr)

	balances, err := h.service.GetLeaveBalance(uint(empID), year)
	if err != nil {
		util.InternalError(c, "Bakiye getirilemedi", err)
		return
	}

	c.JSON(http.StatusOK, balances)
}

// --- LeaveType CRUD ---

// CreateLeaveType handles POST /leave-types
func (h *LeaveHandler) CreateLeaveType(c *gin.Context) {
	var lt core.LeaveType
	if err := c.ShouldBindJSON(&lt); err != nil {
		util.BadRequest(c, "Geçersiz veri", err)
		return
	}

	if err := h.service.CreateLeaveType(&lt); err != nil {
		util.InternalError(c, "İzin türü oluşturulamadı", err)
		return
	}

	c.JSON(http.StatusCreated, lt)
}

// GetAllLeaveTypes handles GET /leave-types
func (h *LeaveHandler) GetAllLeaveTypes(c *gin.Context) {
	types, err := h.service.GetAllLeaveTypes()
	if err != nil {
		util.InternalError(c, "İzin türleri getirilemedi", err)
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
		util.JSONError(c, http.StatusNotFound, "İzin türü bulunamadı", err)
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
		util.BadRequest(c, "Geçersiz veri", err)
		return
	}
	lt.ID = id

	if err := h.service.UpdateLeaveType(&lt); err != nil {
		util.InternalError(c, "Güncelleme başarısız", err)
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
		util.InternalError(c, "Silme işlemi başarısız", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "İzin türü silindi"})
}

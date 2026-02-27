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
	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	var body struct {
		LeaveTypeID uint      `json:"leave_type_id" binding:"required"`
		StartDate   time.Time `json:"start_date" binding:"required"`
		EndDate     time.Time `json:"end_date" binding:"required"`
		Reason      string    `json:"reason"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		util.BadRequest(c, "Gecersiz veri", err)
		return
	}

	req := core.LeaveRequest{
		EmployeeID:  actor.id,
		LeaveTypeID: body.LeaveTypeID,
		StartDate:   body.StartDate,
		EndDate:     body.EndDate,
		Reason:      body.Reason,
	}

	leave, err := h.service.RequestLeave(req)
	if err != nil {
		util.JSONError(c, http.StatusConflict, "Izin talebi cakisiyor", err)
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

	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	leave, err := h.service.GetLeave(id)
	if err != nil {
		util.JSONError(c, http.StatusNotFound, "Izin bulunamadi", err)
		return
	}
	if !ensureOwnsEmployeeResource(c, actor, leave.EmployeeID) {
		return
	}

	c.JSON(http.StatusOK, leave)
}

// ListLeaves handles GET /leaves with query params
func (h *LeaveHandler) ListLeaves(c *gin.Context) {
	actor, ok := getActingUser(c)
	if !ok {
		return
	}

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

	employeeID, ok := resolveEmployeeAccess(c, actor, "employee_id", false)
	if !ok {
		return
	}

	departmentID, ok := parseOptionalUintQuery(c, "department_id")
	if !ok {
		return
	}
	if !actor.isAdmin() && departmentID != 0 {
		util.Forbidden(c, "Bu kaynaga erisim izniniz yok")
		return
	}

	result, err := h.service.GetPaginatedLeaves(params, employeeID, departmentID, start, end)
	if err != nil {
		util.InternalError(c, "Izinler getirilemedi", err)
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

	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	leave, err := h.service.ApproveLeave(id, actor.id)
	if err != nil {
		util.InternalError(c, "Izin onaylanamadi", err)
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

	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	leave, err := h.service.RejectLeave(id, actor.id)
	if err != nil {
		util.InternalError(c, "Izin reddedilemedi", err)
		return
	}

	c.JSON(http.StatusOK, leave)
}

// GetLeaveBalance handles GET /leaves/balance?employee_id=&year=
func (h *LeaveHandler) GetLeaveBalance(c *gin.Context) {
	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))

	employeeID, ok := resolveEmployeeAccess(c, actor, "employee_id", false)
	if !ok {
		return
	}

	year, _ := strconv.Atoi(yearStr)

	balances, err := h.service.GetLeaveBalance(employeeID, year)
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
		util.BadRequest(c, "Gecersiz veri", err)
		return
	}

	if err := h.service.CreateLeaveType(&lt); err != nil {
		util.InternalError(c, "Izin turu olusturulamadi", err)
		return
	}

	c.JSON(http.StatusCreated, lt)
}

// GetAllLeaveTypes handles GET /leave-types
func (h *LeaveHandler) GetAllLeaveTypes(c *gin.Context) {
	types, err := h.service.GetAllLeaveTypes()
	if err != nil {
		util.InternalError(c, "Izin turleri getirilemedi", err)
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
		util.JSONError(c, http.StatusNotFound, "Izin turu bulunamadi", err)
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
		util.BadRequest(c, "Gecersiz veri", err)
		return
	}
	lt.ID = id

	if err := h.service.UpdateLeaveType(&lt); err != nil {
		util.InternalError(c, "Guncelleme basarisiz", err)
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
		util.InternalError(c, "Silme islemi basarisiz", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Izin turu silindi"})
}

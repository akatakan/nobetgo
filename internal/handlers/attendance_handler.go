package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
)

// TimeEntryHandler handles HTTP requests for time entries.
type TimeEntryHandler struct {
	service *services.TimekeepingService
}

// NewTimeEntryHandler creates a new TimeEntryHandler.
func NewTimeEntryHandler(service *services.TimekeepingService) *TimeEntryHandler {
	return &TimeEntryHandler{service: service}
}

// ClockIn handles POST /time-entries/clock-in
func (h *TimeEntryHandler) ClockIn(c *gin.Context) {
	var req core.ClockInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.service.ClockIn(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// ClockOut handles POST /time-entries/clock-out
func (h *TimeEntryHandler) ClockOut(c *gin.Context) {
	var req core.ClockOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.service.ClockOut(req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// CreateTimeEntry handles POST /time-entries
func (h *TimeEntryHandler) CreateTimeEntry(c *gin.Context) {
	var req core.TimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.service.CreateTimeEntry(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// UpdateTimeEntry handles PUT /time-entries/:id
func (h *TimeEntryHandler) UpdateTimeEntry(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var req core.TimeEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.service.UpdateTimeEntry(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// DeleteTimeEntry handles DELETE /time-entries/:id
func (h *TimeEntryHandler) DeleteTimeEntry(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.service.DeleteTimeEntry(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kayıt silindi"})
}

// GetTimeEntry handles GET /time-entries/:id
func (h *TimeEntryHandler) GetTimeEntry(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	entry, err := h.service.GetTimeEntry(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// ListTimeEntries handles GET /time-entries with optional filters and pagination
func (h *TimeEntryHandler) ListTimeEntries(c *gin.Context) {
	var params core.PaginationParams
	params.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	params.Search = c.Query("search")

	start, end, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employeeID, _ := strconv.ParseUint(c.Query("employee_id"), 10, 32)
	departmentID, _ := strconv.ParseUint(c.Query("department_id"), 10, 32)

	result, err := h.service.GetPaginatedTimeEntries(params, uint(employeeID), uint(departmentID), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// --- Helper functions ---

func parseUintParam(c *gin.Context, name string) (uint, error) {
	idStr := c.Param(name)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz " + name})
		return 0, err
	}
	return uint(id), nil
}

func parseDateRange(c *gin.Context) (time.Time, time.Time, error) {
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		// Default to current month
		now := time.Now()
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		end := start.AddDate(0, 1, 0)
		return start, end, nil
	}

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

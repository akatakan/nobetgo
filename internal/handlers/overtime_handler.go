package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
)

// OvertimeHandler handles HTTP requests for overtime calculation, rules, and holidays.
type OvertimeHandler struct {
	service *services.OvertimeService
}

// NewOvertimeHandler creates a new OvertimeHandler.
func NewOvertimeHandler(service *services.OvertimeService) *OvertimeHandler {
	return &OvertimeHandler{service: service}
}

// CalculateOvertime handles GET /overtime/calculate?employee_id=&month=&year=
func (h *OvertimeHandler) CalculateOvertime(c *gin.Context) {
	actor, ok := getActingUser(c)
	if !ok {
		return
	}

	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id, month ve year gerekli"})
		return
	}

	employeeID, ok := resolveEmployeeAccess(c, actor, "employee_id", true)
	if !ok {
		return
	}

	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	summary, err := h.service.CalculateOvertime(employeeID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetDepartmentSummary handles GET /overtime/summary?department_id=&month=&year=
func (h *OvertimeHandler) GetDepartmentSummary(c *gin.Context) {
	actor, ok := getActingUser(c)
	if !ok {
		return
	}
	if !requireAdminAccess(c, actor) {
		return
	}

	deptIDStr := c.Query("department_id")
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month ve year gerekli"})
		return
	}

	deptID, _ := strconv.ParseUint(deptIDStr, 10, 32)
	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	summaries, err := h.service.GetDepartmentOvertimeSummary(uint(deptID), month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summaries)
}

// --- OvertimeRule CRUD ---

// CreateRule handles POST /overtime-rules
func (h *OvertimeHandler) CreateRule(c *gin.Context) {
	var rule core.OvertimeRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateRule(&rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// GetAllRules handles GET /overtime-rules
func (h *OvertimeHandler) GetAllRules(c *gin.Context) {
	rules, err := h.service.GetAllRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

// GetRule handles GET /overtime-rules/:id
func (h *OvertimeHandler) GetRule(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	rule, err := h.service.GetRule(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// UpdateRule handles PUT /overtime-rules/:id
func (h *OvertimeHandler) UpdateRule(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var rule core.OvertimeRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule.ID = id

	if err := h.service.UpdateRule(&rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// DeleteRule handles DELETE /overtime-rules/:id
func (h *OvertimeHandler) DeleteRule(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.service.DeleteRule(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kural silindi"})
}

// --- PublicHoliday CRUD ---

// CreateHoliday handles POST /public-holidays
func (h *OvertimeHandler) CreateHoliday(c *gin.Context) {
	var holiday core.PublicHoliday
	if err := c.ShouldBindJSON(&holiday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateHoliday(&holiday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, holiday)
}

// GetHolidays handles GET /public-holidays?year=
func (h *OvertimeHandler) GetHolidays(c *gin.Context) {
	yearStr := c.DefaultQuery("year", strconv.Itoa(time.Now().Year()))
	year, _ := strconv.Atoi(yearStr)

	holidays, err := h.service.GetHolidays(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, holidays)
}

// UpdateHoliday handles PUT /public-holidays/:id
func (h *OvertimeHandler) UpdateHoliday(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	var holiday core.PublicHoliday
	if err := c.ShouldBindJSON(&holiday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	holiday.ID = id

	if err := h.service.UpdateHoliday(&holiday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, holiday)
}

// DeleteHoliday handles DELETE /public-holidays/:id
func (h *OvertimeHandler) DeleteHoliday(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return
	}

	if err := h.service.DeleteHoliday(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tatil silindi"})
}

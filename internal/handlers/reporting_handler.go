package handlers

import (
	"net/http"
	"strconv"

	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
)

// ReportingHandler handles HTTP requests for reporting and analytics.
type ReportingHandler struct {
	service *services.ReportingService
}

// NewReportingHandler creates a new ReportingHandler.
func NewReportingHandler(service *services.ReportingService) *ReportingHandler {
	return &ReportingHandler{service: service}
}

// GetWorkHoursReport handles GET /reports/work-hours?department_id=&month=&year=
func (h *ReportingHandler) GetWorkHoursReport(c *gin.Context) {
	deptID, month, year, err := parseReportParams(c)
	if err != nil {
		return
	}

	report, err := h.service.GetWorkHoursReport(deptID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetAbsenceReport handles GET /reports/absences?department_id=&month=&year=
func (h *ReportingHandler) GetAbsenceReport(c *gin.Context) {
	deptID, month, year, err := parseReportParams(c)
	if err != nil {
		return
	}

	report, err := h.service.GetAbsenceReport(deptID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// GetEmployeeSummary handles GET /reports/employee-summary?employee_id=&month=&year=
func (h *ReportingHandler) GetEmployeeSummary(c *gin.Context) {
	empIDStr := c.Query("employee_id")
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if empIDStr == "" || monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_id, month ve year gerekli"})
		return
	}

	empID, _ := strconv.ParseUint(empIDStr, 10, 32)
	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	summary, err := h.service.GetEmployeeSummary(uint(empID), month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}

// GetTrendAnalysis handles GET /reports/trends?department_id=&start_month=&end_month=&year=
func (h *ReportingHandler) GetTrendAnalysis(c *gin.Context) {
	deptIDStr := c.DefaultQuery("department_id", "0")
	startMonthStr := c.Query("start_month")
	endMonthStr := c.Query("end_month")
	yearStr := c.Query("year")

	if startMonthStr == "" || endMonthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_month, end_month ve year gerekli"})
		return
	}

	deptID, _ := strconv.ParseUint(deptIDStr, 10, 32)
	startMonth, _ := strconv.Atoi(startMonthStr)
	endMonth, _ := strconv.Atoi(endMonthStr)
	year, _ := strconv.Atoi(yearStr)

	trends, err := h.service.GetTrendAnalysis(uint(deptID), startMonth, endMonth, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trends)
}

func parseReportParams(c *gin.Context) (uint, int, int, error) {
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month ve year gerekli"})
		return 0, 0, 0, errMissingParams
	}

	deptIDStr := c.DefaultQuery("department_id", "0")
	deptID, _ := strconv.ParseUint(deptIDStr, 10, 32)
	month, _ := strconv.Atoi(monthStr)
	year, _ := strconv.Atoi(yearStr)

	return uint(deptID), month, year, nil
}

var errMissingParams = &missingParamsError{}

type missingParamsError struct{}

func (e *missingParamsError) Error() string { return "eksik parametreler" }

// GetCostAnalysis handles GET /reports/cost-analysis?department_id=&month=&year=
func (h *ReportingHandler) GetCostAnalysis(c *gin.Context) {
	deptID, month, year, err := parseReportParams(c)
	if err != nil {
		return
	}

	report, err := h.service.GetCostAnalysis(deptID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

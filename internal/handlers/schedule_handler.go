package handlers

import (
	"net/http"
	"strconv"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	service *services.SchedulerService
}

func NewScheduleHandler(service *services.SchedulerService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

func (h *ScheduleHandler) GenerateSchedule(c *gin.Context) {
	var req core.ScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedules, err := h.service.GenerateSchedule(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedules)
}

func (h *ScheduleHandler) GetSchedule(c *gin.Context) {
	monthStr := c.Query("month")
	yearStr := c.Query("year")

	if monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month and year query parameters are required"})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	// Assuming SchedulerService also has a GetSchedule method, but currently it only has GenerateSchedule.
	// Since Generate return created schedules, maybe we should add Get to Service too?
	// But Repository has GetCombinedSchedule.
	// Let's add simple direct call if service doesn't have it, or expand service.
	// Wait, SchedulerService interface for repo has GetCombinedSchedule.
	// Let's assume we should add GetSchedule to SchedulerService as well?
	// But I implemented `GenerateSchedule` only in previous steps.
	// I should probably add `GetSchedule` to service or expose repo method via service.
	// For now, let's leave this method empty or implement it via a new Service method I'll add.
	// Or I can add `GetMonthlySchedule` to `SchedulerService`.

	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Not implemented yet",
		"month": month,
		"year":  year,
	})
}

func (h *ScheduleHandler) UpdateSchedule(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req core.Schedule
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedSchedule, err := h.service.UpdateSchedule(uint(id), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedSchedule)
}

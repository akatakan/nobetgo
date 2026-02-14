package handlers

import (
	"net/http"
	"strconv"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
)

type ShiftTypeHandler struct {
	service *services.ShiftTypeService
}

func NewShiftTypeHandler(service *services.ShiftTypeService) *ShiftTypeHandler {
	return &ShiftTypeHandler{service: service}
}

func (h *ShiftTypeHandler) CreateShiftType(c *gin.Context) {
	var shiftType core.ShiftType
	if err := c.ShouldBindJSON(&shiftType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateShiftType(&shiftType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shiftType)
}

func (h *ShiftTypeHandler) GetShiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	shiftType, err := h.service.GetShiftTypeByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ShiftType not found"})
		return
	}

	c.JSON(http.StatusOK, shiftType)
}

func (h *ShiftTypeHandler) GetAllShiftTypes(c *gin.Context) {
	shiftTypes, err := h.service.GetAllShiftTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shiftTypes)
}

func (h *ShiftTypeHandler) UpdateShiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var shiftType core.ShiftType
	if err := c.ShouldBindJSON(&shiftType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shiftType.ID = uint(id)

	if err := h.service.UpdateShiftType(&shiftType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shiftType)
}

func (h *ShiftTypeHandler) DeleteShiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteShiftType(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/akatakan/nobetgo/util"
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
		util.InternalError(c, "Çalışma tipi oluşturulamadı", err)
		return
	}

	c.JSON(http.StatusCreated, shiftType)
}

func (h *ShiftTypeHandler) GetShiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "Geçersiz ID", err)
		return
	}

	shiftType, err := h.service.GetShiftTypeByID(uint(id))
	if err != nil {
		util.JSONError(c, http.StatusNotFound, "Çalışma tipi bulunamadı", err)
		return
	}

	c.JSON(http.StatusOK, shiftType)
}

func (h *ShiftTypeHandler) GetAllShiftTypes(c *gin.Context) {
	shiftTypes, err := h.service.GetAllShiftTypes()
	if err != nil {
		util.InternalError(c, "Çalışma tipleri getirilemedi", err)
		return
	}

	c.JSON(http.StatusOK, shiftTypes)
}

func (h *ShiftTypeHandler) UpdateShiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "Geçersiz ID", err)
		return
	}

	var shiftType core.ShiftType
	if err := c.ShouldBindJSON(&shiftType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shiftType.ID = uint(id)

	if err := h.service.UpdateShiftType(&shiftType); err != nil {
		util.InternalError(c, "Güncelleme başarısız", err)
		return
	}

	c.JSON(http.StatusOK, shiftType)
}

func (h *ShiftTypeHandler) DeleteShiftType(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "Geçersiz ID", err)
		return
	}

	if err := h.service.DeleteShiftType(uint(id)); err != nil {
		util.InternalError(c, "Silme işlemi başarısız", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

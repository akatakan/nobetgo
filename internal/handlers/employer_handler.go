package handlers

import (
	"net/http"
	"strconv"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/akatakan/nobetgo/util"
	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	service *services.EmployeeService
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var employee core.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		util.BadRequest(c, "Geçersiz veri", err)
		return
	}

	if err := h.service.CreateEmployee(&employee); err != nil {
		util.InternalError(c, "Çalışan oluşturulamadı", err)
		return
	}

	c.JSON(http.StatusCreated, employee)
}

func (h *EmployeeHandler) GetEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "Geçersiz ID", err)
		return
	}

	employee, err := h.service.GetEmployeeByID(uint(id))
	if err != nil {
		util.JSONError(c, http.StatusNotFound, "Çalışan bulunamadı", err)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) GetAllEmployees(c *gin.Context) {
	var params core.PaginationParams
	params.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	params.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	params.Search = c.Query("search")

	// If no pagination params are provided and it's a simple list request, we might want to keep the old behavior
	// but for standard table view, pagination is better.
	// Let's assume the UI will always send at least defaults.

	result, err := h.service.GetPaginatedEmployees(c.Request.Context(), params)
	if err != nil {
		util.InternalError(c, "Çalışanlar getirilemedi", err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "Geçersiz ID", err)
		return
	}

	var employee core.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		util.BadRequest(c, "Geçersiz veri", err)
		return
	}

	// Ensure ID is set from URL param to avoid overwriting wrong record or creating new one if ID is missing in body
	employee.ID = uint(id)

	if err := h.service.UpdateEmployee(&employee); err != nil {
		util.InternalError(c, "Güncelleme başarısız", err)
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "Geçersiz ID", err)
		return
	}

	if err := h.service.DeleteEmployee(uint(id)); err != nil {
		util.InternalError(c, "Silme işlemi başarısız", err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *EmployeeHandler) ImportEmployees(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		util.BadRequest(c, "Dosya gerekli", err)
		return
	}
	defer file.Close()

	if err := h.service.ImportEmployees(file); err != nil {
		util.InternalError(c, "İçe aktarma başarısız", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Import successful"})
}

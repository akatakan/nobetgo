package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/akatakan/nobetgo/internal/core"
	"github.com/akatakan/nobetgo/internal/services"
	"github.com/gin-gonic/gin"
)

type TitleHandler struct {
	service *services.TitleService
}

func NewTitleHandler(service *services.TitleService) *TitleHandler {
	return &TitleHandler{service: service}
}

func (h *TitleHandler) CreateTitle(c *gin.Context) {
	var title core.Title
	if err := c.ShouldBindJSON(&title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateTitle(&title); err != nil {
		slog.Error("Failed to create title", "error", err, "name", title.Name)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Title created", "id", title.ID, "name", title.Name)
	c.JSON(http.StatusCreated, title)
}

func (h *TitleHandler) GetTitle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	title, err := h.service.GetTitleByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Title not found"})
		return
	}

	c.JSON(http.StatusOK, title)
}

func (h *TitleHandler) GetAllTitles(c *gin.Context) {
	titles, err := h.service.GetAllTitles()
	if err != nil {
		slog.Error("Failed to list titles", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, titles)
}

func (h *TitleHandler) UpdateTitle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var title core.Title
	if err := c.ShouldBindJSON(&title); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	title.ID = uint(id)

	if err := h.service.UpdateTitle(&title); err != nil {
		slog.Error("Failed to update title", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Title updated", "id", id)
	c.JSON(http.StatusOK, title)
}

func (h *TitleHandler) DeleteTitle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteTitle(uint(id)); err != nil {
		slog.Error("Failed to delete title", "error", err, "id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Title deleted", "id", id)
	c.JSON(http.StatusNoContent, nil)
}

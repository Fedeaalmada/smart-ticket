package controllers

import (
	"net/http"
	"strconv"

	"ticketek/backend/domain"
	"ticketek/backend/middleware"
	"ticketek/backend/services"

	"github.com/gin-gonic/gin"
)

// GET /eventos  — sin token (Funcionalidad A cliente)
func GetEventos(c *gin.Context) {
	filtros := map[string]string{
		"categoria": c.Query("categoria"),
		"fecha":     c.Query("fecha"),
		"titulo":    c.Query("titulo"),
	}

	eventos, err := services.ObtenerCatalogoEventos(filtros)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: "error al obtener eventos"})
		return
	}

	c.JSON(http.StatusOK, eventos)
}

// GET /eventos/:id  — sin token (Funcionalidad B cliente)
func GetEventoByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "id inválido"})
		return
	}

	evento, err := services.ObtenerDetalleEvento(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Error: "evento no encontrado"})
		return
	}

	c.JSON(http.StatusOK, evento)
}

// POST /admin/eventos  — solo admin (Funcionalidad A admin)
func CrearEvento(c *gin.Context) {
	var req domain.CrearEventoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	adminID := middleware.ObtenerUsuarioID(c)
	evento, err := services.CrearEvento(req, adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, evento)
}

// PUT /admin/eventos/:id  — solo admin (Funcionalidad B admin)
func ActualizarEvento(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "id inválido"})
		return
	}

	var req domain.ActualizarEventoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	if err := services.ActualizarEvento(uint(id), req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.MessageResponse{Message: "evento actualizado correctamente"})
}

// DELETE /admin/eventos/:id  — solo admin (Funcionalidad C admin)
func CancelarEvento(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "id inválido"})
		return
	}

	if err := services.CancelarEvento(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.MessageResponse{Message: "evento cancelado correctamente"})
}

// GET /admin/eventos/:id/reporte  — solo admin (Funcionalidad D admin)
func GetReporteEvento(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "id inválido"})
		return
	}

	reporte, err := services.ObtenerReporteEvento(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, reporte)
}

package controllers

import (
	"net/http"
	"strconv"

	"ticketek/backend/domain"
	"ticketek/backend/middleware"
	"ticketek/backend/services"

	"github.com/gin-gonic/gin"
)

// POST /entradas/comprar  — requiere token (Funcionalidad C cliente)
func ComprarEntrada(c *gin.Context) {
	var req domain.ComprarEntradaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	usuarioID := middleware.ObtenerUsuarioID(c)
	entrada, err := services.ComprarEntrada(usuarioID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entrada)
}

// GET /entradas/mis-entradas  — requiere token (Funcionalidad D cliente)
func GetMisEntradas(c *gin.Context) {
	usuarioID := middleware.ObtenerUsuarioID(c)

	entradas, err := services.ObtenerMisEntradas(usuarioID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: "error al obtener entradas"})
		return
	}

	c.JSON(http.StatusOK, entradas)
}

// DELETE /entradas/:id  — requiere token (Funcionalidad E cliente)
func CancelarEntrada(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "id inválido"})
		return
	}

	usuarioID := middleware.ObtenerUsuarioID(c)
	if err := services.CancelarEntrada(uint(id), usuarioID); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.MessageResponse{Message: "entrada cancelada correctamente"})
}

// POST /entradas/:id/transferir  — requiere token (Funcionalidad F cliente)
func TransferirEntrada(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "id inválido"})
		return
	}

	var req domain.TransferirEntradaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	usuarioID := middleware.ObtenerUsuarioID(c)
	if err := services.TransferirEntrada(uint(id), usuarioID, req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, domain.MessageResponse{Message: "entrada transferida correctamente"})
}

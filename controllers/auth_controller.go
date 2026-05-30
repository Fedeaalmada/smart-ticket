package controllers

import (
	"net/http"

	"ticketek/backend/domain"
	"ticketek/backend/services"

	"github.com/gin-gonic/gin"
)

// POST /auth/register
func Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := services.Registrar(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// POST /auth/login
func Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := services.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

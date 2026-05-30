package middleware

import (
	"net/http"
	"strings"

	"ticketek/backend/domain"
	"ticketek/backend/utils"

	"github.com/gin-gonic/gin"
)

// AuthRequerido valida el JWT en el header Authorization.
// Guarda los claims en el contexto de Gin para uso posterior.
func AuthRequerido() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{
				Error: "token requerido",
			})
			return
		}

		tokenString := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParsearToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorResponse{
				Error: "token inválido o expirado",
			})
			return
		}

		// Guardamos los datos del usuario en el contexto
		c.Set("usuario_id", claims.UsuarioID)
		c.Set("email", claims.Email)
		c.Set("rol", claims.Rol)
		c.Next()
	}
}

// SoloAdmin verifica que el usuario autenticado sea administrador.
// Debe usarse DESPUÉS de AuthRequerido.
func SoloAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		rol, _ := c.Get("rol")
		if rol != domain.RolAdministrador {
			c.AbortWithStatusJSON(http.StatusForbidden, domain.ErrorResponse{
				Error: "acceso denegado: se requiere rol administrador",
			})
			return
		}
		c.Next()
	}
}

// ObtenerUsuarioID extrae el ID del usuario del contexto de Gin.
func ObtenerUsuarioID(c *gin.Context) uint {
	id, _ := c.Get("usuario_id")
	return id.(uint)
}

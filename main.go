package main

import (
	"log"
	"os"

	"ticketek/backend/clients"
	"ticketek/backend/controllers"
	"ticketek/backend/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload" // carga el .env automáticamente
)

func main() {
	// Conectar a la base de datos, crear tablas y cargar datos iniciales
	clients.ConnectDB()
	clients.MigrateDB()
	clients.SeedDB()

	// Configurar router
	r := gin.Default()

	// CORS básico para desarrollo
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ── Rutas públicas (sin token) ────────────────────────────
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	// Catálogo público
	r.GET("/eventos", controllers.GetEventos)
	r.GET("/eventos/:id", controllers.GetEventoByID)

	// ── Rutas del cliente (requieren token) ───────────────────
	cliente := r.Group("/entradas")
	cliente.Use(middleware.AuthRequerido())
	{
		cliente.POST("/comprar", controllers.ComprarEntrada)
		cliente.GET("/mis-entradas", controllers.GetMisEntradas)
		cliente.DELETE("/:id", controllers.CancelarEntrada)
		cliente.POST("/:id/transferir", controllers.TransferirEntrada)
	}

	// ── Rutas del admin (requieren token + rol admin) ─────────
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequerido(), middleware.SoloAdmin())
	{
		admin.POST("/eventos", controllers.CrearEvento)
		admin.PUT("/eventos/:id", controllers.ActualizarEvento)
		admin.DELETE("/eventos/:id", controllers.CancelarEvento)
		admin.GET("/eventos/:id/reporte", controllers.GetReporteEvento)
	}

	// Iniciar servidor
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Servidor corriendo en http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}

package clients

import (
	"fmt"
	"log"
	"os"
	"time"

	"ticketek/backend/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	host     := os.Getenv("DB_HOST")
	port     := os.Getenv("DB_PORT")
	user     := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname   := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Error conectando a la base de datos: ", err)
	}

	log.Println("Conexión a MySQL exitosa")
	DB = db
}

// MigrateDB crea/actualiza las tablas según los modelos.
// Usalo solo en desarrollo; en producción usá migraciones manuales.
func MigrateDB() {
	err := DB.AutoMigrate(
		&domain.Usuario{},
		&domain.Evento{},
		&domain.Sector{},
		&domain.Entrada{},
		&domain.Transferencia{},
	)
	if err != nil {
		log.Fatal("Error en AutoMigrate: ", err)
	}
	log.Println("Migración completada")
}

// SeedDB carga datos iniciales si la base está vacía.
func SeedDB() {
	var count int64
	DB.Model(&domain.Usuario{}).Count(&count)
	if count > 0 {
		log.Println("Base de datos ya tiene datos, saltando seed")
		return
	}

	log.Println("Cargando datos iniciales...")

	// Usuarios
	DB.Create(&domain.Usuario{Nombre: "Admin", Email: "admin@smartticket.com", PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Rol: domain.RolAdministrador, Activo: true})
	DB.Create(&domain.Usuario{Nombre: "Juan Perez", Email: "juan@email.com", PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Rol: domain.RolCliente, Activo: true})
	DB.Create(&domain.Usuario{Nombre: "Maria Garcia", Email: "maria@email.com", PasswordHash: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Rol: domain.RolCliente, Activo: true})

	// Eventos
	fecha1, _ := time.Parse("2006-01-02", "2026-08-10")
	fecha2, _ := time.Parse("2006-01-02", "2026-09-05")
	fecha3, _ := time.Parse("2006-01-02", "2026-10-09")
	eventos := []domain.Evento{
		{Titulo: "Tan Bionica - La Ultima Noche Magica Tour", Descripcion: "La banda mas convocante del rock nacional regresa con su gira de despedida.", FotoURL: "https://cdn.efstatic.com/VistasNew/Contenido/tanbioonica.jpeg", Fecha: fecha1, Horario: "21:00:00", DuracionMinutos: 150, Categoria: "Concierto", CapacidadTotal: 5000, Estado: domain.EstadoActivo, CreadoPor: 1},
		{Titulo: "River vs Boca - Superclasico", Descripcion: "El partido mas esperado del ano en el Monumental.", FotoURL: "https://media.elpatagonico.com/p/9b3c2a0383ab1a7187a54ea268c24ac7/adjuntos/193/imagenes/036/334/0036334699/rivjpeg.jpeg", Fecha: fecha2, Horario: "17:00:00", DuracionMinutos: 120, Categoria: "Deporte", CapacidadTotal: 3000, Estado: domain.EstadoActivo, CreadoPor: 1},
		{Titulo: "Comic Con Argentina 2026", Descripcion: "La convencion de cultura pop mas grande del pais.", FotoURL: "https://turismo.buenosaires.gob.ar/sites/turismo/files/field/image/comic-con-arg-1500x610_0.png", Fecha: fecha3, Horario: "10:00:00", DuracionMinutos: 480, Categoria: "Feria", CapacidadTotal: 2000, Estado: domain.EstadoActivo, CreadoPor: 1},
	}
	DB.Create(&eventos)

	// Sectores
	sectores := []domain.Sector{
		{EventoID: eventos[0].ID, Nombre: "Campo", CapacidadMaxima: 2000, CapacidadDisponible: 2000, Precio: 15000},
		{EventoID: eventos[0].ID, Nombre: "Platea Baja", CapacidadMaxima: 1500, CapacidadDisponible: 1500, Precio: 25000},
		{EventoID: eventos[0].ID, Nombre: "Platea Alta", CapacidadMaxima: 1000, CapacidadDisponible: 1000, Precio: 18000},
		{EventoID: eventos[0].ID, Nombre: "VIP", CapacidadMaxima: 500, CapacidadDisponible: 500, Precio: 45000},
		{EventoID: eventos[1].ID, Nombre: "Popular", CapacidadMaxima: 1500, CapacidadDisponible: 1500, Precio: 8000},
		{EventoID: eventos[1].ID, Nombre: "Platea", CapacidadMaxima: 1200, CapacidadDisponible: 1200, Precio: 20000},
		{EventoID: eventos[1].ID, Nombre: "Platea Preferencial", CapacidadMaxima: 300, CapacidadDisponible: 300, Precio: 35000},
		{EventoID: eventos[2].ID, Nombre: "General", CapacidadMaxima: 1500, CapacidadDisponible: 1500, Precio: 5000},
		{EventoID: eventos[2].ID, Nombre: "VIP Pass", CapacidadMaxima: 500, CapacidadDisponible: 500, Precio: 12000},
	}
	DB.Create(&sectores)

	log.Println("Datos iniciales cargados correctamente")
}

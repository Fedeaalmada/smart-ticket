package clients

import (
	"fmt"
	"log"
	"os"

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

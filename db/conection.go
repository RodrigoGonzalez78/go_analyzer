package db

import (
	"fmt"
	"log"

	"github.com/RodrigoGonzalez78/go_analyzer/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB

func StartDB() {

	var err error

	database, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Fallo al conectar a la base de datos SQLite:", err)
	}

	fmt.Println("Conexión a la base de datos SQLite (GORM) establecida.")
}

func MigrateModels() {
	database.AutoMigrate(models.User{})
}

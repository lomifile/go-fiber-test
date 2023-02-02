package database

import (
	"github.com/auth-api/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	conn, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{
    DisableForeignKeyConstraintWhenMigrating: true,
  })

	if err != nil {
		panic("Could not connect to db")
	}

  DB = conn
  conn.AutoMigrate(&models.User{})
}

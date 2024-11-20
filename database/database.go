package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	fmt.Println("Attempting database connection...")
	var err error
	DB, err = gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(2)
	}

	fmt.Println("Connected to the database successfully")
	DB.Logger = logger.Default.LogMode(logger.Info)
}

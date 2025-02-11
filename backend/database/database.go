package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlog "gorm.io/gorm/logger"
)

// Logger for the database
var logger = log.New(os.Stdout, "database ", 0)

// The connection to the database
var DBConn *gorm.DB

func Connect() {
	url := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_DATABASE") + " port=" + os.Getenv("DB_PORT")

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: gormlog.Default.LogMode(gormlog.Warn),
	})

	if err != nil {
		logger.Fatal("Something went wrong during the connection with the database.", err)
	}

	logger.Println("Successfully connected to the database.")

	// Configure the database driver
	driver, _ := db.DB()
	driver.SetMaxIdleConns(10)
	driver.SetMaxOpenConns(100)
	driver.SetConnMaxLifetime(time.Hour)

	// Add the uuid extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		logger.Fatal("uuid extension 'uuid-ossp' not found:", err)
	}

	// Migrate the schema

	// Assign the database to the global variable
	DBConn = db

	// Create the default admin account for the database
	CreateDefaultAccount()
}

// Create the default account
func CreateDefaultAccount() {
	username := os.Getenv("MAGIC_DEFAULT_USERNAME")
	password := os.Getenv("MAGIC_DEFAULT_PASSWORD")

	// Make sure the default user is set
	if username == "" || password == "" {
		logger.Fatal("MAGIC_DEFAULT_(USERNAME|PASSWORD) not set. Can't start the server.")
	}

}

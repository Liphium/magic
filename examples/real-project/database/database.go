package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB = nil

func Connect() {
	if DBConn != nil {
		return
	}

	url := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_DATABASE") + " port=" + os.Getenv("DB_PORT")

	db, err := gorm.Open(postgres.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		log.Fatal("Something went wrong during the connection with the database.", err)
	}

	log.Println("Successfully connected to the database.")

	// Configure the database driver
	driver, _ := db.DB()

	driver.SetMaxIdleConns(10)
	driver.SetMaxOpenConns(100)
	driver.SetConnMaxLifetime(time.Hour)

	// Add the uuid extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error; err != nil {
		fmt.Println("uuid extension 'uuid-ossp' not found.")
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(
		&Post{},
	)

	// Assign the database to the global variable
	DBConn = db
}

type Post struct {
	// The ID of a post. The magic: ignore tag here makes sure this isn't asked for when we use the Post in scripts.
	ID uuid.UUID `json:"id" gorm:"primaryKey,type:uuid;default:uuid_generate_v4()" magic:"ignore"`

	// With the prompt: "" struct tag, we can tell Magic what it should ask for (try running go run . -r create-post c:).
	//
	// Magic additionally supports validation using the validator package. Refer to https://github.com/go-playground/validator for everything possible.
	Author  string `json:"author" prompt:"Author name" validate:"required"`
	Content string `json:"content" prompt:"Content of the post" validate:"required,max=256"`
}

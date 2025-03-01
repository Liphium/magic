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
var logger = log.New(os.Stdout, "database: ", 0)

// The connection to the database
var DBConn *gorm.DB

func Connect() {
	url := "host=" + os.Getenv("MAGIC_DB_HOST") + " user=" + os.Getenv("MAGIC_DB_USERNAME") + " password=" + os.Getenv("MAGIC_DB_PASSWORD") + " dbname=" + os.Getenv("MAGIC_DB_NAME") + " port=" + os.Getenv("MAGIC_DB_PORT")

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
	if err := db.AutoMigrate(
		// All tables related to the accounts
		&Account{},
		&Credential{},
		&Rank{},

		// All tables for Forge and the build pipeline
		&Forge{},
		&Build{},
		&Asset{},
		&Target{},

		// All tables for Preview and stuff
		&Preview{},
		&PreviewSecret{},
		&Environment{},
		&EnvironmentFile{},

		// All tables that Wizard needs to deploy new jobs
		&Wizard{},
		&WizardCreationToken{},
		&Job{},
	); err != nil {
		logger.Fatal("Something went wrong during the migration.", err)
	}

	// Assign the database to the global variable
	DBConn = db

	// Create the default admin account for the database
	createDefaultThings()
}

// Create the default account
func createDefaultThings() {

	// Create default ranks
	if err := DBConn.FirstOrCreate(&Rank{
		ID:              1,
		Label:           "Default",
		PermissionLevel: 0,
	}).Error; err != nil {
		logger.Fatal("couldn't create default ranks:", err)
	}
	if err := DBConn.FirstOrCreate(&Rank{
		ID:              2,
		Label:           "Admin",
		PermissionLevel: 100,
	}).Error; err != nil {
		logger.Fatal("couldn't create default ranks:", err)
	}
}

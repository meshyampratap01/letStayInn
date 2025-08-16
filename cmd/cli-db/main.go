package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/container"
	"github.com/meshyampratap01/letStayInn/internal/db"
	"github.com/meshyampratap01/letStayInn/internal/models"
)

func main() {

	// Connect to NeonDB
	db.ConnectDB()

	// Get the underlying sql.DB to manage connection pool and close on exit
	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get sql.DB from gorm: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("⚠️  Error closing database connection: %v", err)
		}
	}()

	// Migrate all models
	if err := db.DB.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.Booking{},
		&models.Feedback{},
		&models.ServiceRequest{},
	); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	} else {
		fmt.Println("✅ Database migrated successfully!")
	}

	CLIUserHandler := container.InitHandlers()

	for {
		color.Cyan(config.WelcomeMsg)
		color.Cyan(" " + config.AppDescription + " ")
		color.Yellow("1. Signup")
		color.Yellow("2. Login")
		color.Yellow("3. Exit")

		fmt.Print(color.HiWhiteString("Select Option: "))
		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			CLIUserHandler.SignupHandler()
		case 2:
			CLIUserHandler.LoginHandler()
		case 3:
			color.Green("Exiting...")
			os.Exit(0)
		default:
			color.Red(config.InvalidOption)
		}
	}
}

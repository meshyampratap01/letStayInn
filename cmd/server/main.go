package main

import (
	"log"

	"github.com/meshyampratap01/letStayInn/internal/config"
	"github.com/meshyampratap01/letStayInn/internal/db"
	"github.com/meshyampratap01/letStayInn/internal/logger"
	"github.com/meshyampratap01/letStayInn/internal/models"
	"github.com/meshyampratap01/letStayInn/internal/server"
)

func main() {
	logger.InitLogger()
	defer logger.SyncLogger()

	db.ConnectDB()

	sqlDB, err := db.DB.DB()
	if err != nil {
		log.Fatalf(config.ErrFailedToGetSQLDB, err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf(config.ErrDBClose, err)
		}
	}()

	if err := db.DB.AutoMigrate(
		&models.User{},
		&models.Room{},
		&models.Booking{},
		&models.Feedback{},
		&models.ServiceRequest{},
	); err != nil {
		log.Fatalf(config.ErrMigrationFailed, err)
	}
	server.StartServer()
}

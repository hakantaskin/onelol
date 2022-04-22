package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/hakantaskin/onelol/app/domain/validators"
	"github.com/hakantaskin/onelol/app/infrastructure/logger"
	"github.com/hakantaskin/onelol/app/infrastructure/shared_db"
	crud_profile "github.com/hakantaskin/onelol/app/presentation/echo/crud/profile"
	service_profile "github.com/hakantaskin/onelol/app/usecases/profile"
	"github.com/hakantaskin/onelol/app/usecases/transformer"
)

var AppConfig struct {
	Environment string `json:"environment"`
	LogLevel    string `json:"logLevel"`
	AppName     string `json:"appName"`
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig.Environment = os.Getenv("ENVIRONMENT")
	AppConfig.LogLevel = os.Getenv("LOG_LEVEL")
	AppConfig.AppName = os.Getenv("APP_NAME")
}

func main() {
	var err error
	l := setupLoggerEntry()

	dbConfig := shared_db.DatabaseConfig{
		Name: os.Getenv("MYSQL_USER_NAME"),
		Host: os.Getenv("MYSQL_HOST_NAME"),
		Port: os.Getenv("MYSQL_PORT"),
		Pass: os.Getenv("MYSQL_PASS"),
		DB:   os.Getenv("MYSQL_DB"),
	}

	db := shared_db.NewConnector(dbConfig)
	err = shared_db.AutoMigrate(db)
	if err != nil {
		l.Error(err)
	}

	e := echo.New()
	v1 := e.Group("/api/v1")

	// Api && Service && Repository

	profileTransformer := transformer.NewProfileTransformer()
	profileRepository := shared_db.NewProfileRepository(db)
	profileCreatorService := service_profile.NewCreateService(profileRepository, validators.ValidateProfileEntity, profileTransformer)
	profileGetterService := service_profile.NewGetService(profileRepository)
	profileListService := service_profile.NewListService(profileRepository)

	profilePresentation := crud_profile.NewController(profileCreatorService, profileGetterService, profileListService)
	profilePresentation.Init(v1)

	// Api && Service && Repository

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func setupLoggerEntry() *logrus.Entry {
	// initialize the logger instance and logger entry using it
	l := logger.New(os.Stdout, AppConfig.Environment, AppConfig.LogLevel)

	return l.WithFields(logrus.Fields{
		"env":          AppConfig.Environment,
		"program":      AppConfig.AppName,
		"channel":      "http",
		"request_path": "",
		"remote_addr":  "",
		"status_code":  "",
	})
}

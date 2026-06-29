package main

import (
	"log"
	"os"

	"spotsync/handler"
	"spotsync/repository"
	"spotsync/service"
	"spotsync/utils"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Initialize Database
	db := repository.InitDB()

	// Initialize Validator
	v := validator.New()
	customValidator := &utils.CustomValidator{Validator: v}

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health Check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Manual Dependency Injection

	// 1. User & Auth
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo)
	handler.NewAuthHandler(e, authService, customValidator)

	// 2. Parking Zones
	zoneRepo := repository.NewZoneRepository(db)
	zoneService := service.NewZoneService(zoneRepo)
	handler.NewZoneHandler(e, zoneService, customValidator)

	// 3. Reservations
	reservationRepo := repository.NewReservationRepository(db)
	reservationService := service.NewReservationService(reservationRepo)
	handler.NewReservationHandler(e, reservationService, customValidator)

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
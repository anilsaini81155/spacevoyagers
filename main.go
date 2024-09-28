package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anilsaini81155/spacevoyagers/db"
	"github.com/anilsaini81155/spacevoyagers/handlers"
	"github.com/anilsaini81155/spacevoyagers/middleware"
	"github.com/anilsaini81155/spacevoyagers/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// Load environment variables

	// if os.Getenv("GO_ENV") != "TEST" {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file0")
	}
	// } else {
	// 	err := godotenv.Load(".envtest")
	// 	if err != nil {
	// 		log.Fatalf("Error loading .env test file")
	// 	}
	// })
	// Get the application port from the environment variable
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		log.Fatalf("APP_PORT not set in .env file")
	}

	// Initialize the database connection (singleton)
	dbConn, connerr := db.GetDB()
	if connerr != nil {
		log.Fatalf("Error initializing database: %v", connerr)
	}

	models.SetDB(dbConn)

	// Run migrations (create tables)
	models.RunMigrations()

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.LoggingMiddleware) // Use logging middleware
	r.Use(middleware.CORSMiddleware)    // Use CORS middleware

	r.HandleFunc("/exoplanets", handlers.CreateExoplanet).Methods("POST")
	r.HandleFunc("/exoplanets", handlers.ListExoplanets).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handlers.GetExoplanetByID).Methods("GET")
	r.HandleFunc("/exoplanets/{id}", handlers.UpdateExoplanet).Methods("PUT")
	r.HandleFunc("/exoplanets/{id}", handlers.DeleteExoplanet).Methods("DELETE")
	r.HandleFunc("/exoplanets/{id}/fuel", handlers.FuelEstimation).Methods("GET")

	log.Printf("Starting server on port %s...", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, r))
}

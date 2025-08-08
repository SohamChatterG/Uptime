package main

import (
	"log"
	"net/http"

	"github.com/SohamChatterG/uptime/auth"
	"github.com/SohamChatterG/uptime/config"
	"github.com/SohamChatterG/uptime/db"
	"github.com/SohamChatterG/uptime/handler"
	"github.com/SohamChatterG/uptime/middleware"
	"github.com/SohamChatterG/uptime/repository"
	"github.com/SohamChatterG/uptime/router"
	"github.com/SohamChatterG/uptime/service"
	"github.com/SohamChatterG/uptime/worker"
	"github.com/gorilla/handlers" // <-- 1. IMPORT THE NEW PACKAGE
	"github.com/gorilla/mux"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Connect to Database
	mongoClient, err := db.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	database := mongoClient.Database(cfg.DBName)

	// 3. Initialize all layers (Dependency Injection)
	jwtService := auth.NewJWTService(cfg.JWTSecret)

	// Repositories
	userRepo := repository.NewUserRepository(database)
	urlRepo := repository.NewURLRepository(database)
	checkRepo := repository.NewCheckRepository(database)

	// Services
	userService := service.NewUserService(userRepo, jwtService)
	urlService := service.NewURLService(urlRepo, checkRepo)

	// Handlers
	userHandler := handler.NewUserHandler(userService)
	urlHandler := handler.NewURLHandler(urlService)

	// Middleware
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// 4. Setup Router
	mainRouter := mux.NewRouter()
	router.SetupRoutes(mainRouter, userHandler, urlHandler, authMiddleware)

	// 5. Start the Background Worker
	checker := worker.NewChecker(urlRepo, checkRepo, cfg.CheckInterval)
	go checker.Start()
	log.Println("Background uptime checker started.")

	// --- 6. CONFIGURE CORS ---
	// This tells the backend to allow requests from your React frontend.
	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:5173"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	// 7. Start Server with the CORS middleware wrapping the main router
	log.Printf("Server starting on port %s", cfg.Port)
	// The handlers.CORS middleware will add the required headers to each response.
	if err := http.ListenAndServe(":"+cfg.Port, handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(mainRouter)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

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
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	mongoClient, err := db.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	database := mongoClient.Database(cfg.DBName)

	jwtService := auth.NewJWTService(cfg.JWTSecret)

	// Repositories
	userRepo := repository.NewUserRepository(database)
	urlRepo := repository.NewURLRepository(database)
	checkRepo := repository.NewCheckRepository(database)

	// Services
	userService := service.NewUserService(userRepo, jwtService)
	urlService := service.NewURLService(urlRepo, checkRepo)
	notificationService := service.NewGmailService(cfg.EmailUser, cfg.EmailPass)

	userHandler := handler.NewUserHandler(userService)
	urlHandler := handler.NewURLHandler(urlService)
	oauthHandler := handler.NewOAuthHandler(cfg, userService)

	authMiddleware := middleware.AuthMiddleware(jwtService)

	mainRouter := mux.NewRouter()
	router.SetupRoutes(mainRouter, userHandler, urlHandler, oauthHandler, authMiddleware)

	checker := worker.NewChecker(urlRepo, userRepo, checkRepo, notificationService, cfg.CheckInterval)
	go checker.Start()
	log.Println("Background uptime checker started.")

	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:5173"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(mainRouter)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

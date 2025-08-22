package router

import (
	"net/http"

	"github.com/SohamChatterG/uptime/handler"
	"github.com/gorilla/mux"
)

func SetupRoutes(
	r *mux.Router,
	userHandler *handler.UserHandler,
	urlHandler *handler.URLHandler,
	oauthHandler *handler.OAuthHandler,
	authMiddleware func(http.Handler) http.Handler,
) {
	// routes for authentication
	authRouter := r.PathPrefix("/auth").Subrouter()
	authRouter.HandleFunc("/register", userHandler.Register).Methods("POST")
	authRouter.HandleFunc("/login", userHandler.Login).Methods("POST")

	authRouter.HandleFunc("/google/login", oauthHandler.GoogleLogin).Methods("GET")
	authRouter.HandleFunc("/google/callback", oauthHandler.GoogleCallback).Methods("GET")
	authRouter.HandleFunc("/github/login", oauthHandler.GithubLogin).Methods("GET")
	authRouter.HandleFunc("/github/callback", oauthHandler.GithubCallback).Methods("GET")

	// Protected API routes
	apiRouter := r.PathPrefix("/api").Subrouter()
	apiRouter.Use(authMiddleware)

	// URL management routes
	apiRouter.HandleFunc("/urls", urlHandler.CreateURLHandler).Methods("POST")
	apiRouter.HandleFunc("/urls", urlHandler.GetURLsHandler).Methods("GET")
	apiRouter.HandleFunc("/urls/{id}", urlHandler.DeleteURLHandler).Methods("DELETE")

	// Reporting routes
	apiRouter.HandleFunc("/urls/{id}/history", urlHandler.GetURLHistoryHandler).Methods("GET")
}

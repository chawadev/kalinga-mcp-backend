package api

import (
	"net/http"

	"github.com/chawadev/kalinga-backend/internal/auth"
	"github.com/chawadev/kalinga-backend/internal/handlers"
)

type Router struct {
	authHandler *handlers.AuthHandler
}

func NewRouter(authService *auth.Service) *Router {
	return &Router{
		authHandler: handlers.NewAuthHandler(authService),
	}
}

func (r *Router) SetupRoutes(mux *http.ServeMux) {
	// Auth routes with CORS
	mux.HandleFunc("/api/auth/register", corsMiddleware(r.authHandler.Register))
	mux.HandleFunc("/api/auth/login", corsMiddleware(r.authHandler.Login))
	mux.HandleFunc("/api/auth/user", corsMiddleware(r.authHandler.GetUser))
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

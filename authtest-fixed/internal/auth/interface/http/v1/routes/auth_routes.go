package routes

import (

	"github.com/go-chi/chi/v5"

	"github.com/test/authtest-fixed/internal/auth/interface/http/v1/handlers"
)


func RegisterAuthRoutes(router chi.Router, authHandler *handlers.AuthHandler) {
	router.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.Post("/refresh", authHandler.RefreshToken)
	})
}


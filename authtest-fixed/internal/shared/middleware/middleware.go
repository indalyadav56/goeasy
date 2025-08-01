package middleware

import (

	"context"
	"net/http"

	

	"github.com/test/authtest-fixed/internal/auth/application"
	"github.com/test/authtest-fixed/pkg/auth"

)

type AuthMiddleware struct {
	authService *application.AuthService
	jwtSecret   string
}

func NewAuthMiddleware(authService *application.AuthService, jwtSecret string) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		jwtSecret:   jwtSecret,
	}
}


func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		token, err := auth.ExtractTokenFromHeader(authHeader)
		if err != nil {
			http.Error(w, "Authorization token required", http.StatusUnauthorized)
			return
		}

		// Validate token and get user
		user, err := m.authService.ValidateToken(token)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Set user in context
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "user_id", user.ID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user from context
			// user, ok := r.Context().Value("user").(*entity.User)
			// if !ok {
			// 	http.Error(w, "User not found in context", http.StatusUnauthorized)
			// 	return
			// }

			// // Check if user has required role
			// if !user.HasRole(roleName) {
			// 	http.Error(w, "Insufficient permissions", http.StatusForbidden)
			// 	return
			// }

			next.ServeHTTP(w, r)
		}))
	}
}

func (m *AuthMiddleware) RequirePermission(permissionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return m.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// // Get user from context
			// user, ok := r.Context().Value("user").(*entity.User)
			// if !ok {
			// 	http.Error(w, "User not found in context", http.StatusUnauthorized)
			// 	return
			// }

			// // Check if user has required permission
			// if !user.HasPermission(permissionName) {
			// 	http.Error(w, "Insufficient permissions", http.StatusForbidden)
			// 	return
			// }


			next.ServeHTTP(w, r)
		}))
	}
}



func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token, err := auth.ExtractTokenFromHeader(authHeader)
			if err == nil {
				user, err := m.authService.ValidateToken(token)
				if err == nil {
					ctx := context.WithValue(r.Context(), "user", user)
					ctx = context.WithValue(ctx, "user_id", user.ID)
					r = r.WithContext(ctx)
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}


// CORS middleware for handling cross-origin requests

func (m *AuthMiddleware) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}


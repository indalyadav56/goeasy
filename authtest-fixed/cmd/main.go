package main

import (
	"fmt"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/test/authtest-fixed/pkg/db"
	"github.com/test/authtest-fixed/pkg/logger"
	
	product_app "github.com/test/authtest-fixed/internal/product/application"
	product_handlers "github.com/test/authtest-fixed/internal/product/interface/http/v1/handlers"
	product_routes "github.com/test/authtest-fixed/internal/product/interface/http/v1/routes"
	product_postgres "github.com/test/authtest-fixed/internal/product/infrastructure/postgres"
	
	// Auth-related imports
	auth_app "github.com/test/authtest-fixed/internal/auth/application" 
	auth_handlers "github.com/test/authtest-fixed/internal/auth/interface/http/v1/handlers"
	auth_routes "github.com/test/authtest-fixed/internal/auth/interface/http/v1/routes"
	user_app "github.com/test/authtest-fixed/internal/user/application" 
	user_handlers "github.com/test/authtest-fixed/internal/user/interface/http/v1/handlers"
	user_routes "github.com/test/authtest-fixed/internal/user/interface/http/v1/routes"
	user_postgres "github.com/test/authtest-fixed/internal/user/infrastructure/postgres"
	role_app "github.com/test/authtest-fixed/internal/role/application" 
	role_handlers "github.com/test/authtest-fixed/internal/role/interface/http/v1/handlers"
	role_routes "github.com/test/authtest-fixed/internal/role/interface/http/v1/routes"
	role_postgres "github.com/test/authtest-fixed/internal/role/infrastructure/postgres"
	permission_app "github.com/test/authtest-fixed/internal/permission/application" 
	permission_handlers "github.com/test/authtest-fixed/internal/permission/interface/http/v1/handlers"
	permission_routes "github.com/test/authtest-fixed/internal/permission/interface/http/v1/routes"
	permission_postgres "github.com/test/authtest-fixed/internal/permission/infrastructure/postgres"
	auth_middleware "github.com/test/authtest-fixed/internal/shared/middleware"
	
	
)

func main() {
	// init db
	dbConn, err := db.InitDB()
	if err != nil {
		fmt.Println(err)
	}
	defer dbConn.Close()

	// init logger
	logger.InitLogger()

	// init router
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	
	// Initialize product bounded context
	productRepo := product_postgres.NewProductRepository(dbConn)
	productService := product_app.NewProductService(productRepo)
	productHandler := product_handlers.NewProductHandler(productService)

	// Register product routes
	product_routes.SetupProductRoutes(r, productHandler)

	
	// Initialize auth bounded contexts
	// User bounded context
	userRepo := user_postgres.NewUserRepository(dbConn)
	userService := user_app.NewUserService(userRepo)
	userHandler := user_handlers.NewUserHandler(userService)

	// Role bounded context
	roleRepo := role_postgres.NewRoleRepository(dbConn)
	roleService := role_app.NewRoleService(roleRepo)
	roleHandler := role_handlers.NewRoleHandler(roleService)

	// Permission bounded context
	permissionRepo := permission_postgres.NewPermissionRepository(dbConn)
	permissionService := permission_app.NewPermissionService(permissionRepo)
	permissionHandler := permission_handlers.NewPermissionHandler(permissionService)

	// Auth service (depends on user, role, permission services)
	authService := auth_app.NewAuthService(userService, roleService, permissionService)
	authHandler := auth_handlers.NewAuthHandler(authService)

	// Initialize middleware
	authMW := auth_middleware.NewAuthMiddleware(authService, "your-jwt-secret-key")

	// Register auth routes
	auth_routes.RegisterAuthRoutes(r, authHandler)
	user_routes.RegisterUserRoutes(r, userHandler, authMW)
	role_routes.RegisterRoleRoutes(r, roleHandler, authMW)
	permission_routes.RegisterPermissionRoutes(r, permissionHandler, authMW)
	
	

	fmt.Println("Server starting on :8080")
	http.ListenAndServe(":8080", r)
}

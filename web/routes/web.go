package routes

import (
	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	"github.com/MarcelArt/passwordless/web/handlers"
	"github.com/gin-gonic/gin"
)

// SetupWebRoutes initializes dependencies and defines UI routes on the Gin engine
func SetupWebRoutes(r *gin.Engine, authService services.IAuthService) {
	userRepo := repositories.NewUserRepo(configs.DB)
	userService := services.NewUserService(userRepo)

	h := handlers.NewWebAuthHandler(userRepo, authService, userService)

	r.GET("/.well-known/assetlinks.json", handlers.AssetLinks)

	// Page routes
	r.GET("/register", h.ShowRegister)
	r.GET("/login", h.ShowLogin)
	r.GET("/profile", h.ShowProfile)
	r.GET("/logout", h.Logout)

	// Validation routes (for HTMX inline feedback)
	r.POST("/register/validate/username", h.ValidateUsername)
	r.POST("/register/validate/email", h.ValidateEmail)
}

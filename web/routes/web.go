package routes

import (
	"log"

	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	"github.com/MarcelArt/passwordless/web/handlers"
	"github.com/gin-gonic/gin"
)

// SetupWebRoutes initializes dependencies and defines UI routes on the Gin engine
func SetupWebRoutes(r *gin.Engine) {
	userRepo := repositories.NewUserRepo(configs.DB)
	authService, err := services.NewAuthService(userRepo)
	if err != nil {
		log.Fatalf("failed to construct auth service for web: %s", err.Error())
		return
	}

	h := handlers.NewWebAuthHandler(authService, userRepo)

	// Page routes
	r.GET("/register", h.ShowRegister)
	r.GET("/login", h.ShowLogin)

	// Validation routes (for HTMX inline feedback)
	r.POST("/web/register/validate/username", h.ValidateUsername)
	r.POST("/web/register/validate/email", h.ValidateEmail)

	// WebAuthn API endpoint wrappers (for Alpine.js calls)
	r.POST("/web/register/begin", h.RegisterBegin)
	r.POST("/web/register/finish", h.RegisterFinish)
	r.POST("/web/login/begin", h.LoginBegin)
	r.POST("/web/login/finish", h.LoginFinish)
}

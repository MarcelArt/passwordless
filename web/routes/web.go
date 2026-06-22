package routes

import (
	"github.com/MarcelArt/passwordless/web/handlers"
	"github.com/gin-gonic/gin"
)

// SetupWebRoutes initializes dependencies and defines UI routes on the Gin engine
func SetupWebRoutes(r *gin.Engine, waHandler *handlers.WebAuthHandler) {

	r.GET("/.well-known/assetlinks.json", handlers.AssetLinks)

	// Page routes
	r.GET("/register", waHandler.ShowRegister)
	r.GET("/login", waHandler.ShowLogin)
	r.GET("/profile", waHandler.ShowProfile)
	r.GET("/logout", waHandler.Logout)

	// Validation routes (for HTMX inline feedback)
	r.POST("/register/validate/username", waHandler.ValidateUsername)
	r.POST("/register/validate/email", waHandler.ValidateEmail)
}

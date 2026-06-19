package routes

import (
	"github.com/MarcelArt/passwordless/internal/v1/middlewares"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(api *gin.RouterGroup, authService services.IAuthService) {
	authM := middlewares.NewAuthMiddleware()

	v1 := api.Group("/v1")
	setupOauthRoutes(v1)
	setupAuthRoutes(v1, authM, authService)
	setupUserRoutes(v1, authM)
}

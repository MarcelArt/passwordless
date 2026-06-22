package routes

import (
	"github.com/MarcelArt/passwordless/internal/v1/handlers"
	"github.com/MarcelArt/passwordless/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	api *gin.RouterGroup,
	oHandler *handlers.OauthHandler,
	aHandler *handlers.AuthHandler,
	uHandler *handlers.UserHandler,
) {
	authM := middlewares.NewAuthMiddleware()

	v1 := api.Group("/v1")
	setupOauthRoutes(v1, oHandler)
	setupAuthRoutes(v1, authM, aHandler)
	setupUserRoutes(v1, authM, uHandler)
}

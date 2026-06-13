package routes

import (
	"github.com/MarcelArt/passwordless/internal/v1/handlers"
	"github.com/gin-gonic/gin"
)

func setupOauthRoutes(v1 *gin.RouterGroup) {
	h := handlers.NewOauthHandler()
	oauth := v1.Group("/oauth")

	oauth.GET("/authorize", h.Authorize)
	oauth.POST("/token", h.Token)
}

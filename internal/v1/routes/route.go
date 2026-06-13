package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	setupOauthRoutes(v1)
	setupAuthRoutes(v1)
}

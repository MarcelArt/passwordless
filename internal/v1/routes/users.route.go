package routes

import (
	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/handlers"
	"github.com/MarcelArt/passwordless/internal/v1/middlewares"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	"github.com/gin-gonic/gin"
)

func setupUserRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware) {
	r := repositories.NewUserRepo(configs.DB)
	s := services.NewUserService(r)
	h := handlers.NewUserHandler(s)

	users := v1.Group("/users")

	users.POST("/", h.Create)
	users.POST("/login", h.Login)
	users.POST("/refresh", authM.Refresh, h.Refresh)

	users.GET("/", h.Read)
	users.GET("/current", authM.Authn, h.GetCurrent)
	users.GET("/:id", h.GetByID)

	users.PUT("/:id", h.Update)

	users.DELETE("/:id", h.Delete)
}

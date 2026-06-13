package routes

import (
	"log"

	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/handlers"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(v1 *gin.RouterGroup) {
	r := repositories.NewUserRepo(configs.DB)
	s, err := services.NewAuthService(r)
	if err != nil {
		log.Fatalf("failed constructing auth service: %s", err.Error())
		return
	}
	h := handlers.NewAuthHandler(s)

	auth := v1.Group("/auth")
	auth.POST("/register/begin", h.RegisterBegin)
	auth.POST("/register/finish", h.RegisterFinish)
	auth.POST("/login/begin", h.LoginBegin)
	auth.POST("/login/finish", h.LoginFinish)
}

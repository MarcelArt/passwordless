package routes

import (
	"log"

	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/handlers"
	"github.com/MarcelArt/passwordless/internal/v1/middlewares"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware) {
	uRepo := repositories.NewUserRepo(configs.DB)
	service, err := services.NewAuthService(uRepo)
	if err != nil {
		log.Fatalf("failed constructing auth service: %s", err.Error())
		return
	}
	uService := services.NewUserService(uRepo)
	h := handlers.NewAuthHandler(service, uService)

	auth := v1.Group("/auth")
	auth.POST("/register/begin", h.RegisterBegin)
	auth.POST("/register/finish", h.RegisterFinish)
	auth.POST("/login/begin", h.LoginBegin)
	auth.POST("/login/finish", h.LoginFinish)
	auth.POST("/qr/start", h.QrStart)
	auth.POST("/qr/approve/:session_id", authM.Authn, h.QrApprove)

	auth.GET("/qr/poll/:session_id", h.QrPoll)
}

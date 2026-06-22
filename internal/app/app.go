package app

import (
	"fmt"
	"time"

	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/handlers"
	"github.com/MarcelArt/passwordless/internal/v1/middlewares"
	"github.com/MarcelArt/passwordless/internal/v1/routes"
	webHandlers "github.com/MarcelArt/passwordless/web/handlers"
	webroutes "github.com/MarcelArt/passwordless/web/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	uHandler *handlers.UserHandler
	aHandler *handlers.AuthHandler
	oHandler *handlers.OauthHandler

	authM *middlewares.AuthMiddleware

	waHandler *webHandlers.WebAuthHandler
}

func New(
	uHandler *handlers.UserHandler,
	aHandler *handlers.AuthHandler,
	oHandler *handlers.OauthHandler,

	authM *middlewares.AuthMiddleware,

	waHandler *webHandlers.WebAuthHandler,
) *App {
	return &App{
		uHandler: uHandler,
		aHandler: aHandler,
		oHandler: oHandler,

		authM: authM,

		waHandler: waHandler,
	}
}

func (a *App) Run() error {
	if configs.Env.ServerENV == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST, OPTIONS, GET, PUT, PATCH, DELETE"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With", "X-Refresh-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group("/api")
	routes.SetupRoutes(api, a.oHandler, a.aHandler, a.uHandler)

	webroutes.SetupWebRoutes(r, a.waHandler)

	port := fmt.Sprintf(":%s", configs.Env.PORT)
	return r.Run(port)
}

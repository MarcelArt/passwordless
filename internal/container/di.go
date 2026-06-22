package container

import (
	"github.com/MarcelArt/passwordless/internal/app"
	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/handlers"
	"github.com/MarcelArt/passwordless/internal/v1/middlewares"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	webHandlers "github.com/MarcelArt/passwordless/web/handlers"
	"go.uber.org/dig"
)

func New() *dig.Container {
	c := dig.New()

	c.Provide(configs.ConnectDB)

	c.Provide(repositories.NewUserRepo, dig.As(new(repositories.IUserRepo)))

	c.Provide(services.NewAuthService, dig.As(new(services.IAuthService)))
	c.Provide(services.NewUserService, dig.As(new(services.IUserService)))

	c.Provide(middlewares.NewAuthMiddleware)

	c.Provide(handlers.NewUserHandler)
	c.Provide(handlers.NewAuthHandler)
	c.Provide(handlers.NewOauthHandler)

	c.Provide(webHandlers.NewWebAuthHandler)

	c.Provide(app.New)

	return c
}

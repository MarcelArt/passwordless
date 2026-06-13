package configs

import (
	ginserver "github.com/go-oauth2/gin-server"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

func SetupOauth2() {
	manager := manage.NewDefaultManager()
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	clientStore := store.NewClientStore()
	clientStore.Set("client_id_123", &models.Client{
		ID:     "client_id_123",
		Secret: "super_secret_456",
		Domain: "http://localhost:9000", // The allowed third-party redirect URI
	})
	manager.MapClientStorage(clientStore)
	ginserver.InitServer(manager)

	ginserver.SetClientInfoHandler(server.ClientFormHandler)
}

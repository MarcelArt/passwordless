package handlers

import (
	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"
)

type OauthHandler struct {
}

func NewOauthHandler() *OauthHandler {
	return &OauthHandler{}
}

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"2YotnFZFEjr1zCsicMWpAA"`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int64  `json:"expires_in" example:"3600"`
	RefreshToken string `json:"refresh_token,omitempty" example:"tGzv3JOkF0XG5Qx2TlKWIA"`
	Scope        string `json:"scope,omitempty" example:"read"`
}

// OAuth2ErrorResponse represents the OAuth2 error response
type OAuth2ErrorResponse struct {
	Error            string `json:"error" example:"invalid_grant"`
	ErrorDescription string `json:"error_description,omitempty" example:"The provided authorization grant or refresh token is invalid"`
}

// Authorize handles the OAuth2 authorization request
// @Summary      OAuth2 Authorization Endpoint
// @Description  Handles the authorization request, redirects the user to login or consent page.
// @Tags         OAuth
// @Accept       x-www-form-urlencoded
// @Produce      html
// @Param        response_type  query  string  true  "Response type, must be 'code'"
// @Param        client_id      query  string  true  "Client ID"
// @Param        redirect_uri   query  string  true  "Redirect URI"
// @Param        scope          query  string  false "Requested scope"
// @Param        state          query  string  false "State parameter to prevent CSRF"
// @Success      302            {string} string "Redirects to redirect_uri"
// @Failure      400            {object} OAuth2ErrorResponse "Bad Request"
// @Router       /v1/oauth/authorize [get]
func (h *OauthHandler) Authorize(c *gin.Context) {
	ginserver.HandleAuthorizeRequest(c)
}

// Token handles the OAuth2 token request
// @Summary      OAuth2 Token Endpoint
// @Description  Exchanges authorization code, refresh token, client credentials, or password for an access token.
// @Tags         OAuth
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        grant_type     formData  string  true   "Grant type (e.g., authorization_code, client_credentials, password, refresh_token)"
// @Param        client_id      formData  string  false  "Client ID"
// @Param        client_secret  formData  string  false  "Client Secret"
// @Param        code           formData  string  false  "Authorization code (required for authorization_code)"
// @Param        redirect_uri   formData  string  false  "Redirect URI (required for authorization_code)"
// @Param        scope          formData  string  false  "Requested scope"
// @Param        refresh_token  formData  string  false  "Refresh token (required for refresh_token)"
// @Param        username       formData  string  false  "Username (required for password)"
// @Param        password       formData  string  false  "Password (required for password)"
// @Success      200            {object}  TokenResponse "Token response"
// @Failure      400            {object}  OAuth2ErrorResponse "Bad Request"
// @Failure      401            {object}  OAuth2ErrorResponse "Unauthorized"
// @Router       /v1/oauth/token [post]
func (h *OauthHandler) Token(c *gin.Context) {
	ginserver.HandleTokenRequest(c)
}


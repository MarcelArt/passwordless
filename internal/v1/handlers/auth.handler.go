package handlers

import (
	"errors"
	"net/http"

	"github.com/MarcelArt/passwordless/internal/common"
	"github.com/MarcelArt/passwordless/internal/v1/models"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	"github.com/gin-gonic/gin"
	_ "github.com/go-webauthn/webauthn/protocol"
)

type AuthHandler struct {
	service services.IAuthService
}

func NewAuthHandler(service services.IAuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// RegisterBegin handles the initiation of WebAuthn registration
// @Summary      Begin WebAuthn Registration
// @Description  Initiates WebAuthn registration by checking username availability and returning credential creation options.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input  body      models.UserInput  true  "User Registration Input"
// @Success      200    {object}  common.JSONResponse{items=models.BeginRegisterWebAuthn} "Registration started successfully"
// @Failure      400    {object}  common.JSONResponse "Invalid input data"
// @Failure      500    {object}  common.JSONResponse "Internal server error"
// @Router       /v1/auth/register/begin [post]
func (h *AuthHandler) RegisterBegin(c *gin.Context) {
	var input models.UserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, common.NewJSONResponse(err, "invalid json"))
		return
	}

	res, err := h.service.RegisterBegin(c, input)
	if err != nil {
		c.JSON(common.StatusCodeFromError(err), common.NewJSONResponse(err, "failed to register"))
		return
	}

	c.JSON(http.StatusOK, common.NewJSONResponse(res, "registration started"))
}

// RegisterFinish handles the completion of WebAuthn registration
// @Summary      Finish WebAuthn Registration
// @Description  Validates the credential creation response and saves the credential to database.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        session_id  query     string                            true  "Session ID"
// @Param        username    query     string                            true  "Username"
// @Param        email    query     string                            true  "Email"
// @Param        input       body      protocol.CredentialCreationResponse true  "Credential Creation Response"
// @Success      200         {object}  common.JSONResponse               "Registration finished"
// @Failure      400         {object}  common.JSONResponse               "Invalid input data"
// @Failure      500         {object}  common.JSONResponse               "Internal server error"
// @Router       /v1/auth/register/finish [post]
func (h *AuthHandler) RegisterFinish(c *gin.Context) {
	sessionID := c.Query("session_id")
	username := c.Query("username")
	email := c.Query("email")
	user := models.UserInput{
		Username: username,
		Email:    email,
	}

	if sessionID == "" || username == "" {
		c.JSON(http.StatusBadRequest, common.NewJSONResponse(errors.New("missing session_id or username parameters"), "missing session_id or username parameters"))
		return
	}

	if err := h.service.RegisterFinish(c, user, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, common.NewJSONResponse(err, "failed to register"))
		return
	}

	c.JSON(http.StatusOK, common.NewJSONResponse(nil, "Registration finished"))
}

// LoginBegin handles the initiation of WebAuthn login
// @Summary      Begin WebAuthn Login
// @Description  Initiates WebAuthn login by retrieving credential assertion options.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        username  query     string  true  "Username"
// @Success      200       {object}  common.JSONResponse{items=models.BeginLoginWebAuthn} "Login started successfully"
// @Failure      500       {object}  common.JSONResponse "Internal server error"
// @Router       /v1/auth/login/begin [post]
func (h *AuthHandler) LoginBegin(c *gin.Context) {
	username := c.Query("username")

	res, err := h.service.LoginBegin(c, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.NewJSONResponse(err, "failed to start login session"))
		return
	}

	c.JSON(http.StatusOK, common.NewJSONResponse(res, "Login started"))
}

// LoginFinish handles the completion of WebAuthn login
// @Summary      Finish WebAuthn Login
// @Description  Validates the credential assertion response and returns access and refresh tokens.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        session_id  query     string                              true  "Session ID"
// @Param        username    query     string                              true  "Username"
// @Param        input       body      protocol.CredentialAssertionResponse true  "Credential Assertion Response"
// @Success      200         {object}  common.JSONResponse{items=models.LoginResponse} "Login finished successfully"
// @Failure      400         {object}  common.JSONResponse                 "Invalid parameters"
// @Failure      401         {object}  common.JSONResponse                 "Unauthorized"
// @Router       /v1/auth/login/finish [post]
func (h *AuthHandler) LoginFinish(c *gin.Context) {
	sessionID := c.Query("session_id")
	username := c.Query("username")

	if sessionID == "" || username == "" {
		c.JSON(http.StatusBadRequest, common.NewJSONResponse(errors.New("missing session_id or username parameters"), "missing session_id or username parameters"))
		return
	}

	res, err := h.service.LoginFinish(c, username, sessionID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.NewJSONResponse(err, "failed to login"))
		return
	}

	c.JSON(http.StatusOK, common.NewJSONResponse(res, "Login finished"))
}

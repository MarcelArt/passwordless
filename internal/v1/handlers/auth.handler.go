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
	service  services.IAuthService
	uService services.IUserService
}

func NewAuthHandler(service services.IAuthService, uService services.IUserService) *AuthHandler {
	return &AuthHandler{
		service:  service,
		uService: uService,
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
		_, res := common.ResultErr(err, "invalid json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res, err := h.service.RegisterBegin(c.Request.Context(), input)
	if err != nil {
		c.JSON(common.ResultErr(err, "invalid json"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(res, "registration started"))
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
		_, res := common.ResultErr(errors.New("missing session_id or username parameters"), "")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := h.service.RegisterFinish(c, user, sessionID); err != nil {
		c.JSON(common.ResultErr(err, "failed to register"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Registration finished"))
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

	res, err := h.service.LoginBegin(c.Request.Context(), username)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed to start login session"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(res, "Login started"))
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
		_, res := common.ResultErr(errors.New("missing session_id or username parameters"), "")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	res, err := h.service.LoginFinish(c, username, sessionID)
	if err != nil {
		_, res := common.ResultErr(err, "failed to login")
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(res, "Login finished"))
}

// QrStart handles the initiation of QR code authentication
// @Summary      Start QR Authentication Session
// @Description  Initiates a QR code authentication session and returns the QR code image.
// @Tags         Auth
// @Produce      image/png
// @Success      200  {file}    string  "QR Code image in PNG format"
// @Failure      500  {object}  common.Result[string]  "Internal server error"
// @Router       /v1/auth/qr/start [post]
func (h *AuthHandler) QrStart(c *gin.Context) {
	_, png, err := h.service.QrStart()
	if err != nil {
		c.JSON(common.ResultErr(err, "failed generating qr"))
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}

// QrPoll handles checking the status of QR code authentication
// @Summary      Poll QR Authentication Status
// @Description  Checks the authentication status of the QR session using the session ID. Returns token pair if logged in.
// @Tags         Auth
// @Produce      json
// @Param        session_id  path      string  true  "Session ID"
// @Success      200         {object}  common.Result[models.LoginResponse] "Logged in or not logged in yet (item is nil when not logged in)"
// @Failure      401         {object}  common.Result[string] "Unauthorized"
// @Router       /v1/auth/qr/poll/{session_id} [get]
func (h *AuthHandler) QrPoll(c *gin.Context) {
	sessionID := c.Param("session_id")

	userID := h.service.QrCheck(sessionID)

	if userID == 0 {
		c.JSON(http.StatusOK, common.ResultOk[any](nil, "not logged in yet"))
		return
	}

	user, err := h.uService.RegenerateTokenPair(c, userID, true)
	if err != nil {
		_, res := common.ResultErr(err, "")
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(user, "Logged in"))
}

// QrApprove handles approving a QR code authentication session
// @Summary      Approve QR Authentication
// @Description  Approves a QR code authentication session by setting the associated user ID.
// @Tags         Auth
// @Produce      json
// @Param        session_id  path      string  true  "Session ID"
// @Success      200         {object}  common.Result[bool]      "QR authentication approved"
// @Failure      401         {object}  common.Result[any]       "Invalid token or session invalid"
// @Security     ApiKeyAuth
// @Router       /v1/auth/qr/approve/{session_id} [post]
func (h *AuthHandler) QrApprove(c *gin.Context) {
	sessionID := c.Param("session_id")

	userID, err := common.MustGet[float64](c, "userId")
	if err != nil {
		_, res := common.ResultErr(err, "invalid token")
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	if ok := h.service.QrApprove(sessionID, uint(userID)); !ok {
		c.JSON(http.StatusUnauthorized, common.ResultOk(false, "qr invalid"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(true, "qr approved"))
}

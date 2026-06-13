package handlers

import (
	"errors"
	"net/http"

	"github.com/MarcelArt/passwordless/internal/v1/models"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	"github.com/MarcelArt/passwordless/web/viewmodels"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebAuthHandler struct {
	authService services.IAuthService
	userRepo    repositories.IUserRepo
}

func NewWebAuthHandler(authService services.IAuthService, userRepo repositories.IUserRepo) *WebAuthHandler {
	return &WebAuthHandler{
		authService: authService,
		userRepo:    userRepo,
	}
}

// ShowRegister renders the WebAuthn registration page template wrapped inside the main layout
func (h *WebAuthHandler) ShowRegister(c *gin.Context) {
	Render(c, http.StatusOK, "register.html", gin.H{
		"Title": "Register",
	})
}

// ValidateUsername checks if a username is available in the database and renders inline validation UI
func (h *WebAuthHandler) ValidateUsername(c *gin.Context) {
	username := c.PostForm("username")
	if username == "" {
		RenderPartial(c, http.StatusOK, "validation.html", gin.H{
			"Target":  "username",
			"IsValid": false,
			"Message": "Username is required",
		})
		return
	}

	_, err := h.userRepo.GetByUsernameOrEmail(c.Request.Context(), username)
	if err == nil {
		RenderPartial(c, http.StatusOK, "validation.html", gin.H{
			"Target":  "username",
			"IsValid": false,
			"Message": "Username is already taken",
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		RenderPartial(c, http.StatusOK, "validation.html", gin.H{
			"Target":  "username",
			"IsValid": true,
			"Message": "Username is available",
		})
		return
	}

	RenderPartial(c, http.StatusOK, "validation.html", gin.H{
		"Target":  "username",
		"IsValid": false,
		"Message": "Error checking username: " + err.Error(),
	})
}

// ValidateEmail checks if an email is available in the database and renders inline validation UI
func (h *WebAuthHandler) ValidateEmail(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		RenderPartial(c, http.StatusOK, "validation.html", gin.H{
			"Target":  "email",
			"IsValid": false,
			"Message": "Email is required",
		})
		return
	}

	_, err := h.userRepo.GetByUsernameOrEmail(c.Request.Context(), email)
	if err == nil {
		RenderPartial(c, http.StatusOK, "validation.html", gin.H{
			"Target":  "email",
			"IsValid": false,
			"Message": "Email is already in use",
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		RenderPartial(c, http.StatusOK, "validation.html", gin.H{
			"Target":  "email",
			"IsValid": true,
			"Message": "Email is available",
		})
		return
	}

	RenderPartial(c, http.StatusOK, "validation.html", gin.H{
		"Target":  "email",
		"IsValid": false,
		"Message": "Error checking email: " + err.Error(),
	})
}

// RegisterBegin handles initiating registration: inserts user skeleton and returns credential options
func (h *WebAuthHandler) RegisterBegin(c *gin.Context) {
	var form viewmodels.RegisterForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid registration data: " + err.Error()})
		return
	}

	input := models.UserInput{
		Username: form.Username,
		Email:    form.Email,
		Password: "", // Passwordless registration: no initial password
	}

	res, err := h.authService.RegisterBegin(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to initiate registration: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// RegisterFinish finishes the WebAuthn flow by validating and storing the device credential in the database
func (h *WebAuthHandler) RegisterFinish(c *gin.Context) {
	sessionID := c.Query("session_id")
	username := c.Query("username")

	if sessionID == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing session_id or username parameters"})
		return
	}

	// Pass context and parameters to finish registration.
	// Since we didn't read the Request Body here, it will be fully read inside authService.RegisterFinish.
	if err := h.authService.RegisterFinish(c, username, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to complete registration: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "WebAuthn Registration complete"})
}

// ShowLogin renders the login page template
func (h *WebAuthHandler) ShowLogin(c *gin.Context) {
	Render(c, http.StatusOK, "login.html", gin.H{
		"Title": "Sign In",
	})
}

// LoginBegin initiates passwordless WebAuthn login assertion
func (h *WebAuthHandler) LoginBegin(c *gin.Context) {
	var form viewmodels.LoginForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid login data: " + err.Error()})
		return
	}

	res, err := h.authService.LoginBegin(c, form.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to initiate login: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// LoginFinish completes passwordless WebAuthn login assertion and returns user details
func (h *WebAuthHandler) LoginFinish(c *gin.Context) {
	sessionID := c.Query("session_id")
	username := c.Query("username")

	if sessionID == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Missing session_id or username parameters"})
		return
	}

	// Forward context directly so the request body remains intact.
	res, err := h.authService.LoginFinish(c, username, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to authenticate: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

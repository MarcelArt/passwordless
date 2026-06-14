package handlers

import (
	"errors"
	"net/http"

	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WebAuthHandler struct {
	userRepo    repositories.IUserRepo
}

func NewWebAuthHandler(userRepo repositories.IUserRepo) *WebAuthHandler {
	return &WebAuthHandler{
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

// ShowLogin renders the login page template
func (h *WebAuthHandler) ShowLogin(c *gin.Context) {
	Render(c, http.StatusOK, "login.html", gin.H{
		"Title": "Sign In",
	})
}

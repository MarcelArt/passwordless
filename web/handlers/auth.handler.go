package handlers

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/internal/v1/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type WebAuthHandler struct {
	userRepo    repositories.IUserRepo
	authService services.IAuthService
	userService services.IUserService
}

func NewWebAuthHandler(userRepo repositories.IUserRepo, authService services.IAuthService, userService services.IUserService) *WebAuthHandler {
	return &WebAuthHandler{
		userRepo:    userRepo,
		authService: authService,
		userService: userService,
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
	sessionID, png, err := h.authService.QrStart()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to start QR login session: %s", err.Error())
		return
	}

	qrBase64 := base64.StdEncoding.EncodeToString(png)

	Render(c, http.StatusOK, "login.html", gin.H{
		"Title":     "Sign In",
		"SessionID": sessionID,
		"QrCode":    qrBase64,
	})
}

// ShowProfile renders the authenticated user profile page template
func (h *WebAuthHandler) ShowProfile(c *gin.Context) {
	tokenString, err := c.Cookie("at")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(configs.Env.JwtSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	userIDVal, ok := claims["userId"]
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	userIDFloat, ok := userIDVal.(float64)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), uint(userIDFloat))
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	Render(c, http.StatusOK, "profile.html", gin.H{
		"Title": "Profile",
		"User":  user,
	})
}

// Logout clears the cookies and redirects to /login
func (h *WebAuthHandler) Logout(c *gin.Context) {
	isProd := configs.Env.ServerENV == "prod"
	c.SetCookie("at", "", -1, "/", "", isProd, true)
	c.SetCookie("rt", "", -1, "/", "", isProd, true)
	c.Redirect(http.StatusSeeOther, "/login")
}

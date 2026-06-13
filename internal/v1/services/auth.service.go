package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/v1/models"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/pkg/jsonb"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type IAuthService interface {
	RegisterBegin(c context.Context, input models.UserInput) (models.BeginRegisterWebAuthn, error)
	RegisterFinish(c *gin.Context, username string, sessionID string) error
	LoginBegin(c *gin.Context, username string) (models.BeginLoginWebAuthn, error)
	LoginFinish(c *gin.Context, username string, sessionID string) (models.LoginResponse, error)
}

type AuthService struct {
	uRepo           repositories.IUserRepo
	webAuthnHandler *webauthn.WebAuthn
	sessionDB       map[string]*webauthn.SessionData
}

func NewAuthService(uRepo repositories.IUserRepo) (*AuthService, error) {
	webAuthnHandler, err := webauthn.New(&webauthn.Config{
		RPDisplayName: configs.Env.RPDisplayName,
		RPID:          configs.Env.RPID, // Must match your domain
		RPOrigins:     []string{configs.Env.RPOrigins},
	})
	if err != nil {
		return nil, err
	}

	return &AuthService{
		uRepo:           uRepo,
		webAuthnHandler: webAuthnHandler,
		sessionDB:       make(map[string]*webauthn.SessionData),
	}, nil
}

func (s *AuthService) RegisterBegin(c context.Context, input models.UserInput) (models.BeginRegisterWebAuthn, error) {
	var res models.BeginRegisterWebAuthn
	user, err := s.uRepo.GetByUsernameOrEmail(c, input.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return res, fmt.Errorf("unexpected error: %w", err)
	}
	if user.ID == 0 {
		id, err := s.uRepo.Create(c, input)
		if err != nil {
			return res, fmt.Errorf("error creating user: %w", err)
		}

		user.ID = id
		user.Email = input.Email
		user.Username = input.Username
		user.Password = input.Password
	}

	options, sessionData, err := s.webAuthnHandler.BeginRegistration(&user)
	sessionID := base64.StdEncoding.EncodeToString([]byte(sessionData.Challenge))

	s.sessionDB[sessionID] = sessionData

	res.Options = options
	res.SessionID = sessionID

	return res, nil
}

func (s *AuthService) RegisterFinish(c *gin.Context, username string, sessionID string) error {
	sessionData := s.sessionDB[sessionID]
	if sessionData == nil {
		return errors.New("session not found")
	}

	user, err := s.uRepo.GetByUsernameOrEmail(c, username)
	if err != nil {
		return err
	}

	credential, err := s.webAuthnHandler.FinishRegistration(&user, *sessionData, c.Request)
	if err != nil {
		return fmt.Errorf("failed parsing credential: %w", err)
	}

	credentials := user.WebAuthnCredentials()
	credentials = append(credentials, *credential)
	jCredentials, err := jsonb.New(credentials)
	if err != nil {
		return fmt.Errorf("failed serializing credentials: %w", err)
	}
	userUpdate := models.UserInput{
		Credentials: jCredentials,
	}
	if err := s.uRepo.Update(c, user.ID, userUpdate); err != nil {
		return fmt.Errorf("failed updating user: %w", err)
	}
	delete(s.sessionDB, sessionID)

	return nil
}

func (s *AuthService) LoginBegin(c *gin.Context, username string) (models.BeginLoginWebAuthn, error) {
	var res models.BeginLoginWebAuthn
	user, err := s.uRepo.GetByUsernameOrEmail(c, username)
	if err != nil {
		return res, err
	}

	options, sessionData, err := s.webAuthnHandler.BeginLogin(&user)
	if err != nil {
		return res, fmt.Errorf("failed starting login session: %w", err)
	}

	sessionID := base64.StdEncoding.EncodeToString([]byte(sessionData.Challenge))
	s.sessionDB[sessionID] = sessionData

	res.Options = options
	res.SessionID = sessionID

	return res, nil
}

func (s *AuthService) LoginFinish(c *gin.Context, username string, sessionID string) (models.LoginResponse, error) {
	var res models.LoginResponse

	sessionData := s.sessionDB[sessionID]
	if sessionData == nil {
		return res, errors.New("session not found")
	}

	user, err := s.uRepo.GetByUsernameOrEmail(c, username)
	if err != nil {
		return res, err
	}

	if _, err := s.webAuthnHandler.FinishLogin(&user, *sessionData, c.Request); err != nil {
		return res, fmt.Errorf("failed parsing credential: %w", err)
	}

	res.User = user
	delete(s.sessionDB, sessionID)

	return res, nil
}

package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MarcelArt/passwordless/internal/common"
	"github.com/MarcelArt/passwordless/internal/configs"
	"github.com/MarcelArt/passwordless/internal/enums"
	"github.com/MarcelArt/passwordless/internal/v1/models"
	"github.com/MarcelArt/passwordless/internal/v1/repositories"
	"github.com/MarcelArt/passwordless/pkg/jsonb"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

type IAuthService interface {
	RegisterBegin(c context.Context, input models.UserInput) (models.BeginRegisterWebAuthn, error)
	RegisterFinish(c *gin.Context, user models.UserInput, sessionID string) error
	LoginBegin(c context.Context, username string) (models.BeginLoginWebAuthn, error)
	LoginFinish(c *gin.Context, username string, sessionID string) (models.LoginResponse, error)
	QrStart() (string, []byte, error)
	QrCheck(sessionID string) uint
	QrApprove(sessionID string, userID uint) bool
}

type AuthService struct {
	uRepo           repositories.IUserRepo
	webAuthnHandler *webauthn.WebAuthn
	sessionDB       map[string]*webauthn.SessionData
	qrSession       map[string]*models.QrSession
}

func NewAuthService(uRepo repositories.IUserRepo) (*AuthService, error) {
	rpAPK, err := common.GetAndroidWebAuthnOrigin(enums.RPOriginAPK)
	if err != nil {
		return nil, err
	}

	webAuthnHandler, err := webauthn.New(&webauthn.Config{
		RPDisplayName: configs.Env.RPDisplayName,
		RPID:          configs.Env.RPID, // Must match your domain
		RPOrigins: []string{
			configs.Env.RPOrigins,
			rpAPK,
		},
	})
	if err != nil {
		return nil, err
	}

	return &AuthService{
		uRepo:           uRepo,
		webAuthnHandler: webAuthnHandler,
		sessionDB:       make(map[string]*webauthn.SessionData),
		qrSession:       make(map[string]*models.QrSession),
	}, nil
}

func (s *AuthService) RegisterBegin(c context.Context, input models.UserInput) (models.BeginRegisterWebAuthn, error) {
	var res models.BeginRegisterWebAuthn
	user, err := s.uRepo.GetByUsernameOrEmail(c, input.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return res, fmt.Errorf("unexpected error: %w", err)
	}
	if user.ID != 0 {
		return res, enums.ErrAlreadyRegsitered
	}
	user.Email = input.Email
	user.Username = input.Username
	user.Password = input.Password

	options, sessionData, err := s.webAuthnHandler.BeginRegistration(&user)
	if err != nil {
		return res, fmt.Errorf("error creating session: %w", err)
	}
	sessionID := base64.StdEncoding.EncodeToString([]byte(sessionData.Challenge))

	s.sessionDB[sessionID] = sessionData

	res.Options = options
	res.SessionID = sessionID

	return res, nil
}

func (s *AuthService) RegisterFinish(c *gin.Context, user models.UserInput, sessionID string) error {
	sessionData := s.sessionDB[sessionID]
	if sessionData == nil {
		return errors.New("session not found")
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
	user.Credentials = jCredentials
	if _, err := s.uRepo.Create(c, user); err != nil {
		return fmt.Errorf("failed updating user: %w", err)
	}
	delete(s.sessionDB, sessionID)

	return nil
}

func (s *AuthService) LoginBegin(c context.Context, username string) (models.BeginLoginWebAuthn, error) {
	var res models.BeginLoginWebAuthn
	user, err := s.uRepo.GetByUsernameOrEmail(c, username)
	if err != nil {
		return res, err
	}

	options, sessionData, err := s.webAuthnHandler.BeginLogin(&user)
	if err != nil {
		return res, fmt.Errorf("failed starting login session: %w", err)
	}
	options.Response.AllowedCredentials = nil

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

func (s *AuthService) QrStart() (string, []byte, error) {
	today := time.Now()
	id, err := uuid.NewV7()
	if err != nil {
		return "", nil, fmt.Errorf("failed generating uuid: %w", err)
	}

	sID := id.String()
	exp := today.Add(5 * time.Minute)
	s.qrSession[sID] = &models.QrSession{Exp: exp}

	png, err := qrcode.Encode(sID, qrcode.Medium, 256)

	log.Println("sID :>> ", sID)
	return sID, png, err
}

func (s *AuthService) QrApprove(sessionID string, userID uint) bool {
	session, ok := s.qrSession[sessionID]
	if !ok {
		return false
	}

	today := time.Now()
	isApproved := today.Before(session.Exp)
	if isApproved {
		session.UserID = userID
	}

	return isApproved
}

func (s *AuthService) QrCheck(sessionID string) uint {
	session, ok := s.qrSession[sessionID]
	if !ok {
		return 0
	}

	if time.Now().After(session.Exp) {
		delete(s.qrSession, sessionID)
		return 0
	}
	userID := session.UserID

	// delete(s.qrSession, sessionID)
	return userID
}

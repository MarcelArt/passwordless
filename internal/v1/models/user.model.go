package models

import (
	"fmt"

	"github.com/MarcelArt/passwordless/internal/common"
	"github.com/MarcelArt/passwordless/internal/entities"
	"github.com/MarcelArt/passwordless/pkg/jsonb"
	"github.com/go-webauthn/webauthn/webauthn"
)

// UserInput
type UserInput struct {
	common.InputModel
	Username    string                             `gorm:"not null;unique" json:"username"`
	Email       string                             `gorm:"not null;unique" json:"email"`
	Password    string                             `gorm:"not null" json:"password"`
	Credentials jsonb.JSONB[[]webauthn.Credential] `json:"-"`
}

func (m *UserInput) WebAuthnID() []byte {
	return []byte(m.Username)
}

func (m *UserInput) WebAuthnName() string {
	return m.Username
}

func (m *UserInput) WebAuthnDisplayName() string {
	return m.Username
}

func (m *UserInput) WebAuthnIcon() string {
	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s", m.Username)
}

func (m *UserInput) WebAuthnCredentials() []webauthn.Credential {
	credentials, _ := m.Credentials.Deserialize()
	return credentials
}

var _ webauthn.User = &UserInput{}

// UserInput end

type UserPage struct {
	ID       uint                  `json:"ID"`
	Username string                `gorm:"not null;unique" json:"username"`
	Email    string                `gorm:"not null;unique" json:"email"`
	Roles    jsonb.JSONB[[]string] `json:"roles"`
}

type LoginInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsRemember bool   `json:"isRemember"`
}

type LoginResponse struct {
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	User         entities.User `json:"user"`
}

type UserRole struct {
	ID          uint   `json:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

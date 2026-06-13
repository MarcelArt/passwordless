package entities

import (
	"fmt"
	"strconv"

	"github.com/MarcelArt/passwordless/pkg/jsonb"
	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string                             `gorm:"not null;unique" json:"username"`
	Email       string                             `gorm:"not null;unique" json:"email"`
	Password    string                             `json:"password"`
	Credentials jsonb.JSONB[[]webauthn.Credential] `json:"credentials"`
}

func (m *User) WebAuthnID() []byte {
	return []byte(strconv.Itoa(int(m.ID)))
}

func (m *User) WebAuthnName() string {
	return m.Username
}

func (m *User) WebAuthnDisplayName() string {
	return m.Username
}

func (m *User) WebAuthnIcon() string {
	return fmt.Sprintf("https://ui-avatars.com/api/?name=%s", m.Username)
}

func (m *User) WebAuthnCredentials() []webauthn.Credential {
	credentials, _ := m.Credentials.Deserialize()
	return credentials
}

var _ webauthn.User = &User{}

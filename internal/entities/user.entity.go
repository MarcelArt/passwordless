package entities

import (
	"encoding/binary"
	"fmt"

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
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(m.ID)) // Converts ID 3 into a stable 8-byte sequence
	return buf
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

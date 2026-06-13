package models

import "github.com/go-webauthn/webauthn/protocol"

type BeginRegisterWebAuthn struct {
	Options   *protocol.CredentialCreation `json:"options"`
	SessionID string                       `json:"sessionId"`
}

type BeginLoginWebAuthn struct {
	Options   *protocol.CredentialAssertion `json:"options"`
	SessionID string                        `json:"sessionId"`
}

package common

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

func GetAndroidWebAuthnOrigin(sha256Hex string) (string, error) {
	// 1. Clean the colons out of your keytool SHA-256 output
	cleanHex := strings.ReplaceAll(sha256Hex, ":", "")

	// 2. Decode the hex string into raw bytes
	bytes, err := hex.DecodeString(cleanHex)
	if err != nil {
		return "", err
	}

	// 3. Encode to Base64URL string strictly WITHOUT padding (=)
	base64UrlHash := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes)

	// 4. Wrap it in the official Android protocol format
	return fmt.Sprintf("android:apk-key-hash:%s", base64UrlHash), nil
}

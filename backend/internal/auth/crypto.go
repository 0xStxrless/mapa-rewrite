package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/scrypt"
)

const (
	scryptKeyLen = 64
	scryptN      = 32768
	scryptR      = 8
	scryptP      = 1
)

type SessionPayload struct {
	UserID    int32  `json:"userId"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"expiresAt"`
}

func GetSessionSecret() (string, error) {
	s := os.Getenv("SESSION_SECRET")
	if s == "" {
		if os.Getenv("APP_ENV") == "production" {
			return "", errors.New("SESSION_SECRET env var must be set in production")
		}
		return "dev-only-insecure-secret-change-me", nil
	}
	return s, nil
}

// HashPassword hashes a password with scrypt + random salt.
// Output format: "scrypt:<salt_hex>:<hash_hex>"
func HashPassword(password string) (string, error) {
	saltBytes := make([]byte, 16)
	if _, err := rand.Read(saltBytes); err != nil {
		return "", fmt.Errorf("generating salt: %w", err)
	}
	salt := hex.EncodeToString(saltBytes)

	derived, err := scrypt.Key([]byte(password), []byte(salt), scryptN, scryptR, scryptP, scryptKeyLen)
	if err != nil {
		return "", fmt.Errorf("scrypt: %w", err)
	}

	return fmt.Sprintf("scrypt:%s:%s", salt, hex.EncodeToString(derived)), nil
}

func VerifyPassword(password, stored string) bool {
	if strings.HasPrefix(stored, "scrypt:") {
		parts := strings.Split(stored, ":")
		if len(parts) != 3 {
			return false
		}
		salt := parts[1]
		hashHex := parts[2]

		derived, err := scrypt.Key([]byte(password), []byte(salt), scryptN, scryptR, scryptP, scryptKeyLen)
		if err != nil {
			return false
		}

		storedBytes, err := hex.DecodeString(hashHex)
		if err != nil || len(storedBytes) != scryptKeyLen {
			return false
		}

		return subtle.ConstantTimeCompare(derived, storedBytes) == 1
	}

	sum := sha256.Sum256([]byte(password))
	sha := hex.EncodeToString(sum[:])

	storedBytes, err := hex.DecodeString(stored)
	if err != nil || len(storedBytes) != 32 {
		return false
	}
	shaBytes, err := hex.DecodeString(sha)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(shaBytes, storedBytes) == 1
}

// SignSession produces: base64url(json).hmac-sha256
func SignSession(payload SessionPayload, secret string) (string, error) {
	if payload.ExpiresAt == 0 {
		payload.ExpiresAt = time.Now().Add(7 * 24 * time.Hour).Unix()
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshalling payload: %w", err)
	}

	data := base64.RawURLEncoding.EncodeToString(raw)
	sig := hmacHex(data, secret)

	return data + "." + sig, nil
}

// VerifyAndParseSession verifies the HMAC and returns the parsed payload.
func VerifyAndParseSession(cookie, secret string) (*SessionPayload, error) {
	dot := strings.LastIndex(cookie, ".")
	if dot == -1 {
		return nil, errors.New("malformed cookie: no signature")
	}

	data := cookie[:dot]
	sig := cookie[dot+1:]
	expected := hmacHex(data, secret)

	// timing-safe comparison both sides must be equal-length hex strings
	sigBytes, err := hex.DecodeString(sig)
	if err != nil {
		return nil, errors.New("malformed signature encoding")
	}
	expectedBytes, err := hex.DecodeString(expected)
	if err != nil {
		return nil, errors.New("internal hmac encoding error")
	}
	if subtle.ConstantTimeCompare(sigBytes, expectedBytes) == 0 {
		return nil, errors.New("signature mismatch")
	}

	raw, err := base64.RawURLEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}

	var p SessionPayload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("json decode: %w", err)
	}

	if p.ExpiresAt != 0 && time.Now().Unix() > p.ExpiresAt {
		return nil, errors.New("token expired")
	}

	return &p, nil
}

func hmacHex(data, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

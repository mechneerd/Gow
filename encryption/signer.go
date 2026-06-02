package encryption

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
)

// Signer provides message signing and verification.
type Signer struct {
	key []byte
}

// NewSigner creates a new Signer with the given secret key.
func NewSigner(secret string) (*Signer, error) {
	if secret == "" {
		return nil, errors.New("signing key cannot be empty")
	}
	hash := sha256.Sum256([]byte(secret))
	return &Signer{key: hash[:]}, nil
}

// Sign signs a message using HMAC-SHA256.
func (s *Signer) Sign(message string) string {
	mac := hmac.New(sha256.New, s.key)
	mac.Write([]byte(message))
	signature := mac.Sum(nil)
	return base64.URLEncoding.EncodeToString(signature)
}

// Verify verifies a message signature using HMAC-SHA256.
func (s *Signer) Verify(message, signature string) bool {
	expected := s.Sign(message)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// SignHMAC signs a message using HMAC-SHA512.
func (s *Signer) SignHMAC(message string) string {
	mac := hmac.New(sha512.New, s.key)
	mac.Write([]byte(message))
	signature := mac.Sum(nil)
	return base64.URLEncoding.EncodeToString(signature)
}

// VerifyHMAC verifies a message signature using HMAC-SHA512.
func (s *Signer) VerifyHMAC(message, signature string) bool {
	expected := s.SignHMAC(message)
	return hmac.Equal([]byte(expected), []byte(signature))
}

// SignedMessage represents a signed message with timestamp.
type SignedMessage struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
}

// SignWithTimestamp signs a message with a timestamp.
func (s *Signer) SignWithTimestamp(message string, timestamp int64) *SignedMessage {
	data := message + "|" + string(rune(timestamp))
	return &SignedMessage{
		Payload:   message,
		Signature: s.Sign(data),
		Timestamp: timestamp,
	}
}

// VerifySignedMessage verifies a signed message with timestamp.
func (s *Signer) VerifySignedMessage(sm *SignedMessage, maxAge int64) bool {
	if sm == nil {
		return false
	}
	data := sm.Payload + "|" + string(rune(sm.Timestamp))
	return s.Verify(data, sm.Signature)
}

// SignedURL represents a URL with a signature.
type SignedURL struct {
	URL       string
	Signature string
}

// SignURL signs a URL path.
func (s *Signer) SignURL(url string) *SignedURL {
	return &SignedURL{
		URL:       url,
		Signature: s.Sign(url),
	}
}

// VerifyURL verifies a signed URL.
func (s *Signer) VerifyURL(url, signature string) bool {
	return s.Verify(url, signature)
}

// Hash creates a SHA-256 hash of the message (one-way).
func (s *Signer) Hash(message string) string {
	h := sha256.Sum256([]byte(message))
	return base64.URLEncoding.EncodeToString(h[:])
}

// HashHMAC creates an HMAC-SHA256 hash (keyed, one-way).
func (s *Signer) HashHMAC(message string) string {
	mac := hmac.New(sha256.New, s.key)
	mac.Write([]byte(message))
	return base64.URLEncoding.EncodeToString(mac.Sum(nil))
}

// ConstantTimeEqual performs constant-time comparison of two byte slices.
func ConstantTimeEqual(a, b []byte) bool {
	return hmac.Equal(a, b)
}

// GenerateRandomKey generates a cryptographically secure random key.
func GenerateRandomKey(length int) string {
	key := make([]byte, length)
	_, _ = rand.Read(key)
	return base64.URLEncoding.EncodeToString(key)
}

// GenerateRandomBytes generates cryptographically secure random bytes.
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	return bytes, err
}

package cookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
)

// Manager handles signing and encrypting HTTP cookies.
type Manager struct {
	key []byte
}

// NewManager creates a new Cookie Manager. The key should be exactly 32 bytes for AES-256.
func NewManager(appKey string) *Manager {
	// Pad or truncate key to 32 bytes
	key := []byte(appKey)
	if len(key) < 32 {
		padded := make([]byte, 32)
		copy(padded, key)
		key = padded
	} else if len(key) > 32 {
		key = key[:32]
	}
	return &Manager{key: key}
}

// Set adds an encrypted and signed cookie to the response.
func (m *Manager) Set(w http.ResponseWriter, cookie *http.Cookie) error {
	encrypted, err := m.Encrypt(cookie.Value)
	if err != nil {
		return err
	}
	cookie.Value = encrypted
	http.SetCookie(w, cookie)
	return nil
}

// Get retrieves and decrypts a cookie from the request.
func (m *Manager) Get(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return m.Decrypt(cookie.Value)
}

// Encrypt encrypts and signs a string value.
func (m *Manager) Encrypt(value string) (string, error) {
	block, err := aes.NewCipher(m.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)

	// Sign the ciphertext
	mac := hmac.New(sha256.New, m.key)
	mac.Write(ciphertext)
	signature := mac.Sum(nil)

	// Format: signature + ciphertext
	final := append(signature, ciphertext...)
	return base64.URLEncoding.EncodeToString(final), nil
}

// Decrypt verifies the signature and decrypts the value.
func (m *Manager) Decrypt(value string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	if len(data) < 32 {
		return "", errors.New("invalid cookie payload length")
	}

	signature := data[:32]
	ciphertext := data[32:]

	mac := hmac.New(sha256.New, m.key)
	mac.Write(ciphertext)
	expectedMAC := mac.Sum(nil)

	if !hmac.Equal(signature, expectedMAC) {
		return "", errors.New("cookie signature mismatch")
	}

	block, err := aes.NewCipher(m.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextPart := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextPart, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}


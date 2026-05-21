package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

// Encrypter manages data encryption and decryption.
type Encrypter struct {
	key []byte
}

// NewEncrypter creates a new encrypter instance.
func NewEncrypter(appKey string) (*Encrypter, error) {
	key := []byte(appKey)
	if len(key) < 32 {
		padded := make([]byte, 32)
		copy(padded, key)
		key = padded
	} else if len(key) > 32 {
		key = key[:32]
	}

	return &Encrypter{key: key}, nil
}

// Encrypt encrypts a plaintext string using AES-256-GCM.
func (e *Encrypter) Encrypt(value string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)

	// We base64 encode the combined nonce and ciphertext
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a previously encrypted string.
func (e *Encrypter) Decrypt(payload string) (string, error) {
	// If the payload has unexpected padding or prefix, clean it.
	payload = strings.TrimSpace(payload)
	
	data, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("invalid payload length")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

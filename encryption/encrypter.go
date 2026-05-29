package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypter manages data encryption and decryption.
type Encrypter struct {
	key []byte
}

// NewEncrypter creates a new encrypter instance.
// Keys are hashed with SHA-256 to produce a consistent 32-byte key,
// preventing collision from zero-padding short keys.
func NewEncrypter(appKey string) (*Encrypter, error) {
	if appKey == "" {
		return nil, errors.New("encryption key cannot be empty")
	}
	// Hash the key to get exactly 32 bytes (AES-256)
	hash := sha256.Sum256([]byte(appKey))
	return &Encrypter{key: hash[:]}, nil
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


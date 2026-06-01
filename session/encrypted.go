package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

// EncryptedStore wraps a Store with AES-GCM encryption.
type EncryptedStore struct {
	store Store
	key   []byte
}

// NewEncryptedStore creates a new encrypted session store.
// The key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
func NewEncryptedStore(store Store, key []byte) (*EncryptedStore, error) {
	switch len(key) {
	case 16, 24, 32:
		// valid key size
	default:
		return nil, errors.New("encryption key must be 16, 24, or 32 bytes")
	}
	return &EncryptedStore{store: store, key: key}, nil
}

func (e *EncryptedStore) Read(id string) (map[string]any, error) {
	data, err := e.store.Read(id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}

	// Check if data is encrypted (stored as "_encrypted" key)
	if encryptedRaw, ok := data["_encrypted"]; ok {
		if encryptedStr, ok := encryptedRaw.(string); ok {
			decrypted, err := e.decrypt(encryptedStr)
			if err != nil {
				return nil, err
			}
			var result map[string]any
			if err := json.Unmarshal(decrypted, &result); err != nil {
				return nil, err
			}
			return result, nil
		}
	}

	return data, nil
}

func (e *EncryptedStore) Write(id string, data map[string]any) error {
	encrypted, err := e.encrypt(data)
	if err != nil {
		return err
	}

	wrapped := map[string]any{
		"_encrypted": encrypted,
	}

	return e.store.Write(id, wrapped)
}

func (e *EncryptedStore) Destroy(id string) error {
	return e.store.Destroy(id)
}

func (e *EncryptedStore) encrypt(data map[string]any) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, jsonData, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *EncryptedStore) decrypt(encrypted string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

// GenerateKey generates a random 32-byte encryption key.
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	return key, err
}

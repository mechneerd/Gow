package hashing

import (
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// ScryptHasher implements the Hasher interface using the scrypt algorithm.
type ScryptHasher struct {
	N       int
	r       int
	p       int
	keyLen  int
	saltLen int
}

// NewScryptHasher creates a new ScryptHasher with OWASP-recommended parameters.
func NewScryptHasher() *ScryptHasher {
	return &ScryptHasher{
		N:       16384,
		r:       8,
		p:       1,
		keyLen:  32,
		saltLen: 16,
	}
}

// Make hashes the given string using scrypt.
func (s *ScryptHasher) Make(value string) (string, error) {
	salt := make([]byte, s.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash, err := scrypt.Key([]byte(value), salt, s.N, s.r, s.p, s.keyLen)
	if err != nil {
		return "", err
	}

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$scrypt$N=%d,r=%d,p=%d$%s$%s",
		s.N, s.r, s.p, b64Salt, b64Hash)

	return encoded, nil
}

// Check verifies that the given string matches the hash.
func (s *ScryptHasher) Check(value, hashed string) bool {
	vals := strings.Split(hashed, "$")
	if len(vals) != 5 || vals[1] != "scrypt" {
		return false
	}

	var N, r, p int
	_, err := fmt.Sscanf(vals[2], "N=%d,r=%d,p=%d", &N, &r, &p)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[3])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return false
	}

	keyLen := len(decodedHash)
	comparisonHash, err := scrypt.Key([]byte(value), salt, N, r, p, keyLen)
	if err != nil {
		return false
	}

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1
}

// NeedsRehash determines if the hash needs to be rehashed.
func (s *ScryptHasher) NeedsRehash(hashed string) bool {
	vals := strings.Split(hashed, "$")
	if len(vals) != 5 {
		return true
	}

	var N, r, p int
	fmt.Sscanf(vals[2], "N=%d,r=%d,p=%d", &N, &r, &p)

	return N != s.N || r != s.r || p != s.p
}

// SHA512Hasher implements the Hasher interface using SHA-512 (for non-password hashing).
type SHA512Hasher struct {
	saltLen int
}

// NewSHA512Hasher creates a new SHA512Hasher.
func NewSHA512Hasher() *SHA512Hasher {
	return &SHA512Hasher{saltLen: 16}
}

// Make hashes the given string using SHA-512 with salt.
func (h *SHA512Hasher) Make(value string) (string, error) {
	salt := make([]byte, h.saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := sha512Append(salt, []byte(value))
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash[:])

	return fmt.Sprintf("$sha512$%s$%s", b64Salt, b64Hash), nil
}

// Check verifies that the given string matches the hash.
func (h *SHA512Hasher) Check(value, hashed string) bool {
	vals := strings.Split(hashed, "$")
	if len(vals) != 4 || vals[1] != "sha512" {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[2])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(vals[3])
	if err != nil {
		return false
	}

	hash := sha512Append(salt, []byte(value))
	return subtle.ConstantTimeCompare(decodedHash, hash[:]) == 1
}

// NeedsRehash always returns false for SHA-512 (no configurable cost).
func (h *SHA512Hasher) NeedsRehash(hashed string) bool {
	return false
}

func sha512Append(salt, data []byte) [64]byte {
	h := sha512.New()
	h.Write(salt)
	h.Write(data)
	var result [64]byte
	copy(result[:], h.Sum(nil))
	return result
}

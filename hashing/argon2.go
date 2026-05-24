package hashing

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Argon2Hasher implements the Hasher interface using the Argon2id algorithm.
type Argon2Hasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// NewArgon2Hasher creates a new Argon2Hasher with OWASP recommended safe defaults.
func NewArgon2Hasher(memory uint32, iterations uint32, parallelism uint8, saltLength uint32, keyLength uint32) *Argon2Hasher {
	if memory == 0 { memory = 64 * 1024 } // 64 MB
	if iterations == 0 { iterations = 3 }
	if parallelism == 0 { parallelism = 2 }
	if saltLength == 0 { saltLength = 16 }
	if keyLength == 0 { keyLength = 32 }

	return &Argon2Hasher{
		memory:      memory,
		iterations:  iterations,
		parallelism: parallelism,
		saltLength:  saltLength,
		keyLength:   keyLength,
	}
}

// Make hashes the given string using Argon2id.
func (a *Argon2Hasher) Make(value string) (string, error) {
	salt := make([]byte, a.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(value), salt, a.iterations, a.memory, a.parallelism, a.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, a.memory, a.iterations, a.parallelism, b64Salt, b64Hash)

	return encoded, nil
}

// Check verifies that the given string matches the hash.
func (a *Argon2Hasher) Check(value, hashed string) bool {
	vals := strings.Split(hashed, "$")
	if len(vals) != 6 || vals[1] != "argon2id" {
		return false
	}

	var version int
	_, err := fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return false
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return false
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return false
	}

	keyLength := uint32(len(decodedHash))
	comparisonHash := argon2.IDKey([]byte(value), salt, iterations, memory, parallelism, keyLength)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1
}

// NeedsRehash determines if the hash needs to be recomputed due to config changes.
func (a *Argon2Hasher) NeedsRehash(hashed string) bool {
	vals := strings.Split(hashed, "$")
	if len(vals) != 6 {
		return true
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8
	fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)

	return memory != a.memory || iterations != a.iterations || parallelism != a.parallelism
}


package hashing

import "golang.org/x/crypto/bcrypt"

// BcryptHasher implements the Hasher interface using the bcrypt algorithm.
type BcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new BcryptHasher with the given cost factor.
// Defaults to 12 if cost is 0.
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost <= 0 {
		cost = 12
	}
	return &BcryptHasher{cost: cost}
}

// Make hashes the given string using bcrypt.
func (b *BcryptHasher) Make(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), b.cost)
	return string(bytes), err
}

// Check verifies that the given string matches the hash.
func (b *BcryptHasher) Check(value, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(value))
	return err == nil
}

// NeedsRehash determines if the given hash needs to be rehashed based on current cost.
func (b *BcryptHasher) NeedsRehash(hashed string) bool {
	cost, err := bcrypt.Cost([]byte(hashed))
	if err != nil {
		return true
	}
	return cost != b.cost
}

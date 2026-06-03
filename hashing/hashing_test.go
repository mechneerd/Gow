package hashing

import (
	"strings"
	"testing"
)

func TestBcryptHasher_Make(t *testing.T) {
	hasher := NewBcryptHasher(10)

	hash, err := hasher.Make("password123")
	if err != nil {
		t.Fatalf("Make failed: %v", err)
	}

	if hash == "" {
		t.Error("expected non-empty hash")
	}

	if !strings.HasPrefix(hash, "$2") {
		t.Error("expected bcrypt hash to start with $2")
	}
}

func TestBcryptHasher_Check(t *testing.T) {
	hasher := NewBcryptHasher(10)

	hash, _ := hasher.Make("password123")

	if !hasher.Check("password123", hash) {
		t.Error("expected Check to return true for correct password")
	}

	if hasher.Check("wrongpassword", hash) {
		t.Error("expected Check to return false for wrong password")
	}
}

func TestBcryptHasher_NeedsRehash(t *testing.T) {
	hasher := NewBcryptHasher(10)

	hash, _ := hasher.Make("password123")

	// Same cost should not need rehash
	if hasher.NeedsRehash(hash) {
		t.Error("expected NeedsRehash to return false for same cost")
	}
}

func TestScryptHasher_Make(t *testing.T) {
	hasher := NewScryptHasher()

	hash, err := hasher.Make("password123")
	if err != nil {
		t.Fatalf("Make failed: %v", err)
	}

	if hash == "" {
		t.Error("expected non-empty hash")
	}

	if !strings.HasPrefix(hash, "$scrypt$") {
		t.Error("expected scrypt hash to start with $scrypt$")
	}
}

func TestScryptHasher_Check(t *testing.T) {
	hasher := NewScryptHasher()

	hash, _ := hasher.Make("password123")

	if !hasher.Check("password123", hash) {
		t.Error("expected Check to return true for correct password")
	}

	if hasher.Check("wrongpassword", hash) {
		t.Error("expected Check to return false for wrong password")
	}
}

func TestScryptHasher_NeedsRehash(t *testing.T) {
	hasher := NewScryptHasher()

	hash, _ := hasher.Make("password123")

	// Scrypt with current params should not need rehash
	if hasher.NeedsRehash(hash) {
		t.Error("expected NeedsRehash to return false")
	}
}

func TestSHA512Hasher_Make(t *testing.T) {
	hasher := NewSHA512Hasher()

	hash, err := hasher.Make("password123")
	if err != nil {
		t.Fatalf("Make failed: %v", err)
	}

	if hash == "" {
		t.Error("expected non-empty hash")
	}

	if !strings.HasPrefix(hash, "$sha512$") {
		t.Error("expected SHA-512 hash to start with $sha512$")
	}
}

func TestSHA512Hasher_Check(t *testing.T) {
	hasher := NewSHA512Hasher()

	hash, _ := hasher.Make("password123")

	if !hasher.Check("password123", hash) {
		t.Error("expected Check to return true for correct password")
	}

	if hasher.Check("wrongpassword", hash) {
		t.Error("expected Check to return false for wrong password")
	}
}

func TestSHA512Hasher_NeedsRehash(t *testing.T) {
	hasher := NewSHA512Hasher()

	hash, _ := hasher.Make("password123")

	// SHA-512 never needs rehash
	if hasher.NeedsRehash(hash) {
		t.Error("expected NeedsRehash to return false for SHA-512")
	}
}

func TestSHA512Hasher_ConsistentOutput(t *testing.T) {
	hasher := NewSHA512Hasher()

	hash1, _ := hasher.Make("same-password")
	hash2, _ := hasher.Make("same-password")

	// Different salts, so hashes should be different
	if hash1 == hash2 {
		t.Error("expected different hashes due to random salt")
	}
}

func TestHasherInterface(t *testing.T) {
	var _ Hasher = NewBcryptHasher(10)
	var _ Hasher = NewScryptHasher()
	var _ Hasher = NewSHA512Hasher()
}

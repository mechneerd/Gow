package hashing

// Hasher provides an interface for password hashing mechanisms.
type Hasher interface {
	// Make hashes the given string.
	Make(value string) (string, error)
	// Check verifies that the given plain text string matches the hash.
	Check(value, hashed string) bool
	// NeedsRehash determines if the given hash has been hashed using the configured options.
	NeedsRehash(hashed string) bool
}


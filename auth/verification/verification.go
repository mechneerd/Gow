package verification

import (
	"fmt"
	"time"

	"github.com/mechneerd/gow/auth"
	"github.com/mechneerd/gow/routing"
)

// GenerateVerificationURL creates a signed URL for email verification.
// The URL will be valid for the given duration.
func GenerateVerificationURL(user auth.Authenticatable, expiresIn time.Duration) string {
	// We use the signed route system
	// Assume a route named "verification.verify" exists and accepts id + hash (or email + signature)
	id := user.GetAuthIdentifier()

	// For simplicity, we sign the email + timestamp
	payload := fmt.Sprintf("%s:%d", id, time.Now().Add(expiresIn).Unix())

	signedURL, _ := routing.TemporarySignedRoute("verification.verify", expiresIn, map[string]string{
		"id":    id,
		"hash":  payload, // in real use we'd hash the payload
	})

	return signedURL
}

// MarkAsVerified sets the email_verified_at timestamp on the user.
// In a real app, this would update the user record via ORM.
func MarkAsVerified(user auth.Authenticatable) {
	// Placeholder — the actual User model should have this field and be persisted
	// Example: user.EmailVerifiedAt = time.Now()
	// Then save via your UserProvider or ORM
}


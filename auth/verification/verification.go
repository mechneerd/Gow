package verification

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/mechneerd/gow/database/orm"
)

// Verifyable is an interface that users must implement for email verification.
type Verifyable interface {
	GetAuthIdentifier() string
}

// VerifyableModel extends Verifyable with verification methods.
type VerifyableModel interface {
	Verifyable
	IsEmailVerified() bool
}

// GenerateVerificationURL creates a signed URL for email verification.
// The URL will be valid for the given duration.
func GenerateVerificationURL(userID string, secret []byte, expiresIn time.Duration) string {
	expiry := time.Now().Add(expiresIn).Unix()
	payload := fmt.Sprintf("%s:%d", userID, expiry)

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("/verify-email?id=%s&expires=%d&signature=%s", userID, expiry, signature)
}

// VerifyEmail validates the verification URL parameters.
func VerifyEmail(userID, expiresStr, signature string, secret []byte) error {
	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid expiration")
	}

	if time.Now().Unix() > expires {
		return fmt.Errorf("verification link has expired")
	}

	payload := fmt.Sprintf("%s:%d", userID, expires)
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	expectedSig := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return fmt.Errorf("invalid verification signature")
	}

	return nil
}

// MarkAsVerified updates the user's email_verified_at timestamp.
func MarkAsVerified(db *orm.DB, userID string, userTable string) error {
	if userTable == "" {
		userTable = "users"
	}
	_, err := db.Builder.Clone().
		Table(userTable).
		Where("id", "=", userID).
		Update(map[string]any{
			"email_verified_at": time.Now(),
		})
	return err
}

// IsEmailVerified checks if a user has verified their email.
func IsEmailVerified(db *orm.DB, userID string, userTable string) bool {
	if userTable == "" {
		userTable = "users"
	}
	rows, err := db.Builder.Clone().
		Table(userTable).
		Select("email_verified_at").
		Where("id", "=", userID).
		Limit(1).
		Get()
	if err != nil {
		return false
	}
	defer rows.Close()

	if !rows.Next() {
		return false
	}

	var verifiedAt *time.Time
	if err := rows.Scan(&verifiedAt); err != nil {
		return false
	}

	return verifiedAt != nil
}

// GenerateToken creates a verification token and stores it.
func GenerateToken(db *orm.DB, userID string, expiresIn time.Duration) (string, error) {
	token, err := generateRandomToken(32)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(expiresIn)

	_, err = db.Builder.Clone().
		Table("email_verifications").
		Insert(map[string]any{
			"user_id":    userID,
			"token":      token,
			"expires_at": expiresAt,
			"created_at": time.Now(),
		})
	if err != nil {
		return "", err
	}

	return token, nil
}

// VerifyToken validates a verification token.
func VerifyToken(db *orm.DB, userID, token string) error {
	rows, err := db.Builder.Clone().
		Table("email_verifications").
		Where("user_id", "=", userID).
		Where("token", "=", token).
		Limit(1).
		Get()
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return fmt.Errorf("invalid verification token")
	}

	// Delete the token after successful verification
	_, err = db.Builder.Clone().
		Table("email_verifications").
		Where("user_id", "=", userID).
		Where("token", "=", token).
		Delete()
	return err
}

func generateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to time-based
		return fmt.Sprintf("%d", time.Now().UnixNano()), nil
	}
	return hex.EncodeToString(b), nil
}


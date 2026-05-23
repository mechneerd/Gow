package fortify

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// TwoFactor provides TOTP-based 2FA (Laravel Fortify style).
// Pure Go implementation — no external libraries needed.

const (
	totpDigits = 6
	totpPeriod = 30 // seconds
)

// GenerateSecret creates a new base32 secret for the user (20 bytes recommended).
func GenerateSecret() (string, error) {
	secret := make([]byte, 20)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(secret), nil
}

// GenerateTOTPCode returns the current 6-digit TOTP code for the given secret.
func GenerateTOTPCode(secret string) string {
	return generateTOTP(secret, time.Now().Unix()/totpPeriod)
}

// VerifyTOTP checks if the provided code matches the current or previous window (allows clock skew).
func VerifyTOTP(secret, code string, skew int) bool {
	if len(code) != totpDigits {
		return false
	}
	now := time.Now().Unix() / totpPeriod

	for i := -skew; i <= skew; i++ {
		if generateTOTP(secret, now+int64(i)) == code {
			return true
		}
	}
	return false
}

func generateTOTP(secret string, counter int64) string {
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(secret))
	if err != nil {
		return "000000"
	}

	msg := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		msg[i] = byte(counter & 0xff)
		counter >>= 8
	}

	h := hmac.New(sha1.New, key)
	h.Write(msg)
	hmacHash := h.Sum(nil)

	offset := hmacHash[len(hmacHash)-1] & 0x0f
	binary := int64(hmacHash[offset]&0x7f)<<24 |
		int64(hmacHash[offset+1])<<16 |
		int64(hmacHash[offset+2])<<8 |
		int64(hmacHash[offset+3])

	otp := binary % int64(math.Pow10(totpDigits))
	return fmt.Sprintf("%0"+strconv.Itoa(totpDigits)+"d", otp)
}

// GetQRCodeURL returns a otpauth:// URL for Google Authenticator / Authy etc.
// Example: otpauth://totp/Company:email?secret=XXXX&issuer=Company
func GetQRCodeURL(issuer, accountName, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=SHA1&digits=%d&period=%d",
		issuer, accountName, secret, issuer, totpDigits, totpPeriod)
}

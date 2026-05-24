package routing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// URLGenerator provides methods to generate URLs to named routes.
type URLGenerator struct {
	router *Router
	appKey string // Used for HMAC signatures
}

// NewURLGenerator creates a new URL Generator.
func NewURLGenerator(router *Router, appKey string) *URLGenerator {
	return &URLGenerator{
		router: router,
		appKey: appKey,
	}
}

// SignedRoute generates a cryptographically signed URL for a named route.
func (u *URLGenerator) SignedRoute(name string, parameters map[string]string, expiration time.Time) (string, error) {
	route, exists := u.router.namedRoutes[name]
	if !exists {
		return "", fmt.Errorf("route %s not found", name)
	}

	path := route.Path
	for k, v := range parameters {
		path = strings.ReplaceAll(path, "{"+k+"}", url.PathEscape(v))
	}

	queryParams := url.Values{}
	if !expiration.IsZero() {
		queryParams.Set("expires", strconv.FormatInt(expiration.Unix(), 10))
	}

	// Generate signature
	signatureString := path
	if len(queryParams) > 0 {
		signatureString += "?" + queryParams.Encode()
	}

	mac := hmac.New(sha256.New, []byte(u.appKey))
	mac.Write([]byte(signatureString))
	signature := hex.EncodeToString(mac.Sum(nil))

	queryParams.Set("signature", signature)

	return path + "?" + queryParams.Encode(), nil
}

// HasValidSignature verifies if the request has a valid signature.
func (u *URLGenerator) HasValidSignature(r *http.Request) bool {
	originalSignature := r.URL.Query().Get("signature")
	if originalSignature == "" {
		return false
	}

	// Remove signature from query to verify the rest
	q := r.URL.Query()
	q.Del("signature")

	// Check expiration
	expiresStr := q.Get("expires")
	if expiresStr != "" {
		expiresInt, err := strconv.ParseInt(expiresStr, 10, 64)
		if err == nil {
			if time.Now().Unix() > expiresInt {
				return false
			}
		}
	}

	signatureString := r.URL.Path
	if len(q) > 0 {
		signatureString += "?" + q.Encode()
	}

	mac := hmac.New(sha256.New, []byte(u.appKey))
	mac.Write([]byte(signatureString))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(originalSignature), []byte(expectedSignature))
}

// TemporarySignedRoute is a convenience helper that generates a signed URL
// that expires after the given duration. It requires a global router to be set
// via SetGlobalRouterForSignedURLs (or similar) in real apps.
var globalURLGenerator *URLGenerator

// SetGlobalURLGenerator allows setting a global URL generator for convenience helpers.
func SetGlobalURLGenerator(gen *URLGenerator) {
	globalURLGenerator = gen
}

// TemporarySignedRoute generates a temporary signed URL (convenience wrapper).
func TemporarySignedRoute(name string, expiresIn time.Duration, params map[string]string) (string, error) {
	if globalURLGenerator == nil {
		return "", fmt.Errorf("no global URL generator set. Use routing.SetGlobalURLGenerator")
	}
	expiresAt := time.Now().Add(expiresIn)
	return globalURLGenerator.SignedRoute(name, params, expiresAt)
}


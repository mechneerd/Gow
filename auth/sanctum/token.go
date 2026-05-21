package sanctum

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"gow/database/orm"
	"gow/database/query"
	"gow/support/str"
)

// PersonalAccessToken represents the database model.
type PersonalAccessToken struct {
	ID            int
	TokenableType string
	TokenableID   string
	Name          string
	Token         string
	Abilities     []string
	LastUsedAt    *time.Time
	ExpiresAt     *time.Time
}

// TokenManager handles issuing and verifying tokens.
type TokenManager struct {
	db *orm.DB
}

// NewTokenManager creates a new token manager.
func NewTokenManager(db *orm.DB) *TokenManager {
	return &TokenManager{db: db}
}

// CreateToken issues a new token for the given entity.
func (tm *TokenManager) CreateToken(tokenableType, tokenableID, name string, abilities []string) (string, error) {
	// Generate plain text token
	plainTextToken := str.Random(40)
	
	// Hash it for database storage
	hash := sha256.Sum256([]byte(plainTextToken))
	hashedToken := hex.EncodeToString(hash[:])
	
	abilitiesJSON, _ := json.Marshal(abilities)

	builder := query.NewBuilder(tm.db.RawDB(), tm.db.Dialect())
	builder.Table("personal_access_tokens")
	
	// Insert into DB
	// In a real framework, we'd retrieve the inserted ID and return `ID|plainTextToken`
	_, err := builder.Insert(map[string]any{
		"tokenable_type": tokenableType,
		"tokenable_id":   tokenableID,
		"name":           name,
		"token":          hashedToken,
		"abilities":      string(abilitiesJSON),
		"created_at":     time.Now().Unix(),
		"updated_at":     time.Now().Unix(),
	})

	if err != nil {
		return "", err
	}

	// For simplicity, we just return the plain token here without the ID prefix
	return plainTextToken, nil
}

// ContextKey used for storing token/user info in request context
type ContextKey string
const TokenKey ContextKey = "sanctum_token"

// Middleware protects routes requiring API authentication.
func (tm *TokenManager) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			
			plainToken := strings.TrimPrefix(authHeader, "Bearer ")
			hash := sha256.Sum256([]byte(plainToken))
			hashedToken := hex.EncodeToString(hash[:])

			// Query DB for token
			builder := query.NewBuilder(tm.db.RawDB(), tm.db.Dialect())
			builder.Table("personal_access_tokens")
			builder.Where("token", "=", hashedToken)
			builder.Limit(1)
			
			rows, err := builder.Get()
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			defer rows.Close()
			
			if !rows.Next() {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// We could scan the rows to populate PersonalAccessToken
			// and verify ExpiresAt
			// Here we just accept it if it exists
			
			// Inject token identity into context
			ctx := context.WithValue(r.Context(), TokenKey, plainToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

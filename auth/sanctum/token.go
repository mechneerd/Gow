package sanctum

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mechneerd/gow/database/orm"
	"github.com/mechneerd/gow/support/str"
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
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TableName returns the table name for the model.
func (PersonalAccessToken) TableName() string {
	return "personal_access_tokens"
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
func (tm *TokenManager) CreateToken(tokenableType, tokenableID, name string, abilities []string, expiry ...time.Duration) (string, error) {
	// Generate plain text token
	plainTextToken := str.Random(40)

	// Hash it for database storage
	hash := sha256.Sum256([]byte(plainTextToken))
	hashedToken := hex.EncodeToString(hash[:])

	abilitiesJSON, _ := json.Marshal(abilities)

	now := time.Now()
	var expiresAt *time.Time
	if len(expiry) > 0 {
		exp := now.Add(expiry[0])
		expiresAt = &exp
	}

	// Insert into DB and get the ID
	builder := tm.db.Builder.Clone()
	builder.Table("personal_access_tokens")

	_, err := builder.Insert(map[string]any{
		"tokenable_type": tokenableType,
		"tokenable_id":   tokenableID,
		"name":           name,
		"token":          hashedToken,
		"abilities":      string(abilitiesJSON),
		"last_used_at":   nil,
		"expires_at":     expiresAt,
		"created_at":     now,
		"updated_at":     now,
	})

	if err != nil {
		return "", err
	}

	return plainTextToken, nil
}

// Revoke deletes a token by its plain text value.
func (tm *TokenManager) Revoke(plainToken string) error {
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	_, err := tm.db.Builder.Clone().
		Table("personal_access_tokens").
		Where("token", "=", hashedToken).
		Delete()
	return err
}

// RevokeAllFor deletes all tokens for a given entity.
func (tm *TokenManager) RevokeAllFor(tokenableType, tokenableID string) error {
	_, err := tm.db.Builder.Clone().
		Table("personal_access_tokens").
		Where("tokenable_type", "=", tokenableType).
		Where("tokenable_id", "=", tokenableID).
		Delete()
	return err
}

// GetToken retrieves a token record by its plain text value.
func (tm *TokenManager) GetToken(plainToken string) (*PersonalAccessToken, error) {
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	rows, err := tm.db.Builder.Clone().
		Table("personal_access_tokens").
		Where("token", "=", hashedToken).
		Limit(1).
		Get()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	cols, _ := rows.Columns()
	token := &PersonalAccessToken{}
	values := make([]any, len(cols))
	valuePtrs := make([]any, len(cols))
	for i := range values {
		valuePtrs[i] = &values[i]
	}
	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	// Map columns to struct fields
	colMap := make(map[string]int)
	for i, col := range cols {
		colMap[col] = i
	}

	if idx, ok := colMap["id"]; ok {
		if v, ok := values[idx].(*any); ok {
			token.ID = int((*v).(int64))
		}
	}
	if idx, ok := colMap["tokenable_type"]; ok {
		if v, ok := values[idx].(*any); ok {
			token.TokenableType = (*v).(string)
		}
	}
	if idx, ok := colMap["tokenable_id"]; ok {
		if v, ok := values[idx].(*any); ok {
			token.TokenableID = (*v).(string)
		}
	}
	if idx, ok := colMap["name"]; ok {
		if v, ok := values[idx].(*any); ok {
			token.Name = (*v).(string)
		}
	}
	if idx, ok := colMap["abilities"]; ok {
		if v, ok := values[idx].(*any); ok {
			var abilities []string
			json.Unmarshal([]byte((*v).(string)), &abilities)
			token.Abilities = abilities
		}
	}
	if idx, ok := colMap["expires_at"]; ok {
		if v, ok := values[idx].(*any); ok && *v != nil {
			if t, ok := (*v).(time.Time); ok {
				token.ExpiresAt = &t
			}
		}
	}
	if idx, ok := colMap["last_used_at"]; ok {
		if v, ok := values[idx].(*any); ok && *v != nil {
			if t, ok := (*v).(time.Time); ok {
				token.LastUsedAt = &t
			}
		}
	}

	return token, nil
}

// HasAbility checks if a token has a specific ability.
func (t *PersonalAccessToken) HasAbility(ability string) bool {
	for _, a := range t.Abilities {
		if a == ability || a == "*" {
			return true
		}
	}
	return false
}

// Can checks if the token has ALL of the given abilities.
func (t *PersonalAccessToken) Can(abilities ...string) bool {
	for _, ability := range abilities {
		if !t.HasAbility(ability) {
			return false
		}
	}
	return true
}

// Cannot checks if the token is missing any of the given abilities.
func (t *PersonalAccessToken) Cannot(abilities ...string) bool {
	return !t.Can(abilities...)
}

// IsExpired checks if the token has expired.
func (t *PersonalAccessToken) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false // no expiry set
	}
	return time.Now().After(*t.ExpiresAt)
}

// UpdateLastUsedAt updates the last_used_at timestamp.
func (tm *TokenManager) UpdateLastUsedAt(plainToken string) error {
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	_, err := tm.db.Builder.Clone().
		Table("personal_access_tokens").
		Where("token", "=", hashedToken).
		Update(map[string]any{
			"last_used_at": time.Now(),
			"updated_at":   time.Now(),
		})
	return err
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

			// Retrieve and validate token
			token, err := tm.GetToken(plainToken)
			if err != nil || token == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check expiration
			if token.IsExpired() {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}

			// Update last used timestamp
			_ = tm.UpdateLastUsedAt(plainToken)

			// Inject token into context
			ctx := context.WithValue(r.Context(), TokenKey, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAbilities returns middleware that checks for specific abilities.
func (tm *TokenManager) RequireAbilities(abilities ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := r.Context().Value(TokenKey).(*PersonalAccessToken)
			if !ok || token == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !token.Can(abilities...) {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// TokenFromContext extracts the token from the request context.
func TokenFromContext(r *http.Request) (*PersonalAccessToken, bool) {
	token, ok := r.Context().Value(TokenKey).(*PersonalAccessToken)
	return token, ok
}


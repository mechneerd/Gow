package password

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/mechneerd/gow/auth"
	"github.com/mechneerd/gow/database/orm"
	"github.com/mechneerd/gow/hashing"
	"github.com/mechneerd/gow/mail"
)

// Broker handles password reset logic.
type Broker struct {
	db      *orm.DB
	mailer  *mail.Mailer
	hasher  hashing.Hasher
	userTab string // table name for users (default: "users")
}

func NewBroker(db *orm.DB, mailer *mail.Mailer) *Broker {
	return &Broker{
		db:      db,
		mailer:  mailer,
		hasher:  hashing.NewBcryptHasher(10),
		userTab: "users",
	}
}

// NewBrokerWithHasher creates a broker with a custom hasher.
func NewBrokerWithHasher(db *orm.DB, mailer *mail.Mailer, hasher hashing.Hasher) *Broker {
	return &Broker{
		db:      db,
		mailer:  mailer,
		hasher:  hasher,
		userTab: "users",
	}
}

// SetUserTable sets the users table name.
func (b *Broker) SetUserTable(table string) {
	b.userTab = table
}

// SendResetLink generates a token and sends the reset email.
func (b *Broker) SendResetLink(email string) error {
	token, err := generateToken(32)
	if err != nil {
		return err
	}

	// Store token
	if err := Create(b.db, email, token); err != nil {
		return err
	}

	// Send email
	mailable := NewResetPasswordMailable(email, token)
	return b.mailer.Send(mailable)
}

// Reset validates the token and updates the user's password.
func (b *Broker) Reset(email, token, newPassword string) error {
	reset, err := Find(b.db, email, token)
	if err != nil || reset == nil {
		return auth.ErrInvalidToken
	}

	// Hash the new password
	hashedPassword, err := b.hasher.Make(newPassword)
	if err != nil {
		return err
	}

	// Update the user's password in the database
	_, err = b.db.Builder.Clone().
		Table(b.userTab).
		Where("email", "=", email).
		Update(map[string]any{
			"password": hashedPassword,
		})
	if err != nil {
		return err
	}

	// Delete the reset token
	if err := Delete(b.db, email, token); err != nil {
		return err
	}

	return nil
}

// ValidateToken checks if a reset token is valid without consuming it.
func (b *Broker) ValidateToken(email, token string) error {
	reset, err := Find(b.db, email, token)
	if err != nil {
		return err
	}
	if reset == nil {
		return errors.New("invalid or expired reset token")
	}
	return nil
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}


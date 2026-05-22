package password

import (
	"crypto/rand"
	"encoding/hex"
	"gow/auth"
	"gow/database/orm"
	"gow/mail"
)

// Broker handles password reset logic.
type Broker struct {
	db     *orm.DB
	mailer *mail.Mailer
}

func NewBroker(db *orm.DB, mailer *mail.Mailer) *Broker {
	return &Broker{db: db, mailer: mailer}
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
		return auth.ErrInvalidToken // we'll define this
	}

	// In a real app, you would look up the user and update their password using the UserProvider or ORM.
	// For now, we just delete the token as a placeholder.
	// TODO: Integrate with actual User model and hash the password.

	if err := Delete(b.db, email, token); err != nil {
		return err
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

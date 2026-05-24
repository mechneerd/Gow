package password

import (
	"time"

	"gow/database/orm"
)

// PasswordReset represents a password reset token.
type PasswordReset struct {
	Email     string    `db:"email"`
	Token     string    `db:"token"`
	CreatedAt time.Time `db:"created_at"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}

// Create a new reset token record.
func Create(db *orm.DB, email, token string) error {
	reset := &PasswordReset{
		Email:     email,
		Token:     token,
		CreatedAt: time.Now(),
	}
	return orm.NewQuery[PasswordReset](db).Insert(reset)
}

// Find a valid token (not expired, e.g. within 60 minutes).
func Find(db *orm.DB, email, token string) (*PasswordReset, error) {
	reset, err := orm.NewQuery[PasswordReset](db).
		Where("email", "=", email).
		Where("token", "=", token).
		First()
	if err != nil {
		return nil, err
	}
	if reset == nil {
		return nil, nil
	}
	// Check expiry (60 minutes)
	if time.Since(reset.CreatedAt) > 60*time.Minute {
		// Delete expired
		orm.NewQuery[PasswordReset](db).
			Where("email", "=", email).
			Where("token", "=", token).
			Delete(reset)
		return nil, nil
	}
	return reset, nil
}

// Delete token after successful reset.
func Delete(db *orm.DB, email, token string) error {
	err := orm.NewQuery[PasswordReset](db).
		Where("email", "=", email).
		Where("token", "=", token).
		Delete(&PasswordReset{})
	return err
}

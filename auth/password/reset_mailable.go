package password

import "github.com/mechneerd/gow/mail"

// ResetPasswordMailable sends the password reset link.
type ResetPasswordMailable struct {
	*mail.BaseMailable
	Email string
	Token string
}

func NewResetPasswordMailable(email, token string) *ResetPasswordMailable {
	return &ResetPasswordMailable{
		BaseMailable: mail.NewMailable(),
		Email:        email,
		Token:        token,
	}
}

func (m *ResetPasswordMailable) Build() *mail.Message {
	resetURL := "http://localhost:8080/reset-password?token=" + m.Token + "&email=" + m.Email

	markdown := `
# Reset Your Password

Hello,

You are receiving this email because we received a password reset request for your account.

**Reset Password Link:**

[` + resetURL + `]("` + resetURL + `")

This password reset link will expire in 60 minutes.

If you did not request a password reset, no further action is required.

Regards,
The GoW Team
`

	return m.BaseMailable.
		To(m.Email).
		Subject("Reset Password Notification").
		Markdown(markdown).
		Build()
}


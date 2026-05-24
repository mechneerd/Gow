package access

import (
	"testing"
)

type User struct {
	ID    int
	Admin bool
}

type Post struct {
	ID     int
	UserID int
}

type PostPolicy struct{}

func (p *PostPolicy) Update(user User, post Post) bool {
	return user.ID == post.UserID
}

func (p *PostPolicy) Delete(user User, post Post) bool {
	return user.ID == post.UserID
}

func TestGateClosures(t *testing.T) {
	gate := NewGate()

	gate.Define("view-dashboard", func(user any, args ...any) bool {
		u := user.(User)
		return u.Admin
	})

	admin := User{ID: 1, Admin: true}
	guest := User{ID: 2, Admin: false}

	if !gate.Allows(admin, "view-dashboard") {
		t.Error("Admin should be allowed to view dashboard")
	}

	if gate.Allows(guest, "view-dashboard") {
		t.Error("Guest should be denied from viewing dashboard")
	}

	if !gate.Denies(guest, "view-dashboard") {
		t.Error("Denies should return true for guest")
	}
}

func TestGatePolicies(t *testing.T) {
	gate := NewGate()
	gate.Policy("Post", &PostPolicy{})

	author := User{ID: 1}
	reader := User{ID: 2}
	post := Post{ID: 10, UserID: 1}

	if !gate.Allows(author, "update", post) {
		t.Error("Author should be allowed to update their post")
	}

	if gate.Allows(reader, "update", post) {
		t.Error("Reader should be denied from updating another's post")
	}
}

func TestGateBeforeAfterHooks(t *testing.T) {
	gate := NewGate()

	gate.Define("edit-settings", func(user any, args ...any) bool {
		return false // Normally nobody can edit settings
	})

	// Super admin bypass
	gate.Before(func(user any, ability string, args ...any) *bool {
		u := user.(User)
		if u.Admin {
			res := true
			return &res
		}
		return nil
	})

	admin := User{ID: 1, Admin: true}
	normal := User{ID: 2, Admin: false}

	if !gate.Allows(admin, "edit-settings") {
		t.Error("Admin should bypass closure and be allowed")
	}

	if gate.Allows(normal, "edit-settings") {
		t.Error("Normal user should hit closure and be denied")
	}
}


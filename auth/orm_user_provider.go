package auth

import (
	"gow/database/orm"
	"gow/hashing"
)

// ORMUserProvider is a UserProvider implementation that retrieves users from the ORM.
type ORMUserProvider struct {
	db         *orm.DB
	model      Authenticatable // used as prototype for reflection / table
	tableName  string
}

// NewORMUserProvider creates a new ORM-backed user provider.
func NewORMUserProvider(db *orm.DB, model Authenticatable) *ORMUserProvider {
	return &ORMUserProvider{
		db:    db,
		model: model,
	}
}

func (p *ORMUserProvider) RetrieveByID(identifier string) any {
	// This is a simplified version. In real use, you would use the model's table.
	// For now we assume the model has a table and we query by primary key "id".
	_ = p.model // prototype for future reflection-based instantiation
	// Placeholder implementation - users should replace with their actual model query
	_ = identifier
	return nil
}

func (p *ORMUserProvider) RetrieveByToken(identifier string, token string) any {
	// Remember token support (future)
	return nil
}

func (p *ORMUserProvider) UpdateRememberToken(user any, token string) {
	// Future
}

func (p *ORMUserProvider) RetrieveByCredentials(credentials map[string]any) any {
	// Typically find by email/username
	// Note: This is a simplified implementation. Real apps should register a concrete
	// user model that satisfies Authenticatable and use the model's table/query methods.
	return nil
}

func (p *ORMUserProvider) ValidateCredentials(user any, credentials map[string]any) bool {
	if u, ok := user.(Authenticatable); ok {
		password, hasPass := credentials["password"].(string)
		if !hasPass {
			return false
		}
		hashed := u.GetAuthPassword()
		return hashing.NewBcryptHasher(10).Check(password, hashed)
	}
	return false
}

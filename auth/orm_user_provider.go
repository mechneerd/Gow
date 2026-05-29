package auth

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/mechneerd/gow/database/orm"
	"github.com/mechneerd/gow/hashing"
)

// ORMUserProvider is a UserProvider implementation that retrieves users from the ORM.
type ORMUserProvider struct {
	db    *orm.DB
	model Authenticatable
}

// NewORMUserProvider creates a new ORM-backed user provider.
func NewORMUserProvider(db *orm.DB, model Authenticatable) *ORMUserProvider {
	return &ORMUserProvider{
		db:    db,
		model: model,
	}
}

func (p *ORMUserProvider) getTableName() string {
	t := reflect.TypeOf(p.model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if m, ok := p.model.(orm.Model); ok {
		return m.TableName()
	}
	return fmt.Sprintf("%ss", strings.ToLower(t.Name()))
}

func (p *ORMUserProvider) getSQLDB() (*sql.DB, bool) {
	db, ok := p.db.RawDB().(*sql.DB)
	return db, ok
}

func (p *ORMUserProvider) RetrieveByID(identifier string) any {
	sqlDB, ok := p.getSQLDB()
	if !ok {
		return nil
	}

	table := p.getTableName()
	var id int
	var name, email, password string
	err := sqlDB.QueryRow(
		fmt.Sprintf(`SELECT "id", "name", "email", "password" FROM "%s" WHERE "id" = ? LIMIT 1`, table),
		identifier,
	).Scan(&id, &name, &email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return nil
	}

	return &GenericUser{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
}

func (p *ORMUserProvider) RetrieveByCredentials(credentials map[string]any) any {
	sqlDB, ok := p.getSQLDB()
	if !ok {
		return nil
	}

	table := p.getTableName()
	var conditions []string
	var args []any

	for k, v := range credentials {
		if k == "password" {
			continue
		}
		conditions = append(conditions, fmt.Sprintf(`"%s" = ?`, k))
		args = append(args, v)
	}

	if len(conditions) == 0 {
		return nil
	}

	var id int
	var name, email, password string
	query := fmt.Sprintf(`SELECT "id", "name", "email", "password" FROM "%s" WHERE %s LIMIT 1`, table, joinConditions(conditions))
	err := sqlDB.QueryRow(query, args...).Scan(&id, &name, &email, &password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return nil
	}

	return &GenericUser{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
	}
}

func (p *ORMUserProvider) RetrieveByToken(identifier string, token string) any {
	return nil
}

func (p *ORMUserProvider) UpdateRememberToken(user any, token string) {
	sqlDB, ok := p.getSQLDB()
	if !ok {
		return
	}
	if authUser, ok := user.(Authenticatable); ok {
		table := p.getTableName()
		_, _ = sqlDB.Exec(
			fmt.Sprintf(`UPDATE "%s" SET "remember_token" = ? WHERE "id" = ?`, table),
			token, authUser.GetAuthIdentifier(),
		)
	}
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

// GenericUser is a basic Authenticatable implementation for ORMUserProvider.
type GenericUser struct {
	ID       int
	Name     string
	Email    string
	Password string
}

func (u *GenericUser) GetAuthIdentifier() string {
	return fmt.Sprintf("%d", u.ID)
}

func (u *GenericUser) GetAuthPassword() string {
	return u.Password
}

func joinConditions(conditions []string) string {
	result := ""
	for i, c := range conditions {
		if i > 0 {
			result += " AND "
		}
		result += c
	}
	return result
}


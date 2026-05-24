package rbac

import "github.com/mechneerd/gow/database/orm"

type Permission struct {
	orm.Model
	Name      string
	GuardName string
}

func (Permission) TableName() string { return "permissions" }


package rbac

import "gow/database/orm"

type Role struct {
	orm.Model
	Name      string
	GuardName string
}

func (Role) TableName() string { return "roles" }

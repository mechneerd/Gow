package orm

// Lifecycle events dispatched by the ORM.

type ModelSaving struct{ Model any }
type ModelSaved struct{ Model any }
type ModelCreating struct{ Model any }
type ModelCreated struct{ Model any }
type ModelUpdating struct{ Model any }
type ModelUpdated struct{ Model any }
type ModelDeleting struct{ Model any }
type ModelDeleted struct{ Model any }

// Interface Hooks that models can implement.

type BeforeSaveHook interface {
	BeforeSave() error
}

type AfterSaveHook interface {
	AfterSave() error
}

type BeforeCreateHook interface {
	BeforeCreate() error
}

type AfterCreateHook interface {
	AfterCreate() error
}

type BeforeUpdateHook interface {
	BeforeUpdate() error
}

type AfterUpdateHook interface {
	AfterUpdate() error
}

type BeforeDeleteHook interface {
	BeforeDelete() error
}

type AfterDeleteHook interface {
	AfterDelete() error
}

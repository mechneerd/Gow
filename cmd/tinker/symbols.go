package tinker

import (
	"reflect"

	"github.com/traefik/yaegi/interp"

	"gow/config"
	"gow/database/orm"
	"gow/foundation"
)

// Symbols returns a map of preloaded GoW objects for the Tinker REPL.
func Symbols(app *foundation.Application, db *orm.DB, cfg *config.Repository) interp.Exports {
	return interp.Exports{
		"tinker/tinker": map[string]reflect.Value{
			"App":    reflect.ValueOf(app),
			"DB":     reflect.ValueOf(db),
			"Config": reflect.ValueOf(cfg),
		},
	}
}

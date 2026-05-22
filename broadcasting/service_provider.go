package broadcasting

import (
	"reflect"

	"gow/foundation"
	"log"
)

// ServiceProvider bootstraps the native WebSocket broadcasting system.
// It creates a Hub (started in background), the WebSocketDriver, registers the
// driver with the broadcasting Manager (if present in container), and binds
// both for dependency injection.
//
// To wire the WebSocket endpoint, register the route in your application routes:
//
//   hub := app.Make(broadcasting.Hub{})
//   router.Any("/ws", func(w http.ResponseWriter, r *http.Request) {
//       broadcasting.ServeWs(hub, w, r)
//   })
//
// The /ws path is conventional for native WebSocket connections.
type ServiceProvider struct {
	foundation.BaseServiceProvider
}

// Register sets up the Hub and WebSocket driver.
func (p *ServiceProvider) Register(app *foundation.Application) {
	hub := NewHub()
	go hub.Run()

	app.Instance((*Hub)(nil), hub)

	wsDriver := NewWebSocketDriver(hub)
	app.Instance((*WebSocketDriver)(nil), wsDriver)

	// Attempt to register the websocket driver with the central broadcasting manager
	managerType := reflect.TypeOf((*Manager)(nil))
	if iface, err := app.Resolve(managerType); err == nil {
		if bm, ok := iface.(*Manager); ok {
			bm.Extend("websocket", wsDriver)
			log.Println("[GoW] WebSocket broadcast driver registered")
		}
	}
}

// Boot can be used to attach shutdown hooks etc.
func (p *ServiceProvider) Boot(app *foundation.Application) {
	// Example: if kernel is available, register hub stop
	// kernelIface, _ := app.Resolve((*http.Kernel)(nil))
	// if k, ok := ... ; ok { k.OnShutdown(hub.Stop) }
}

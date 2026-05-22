> ✅ Implemented · 🚧 In Progress · 📋 Planned

# Events

> **Status**: ✅ Implemented


GoW's events provide a simple observer pattern implementation, allowing you to subscribe and listen for various events that occur in your application. 

## Registering Events & Listeners

You can register events and their listeners in your `app/providers/event_service_provider.go`.

```go
func (p *EventServiceProvider) Register(events *events.Manager) {
    events.Listen("OrderShipped", func(e events.Event) {
        orderEvent := e.(*OrderShipped)
        // Send email to customer...
    })
}
```

## Dispatching Events

To dispatch an event, you may pass an instance of the event struct to the `events.Dispatch` method. The event manager will automatically derive the event name from the struct type and dispatch it to all registered listeners.

```go
type OrderShipped struct {
    OrderID int
}

eventsManager.Dispatch(&OrderShipped{OrderID: 123})
```

## Queued Listeners

If your listener performs a slow task like sending an e-mail or making an HTTP request, you may queue it. When using `QueueListen`, GoW will automatically dispatch the listener execution to your configured Queue driver in the background.

```go
events.QueueListen("OrderShipped", SendShipmentNotification)
```

## Event Subscribers

Event subscribers are structs that may subscribe to multiple events from within the subscriber class itself.

```go
type UserEventSubscriber struct {}

func (s *UserEventSubscriber) Subscribe(events *events.Manager) {
    events.Listen("UserLoggedIn", s.handleUserLogin)
    events.Listen("UserLoggedOut", s.handleUserLogout)
}
```

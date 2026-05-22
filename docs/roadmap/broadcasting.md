# Broadcasting

> **Status**: ?? Planned (Currently stubbed or not implemented)


In many modern web applications, WebSockets are used to implement real-time, live-updating user interfaces. When some data is updated on the server, a message is typically sent over a WebSocket connection to be handled by the client.

To assist you in building these types of features, GoW makes it easy to "broadcast" your events over a WebSocket connection.

## Configuration

All of your application's broadcasting configuration is stored in the `config/broadcasting.go` configuration file. 

GoW supports several broadcast drivers out of the box:
- `pusher`: Utilize the Pusher API directly.
- `redis`: Utilize Redis Pub/Sub, which is ideal if you are hosting your own WebSocket server like Laravel Reverb or Soketi.

## Defining Channel Authorization

If your application uses private or presence channels, users must be authorized to join them. You can define your channel authorization callbacks in `routes/channels.go`.

```go
broadcasting.Channel("orders.{orderId}", func(user *models.User, orderId int) bool {
    order := orm.Table("orders").Find(orderId)
    return user.ID == order.CreatorID
})
```

## Broadcasting Events

To broadcast an event, you simply call the `Broadcast` method on the `broadcasting` manager, specifying the channels and the payload.

```go
eventPayload := map[string]any{
    "order_id": 123,
    "status": "shipped",
}

// Broadcast to a public channel
broadcasting.Broadcast([]string{"orders"}, "OrderShipped", eventPayload)

// Broadcast to a private channel
broadcasting.Broadcast([]string{"private-orders.123"}, "OrderShipped", eventPayload)
```

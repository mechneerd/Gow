package broadcasting

// Channel represents an abstract broadcasting channel (e.g. public, private, presence).
type Channel interface {
	Name() string
}

// PublicChannel represents a channel anyone can subscribe to.
type PublicChannel struct {
	name string
}

func NewPublicChannel(name string) *PublicChannel {
	return &PublicChannel{name: name}
}

func (c *PublicChannel) Name() string {
	return c.name
}

// PrivateChannel represents a channel requiring authorization.
type PrivateChannel struct {
	name string
}

func NewPrivateChannel(name string) *PrivateChannel {
	return &PrivateChannel{name: "private-" + name}
}

func (c *PrivateChannel) Name() string {
	return c.name
}

// PresenceChannel represents a channel requiring authorization and tracking users.
type PresenceChannel struct {
	name string
}

func NewPresenceChannel(name string) *PresenceChannel {
	return &PresenceChannel{name: "presence-" + name}
}

func (c *PresenceChannel) Name() string {
	return c.name
}

// Event represents an event that can be broadcasted.
type Event interface {
	BroadcastOn() []Channel
	BroadcastAs() string
	BroadcastWith() map[string]any
}

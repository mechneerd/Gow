package cashier

// Subscription represents a recurring billing subscription.
type Subscription struct {
	ID     string
	UserID string
	Plan   string
	Status string // active, canceled, etc.
}

// Cashier provides billing (Stripe, Paddle, etc.).
type Cashier struct {
	// Add Stripe/Paddle client here in real impl
}

func NewCashier() *Cashier {
	return &Cashier{}
}

func (c *Cashier) NewSubscription(userID, plan string) *Subscription {
	return &Subscription{UserID: userID, Plan: plan, Status: "active"}
}


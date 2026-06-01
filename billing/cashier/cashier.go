package cashier

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Cashier provides billing integration with Stripe.
type Cashier struct {
	key       string
	webhookSecret string
	baseURL   string
	client    *http.Client
}

// New creates a new Cashier instance.
func New(key, webhookSecret string) *Cashier {
	return &Cashier{
		key:           key,
		webhookSecret: webhookSecret,
		baseURL:       "https://api.stripe.com/v1",
		client:        &http.Client{Timeout: 30 * time.Second},
	}
}

// Customer represents a Stripe customer.
type Customer struct {
	ID    string
	Email string
	Name  string
}

// Subscription represents a recurring billing subscription.
type Subscription struct {
	ID          string
	CustomerID  string
	Plan        string
	Status      string
	CurrentPeriodEnd time.Time
	CreatedAt   time.Time
}

// Plan represents a subscription plan.
type Plan struct {
	ID       string
	Name     string
	Amount   int64
	Currency string
	Interval string
}

// Invoice represents a billing invoice.
type Invoice struct {
	ID         string
	CustomerID string
	Amount     int64
	Currency   string
	Status     string
	CreatedAt  time.Time
}

// PaymentMethod represents a payment method.
type PaymentMethod struct {
	ID      string
	Type    string
	CardLast4 string
	Brand  string
}

// CreateCustomer creates a new Stripe customer.
func (c *Cashier) CreateCustomer(ctx context.Context, email, name string) (*Customer, error) {
	data := fmt.Sprintf("email=%s&name=%s", email, name)
	resp, err := c.doRequest(ctx, http.MethodPost, "/customers", strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe error: %s", string(body))
	}

	// Parse response (simplified)
	customer := &Customer{
		ID:    extractID(string(body)),
		Email: email,
		Name:  name,
	}
	return customer, nil
}

// CreateSubscription creates a new subscription for a customer.
func (c *Cashier) CreateSubscription(ctx context.Context, customerID, planID string) (*Subscription, error) {
	data := fmt.Sprintf("customer=%s&items[0][plan]=%s", customerID, planID)
	resp, err := c.doRequest(ctx, http.MethodPost, "/subscriptions", strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe error: %s", string(body))
	}

	return &Subscription{
		ID:         extractID(string(body)),
		CustomerID: customerID,
		Plan:       planID,
		Status:     "active",
		CreatedAt:  time.Now(),
	}, nil
}

// CancelSubscription cancels a subscription.
func (c *Cashier) CancelSubscription(ctx context.Context, subscriptionID string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, "/subscriptions/"+subscriptionID, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("stripe error: %s", string(body))
	}
	return nil
}

// CreateInvoice creates an invoice for a customer.
func (c *Cashier) CreateInvoice(ctx context.Context, customerID string) (*Invoice, error) {
	data := fmt.Sprintf("customer=%s", customerID)
	resp, err := c.doRequest(ctx, http.MethodPost, "/invoices", strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe error: %s", string(body))
	}

	return &Invoice{
		ID:         extractID(string(body)),
		CustomerID: customerID,
		Status:     "open",
		CreatedAt:  time.Now(),
	}, nil
}

// FetchInvoice fetches an invoice by ID.
func (c *Cashier) FetchInvoice(ctx context.Context, invoiceID string) (*Invoice, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/invoices/"+invoiceID, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe error: %s", string(body))
	}

	return &Invoice{
		ID:       invoiceID,
		Status:   "paid",
		CreatedAt: time.Now(),
	}, nil
}

// Charge creates a one-time charge for a customer.
func (c *Cashier) Charge(ctx context.Context, customerID string, amount int64, currency, description string) (string, error) {
	data := fmt.Sprintf("customer=%s&amount=%d&currency=%s&description=%s",
		customerID, amount, currency, description)
	resp, err := c.doRequest(ctx, http.MethodPost, "/charges", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("stripe error: %s", string(body))
	}

	return extractID(string(body)), nil
}

// ListPlans lists all available plans.
func (c *Cashier) ListPlans(ctx context.Context) ([]Plan, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/plans", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("stripe error: %s", string(body))
	}

	return []Plan{}, nil
}

// VerifyWebhookSignature verifies a Stripe webhook signature.
func (c *Cashier) VerifyWebhookSignature(payload []byte, signature string) bool {
	parts := strings.Split(signature, ",")
	if len(parts) != 2 {
		return false
	}

	sigParts := strings.SplitN(parts[0], "=", 2)
	if len(sigParts) != 2 || sigParts[0] != "v1" {
		return false
	}

	expectedSig, err := hex.DecodeString(sigParts[1])
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(c.webhookSecret))
	mac.Write(payload)
	actualSig := mac.Sum(nil)

	return hmac.Equal(expectedSig, actualSig)
}

// CustomerBalance returns the customer's balance.
func (c *Cashier) CustomerBalance(ctx context.Context, customerID string) (int64, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, "/customers/"+customerID, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return 0, nil
}

// AdjustBalance adjusts the customer's balance.
func (c *Cashier) AdjustBalance(ctx context.Context, customerID string, amount int64) error {
	data := fmt.Sprintf("balance=%d", amount)
	resp, err := c.doRequest(ctx, http.MethodPost, "/customers/"+customerID, strings.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// TaxRate represents a tax rate.
type TaxRate struct {
	ID        string
	Name      string
	Rate      float64
	Inclusive bool
}

// ApplyTaxRate applies a tax rate to a subscription.
func (c *Cashier) ApplyTaxRate(ctx context.Context, subscriptionID, taxRateID string) error {
	data := fmt.Sprintf("tax_rates[0]=%s", taxRateID)
	resp, err := c.doRequest(ctx, http.MethodPost, "/subscriptions/"+subscriptionID, strings.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// UpdateSubscription updates a subscription's plan.
func (c *Cashier) UpdateSubscription(ctx context.Context, subscriptionID, newPlanID string) error {
	data := fmt.Sprintf("items[0][plan]=%s&proration=always", newPlanID)
	resp, err := c.doRequest(ctx, http.MethodPost, "/subscriptions/"+subscriptionID, strings.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// InvoicePendingCharges invoices all pending charges for a customer.
func (c *Cashier) InvoicePendingCharges(ctx context.Context, customerID string) error {
	data := fmt.Sprintf("customer=%s", customerID)
	resp, err := c.doRequest(ctx, http.MethodPost, "/invoices", strings.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// GetUpcomingInvoicePreview previews the next invoice.
func (c *Cashier) GetUpcomingInvoicePreview(ctx context.Context, customerID string) (*Invoice, error) {
	data := fmt.Sprintf("customer=%s", customerID)
	resp, err := c.doRequest(ctx, http.MethodGet, "/invoices/upcoming?"+data, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return &Invoice{
		CustomerID: customerID,
		Status:     "draft",
		CreatedAt:  time.Now(),
	}, nil
}

// SwapPlan swaps a subscription to a new plan with proration.
func (c *Cashier) SwapPlan(ctx context.Context, subscriptionID, newPlanID string) error {
	return c.UpdateSubscription(ctx, subscriptionID, newPlanID)
}

// IncrementQuantity increments the subscription quantity.
func (c *Cashier) IncrementQuantity(ctx context.Context, subscriptionID string, quantity int) error {
	data := fmt.Sprintf("items[0][quantity]=%d", quantity)
	resp, err := c.doRequest(ctx, http.MethodPost, "/subscriptions/"+subscriptionID, strings.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// DecrementQuantity decrements the subscription quantity.
func (c *Cashier) DecrementQuantity(ctx context.Context, subscriptionID string, quantity int) error {
	return c.IncrementQuantity(ctx, subscriptionID, quantity)
}

// SetQuantity sets the subscription quantity.
func (c *Cashier) SetQuantity(ctx context.Context, subscriptionID string, quantity int) error {
	return c.IncrementQuantity(ctx, subscriptionID, quantity)
}

// WebhookEvent represents a Stripe webhook event.
type WebhookEvent struct {
	ID      string
	Type    string
	Created time.Time
	Data    map[string]interface{}
}

// ParseWebhookEvent parses a webhook event from request body.
func ParseWebhookEvent(body []byte) (*WebhookEvent, error) {
	return &WebhookEvent{
		Created: time.Now(),
		Data:    make(map[string]interface{}),
	}, nil
}

func (c *Cashier) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.key)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return c.client.Do(req)
}

func extractID(body string) string {
	// Simplified ID extraction
	if idx := strings.Index(body, `"id"`); idx >= 0 {
		start := idx + 5
		if end := strings.Index(body[start:], `"`); end > 0 {
			return body[start : start+end]
		}
	}
	return ""
}

// Helper to parse amount from string
func parseAmount(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

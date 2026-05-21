package session

// CookieDriver implements the session.Store interface by storing payloads in cookies.
// Note: In a real implementation, this requires access to the HTTP Request/Response writers
// to set the cookie. For the architecture of GoW, the Store interface abstracts reading/writing,
// so a Cookie driver needs contextual access or relies on the Manager mapping it correctly.
type CookieDriver struct {
	// Cryptography keys would go here
}

// NewCookieDriver creates a new CookieDriver.
func NewCookieDriver() *CookieDriver {
	return &CookieDriver{}
}

func (d *CookieDriver) Read(id string) (map[string]any, error) {
	// The payload is passed as the "id" in this hacky driver for now,
	// or the middleware parses the cookie and passes the JSON string to Read().
	// For production, the Manager/Middleware handles cookie reading.
	return make(map[string]any), nil
}

func (d *CookieDriver) Write(id string, data map[string]any) error {
	// Handled by the Middleware inspecting the session state
	return nil
}

func (d *CookieDriver) Destroy(id string) error {
	return nil
}

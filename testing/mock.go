package testing

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

// WithoutExceptionHandling runs the given function with exception handling disabled.
func WithoutExceptionHandling(fn func()) {
	fn()
}

// Mock provides a simple mock framework for testing.
type Mock struct {
	calls   []MockCall
	stubs   map[string]MockStub
	mu      sync.RWMutex
}

// MockCall represents a recorded method call.
type MockCall struct {
	Method string
	Args   []any
	Return []any
}

// MockStub represents a stubbed method return value.
type MockStub struct {
	Returns []any
}

// NewMock creates a new Mock instance.
func NewMock() *Mock {
	return &Mock{
		stubs: make(map[string]MockStub),
	}
}

// Stub stubs a method to return specific values.
func (m *Mock) Stub(method string, returns ...any) *Mock {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stubs[method] = MockStub{Returns: returns}
	return m
}

// Call records a method call and returns stubbed values.
func (m *Mock) Call(method string, args ...any) []any {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls = append(m.calls, MockCall{
		Method: method,
		Args:   args,
	})

	if stub, ok := m.stubs[method]; ok {
		return stub.Returns
	}
	return nil
}

// CalledWith checks if a method was called with specific arguments.
func (m *Mock) CalledWith(method string, args ...any) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, call := range m.calls {
		if call.Method == method {
			if len(call.Args) == len(args) {
				match := true
				for i := range args {
					if call.Args[i] != args[i] {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
		}
	}
	return false
}

// CalledTimes returns the number of times a method was called.
func (m *Mock) CalledTimes(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, call := range m.calls {
		if call.Method == method {
			count++
		}
	}
	return count
}

// WasCalled checks if a method was called at all.
func (m *Mock) WasCalled(method string) bool {
	return m.CalledTimes(method) > 0
}

// LastCall returns the last call for a method.
func (m *Mock) LastCall(method string) *MockCall {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for i := len(m.calls) - 1; i >= 0; i-- {
		if m.calls[i].Method == method {
			return &m.calls[i]
		}
	}
	return nil
}

// Reset clears all recorded calls.
func (m *Mock) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = nil
}

// CallHistory returns all recorded calls.
func (m *Mock) CallHistory() []MockCall {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]MockCall, len(m.calls))
	copy(result, m.calls)
	return result
}

// Spy wraps a real implementation and records calls.
type Spy struct {
	*Mock
	realImplementation any
}

// NewSpy creates a new Spy wrapping a real implementation.
func NewSpy(real any) *Spy {
	return &Spy{
		Mock:               NewMock(),
		realImplementation: real,
	}
}

// FakeHTTP creates a fake HTTP server for testing.
func FakeHTTP(handler http.Handler) *httptest.Server {
	return httptest.NewServer(handler)
}

// FakeHTTPHandler returns a handler that records requests.
func FakeHTTPHandler() (*FakeHandler, http.Handler) {
	h := &FakeHandler{}
	return h, h
}

// FakeHandler records HTTP requests for testing.
type FakeHandler struct {
	Requests  []*http.Request
	Responses []FakeResponse
	mu        sync.Mutex
}

// FakeResponse represents a pre-configured response.
type FakeResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// ServeHTTP records the request and returns a configured response.
func (h *FakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.Requests = append(h.Requests, r)
	h.mu.Unlock()

	if len(h.Responses) > 0 {
		resp := h.Responses[0]
		h.mu.Lock()
		h.Responses = h.Responses[1:]
		h.mu.Unlock()

		for key, value := range resp.Headers {
			w.Header().Set(key, value)
		}
		if resp.StatusCode == 0 {
			resp.StatusCode = http.StatusOK
		}
		w.WriteHeader(resp.StatusCode)
		w.Write([]byte(resp.Body))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

// Queue records requests made to the fake handler.
func (h *FakeHandler) Queue(responses ...FakeResponse) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Responses = append(h.Responses, responses...)
}

// LastRequest returns the last recorded request.
func (h *FakeHandler) LastRequest() *http.Request {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.Requests) == 0 {
		return nil
	}
	return h.Requests[len(h.Requests)-1]
}

// RequestCount returns the number of requests made.
func (h *FakeHandler) RequestCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.Requests)
}

// Reset clears all recorded requests and responses.
func (h *FakeHandler) Reset() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Requests = nil
	h.Responses = nil
}

// AssertSent asserts that a request was sent to a specific URL.
func (h *FakeHandler) AssertSent(url string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, req := range h.Requests {
		if req.URL.String() == url {
			return true
		}
	}
	return false
}

// AssertMethod asserts that a specific HTTP method was used.
func (h *FakeHandler) AssertMethod(method string) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, req := range h.Requests {
		if req.Method == method {
			return true
		}
	}
	return false
}

// FakeClock provides time control for testing.
type FakeClock struct {
	current time.Time
	mu      sync.RWMutex
}

// NewFakeClock creates a new FakeClock at the current time.
func NewFakeClock() *FakeClock {
	return &FakeClock{current: time.Now()}
}

// Now returns the fake current time.
func (fc *FakeClock) Now() time.Time {
	fc.mu.RLock()
	defer fc.mu.RUnlock()
	return fc.current
}

// Advance advances the fake clock by the given duration.
func (fc *FakeClock) Advance(d time.Duration) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.current = fc.current.Add(d)
}

// Set sets the fake clock to a specific time.
func (fc *FakeClock) Set(t time.Time) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.current = t
}

// Sub returns the duration from the fake now to the given time.
func (fc *FakeClock) Sub(t time.Time) time.Duration {
	fc.mu.RLock()
	defer fc.mu.RUnlock()
	return fc.current.Sub(t)
}

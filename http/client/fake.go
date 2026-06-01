package client

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

// Fake is a fake HTTP client for testing.
type Fake struct {
	mu       sync.RWMutex
	responses map[string]*FakeResponse
	requests  []*FakeRequest
	calls     int
}

// FakeResponse represents a canned response.
type FakeResponse struct {
	StatusCode int
	Body       any
	Headers    map[string]string
	Error      error
}

// FakeRequest captures a request made to the fake.
type FakeRequest struct {
	Method string
	URL    string
	Body   []byte
	Headers map[string]string
}

// NewFake creates a new fake HTTP client.
func NewFake() *Fake {
	return &Fake{
		responses: make(map[string]*FakeResponse),
		requests:  make([]*FakeRequest, 0),
	}
}

// When sets up a canned response for a URL pattern.
func (f *Fake) When(method, url string, response *FakeResponse) *Fake {
	f.mu.Lock()
	defer f.mu.Unlock()
	key := method + " " + url
	f.responses[key] = response
	return f
}

// WhenGet is a convenience method for GET requests.
func (f *Fake) WhenGet(url string, statusCode int, body any) *Fake {
	return f.When("GET", url, &FakeResponse{
		StatusCode: statusCode,
		Body:       body,
	})
}

// WhenPost is a convenience method for POST requests.
func (f *Fake) WhenPost(url string, statusCode int, body any) *Fake {
	return f.When("POST", url, &FakeResponse{
		StatusCode: statusCode,
		Body:       body,
	})
}

// WhenPut is a convenience method for PUT requests.
func (f *Fake) WhenPut(url string, statusCode int, body any) *Fake {
	return f.When("PUT", url, &FakeResponse{
		StatusCode: statusCode,
		Body:       body,
	})
}

// WhenDelete is a convenience method for DELETE requests.
func (f *Fake) WhenDelete(url string, statusCode int) *Fake {
	return f.When("DELETE", url, &FakeResponse{
		StatusCode: statusCode,
	})
}

// Do executes the fake request and returns the canned response.
func (f *Fake) Do(method, url string, body io.Reader) (*Response, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.calls++

	// Capture request
	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = io.ReadAll(body)
	}

	req := &FakeRequest{
		Method:  method,
		URL:     url,
		Body:    bodyBytes,
		Headers: make(map[string]string),
	}
	f.requests = append(f.requests, req)

	// Find matching response
	key := method + " " + url
	resp, ok := f.responses[key]
	if !ok {
		// Return 404 for unmatched requests
		return &Response{
			Raw: &http.Response{
				StatusCode: http.StatusNotFound,
				Header:     make(http.Header),
			},
			body: []byte(`{"error": "No matching fake response for ` + method + ` ` + url + `"}`),
		}, nil
	}

	if resp.Error != nil {
		return nil, resp.Error
	}

	// Build response body
	var respBody []byte
	if resp.Body != nil {
		switch b := resp.Body.(type) {
		case string:
			respBody = []byte(b)
		case []byte:
			respBody = b
		default:
			respBody, _ = json.Marshal(b)
		}
	}

	// Build response headers
	headers := make(http.Header)
	for k, v := range resp.Headers {
		headers.Set(k, v)
	}

	return &Response{
		Raw: &http.Response{
			StatusCode: resp.StatusCode,
			Header:     headers,
		},
		body: respBody,
	}, nil
}

// GetRequests returns all captured requests.
func (f *Fake) GetRequests() []*FakeRequest {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := make([]*FakeRequest, len(f.requests))
	copy(result, f.requests)
	return result
}

// GetLastRequest returns the last captured request, or nil if none.
func (f *Fake) GetLastRequest() *FakeRequest {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if len(f.requests) == 0 {
		return nil
	}
	return f.requests[len(f.requests)-1]
}

// GetRequestCount returns the number of requests made.
func (f *Fake) GetRequestCount() int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.calls
}

// AssertSent checks if a request was made to the given URL.
func (f *Fake) AssertSent(method, url string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	for _, req := range f.requests {
		if req.Method == method && req.URL == url {
			return true
		}
	}
	return false
}

// AssertNotSent checks that no request was made to the given URL.
func (f *Fake) AssertNotSent(method, url string) bool {
	return !f.AssertSent(method, url)
}

// AssertSentCount checks the exact number of requests made.
func (f *Fake) AssertSentCount(count int) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.calls == count
}

// AssertSentJson checks if a request was made with the given JSON body.
func (f *Fake) AssertSentJson(method, url string, body any) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	expected, _ := json.Marshal(body)

	for _, req := range f.requests {
		if req.Method == method && req.URL == url {
			// Compare JSON bodies
			var expectedMap, actualMap map[string]any
			json.Unmarshal(expected, &expectedMap)
			json.Unmarshal(req.Body, &actualMap)

			if len(expectedMap) == len(actualMap) {
				for k, v := range expectedMap {
					if actualMap[k] != v {
						return false
					}
				}
				return true
			}
		}
	}
	return false
}

// Clear resets all captured requests.
func (f *Fake) Clear() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.requests = make([]*FakeRequest, 0)
	f.calls = 0
}

// Flush removes all configured responses.
func (f *Fake) Flush() {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.responses = make(map[string]*FakeResponse)
}

// Sequence returns a different response for each sequential call to the same URL.
func (f *Fake) Sequence(method, url string, responses ...*FakeResponse) *Fake {
	f.mu.Lock()
	defer f.mu.Unlock()

	for i, resp := range responses {
		key := method + " " + url
		if i > 0 {
			key = method + " " + url
		}
		f.responses[key] = resp
	}
	return f
}
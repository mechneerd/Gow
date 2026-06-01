package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Request2 represents an HTTP request
type Request2 struct {
	Method      string
	URL         string
	Headers     map[string]string
	Body        any
	Timeout     time.Duration
	BasicUser   string
	BasicPass   string
	BearerToken string
	QueryParams map[string]string
}

// Response2 represents an HTTP response
type Response2 struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	Request    *http.Request
}

// Client2 provides HTTP client functionality
type Client2 struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string
}

// NewClient creates a new HTTP client
func NewClient() *Client2 {
	return &Client2{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

// NewWithTimeout creates a new HTTP client with a custom timeout
func NewWithTimeout(timeout time.Duration) *Client2 {
	return &Client2{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		headers: make(map[string]string),
	}
}

// SetBaseURL sets the base URL for all requests
func (c *Client2) SetBaseURL(baseURL string) *Client2 {
	c.baseURL = baseURL
	return c
}

// SetHeader sets a default header for all requests
func (c *Client2) SetHeader(key, value string) *Client2 {
	c.headers[key] = value
	return c
}

// SetHeaders sets multiple default headers
func (c *Client2) SetHeaders(headers map[string]string) *Client2 {
	for key, value := range headers {
		c.headers[key] = value
	}
	return c
}

// Get sends a GET request
func (c *Client2) Get(url string, headers ...map[string]string) (*Response2, error) {
	return c.Do(http.MethodGet, url, nil, headers...)
}

// Post sends a POST request
func (c *Client2) Post(url string, body any, headers ...map[string]string) (*Response2, error) {
	return c.Do(http.MethodPost, url, body, headers...)
}

// Put sends a PUT request
func (c *Client2) Put(url string, body any, headers ...map[string]string) (*Response2, error) {
	return c.Do(http.MethodPut, url, body, headers...)
}

// Patch sends a PATCH request
func (c *Client2) Patch(url string, body any, headers ...map[string]string) (*Response2, error) {
	return c.Do(http.MethodPatch, url, body, headers...)
}

// Delete sends a DELETE request
func (c *Client2) Delete(url string, headers ...map[string]string) (*Response2, error) {
	return c.Do(http.MethodDelete, url, nil, headers...)
}

// Head sends a HEAD request
func (c *Client2) Head(url string, headers ...map[string]string) (*Response2, error) {
	return c.Do(http.MethodHead, url, nil, headers...)
}

// Options sends an OPTIONS request
func (c *Client2) Options(url string, headers ...map[string]string) (*Response2, error) {
	return c.Do(http.MethodOptions, url, nil, headers...)
}

// Do sends an HTTP request
func (c *Client2) Do(method, url string, body any, headers ...map[string]string) (*Response2, error) {
	// Build full URL
	fullURL := url
	if c.baseURL != "" && !strings.HasPrefix(url, "http") {
		fullURL = strings.TrimRight(c.baseURL, "/") + "/" + strings.TrimLeft(url, "/")
	}

	// Build request body
	var bodyReader io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
	}

	// Create request
	req, err := http.NewRequest(method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Set request headers
	if len(headers) > 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}

	// Set content type for JSON body
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return &Response2{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
		Request:    req,
	}, nil
}

// WithBearerToken sets bearer token authentication for the request
func (c *Client2) WithBearerToken(token string) *Client2 {
	c.headers["Authorization"] = "Bearer " + token
	return c
}

// WithTimeout sets the timeout for the client
func (c *Client2) WithTimeout(timeout time.Duration) *Client2 {
	c.httpClient.Timeout = timeout
	return c
}

// JSON unmarshals the response body into the provided interface
func (r *Response2) JSON(v any) error {
	return json.Unmarshal(r.Body, v)
}

// String returns the response body as a string
func (r *Response2) String() string {
	return string(r.Body)
}

// IsSuccessful returns true if the status code is 2xx
func (r *Response2) IsSuccessful() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsRedirect returns true if the status code is 3xx
func (r *Response2) IsRedirect() bool {
	return r.StatusCode >= 300 && r.StatusCode < 400
}

// IsClientError returns true if the status code is 4xx
func (r *Response2) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if the status code is 5xx
func (r *Response2) IsServerError() bool {
	return r.StatusCode >= 500
}

// Header2 returns a header value
func (r *Response2) Header2(key string) string {
	return r.Headers.Get(key)
}

// Retry sends a request with retries
func (c *Client2) Retry(method, url string, body any, maxRetries int, delay time.Duration, headers ...map[string]string) (*Response2, error) {
	var lastErr error
	for i := 0; i <= maxRetries; i++ {
		resp, err := c.Do(method, url, body, headers...)
		if err == nil && resp.IsSuccessful() {
			return resp, nil
		}
		lastErr = err
		if i < maxRetries {
			time.Sleep(delay * time.Duration(i+1))
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}

// Async sends a request asynchronously
func (c *Client2) Async(method, url string, body any, result chan<- *Response2, errChan chan<- error, headers ...map[string]string) {
	go func() {
		resp, err := c.Do(method, url, body, headers...)
		if err != nil {
			errChan <- err
			return
		}
		result <- resp
	}()
}

// Pool represents a pool of HTTP clients
type Pool struct {
	clients chan *Client2
	size    int
}

// NewPool creates a new pool of HTTP clients
func NewPool(size int) *Pool {
	clients := make(chan *Client2, size)
	for i := 0; i < size; i++ {
		clients <- NewClient()
	}
	return &Pool{
		clients: clients,
		size:    size,
	}
}

// Get retrieves a client from the pool
func (p *Pool) Get() *Client2 {
	return <-p.clients
}

// Put returns a client to the pool
func (p *Pool) Put(client *Client2) {
	p.clients <- client
}

// Size returns the pool size
func (p *Pool) Size() int {
	return p.size
}

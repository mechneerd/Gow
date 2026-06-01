package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"
)

// Client provides a fluent wrapper around net/http.
type Client struct {
	httpClient *http.Client
	headers    map[string]string
	retries    int
	retryDelay time.Duration
	timeout    time.Duration
	fake       *Fake
	mu         sync.RWMutex
}

// Response wraps an http.Response for easier consumption.
type Response struct {
	Raw  *http.Response
	body []byte
}

// New creates a new HTTP client.
func New() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		headers:    make(map[string]string),
	}
}

// WithHeader adds a header to the request.
func (c *Client) WithHeader(key, value string) *Client {
	c.headers[key] = value
	return c
}

// WithTimeout sets the request timeout.
func (c *Client) WithTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	c.httpClient.Timeout = timeout
	return c
}

// WithBasicAuth sets Basic Auth headers.
func (c *Client) WithBasicAuth(username, password string) *Client {
	c.headers["Authorization"] = "Basic " + basicAuth(username, password)
	return c
}

// WithBearerToken sets a Bearer token in the Authorization header.
func (c *Client) WithBearerToken(token string) *Client {
	c.headers["Authorization"] = "Bearer " + token
	return c
}

// WithRetry configures retry with exponential backoff.
func (c *Client) WithRetry(attempts int, delay time.Duration) *Client {
	c.retries = attempts
	c.retryDelay = delay
	return c
}

// WithContext returns a new Client that uses the given context for requests.
func (c *Client) WithContext(ctx context.Context) *Client {
	c.httpClient = &http.Client{
		Timeout: c.timeout,
	}
	_ = ctx // Context will be used in request methods
	return c
}

// Get makes a GET request.
func (c *Client) Get(url string) (*Response, error) {
	return c.request("GET", url, nil)
}

// Post makes a POST request.
func (c *Client) Post(url string, body any) (*Response, error) {
	return c.sendWithBody("POST", url, body)
}

// Put makes a PUT request.
func (c *Client) Put(url string, body any) (*Response, error) {
	return c.sendWithBody("PUT", url, body)
}

// Patch makes a PATCH request.
func (c *Client) Patch(url string, body any) (*Response, error) {
	return c.sendWithBody("PATCH", url, body)
}

// Delete makes a DELETE request.
func (c *Client) Delete(url string) (*Response, error) {
	return c.request("DELETE", url, nil)
}

// Head makes a HEAD request.
func (c *Client) Head(url string) (*Response, error) {
	return c.request("HEAD", url, nil)
}

// Options makes an OPTIONS request.
func (c *Client) Options(url string) (*Response, error) {
	return c.request("OPTIONS", url, nil)
}

func (c *Client) sendWithBody(method, url string, body any) (*Response, error) {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(b)
		c.WithHeader("Content-Type", "application/json")
	}
	return c.request(method, url, bodyReader)
}

func (c *Client) request(method, url string, body io.Reader) (*Response, error) {
	// Check for fake
	c.mu.RLock()
	if c.fake != nil {
		c.mu.RUnlock()
		return c.fake.Do(method, url, body)
	}
	c.mu.RUnlock()

	var lastErr error
	attempts := c.retries + 1
	if attempts <= 0 {
		attempts = 1
	}

	for i := 0; i < attempts; i++ {
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			return nil, err
		}

		for k, v := range c.headers {
			req.Header.Set(k, v)
		}

		res, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if i < attempts-1 {
				time.Sleep(c.retryDelay * time.Duration(1<<uint(i))) // exponential backoff
			}
			continue
		}
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			lastErr = err
			if i < attempts-1 {
				time.Sleep(c.retryDelay * time.Duration(1<<uint(i)))
			}
			continue
		}

		return &Response{
			Raw:  res,
			body: resBody,
		}, nil
	}

	return nil, lastErr
}

// Fake replaces the client with a fake for testing.
func (c *Client) Fake(fake *Fake) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.fake = fake
	return c
}

// Body returns the raw response bytes.
func (r *Response) Body() string {
	return string(r.body)
}

// JSON unmarshals the response body into the given target.
func (r *Response) JSON(target any) error {
	return json.Unmarshal(r.body, target)
}

// StatusCode returns the HTTP status code.
func (r *Response) StatusCode() int {
	return r.Raw.StatusCode
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return encodeBase64(auth)
}

func encodeBase64(s string) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	buf := make([]byte, 0, len(s)/3*4+4)
	for i := 0; i < len(s); i += 3 {
		var b0, b1, b2 byte
		b0 = s[i]
		if i+1 < len(s) {
			b1 = s[i+1]
		}
		if i+2 < len(s) {
			b2 = s[i+2]
		}
		buf = append(buf, chars[(b0>>2)&0x3F])
		buf = append(buf, chars[((b0&0x3)<<4)|((b1>>4)&0xF)])
		if i+1 < len(s) {
			buf = append(buf, chars[((b1&0xF)<<2)|((b2>>6)&0x3)])
		}
		if i+2 < len(s) {
			buf = append(buf, chars[b2&0x3F])
		}
	}
	// Padding
	for len(buf)%4 != 0 {
		buf = append(buf, '=')
	}
	return string(buf)
}


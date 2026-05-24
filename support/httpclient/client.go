package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Response wraps an HTTP response with many convenient helpers.
type Response struct {
	*http.Response
	bodyBytes []byte
}

// Body returns the raw response body (cached after first read).
func (r *Response) Body() []byte {
	if r.bodyBytes == nil {
		r.bodyBytes, _ = io.ReadAll(r.Response.Body)
		r.Response.Body.Close()
	}
	return r.bodyBytes
}

// JSON parses the body as JSON into the target.
func (r *Response) JSON(target any) error {
	return json.Unmarshal(r.Body(), target)
}

// StatusCode returns the HTTP status code.
func (r *Response) StatusCode() int {
	return r.Response.StatusCode
}

// Ok returns true if the response status is 2xx.
func (r *Response) Ok() bool {
	return r.StatusCode() >= 200 && r.StatusCode() < 300
}

// Successful is an alias for Ok().
func (r *Response) Successful() bool {
	return r.Ok()
}

// Failed returns true if the response is 4xx or 5xx.
func (r *Response) Failed() bool {
	return r.StatusCode() >= 400
}

// ServerError returns true for 5xx responses.
func (r *Response) ServerError() bool {
	return r.StatusCode() >= 500
}

// ClientError returns true for 4xx responses.
func (r *Response) ClientError() bool {
	return r.StatusCode() >= 400 && r.StatusCode() < 500
}

// Header returns a specific response header.
func (r *Response) Header(key string) string {
	return r.Response.Header.Get(key)
}

// Client wraps a standard HTTP client and provides fluent request building.
type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
	timeout time.Duration
	retries int
	query   map[string]string
}

// Global fakes registry for testing.
var (
	fakesEnabled bool
	fakeResponses map[string]*Response
)

// Fake globally enables faking and registers mock responses.
// URL path pattern -> Mock Response
func Fake(mocks map[string]*Response) {
	fakesEnabled = true
	fakeResponses = mocks
}

// NewClient creates a new fluent HTTP client.
func NewClient() *Client {
	return &Client{
		client:  &http.Client{Timeout: 30 * time.Second},
		headers: make(map[string]string),
	}
}

// BaseURL sets the base URL for the client.
func (c *Client) BaseURL(url string) *Client {
	c.baseURL = url
	return c
}

// WithHeader adds a header to the client.
func (c *Client) WithHeader(key, value string) *Client {
	c.headers[key] = value
	return c
}

// WithHeaders sets multiple headers at once.
func (c *Client) WithHeaders(headers map[string]string) *Client {
	for k, v := range headers {
		c.headers[k] = v
	}
	return c
}

// WithQuery adds query parameters.
func (c *Client) WithQuery(key, value string) *Client {
	if c.query == nil {
		c.query = make(map[string]string)
	}
	c.query[key] = value
	return c
}

// WithToken sets a Bearer token.
func (c *Client) WithToken(token string) *Client {
	return c.WithHeader("Authorization", "Bearer "+token)
}

// Retry sets the number of retries on failure (5xx).
func (c *Client) Retry(times int) *Client {
	c.retries = times
	return c
}

// Get executes a GET request.
func (c *Client) Get(path string) (*Response, error) {
	return c.request(http.MethodGet, path, nil)
}

// Post executes a POST request with JSON payload.
func (c *Client) Post(path string, data any) (*Response, error) {
	var body io.Reader
	if data != nil {
		b, _ := json.Marshal(data)
		body = bytes.NewReader(b)
		c.WithHeader("Content-Type", "application/json")
	}
	return c.request(http.MethodPost, path, body)
}

// PostForm sends a form-urlencoded POST request.
func (c *Client) PostForm(path string, data map[string]string) (*Response, error) {
	form := make(url.Values)
	for k, v := range data {
		form.Set(k, v)
	}
	c.WithHeader("Content-Type", "application/x-www-form-urlencoded")
	return c.request(http.MethodPost, path, strings.NewReader(form.Encode()))
}

// Put executes a PUT request.
func (c *Client) Put(path string, data any) (*Response, error) {
	var body io.Reader
	if data != nil {
		b, _ := json.Marshal(data)
		body = bytes.NewReader(b)
		c.WithHeader("Content-Type", "application/json")
	}
	return c.request(http.MethodPut, path, body)
}

// Patch executes a PATCH request.
func (c *Client) Patch(path string, data any) (*Response, error) {
	var body io.Reader
	if data != nil {
		b, _ := json.Marshal(data)
		body = bytes.NewReader(b)
		c.WithHeader("Content-Type", "application/json")
	}
	return c.request(http.MethodPatch, path, body)
}

// Delete executes a DELETE request.
func (c *Client) Delete(path string) (*Response, error) {
	return c.request(http.MethodDelete, path, nil)
}

// Head executes a HEAD request.
func (c *Client) Head(path string) (*Response, error) {
	return c.request(http.MethodHead, path, nil)
}

func (c *Client) request(method, path string, body io.Reader) (*Response, error) {
	fullURL := c.baseURL + path

	// Append query parameters if any
	if len(c.query) > 0 {
		u, err := url.Parse(fullURL)
		if err == nil {
			q := u.Query()
			for k, v := range c.query {
				q.Set(k, v)
			}
			u.RawQuery = q.Encode()
			fullURL = u.String()
		}
	}

	// Global faking support (basic)
	if fakesEnabled {
		if resp, ok := fakeResponses[path]; ok {
			return resp, nil
		}
		if resp, ok := fakeResponses["*"]; ok {
			return resp, nil
		}
	}

	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	var resp *http.Response
	var reqErr error

	for i := 0; i <= c.retries; i++ {
		resp, reqErr = c.client.Do(req)
		if reqErr == nil && resp.StatusCode < 500 {
			break
		}
		time.Sleep(100 * time.Millisecond * time.Duration(i+1))
	}

	if reqErr != nil {
		return nil, reqErr
	}

	return &Response{Response: resp}, nil
}

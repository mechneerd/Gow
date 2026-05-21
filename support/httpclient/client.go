package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Response wraps an HTTP response.
type Response struct {
	*http.Response
	bodyBytes []byte
}

// Body returns the raw response body.
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

// Client wraps a standard HTTP client and provides fluent request building and faking capabilities.
type Client struct {
	client  *http.Client
	baseURL string
	headers map[string]string
	timeout time.Duration
	retries int
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

// Retry sets the number of retries.
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

func (c *Client) request(method, path string, body io.Reader) (*Response, error) {
	fullURL := c.baseURL + path

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

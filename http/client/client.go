package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Client provides a fluent wrapper around net/http.
type Client struct {
	httpClient *http.Client
	headers    map[string]string
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
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Raw:  res,
		body: resBody,
	}, nil
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


package testing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCase wraps the standard testing.T to provide fluent, Laravel-style assertions.
type TestCase struct {
	T       *testing.T
	handler http.Handler
}

// NewTestCase creates a new TestCase instance.
func NewTestCase(t *testing.T, handler http.Handler) *TestCase {
	return &TestCase{
		T:       t,
		handler: handler,
	}
}

// TestResponse wraps the httptest.ResponseRecorder for fluent assertions.
type TestResponse struct {
	T        *testing.T
	Recorder *httptest.ResponseRecorder
}

// Get simulates an HTTP GET request.
func (tc *TestCase) Get(url string) *TestResponse {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	w := httptest.NewRecorder()
	tc.handler.ServeHTTP(w, req)
	return &TestResponse{T: tc.T, Recorder: w}
}

// AssertStatus asserts the response HTTP status code.
func (tr *TestResponse) AssertStatus(code int) *TestResponse {
	assert.Equal(tr.T, code, tr.Recorder.Code, "Expected status code to match")
	return tr
}

// AssertOk asserts the response status is 200 OK.
func (tr *TestResponse) AssertOk() *TestResponse {
	return tr.AssertStatus(http.StatusOK)
}

// AssertJson asserts the response JSON matches the given map.
func (tr *TestResponse) AssertJson(expected map[string]any) *TestResponse {
	var actual map[string]any
	err := json.Unmarshal(tr.Recorder.Body.Bytes(), &actual)
	assert.NoError(tr.T, err, "Failed to parse JSON response")
	assert.Equal(tr.T, expected, actual, "JSON response did not match expected map")
	return tr
}

// AssertSee asserts the response body contains the given text.
func (tr *TestResponse) AssertSee(text string) *TestResponse {
	assert.Contains(tr.T, tr.Recorder.Body.String(), text, "Response did not contain expected text")
	return tr
}

// AssertHeader asserts the response contains a specific header value.
func (tr *TestResponse) AssertHeader(key, value string) *TestResponse {
	assert.Equal(tr.T, value, tr.Recorder.Header().Get(key), "Header did not match")
	return tr
}

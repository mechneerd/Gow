package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gow/database/orm"
	"gow/database/query"
	"gow/foundation"
	"gow/routing"
)

// TestCase provides a fluent testing API similar to Laravel.
type TestCase struct {
	*testing.T
	App      *foundation.Application
	DB       *orm.DB
	Router   *routing.Router
	client   *http.Client
	authUser any
}

// NewTestCase initializes a test environment.
func NewTestCase(t *testing.T, app *foundation.Application, db *orm.DB, router *routing.Router) *TestCase {
	return &TestCase{
		T:      t,
		App:    app,
		DB:     db,
		Router: router,
		client: &http.Client{},
	}
}

// TestResponse wraps the httptest.ResponseRecorder for fluent assertions.
type TestResponse struct {
	T        *testing.T
	Recorder *httptest.ResponseRecorder
}

// ActingAs sets the authenticated user for subsequent test requests.
func (tc *TestCase) ActingAs(user any) *TestCase {
	tc.authUser = user
	return tc
}

// Get dispatch a GET request to the application.
func (tc *TestCase) Get(uri string) *TestResponse {
	req := httptest.NewRequest(http.MethodGet, uri, nil)
	if tc.authUser != nil {
		req.Header.Set("X-Test-Auth-User", "1")
	}
	w := httptest.NewRecorder()
	
	tc.Router.ServeHTTP(w, req)
	return &TestResponse{T: tc.T, Recorder: w}
}

// Post dispatch a POST request to the application.
func (tc *TestCase) Post(uri string, body io.Reader) *TestResponse {
	req := httptest.NewRequest(http.MethodPost, uri, body)
	if tc.authUser != nil {
		req.Header.Set("X-Test-Auth-User", "1")
	}
	w := httptest.NewRecorder()
	
	tc.Router.ServeHTTP(w, req)
	return &TestResponse{T: tc.T, Recorder: w}
}

// Upload simulates a file upload (multipart).
func (tc *TestCase) Upload(url, fieldName, filename, content string) *TestResponse {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile(fieldName, filename)
	part.Write([]byte(content))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if tc.authUser != nil {
		req.Header.Set("X-Test-Auth-User", "1")
	}
	w := httptest.NewRecorder()
	tc.Router.ServeHTTP(w, req)
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

// AssertJson asserts the response contains the given JSON key/value.
func (tr *TestResponse) AssertJson(key string, value any) *TestResponse {
	var data map[string]any
	json.Unmarshal(tr.Recorder.Body.Bytes(), &data)
	assert.Equal(tr.T, value, data[key], "JSON assertion failed for key: "+key)
	return tr
}

// AssertJsonStructure asserts top-level keys exist.
func (tr *TestResponse) AssertJsonStructure(keys ...string) *TestResponse {
	var data map[string]any
	json.Unmarshal(tr.Recorder.Body.Bytes(), &data)
	for _, k := range keys {
		_, exists := data[k]
		assert.True(tr.T, exists, "Missing JSON key: "+k)
	}
	return tr
}

// AssertJsonMap asserts the response JSON matches the given map.
func (tr *TestResponse) AssertJsonMap(expected map[string]any) *TestResponse {
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

// AssertDatabaseHas asserts that a database table contains a row matching the given constraints.
func (tc *TestCase) AssertDatabaseHas(table string, conditions map[string]any) {
	tc.Helper()
	builder := query.NewBuilder(tc.DB.RawDB(), tc.DB.Dialect())
	builder.Table(table)
	
	for k, v := range conditions {
		builder.Where(k, "=", v)
	}
	
	count, err := builder.Count("*")
	if err != nil {
		tc.Fatalf("Error querying database: %v", err)
	}
	if count == 0 {
		tc.Errorf("Failed asserting that table [%s] has row matching %v", table, conditions)
	}
}

// AssertDatabaseMissing asserts that a database table does NOT contain a row matching the given constraints.
func (tc *TestCase) AssertDatabaseMissing(table string, conditions map[string]any) {
	tc.Helper()
	builder := query.NewBuilder(tc.DB.RawDB(), tc.DB.Dialect())
	builder.Table(table)
	
	for k, v := range conditions {
		builder.Where(k, "=", v)
	}
	
	count, err := builder.Count("*")
	if err != nil {
		tc.Fatalf("Error querying database: %v", err)
	}
	if count > 0 {
		tc.Errorf("Failed asserting that table [%s] is missing row matching %v", table, conditions)
	}
}

// AssertDatabaseCount asserts that a table has exactly `expected` number of rows.
func (tc *TestCase) AssertDatabaseCount(table string, expected int) {
	tc.Helper()
	builder := query.NewBuilder(tc.DB.RawDB(), tc.DB.Dialect())
	builder.Table(table)

	count, err := builder.Count("*")
	if err != nil {
		tc.Fatalf("Error counting rows in table [%s]: %v", table, err)
	}
	if count != expected {
		tc.Errorf("Failed asserting that table [%s] has %d rows. Found %d", table, expected, count)
	}
}

// AssertDatabaseHasNoRecords asserts that a table is completely empty.
func (tc *TestCase) AssertDatabaseHasNoRecords(table string) {
	tc.AssertDatabaseCount(table, 0)
}

// AssertDatabaseHasExactly asserts that a table contains **exactly one** row matching the conditions.
func (tc *TestCase) AssertDatabaseHasExactly(table string, conditions map[string]any) {
	tc.Helper()
	builder := query.NewBuilder(tc.DB.RawDB(), tc.DB.Dialect())
	builder.Table(table)
	
	for k, v := range conditions {
		builder.Where(k, "=", v)
	}
	
	count, err := builder.Count("*")
	if err != nil {
		tc.Fatalf("Error querying database: %v", err)
	}
	if count != 1 {
		tc.Errorf("Failed asserting that table [%s] has exactly 1 row matching %v. Found %d", table, conditions, count)
	}
}

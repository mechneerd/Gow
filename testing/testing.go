package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/mechneerd/gow/database/orm"
	"github.com/mechneerd/gow/database/query"
	"github.com/mechneerd/gow/foundation"
	"github.com/mechneerd/gow/routing"
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

// AssertViewIs asserts the response was rendered with the given view name.
func (tr *TestResponse) AssertViewIs(view string) *TestResponse {
	tr.T.Helper()
	// View name is typically set in the response context or header
	viewHeader := tr.Recorder.Header().Get("X-View")
	if viewHeader == "" {
		tr.T.Errorf("Expected view [%s], but no view was set", view)
	} else if viewHeader != view {
		tr.T.Errorf("Expected view [%s], got [%s]", view, viewHeader)
	}
	return tr
}

// AssertViewHas asserts the view data contains the given key.
func (tr *TestResponse) AssertViewHas(key string) *TestResponse {
	tr.T.Helper()
	viewData := tr.Recorder.Header().Get("X-View-Data-" + key)
	if viewData == "" {
		tr.T.Errorf("Expected view data to contain key [%s]", key)
	}
	return tr
}

// AssertViewHasValue asserts the view data contains the given key with the expected value.
func (tr *TestResponse) AssertViewHasValue(key string, value string) *TestResponse {
	tr.T.Helper()
	viewData := tr.Recorder.Header().Get("X-View-Data-" + key)
	if viewData != value {
		tr.T.Errorf("Expected view data [%s] to be [%s], got [%s]", key, value, viewData)
	}
	return tr
}

// AssertCookie asserts the response contains a cookie with the given name.
func (tr *TestResponse) AssertCookie(name string) *TestResponse {
	tr.T.Helper()
	for _, cookie := range tr.Recorder.Result().Cookies() {
		if cookie.Name == name {
			return tr
		}
	}
	tr.T.Errorf("Expected cookie [%s] not found", name)
	return tr
}

// AssertCookieValue asserts the response contains a cookie with the given value.
func (tr *TestResponse) AssertCookieValue(name, value string) *TestResponse {
	tr.T.Helper()
	for _, cookie := range tr.Recorder.Result().Cookies() {
		if cookie.Name == name && cookie.Value == value {
			return tr
		}
	}
	tr.T.Errorf("Expected cookie [%s] with value [%s] not found", name, value)
	return tr
}

// AssertDontSee asserts the response body does NOT contain the given text.
func (tr *TestResponse) AssertDontSee(text string) *TestResponse {
	tr.T.Helper()
	assert.NotContains(tr.T, tr.Recorder.Body.String(), text, "Response should not contain text")
	return tr
}

// AssertSeeText asserts the response body contains the given text (as plain text).
func (tr *TestResponse) AssertSeeText(text string) *TestResponse {
	tr.T.Helper()
	assert.Contains(tr.T, tr.Recorder.Body.String(), text, "Response did not contain expected text")
	return tr
}

// AssertJsonCount asserts the response JSON array has the given count.
func (tr *TestResponse) AssertJsonCount(key string, count int) *TestResponse {
	tr.T.Helper()
	var data map[string]any
	json.Unmarshal(tr.Recorder.Body.Bytes(), &data)

	if arr, ok := data[key].([]any); ok {
		assert.Equal(tr.T, count, len(arr), "JSON array count mismatch")
	} else {
		tr.T.Errorf("Expected key [%s] to be an array", key)
	}
	return tr
}

// AssertSuccessful asserts the response status is 2xx.
func (tr *TestResponse) AssertSuccessful() *TestResponse {
	tr.T.Helper()
	assert.True(tr.T, tr.Recorder.Code >= 200 && tr.Recorder.Code < 300,
		"Expected successful status code, got %d", tr.Recorder.Code)
	return tr
}

// AssertClientError asserts the response status is 4xx.
func (tr *TestResponse) AssertClientError() *TestResponse {
	tr.T.Helper()
	assert.True(tr.T, tr.Recorder.Code >= 400 && tr.Recorder.Code < 500,
		"Expected client error status code, got %d", tr.Recorder.Code)
	return tr
}

// AssertServerError asserts the response status is 5xx.
func (tr *TestResponse) AssertServerError() *TestResponse {
	tr.T.Helper()
	assert.True(tr.T, tr.Recorder.Code >= 500,
		"Expected server error status code, got %d", tr.Recorder.Code)
	return tr
}

// AssertDatabaseHas asserts that a database table contains a row matching the given constraints.
func (tc *TestCase) AssertDatabaseHas(table string, conditions map[string]any) {
	tc.Helper()
	d, err := tc.DB.Dialect()
	if err != nil {
		tc.Fatalf("Dialect not configured: %v", err)
	}
	builder := query.NewBuilder(tc.DB.RawDB(), d)
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
	d, err := tc.DB.Dialect()
	if err != nil {
		tc.Fatalf("Dialect not configured: %v", err)
	}
	builder := query.NewBuilder(tc.DB.RawDB(), d)
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
	d, err := tc.DB.Dialect()
	if err != nil {
		tc.Fatalf("Dialect not configured: %v", err)
	}
	builder := query.NewBuilder(tc.DB.RawDB(), d)
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
	d, err := tc.DB.Dialect()
	if err != nil {
		tc.Fatalf("Dialect not configured: %v", err)
	}
	builder := query.NewBuilder(tc.DB.RawDB(), d)
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

// AssertSoftDeleted asserts that a row has been soft-deleted (deleted_at is not null).
func (tc *TestCase) AssertSoftDeleted(table string, conditions map[string]any) {
	tc.Helper()
	d, err := tc.DB.Dialect()
	if err != nil {
		tc.Fatalf("Dialect not configured: %v", err)
	}
	builder := query.NewBuilder(tc.DB.RawDB(), d)
	builder.Table(table)
	
	for k, v := range conditions {
		builder.Where(k, "=", v)
	}
	builder.Where("deleted_at", "IS NOT", nil)
	
	count, err := builder.Count("*")
	if err != nil {
		tc.Fatalf("Error querying database: %v", err)
	}
	if count == 0 {
		tc.Errorf("Failed asserting that table [%s] has soft-deleted row matching %v", table, conditions)
	}
}

// AssertNotSoftDeleted asserts that a row has NOT been soft-deleted.
func (tc *TestCase) AssertNotSoftDeleted(table string, conditions map[string]any) {
	tc.Helper()
	d, err := tc.DB.Dialect()
	if err != nil {
		tc.Fatalf("Dialect not configured: %v", err)
	}
	builder := query.NewBuilder(tc.DB.RawDB(), d)
	builder.Table(table)
	
	for k, v := range conditions {
		builder.Where(k, "=", v)
	}
	builder.Where("deleted_at", "IS", nil)
	
	count, err := builder.Count("*")
	if err != nil {
		tc.Fatalf("Error querying database: %v", err)
	}
	if count == 0 {
		tc.Errorf("Failed asserting that table [%s] has non-soft-deleted row matching %v", table, conditions)
	}
}

// AssertDatabaseTable asserts that a database table exists.
func (tc *TestCase) AssertDatabaseTable(table string) {
	tc.Helper()
	d, err := tc.DB.Dialect()
	if err != nil {
		tc.Fatalf("Dialect not configured: %v", err)
	}
	builder := query.NewBuilder(tc.DB.RawDB(), d)
	builder.Table(table)
	
	_, err = builder.Count("*")
	if err != nil {
		tc.Errorf("Failed asserting that table [%s] exists: %v", table, err)
	}
}

// AssertDatabaseHasColumns asserts that a table has the given columns.
func (tc *TestCase) AssertDatabaseHasColumns(table string, columns ...string) {
	tc.Helper()
	// This is a basic check - in production we'd query information_schema
	// For now, just check that we can query the table
	d, err := tc.DB.Dialect()
	if err != nil {
		tc.Fatalf("Dialect not configured: %v", err)
	}
	builder := query.NewBuilder(tc.DB.RawDB(), d)
	builder.Table(table)
	
	// Try to select the columns
	for _, col := range columns {
		builder.Select(col)
	}
	
	_, err = builder.Get()
	if err != nil {
		tc.Errorf("Failed asserting that table [%s] has columns %v: %v", table, columns, err)
	}
}

// ArtisanTestCase provides testing utilities for artisan commands.
type ArtisanTestCase struct {
	*testing.T
	output     bytes.Buffer
	args       []string
	exitCode   int
}

// NewArtisanTestCase creates a new artisan test case.
func NewArtisanTestCase(t *testing.T, args ...string) *ArtisanTestCase {
	return &ArtisanTestCase{
		T:    t,
		args: args,
	}
}

// Args sets the command arguments.
func (atc *ArtisanTestCase) Args(args ...string) *ArtisanTestCase {
	atc.args = args
	return atc
}

// AssertExitCode asserts the command exited with the given code.
func (atc *ArtisanTestCase) AssertExitCode(code int) *ArtisanTestCase {
	atc.Helper()
	if atc.exitCode != code {
		atc.Errorf("Expected exit code %d, got %d", code, atc.exitCode)
	}
	return atc
}

// AssertSuccessful asserts the command was successful (exit code 0).
func (atc *ArtisanTestCase) AssertSuccessful() *ArtisanTestCase {
	atc.Helper()
	return atc.AssertExitCode(0)
}

// AssertOutputContains asserts the output contains the given string.
func (atc *ArtisanTestCase) AssertOutputContains(str string) *ArtisanTestCase {
	atc.Helper()
	output := atc.output.String()
	if !strings.Contains(output, str) {
		atc.Errorf("Expected output to contain [%s], got [%s]", str, output)
	}
	return atc
}

// AssertOutputNotContains asserts the output does not contain the given string.
func (atc *ArtisanTestCase) AssertOutputNotContains(str string) *ArtisanTestCase {
	atc.Helper()
	output := atc.output.String()
	if strings.Contains(output, str) {
		atc.Errorf("Expected output not to contain [%s], got [%s]", str, output)
	}
	return atc
}

// AssertOutputIs asserts the output is exactly the given string.
func (atc *ArtisanTestCase) AssertOutputIs(expected string) *ArtisanTestCase {
	atc.Helper()
	output := atc.output.String()
	if output != expected {
		atc.Errorf("Expected output [%s], got [%s]", expected, output)
	}
	return atc
}

// Output returns the command output.
func (atc *ArtisanTestCase) Output() string {
	return atc.output.String()
}

// ParallelTestCase provides parallel test execution support.
type ParallelTestCase struct {
	*testing.T
	workers int
}

// NewParallelTestCase creates a new parallel test case.
func NewParallelTestCase(t *testing.T, workers int) *ParallelTestCase {
	if workers <= 0 {
		workers = 4
	}
	return &ParallelTestCase{
		T:       t,
		workers: workers,
	}
}

// Run runs the given function in parallel with multiple workers.
func (ptc *ParallelTestCase) Run(fn func(workerID int)) {
	ptc.Helper()

	if ptc.workers <= 1 {
		fn(0)
		return
	}

	done := make(chan struct{}, ptc.workers)

	for i := 0; i < ptc.workers; i++ {
		go func(workerID int) {
			defer func() {
				if r := recover(); r != nil {
					ptc.Errorf("Worker %d panicked: %v", workerID, r)
				}
				done <- struct{}{}
			}()
			fn(workerID)
		}(i)
	}

	for i := 0; i < ptc.workers; i++ {
		<-done
	}
}

// RunParallel runs multiple test functions in parallel.
func (ptc *ParallelTestCase) RunParallel(fns ...func(workerID int)) {
	ptc.Helper()

	if len(fns) == 0 {
		return
	}

	done := make(chan struct{}, len(fns))

	for _, fn := range fns {
		go func(f func(workerID int)) {
			defer func() {
				if r := recover(); r != nil {
					ptc.Errorf("Parallel test panicked: %v", r)
				}
				done <- struct{}{}
			}()
			f(0)
		}(fn)
	}

	for i := 0; i < len(fns); i++ {
		<-done
	}
}

// ConcurrentTestCase manages concurrent test execution with shared state.
type ConcurrentTestCase struct {
	*testing.T
	results   chan TestResult
	errors    chan error
	timeout   time.Duration
}

// TestResult represents the result of a concurrent test.
type TestResult struct {
	WorkerID  int
	Duration  time.Duration
	Error     error
}

// NewConcurrentTestCase creates a new concurrent test case.
func NewConcurrentTestCase(t *testing.T, workers int, timeout time.Duration) *ConcurrentTestCase {
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &ConcurrentTestCase{
		T:       t,
		results: make(chan TestResult, workers),
		errors:  make(chan error, workers),
		timeout: timeout,
	}
}

// Run runs a function concurrently with the specified number of workers.
func (ctc *ConcurrentTestCase) Run(workers int, fn func(workerID int) error) {
	ctc.Helper()

	done := make(chan error, workers)

	for i := 0; i < workers; i++ {
		go func(workerID int) {
			done <- fn(workerID)
		}(i)
	}

	timer := time.NewTimer(ctc.timeout)
	defer timer.Stop()

	for i := 0; i < workers; i++ {
		select {
		case err := <-done:
			if err != nil {
				ctc.Errorf("Worker %d failed: %v", i, err)
			}
		case <-timer.C:
			ctc.Fatalf("Test timed out after %v", ctc.timeout)
		}
	}
}

// RunWithTimeout runs a function with a timeout.
func (ctc *ConcurrentTestCase) RunWithTimeout(timeout time.Duration, fn func() error) {
	ctc.Helper()

	done := make(chan error, 1)
	go func() {
		done <- fn()
	}()

	select {
	case err := <-done:
		if err != nil {
			ctc.Errorf("Test failed: %v", err)
		}
	case <-time.After(timeout):
		ctc.Fatalf("Test timed out after %v", timeout)
	}
}

// AssertNoErrors asserts that no errors occurred in the error channel.
func (ctc *ConcurrentTestCase) AssertNoErrors() {
	ctc.Helper()
	close(ctc.errors)
	for err := range ctc.errors {
		ctc.Errorf("Unexpected error: %v", err)
	}
}


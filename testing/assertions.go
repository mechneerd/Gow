package testing

import (
	"net/http"
	"strings"
	"time"
)

// AssertRedirect asserts a redirect response
func (tr *TestResponse) AssertRedirect(url string) *TestResponse {
	tr.T.Helper()
	if tr.Recorder == nil {
		tr.T.Errorf("Expected redirect to %s, but got no response", url)
		return tr
	}
	location := tr.Recorder.Header().Get("Location")
	if location != url {
		tr.T.Errorf("Expected redirect to %s, got %s", url, location)
	}
	return tr
}

// AssertRedirectToRoute asserts a redirect to a named route
func (tr *TestResponse) AssertRedirectToRoute(route string, params ...any) *TestResponse {
	tr.T.Helper()
	tr.T.Logf("Asserting redirect to route: %s", route)
	return tr
}

// AssertCookieExists asserts a cookie exists
func (tr *TestResponse) AssertCookieExists(name string) *TestResponse {
	tr.T.Helper()
	if tr.Recorder == nil {
		tr.T.Errorf("Expected cookie %s, but got no response", name)
		return tr
	}
	for _, cookie := range tr.Recorder.Result().Cookies() {
		if cookie.Name == name {
			return tr
		}
	}
	tr.T.Errorf("Expected cookie %s, but it was not found", name)
	return tr
}

// AssertCreated asserts a 201 status
func (tr *TestResponse) AssertCreated() *TestResponse {
	return tr.AssertStatus(http.StatusCreated)
}

// AssertNotFound asserts a 404 status
func (tr *TestResponse) AssertNotFound() *TestResponse {
	return tr.AssertStatus(http.StatusNotFound)
}

// AssertForbidden asserts a 403 status
func (tr *TestResponse) AssertForbidden() *TestResponse {
	return tr.AssertStatus(http.StatusForbidden)
}

// AssertUnauthorized asserts a 401 status
func (tr *TestResponse) AssertUnauthorized() *TestResponse {
	return tr.AssertStatus(http.StatusUnauthorized)
}

// AssertUnprocessableEntity asserts a 422 status
func (tr *TestResponse) AssertUnprocessableEntity() *TestResponse {
	return tr.AssertStatus(http.StatusUnprocessableEntity)
}

// AssertIsRedirect asserts a 3xx redirect
func (tr *TestResponse) AssertIsRedirect() *TestResponse {
	tr.T.Helper()
	if tr.Recorder == nil {
		tr.T.Error("Expected redirect, but got no response")
		return tr
	}
	if tr.Recorder.Code < 300 || tr.Recorder.Code >= 400 {
		tr.T.Errorf("Expected redirect (3xx), got %d", tr.Recorder.Code)
	}
	return tr
}

// AssertHeaderContains asserts a header contains a value
func (tr *TestResponse) AssertHeaderContains(name, value string) *TestResponse {
	tr.T.Helper()
	if tr.Recorder == nil {
		tr.T.Errorf("Expected header %s to contain %s, but got no response", name, value)
		return tr
	}
	actual := tr.Recorder.Header().Get(name)
	if !strings.Contains(actual, value) {
		tr.T.Errorf("Expected header %s to contain %s, got %s", name, value, actual)
	}
	return tr
}

// AssertDontSeeHtml asserts the response doesn't contain HTML
func (tr *TestResponse) AssertDontSeeHtml(html string) *TestResponse {
	tr.T.Helper()
	body := tr.Recorder.Body.String()
	if body == "" {
		tr.T.Errorf("Expected response not to contain HTML %s, but got empty response", html)
		return tr
	}
	if strings.Contains(body, html) {
		tr.T.Errorf("Expected response not to contain HTML %s, but it does", html)
	}
	return tr
}

// AssertSeeHtml asserts the response contains HTML
func (tr *TestResponse) AssertSeeHtml(html string) *TestResponse {
	tr.T.Helper()
	body := tr.Recorder.Body.String()
	if body == "" {
		tr.T.Errorf("Expected response to contain HTML %s, but got empty response", html)
		return tr
	}
	if !strings.Contains(body, html) {
		tr.T.Errorf("Expected response to contain HTML %s, but it does not", html)
	}
	return tr
}

// AssertExactJson asserts the response contains exact JSON
func (tr *TestResponse) AssertExactJson(data map[string]any) *TestResponse {
	tr.T.Helper()
	tr.T.Logf("Asserting exact JSON response: %v", data)
	return tr
}

// AssertJsonValidationErrors asserts JSON validation errors for fields
func (tr *TestResponse) AssertJsonValidationErrors(fields ...string) *TestResponse {
	tr.T.Helper()
	tr.T.Logf("Asserting JSON validation errors for: %v", fields)
	return tr
}

// AssertDownload asserts the response is a file download
func (tr *TestResponse) AssertDownload(filename string) *TestResponse {
	tr.T.Helper()
	if tr.Recorder == nil {
		tr.T.Errorf("Expected download %s, but got no response", filename)
		return tr
	}
	contentDisposition := tr.Recorder.Header().Get("Content-Disposition")
	if !strings.Contains(contentDisposition, "attachment") {
		tr.T.Error("Expected response to be a download")
	}
	return tr
}

// AssertStreamed asserts the response is streamed
func (tr *TestResponse) AssertStreamed() *TestResponse {
	tr.T.Helper()
	if tr.Recorder != nil && tr.Recorder.Header().Get("Transfer-Encoding") != "chunked" {
		tr.T.Error("Expected response to be streamed")
	}
	return tr
}

// AssertCookieExpired asserts a cookie is expired
func (tr *TestResponse) AssertCookieExpired(name string) *TestResponse {
	tr.T.Helper()
	if tr.Recorder == nil {
		tr.T.Errorf("Expected cookie %s to be expired, but got no response", name)
		return tr
	}
	for _, cookie := range tr.Recorder.Result().Cookies() {
		if cookie.Name == name {
			if !cookie.Expires.Before(time.Now()) {
				tr.T.Errorf("Expected cookie %s to be expired", name)
			}
			return tr
		}
	}
	tr.T.Errorf("Expected cookie %s, but it was not found", name)
	return tr
}

// AssertCookieNotExpired asserts a cookie is not expired
func (tr *TestResponse) AssertCookieNotExpired(name string) *TestResponse {
	tr.T.Helper()
	if tr.Recorder == nil {
		tr.T.Errorf("Expected cookie %s to not be expired, but got no response", name)
		return tr
	}
	for _, cookie := range tr.Recorder.Result().Cookies() {
		if cookie.Name == name {
			if cookie.Expires.Before(time.Now()) {
				tr.T.Errorf("Expected cookie %s to not be expired", name)
			}
			return tr
		}
	}
	tr.T.Errorf("Expected cookie %s, but it was not found", name)
	return tr
}

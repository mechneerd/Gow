package exception

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpExceptionWrapping(t *testing.T) {
	baseErr := errors.New("database connection failed")
	httpErr := Wrap(http.StatusInternalServerError, "Internal Server Error", baseErr)

	if !errors.Is(httpErr, baseErr) {
		t.Error("errors.Is failed to unwrap the underlying cause")
	}

	var unwrapped *HttpException
	if !errors.As(httpErr, &unwrapped) {
		t.Error("errors.As failed to target HttpException")
	}

	if unwrapped.Code != http.StatusInternalServerError {
		t.Errorf("expected code 500, got %d", unwrapped.Code)
	}
}

func TestHttpExceptionRenderJSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	err := NotFound("User not found")
	err.Render(w, req)

	res := w.Result()
	if res.StatusCode != 404 {
		t.Errorf("expected status 404, got %d", res.StatusCode)
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		t.Errorf("expected JSON content type, got %s", res.Header.Get("Content-Type"))
	}

	body := w.Body.String()
	if !strings.Contains(body, `"error":true`) || !strings.Contains(body, `"message":"User not found"`) {
		t.Errorf("unexpected JSON body: %s", body)
	}
}

func TestHttpExceptionRenderHTML(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	// No Accept header implies HTML fallback
	w := httptest.NewRecorder()

	err := BadRequest("Invalid input")
	err.Render(w, req)

	res := w.Result()
	if res.StatusCode != 400 {
		t.Errorf("expected status 400, got %d", res.StatusCode)
	}

	if !strings.Contains(res.Header.Get("Content-Type"), "text/html") {
		t.Errorf("expected HTML content type, got %s", res.Header.Get("Content-Type"))
	}

	body := w.Body.String()
	if !strings.Contains(body, "<html>") || !strings.Contains(body, "Invalid input") {
		t.Errorf("unexpected HTML body: %s", body)
	}
}

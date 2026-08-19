package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflight(t *testing.T) {
	handlerCalled := false
	handler := CORS("http://localhost:5173", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	}))

	req := httptest.NewRequest(http.MethodOptions, "/users", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("unexpected allowed origin: %q", got)
	}
	if handlerCalled {
		t.Fatal("preflight request should not reach the application handler")
	}
}

func TestCORSRejectsUnknownOrigin(t *testing.T) {
	handler := CORS("http://localhost:5173", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Origin", "https://example.com")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("unknown origin must not be allowed, got %q", got)
	}
}

func TestCORSWildcardAllowsOrigin(t *testing.T) {
	handler := CORS("*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/games", nil)
	req.Header.Set("Origin", "https://game.example.com")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, req)

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "https://game.example.com" {
		t.Fatalf("wildcard should echo the request origin, got %q", got)
	}
}

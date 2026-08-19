package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLoggingPreservesStatusAndAddsRequestID(t *testing.T) {
	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	request := httptest.NewRequest(http.MethodPost, "/games", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}
	if response.Header().Get("X-Request-ID") == "" {
		t.Fatal("expected X-Request-ID response header")
	}
}

func TestTimeoutSetsContextDeadline(t *testing.T) {
	handler := Timeout(time.Second, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Deadline(); !ok {
			t.Fatal("expected request context deadline")
		}
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/games", nil))
}

func TestTimeoutDoesNotBoundWebSocket(t *testing.T) {
	handler := Timeout(time.Second, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Deadline(); ok {
			t.Fatal("websocket request should not receive the API timeout")
		}
		w.WriteHeader(http.StatusOK)
	}))

	request := httptest.NewRequest(http.MethodGet, "/ws", nil).WithContext(context.Background())
	handler.ServeHTTP(httptest.NewRecorder(), request)
}

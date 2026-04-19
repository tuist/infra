package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewConfiguresHTTPServer(t *testing.T) {
	t.Parallel()

	server := New(Config{ListenAddress: ":18090"})

	if server.httpServer.Addr != ":18090" {
		t.Fatalf("http server addr = %q, want %q", server.httpServer.Addr, ":18090")
	}

	if server.httpServer.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf("http server read header timeout = %s, want %s", server.httpServer.ReadHeaderTimeout, 5*time.Second)
	}

	if server.httpServer.Handler == nil {
		t.Fatal("http server handler is nil")
	}
}

func TestHealthAndReadinessEndpoints(t *testing.T) {
	t.Parallel()

	server := New(Config{})

	testCases := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "healthz",
			path:       "/healthz",
			wantStatus: http.StatusOK,
			wantBody:   "ok",
		},
		{
			name:       "readyz",
			path:       "/readyz",
			wantStatus: http.StatusOK,
			wantBody:   "ready",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, testCase.path, nil)

			server.httpServer.Handler.ServeHTTP(recorder, request)

			if recorder.Code != testCase.wantStatus {
				t.Fatalf("%s status = %d, want %d", testCase.path, recorder.Code, testCase.wantStatus)
			}

			if recorder.Body.String() != testCase.wantBody {
				t.Fatalf("%s body = %q, want %q", testCase.path, recorder.Body.String(), testCase.wantBody)
			}
		})
	}
}

func TestInfoEndpoint(t *testing.T) {
	t.Parallel()

	server := New(Config{})
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/info", nil)

	server.httpServer.Handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("/v1/info status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("/v1/info content type = %q, want %q", contentType, "application/json")
	}

	var payload map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("unable to decode /v1/info response: %v", err)
	}

	if payload["name"] != "infra-api" {
		t.Fatalf("/v1/info name = %v, want %q", payload["name"], "infra-api")
	}

	if payload["apiVersion"] != "infra.tuist.dev/v1" {
		t.Fatalf("/v1/info apiVersion = %v, want %q", payload["apiVersion"], "infra.tuist.dev/v1")
	}

	if payload["version"] != "dev" {
		t.Fatalf("/v1/info version = %v, want %q", payload["version"], "dev")
	}
}

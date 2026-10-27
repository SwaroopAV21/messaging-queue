package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHandlerProxiesRoutes(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		upstream string
	}{
		{name: "projects", path: "/api/projects?limit=2", upstream: "/api/projects?limit=2"},
		{name: "featured projects", path: "/api/projects/featured", upstream: "/api/projects/featured"},
		{name: "project by id", path: "/api/projects/42", upstream: "/api/projects/42"},
		{name: "experiences", path: "/api/experiences", upstream: "/api/experiences"},
		{name: "skills", path: "/api/skill", upstream: "/api/skill"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if got := r.URL.RequestURI(); got != tt.upstream {
					t.Errorf("upstream request URI = %q, want %q", got, tt.upstream)
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
				_, _ = fmt.Fprint(w, `{"ok":true}`)
			}))
			defer upstream.Close()

			handler, err := NewGetHandler(upstream.URL)
			if err != nil {
				t.Fatalf("NewGetHandler() error = %v", err)
			}

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, tt.path, nil)
			handler.ServeHTTP(recorder, request)

			if recorder.Code != http.StatusAccepted {
				t.Errorf("status = %d, want %d", recorder.Code, http.StatusAccepted)
			}
			if got := recorder.Body.String(); got != `{"ok":true}` {
				t.Errorf("body = %q, want %q", got, `{"ok":true}`)
			}
		})
	}
}

func TestGetProjectRejectsInvalidID(t *testing.T) {
	handler, err := NewGetHandler("http://localhost:8080")
	if err != nil {
		t.Fatalf("NewGetHandler() error = %v", err)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/projects/not-a-number", nil)
	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

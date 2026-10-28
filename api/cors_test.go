package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	const allowedOrigin = "http://localhost:5173"
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := corsMiddleware(next, allowedOrigin)

	tests := []struct {
		name             string
		method           string
		origin           string
		wantStatus       int
		wantAllowOrigin  string
		wantAllowMethods string
	}{
		{
			name:            "allowed origin",
			method:          http.MethodGet,
			origin:          allowedOrigin,
			wantStatus:      http.StatusOK,
			wantAllowOrigin: allowedOrigin,
		},
		{
			name:             "allowed preflight",
			method:           http.MethodOptions,
			origin:           allowedOrigin,
			wantStatus:       http.StatusNoContent,
			wantAllowOrigin:  allowedOrigin,
			wantAllowMethods: "GET, OPTIONS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, "/", nil)
			request.Header.Set("Origin", tt.origin)

			handler.ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != tt.wantAllowOrigin {
				t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, tt.wantAllowOrigin)
			}
			if got := recorder.Header().Get("Access-Control-Allow-Methods"); got != tt.wantAllowMethods {
				t.Errorf("Access-Control-Allow-Methods = %q, want %q", got, tt.wantAllowMethods)
			}
		})
	}
}

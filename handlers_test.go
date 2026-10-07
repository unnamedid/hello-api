package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func get(t *testing.T, h http.Handler, path string) (int, string) {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	body, _ := io.ReadAll(rec.Body)
	return rec.Code, string(body)
}

func TestHealthz(t *testing.T) {
	code, body := get(t, newHandler("test"), "/healthz")
	if code != http.StatusOK || strings.TrimSpace(body) != "ok" {
		t.Fatalf("got %d %q, want 200 ok", code, body)
	}
}

func TestRootReportsEnvironment(t *testing.T) {
	code, body := get(t, newHandler("test"), "/")
	if code != http.StatusOK || !strings.Contains(body, `"environment":"test"`) {
		t.Fatalf("got %d %q", code, body)
	}
}

func TestVersionReportsEnv(t *testing.T) {
	t.Setenv("VERSION", "1.2.3")
	code, body := get(t, newHandler("test"), "/version")
	if code != http.StatusOK || !strings.Contains(body, `"version":"1.2.3"`) {
		t.Fatalf("got %d %q", code, body)
	}
}

func TestVersionDefaultsToDev(t *testing.T) {
	t.Setenv("VERSION", "")
	code, body := get(t, newHandler("test"), "/version")
	if code != http.StatusOK || !strings.Contains(body, `"version":"dev"`) {
		t.Fatalf("got %d %q", code, body)
	}
}

func TestUnknownPathIs404(t *testing.T) {
	if code, _ := get(t, newHandler("test"), "/nope"); code != http.StatusNotFound {
		t.Fatalf("got %d, want 404", code)
	}
}

func TestMetricsCountResponsesByClass(t *testing.T) {
	h := newHandler("test")
	get(t, h, "/healthz")
	get(t, h, "/healthz")
	get(t, h, "/nope")
	_, body := get(t, h, "/metrics")
	for _, want := range []string{
		`http_requests_total{status="2xx"} 2`,
		`http_requests_total{status="4xx"} 1`,
		`http_requests_total{status="5xx"} 0`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("metrics lack %q; got:\n%s", want, body)
		}
	}
}

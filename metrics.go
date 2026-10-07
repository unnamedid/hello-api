package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

// metrics counts responses by status class; the index is status/100.
type metrics struct {
	byClass [6]atomic.Int64
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (m *metrics) wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, req)
		class := rec.status / 100
		if class < 1 || class > 5 {
			class = 5
		}
		m.byClass[class].Add(1)
	})
}

func (m *metrics) serve(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintln(w, "# HELP http_requests_total Total HTTP requests by status class.")
	fmt.Fprintln(w, "# TYPE http_requests_total counter")
	for class := 2; class <= 5; class++ {
		fmt.Fprintf(w, "http_requests_total{status=\"%dxx\"} %d\n", class, m.byClass[class].Load())
	}
}

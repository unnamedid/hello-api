package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func newHandler(environment string) http.Handler {
	m := &metrics{}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("/metrics", m.serve)
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		version := os.Getenv("VERSION")
		if version == "" {
			version = "dev"
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"version": version})
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"service": "hello-api", "environment": environment})
	})
	return m.wrap(mux)
}

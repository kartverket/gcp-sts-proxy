package main

import (
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
)

func proxyHandler(w http.ResponseWriter, r *http.Request, cfg *config) {
	slog.Info("request", "url", r.URL.String(), "method", r.Method)

	rawURL := r.URL.Query().Get("url")
	if rawURL == "" {
		http.Error(w, "missing url param", http.StatusBadRequest)
		return
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		http.Error(w, "invalid url param", http.StatusBadRequest)
		return
	}

	token, err := getToken(cfg)
	if err != nil {
		slog.Info("failed to get token", "error", err.Error())
		http.Error(w, "failed to get token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), r.Method, rawURL, r.Body)
	if err != nil {
		slog.Info("failed to create request", "error", err.Error())
		http.Error(w, "failed to create request: "+err.Error(), http.StatusInternalServerError)
		return
	}

	maps.Copy(req.Header, r.Header)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "upstream error: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

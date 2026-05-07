package main

import (
	"io"
	"log/slog"
	"maps"
	"net/http"

	"golang.org/x/oauth2"
)

func proxyHandler(tokenSource oauth2.TokenSource) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "url", r.URL.String(), "method", r.Method)

		url := r.URL.Query().Get("url")
		if url == "" {
			http.Error(w, "missing url param", http.StatusBadRequest)
			return
		}

		token, err := tokenSource.Token()
		if err != nil {
			slog.Error("failed to get token", "error", err.Error())
			http.Error(w, "failed to get token", http.StatusInternalServerError)
			return
		}

		req, err := http.NewRequestWithContext(r.Context(), r.Method, url, r.Body)
		if err != nil {
			slog.Error("failed to create request", "error", err.Error())
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}

		maps.Copy(req.Header, r.Header)
		req.Header.Set("Authorization", "Bearer "+token.AccessToken)

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
		_, _ = io.Copy(w, resp.Body)
	}
}

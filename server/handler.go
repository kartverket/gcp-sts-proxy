package server

import (
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/url"
)

func (server *Server) proxyHandler(w http.ResponseWriter, r *http.Request) {
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

	token, err := server.tokenProvider.GetToken()
	if err != nil {
		slog.Info("failed to get token", "error", err.Error())
		http.Error(w, "failed to get token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), r.Method, rawURL, r.Body)
	if err != nil {
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

	// copy response headers and body
	for key, vals := range resp.Header {
		for _, v := range vals {
			w.Header().Add(key, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

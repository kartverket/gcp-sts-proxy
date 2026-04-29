package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
)

type config struct {
	tokenFile        string
	audience         string
	impersonationURL string
	port             string
}

func main() {
	cfg := config{
		tokenFile:        getEnv("TOKEN_FILE", "/var/run/secrets/tokens/gcp-ksa/token", false),
		audience:         getEnv("AUDIENCE", "", true),
		impersonationURL: getEnv("SERVICE_ACCOUNT_IMPERSONATION_URL", "", false),
		port:             getEnv("PORT", "8080", false),
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxyHandler(w, r, &cfg)
	})

	slog.Info("starting server", "port", cfg.port)
	log.Fatal(http.ListenAndServe(":"+cfg.port, mux))
}

func getEnv(key, fallback string, required bool) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if required {
		log.Fatalf("required env var %s not set", key)
	}
	return fallback
}

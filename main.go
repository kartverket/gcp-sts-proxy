package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/kartverket/gcp-sts-proxy/server"
	"github.com/kartverket/gcp-sts-proxy/token"
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
		impersonationURL: getEnv("IMPERSONATION_URL", "", false),
		port:             getEnv("PORT", "8080", false),
	}

	// set log output format to json
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	tokenProvider := token.NewProvider(token.Config{
		TokenFile:        cfg.tokenFile,
		Audience:         cfg.audience,
		ImpersonationURL: cfg.impersonationURL,
	})

	server := server.New(server.Config{
		Port: cfg.port,
	}, tokenProvider)
	server.Start()
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

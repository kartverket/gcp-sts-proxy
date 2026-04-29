package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/kartverket/gcp-sts-proxy/server"
	"github.com/kartverket/gcp-sts-proxy/token"
)

func main() {
	// set log output format to json
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	tokenProvider := token.NewProvider(token.Config{
		TokenFile:        getEnv("TOKEN_FILE", "/var/run/secrets/tokens/gcp-ksa/token", false),
		Audience:         getEnv("AUDIENCE", "", true),
		ImpersonationURL: getEnv("IMPERSONATION_URL", "", false),
	})

	server := server.New(server.Config{
		Port: getEnv("PORT", "8080", false),
	}, tokenProvider)

	err := server.Start()
	if err != nil {
		slog.Error("failed to start server", "error", err.Error())
	}
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

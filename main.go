package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	tokenSource, err := google.DefaultTokenSource(
		context.Background(),
		"https://www.googleapis.com/auth/cloud-platform",
	)
	if err != nil {
		slog.Error("failed to init token provider", "error", err.Error())
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", proxyHandler(tokenSource))

	slog.Info("starting server", "port", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		slog.Error("failed to start server", "error", err.Error())
		os.Exit(1)
	}
}

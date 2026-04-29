package main

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google/externalaccount"
)

type tokenProvider struct {
	once   sync.Once
	source oauth2.TokenSource
	err    error
}

var tokens tokenProvider

func getToken(cfg *config) (string, error) {
	return tokens.token(cfg)
}

func (p *tokenProvider) token(cfg *config) (string, error) {
	p.once.Do(func() {
		ts, err := externalaccount.NewTokenSource(context.Background(), externalaccount.Config{
			TokenURL:         "https://sts.googleapis.com/v1/token",
			Audience:         cfg.audience,
			SubjectTokenType: "urn:ietf:params:oauth:token-type:jwt",
			Scopes:           []string{"https://www.googleapis.com/auth/cloud-platform"},
			CredentialSource: &externalaccount.CredentialSource{
				File: cfg.tokenFile,
			},
			ServiceAccountImpersonationURL: cfg.impersonationURL,
		})
		if err != nil {
			p.err = fmt.Errorf("init token source: %w", err)
			return
		}
		p.source = oauth2.ReuseTokenSource(nil, ts)
	})

	if p.err != nil {
		return "", p.err
	}

	tok, err := p.source.Token()
	if err != nil {
		return "", err
	}
	return tok.AccessToken, nil
}

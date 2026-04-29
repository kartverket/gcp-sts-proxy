package token

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google/externalaccount"
)

type Provider struct {
	config Config
	// once ensures the token source is initialized only once
	once sync.Once
	// source is the underlying oauth2.TokenSource
	source oauth2.TokenSource
	// err captures any error during initialization
	err error
}

type Config struct {
	TokenFile        string
	Audience         string
	ImpersonationURL string
}

func NewProvider(config Config) *Provider {
	return &Provider{
		config: config,
	}
}

func (p *Provider) GetToken() (string, error) {
	p.once.Do(func() {
		ts, err := externalaccount.NewTokenSource(context.Background(), externalaccount.Config{
			TokenURL:         "https://sts.googleapis.com/v1/token",
			Audience:         p.config.Audience,
			SubjectTokenType: "urn:ietf:params:oauth:token-type:jwt",
			Scopes:           []string{"https://www.googleapis.com/auth/cloud-platform"},
			CredentialSource: &externalaccount.CredentialSource{
				File: p.config.TokenFile,
			},
			ServiceAccountImpersonationURL: p.config.ImpersonationURL,
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

package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	User        string
	Key         string
	Secret      string
	Token       string
	TokenSecret string
}

func setup() (*Config, error) {
	var (
		user        = flag.String("user", os.Getenv("TWITTER_USER"), "Twitter user name")
		key         = flag.String("key", os.Getenv("TWITTER_KEY"), "Twitter API Key")
		secret      = flag.String("secret", os.Getenv("TWITTER_SECRET"), "Twitter API Secret")
		token       = flag.String("token", os.Getenv("TWITTER_TOKEN"), "Twitter API Token")
		tokenSecret = flag.String("tokensecret", os.Getenv("TWITTER_TOKEN_SECRET"), "Twitter API Token Secret")
	)

	flag.Parse()
	if *user == "" {
		return nil, fmt.Errorf("user is required")
	}
	if *key == "" {
		return nil, fmt.Errorf("key is required")
	}
	if *secret == "" {
		return nil, fmt.Errorf("secret is required")
	}
	if *token == "" {
		return nil, fmt.Errorf("token is required")
	}
	if *tokenSecret == "" {
		return nil, fmt.Errorf("tokensecret is required")
	}

	return &Config{
		User:        *user,
		Key:         *key,
		Secret:      *secret,
		Token:       *token,
		TokenSecret: *tokenSecret,
	}, nil
}

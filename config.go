package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	TwitterUser        string
	TwitterKey         string
	TwitterSecret      string
	TwitterToken       string
	TwitterTokenSecret string
	TwitterScoreWeight float64
}

func setup() (*Config, error) {
	var (
		twUser        = flag.String("twuser", os.Getenv("TWITTER_USER"), "Twitter user name")
		twKey         = flag.String("twkey", os.Getenv("TWITTER_KEY"), "Twitter API Key")
		twSecret      = flag.String("twsecret", os.Getenv("TWITTER_SECRET"), "Twitter API Secret")
		twToken       = flag.String("twtoken", os.Getenv("TWITTER_TOKEN"), "Twitter API Token")
		twTokenSecret = flag.String("twtokensecret", os.Getenv("TWITTER_TOKEN_SECRET"), "Twitter API Token Secret")
		twScoreWeight = flag.Float64("twweight", 1.0, "Twitter score weight")
	)

	flag.Parse()
	if *twUser == "" {
		return nil, fmt.Errorf("twuser is required")
	}
	if *twKey == "" {
		return nil, fmt.Errorf("twkey is required")
	}
	if *twSecret == "" {
		return nil, fmt.Errorf("twsecret is required")
	}
	if *twToken == "" {
		return nil, fmt.Errorf("twtoken is required")
	}
	if *twTokenSecret == "" {
		return nil, fmt.Errorf("twtokensecret is required")
	}
	// Twitter score weight should be between 0 and 1
	if *twScoreWeight < 0 || *twScoreWeight > 1 {
		return nil, fmt.Errorf("twweight should be between 0 and 1 (1 is default)")
	}

	return &Config{
		TwitterUser:        *twUser,
		TwitterKey:         *twKey,
		TwitterSecret:      *twSecret,
		TwitterToken:       *twToken,
		TwitterTokenSecret: *twTokenSecret,
		TwitterScoreWeight: *twScoreWeight,
	}, nil
}

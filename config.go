package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Business           string
	Owner              string
	TwitterUser        string
	TwitterKey         string
	TwitterSecret      string
	TwitterToken       string
	TwitterTokenSecret string
	TwitterScoreWeight float64
	TwitterDays        int
	Verbose            bool
}

func setup() (*Config, error) {
	var (
		business      = flag.String("business", "unspecified", "Business name")
		owner         = flag.String("owner", "unspecified", "Business owner name")
		twUser        = flag.String("twuser", os.Getenv("TWITTER_USER"), "Twitter user name")
		twKey         = flag.String("twkey", os.Getenv("TWITTER_KEY"), "Twitter API Key")
		twSecret      = flag.String("twsecret", os.Getenv("TWITTER_SECRET"), "Twitter API Secret")
		twToken       = flag.String("twtoken", os.Getenv("TWITTER_TOKEN"), "Twitter API Token")
		twTokenSecret = flag.String("twtokensecret", os.Getenv("TWITTER_TOKEN_SECRET"), "Twitter API Token Secret")
		twScoreWeight = flag.Float64("twweight", 1.0, "Twitter score weight")
		twDays        = flag.Int("twdays", 180, "Number of days to use for twitter")
		verbose       = flag.Bool("verbose", false, "Verbosity flag")
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
	if *twDays <= 0 {
		return nil, fmt.Errorf("twdays must be greater than 0")
	}

	return &Config{
		Business:           *business,
		Owner:              *owner,
		TwitterUser:        *twUser,
		TwitterKey:         *twKey,
		TwitterSecret:      *twSecret,
		TwitterToken:       *twToken,
		TwitterTokenSecret: *twTokenSecret,
		TwitterScoreWeight: *twScoreWeight,
		TwitterDays:        *twDays,
		Verbose:            *verbose,
	}, nil
}

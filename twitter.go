package main

import (
	"net/url"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

func twitterScore(conf *Config) (float64, error) {
	anaconda.SetConsumerKey(conf.TwitterKey)
	anaconda.SetConsumerSecret(conf.TwitterSecret)
	api := anaconda.NewTwitterApi(conf.TwitterToken, conf.TwitterTokenSecret)

	user, err := api.GetUsersShow(conf.TwitterUser, nil)
	if err != nil {
		return 0, err
	}

	tweets, err := getTweets(api, conf.TwitterUser, time.Now().AddDate(0, 0, -conf.TwitterDays))
	if err != nil {
		return 0, err
	}

	tra, err := tweetRateAdjustment(tweets, conf.TwitterDays)
	if err != nil {
		return 0, err
	}
	rra := tweetReplyRatioAdjustment(tweets, user.Id)

	score := 1 - tra - rra
	if score < 0 {
		score = 0
	} else if score > 1 {
		score = 1
	}
	return score, nil
}

// tweetRateAdjustment assumes a gradual increase of tweets over a period of
// time is good. Within 5 - 10%. Anything outside of that is penalized
// 3 percentage point per percentage point outside that range.
func tweetRateAdjustment(tweets []anaconda.Tweet, dayLimit int) (float64, error) {
	out := 0.0

	oneThirdTweets, err := filterOlderTweets(tweets, time.Now().AddDate(0, 0, -dayLimit/3))
	if err != nil {
		return 0, err
	}

	percentIncrease := (float64(len(oneThirdTweets))/float64(len(tweets)))/(1.0/3.0) - 1

	// Adjust out if it's outside of 5 - 10%, otherwise leave it at 0.0 (no penalty).
	switch {
	case percentIncrease < 0.05:
		out = 0.05 - percentIncrease
	case percentIncrease > 0.1:
		out = percentIncrease - 0.1
	}

	return out * 3, nil
}

// tweetReplyRatioAdjustment assumes that a range of original tweets to replies
// are optimal. Being outside of that range is problematic.
//
// Less than 1 / 100 probably means the company isn't messaging enough, -0.3.
// Greater than 1 / 1 probably means no one is engaging with the company, -0.5.
func tweetReplyRatioAdjustment(tweets []anaconda.Tweet, userId int64) float64 {
	out := 0.0
	replyCount := 0.0
	for _, t := range tweets {
		if t.User.Id == userId && t.InReplyToUserID != 0 {
			replyCount++
		}
	}

	replyRatio := (float64(len(tweets)) - replyCount) / replyCount

	switch {
	case replyRatio < 1.0/100.0:
		out = 0.3
	case replyRatio > 1.0/1.0:
		out = 0.5
	}

	return out
}

func getTweets(api *anaconda.TwitterApi, user string, limit time.Time) ([]anaconda.Tweet, error) {
	v := url.Values{}
	v.Add("screen_name", user)
	v.Add("include_rts", "true")
	v.Add("exclude_replies", "false")
	v.Add("trim_user", "true")
	v.Add("count", "200")

	out := []anaconda.Tweet{}
	for {
		tweets, err := api.GetUserTimeline(v)
		if err != nil {
			return nil, err
		}

		if len(tweets) == 0 {
			// We've gotten all the tweets the user has.
			break
		}
		lastTweet := tweets[len(tweets)-1]
		v.Set("max_id", strconv.FormatInt(lastTweet.Id-1, 10))
		out = append(out, tweets...)

		t, err := lastTweet.CreatedAtTime()
		if err != nil {
			return nil, err
		}

		// Stop if the last tweet is older than the limit.
		if t.Before(limit) {
			break
		}
	}

	// We have at least one tweet older than our limit, so filter the tweets.
	return filterOlderTweets(out, limit)
}

// filterOlderTweets removes all tweets from a list that are older than the
// time limit specified.
//
// TODO: use binary search instead of linear search
func filterOlderTweets(ts []anaconda.Tweet, limit time.Time) ([]anaconda.Tweet, error) {
	for i := range ts {
		t, err := ts[i].CreatedAtTime()
		if err != nil {
			return nil, err
		}
		if t.Before(limit) {
			return ts[:i], nil
		}
	}
	return ts, nil
}

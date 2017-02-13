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

	tweets, err := getTweets(api, conf.TwitterUser, time.Now().AddDate(0, 0, -180))
	if err != nil {
		return 0, err
	}

	/*
		tweets60Days, err := filterOlderTweets(tweets, time.Now().AddDate(0, 0, -60))
		if err != nil {
			return 0, err
		}
	*/

	// tweets - 5 - 10% increase, every percent off is a percent off

	// tweet to reply ratio, shouldn't exceet 1 / 100, nor 1 / 1
	// below -.3, above -.5

	tra := tweetRateAdjustment(tweets)
	rra := tweetReplyRatioAdjustment(tweets, user.Id)

	score := 1 - tra - rra
	if score < 0 {
		score = 0
	}
	return score, nil
}

func tweetRateAdjustment(tweets []anaconda.Tweet) float64 {
	out := 0.0
	return out
}

func tweetReplyRatioAdjustment(tweets []anaconda.Tweet, userId int64) float64 {
	out := 0.0
	return out
}

func getTweets(api *anaconda.TwitterApi, user string, limit time.Time) ([]anaconda.Tweet, error) {
	v := url.Values{}
	v.Add("screen_name", user)
	//v.Add("include_rts", "false")
	v.Add("include_rts", "true")
	//v.Add("exclude_replies", "true")
	v.Add("trim_user", "true")
	v.Add("count", "200")

	//nintyDaysAgo := time.Now().AddDate(0, -6, 0)
	out := []anaconda.Tweet{}
	movePastRetweets := false
	for {
		tweets, err := api.GetUserTimeline(v)
		if err != nil {
			return nil, err
		}

		if len(tweets) == 0 {
			// All tweets in this section are retweets or replies.
			// Let's temporarily move past them.
			if movePastRetweets {
				// We've gotten all the tweets the user has.
				break
			}
			//v.Set("include_rts", "true")
			//v.Set("exclude_replies", "false")
			movePastRetweets = true
			continue
		}

		lastTweet := tweets[len(tweets)-1]

		v.Set("max_id", strconv.FormatInt(lastTweet.Id-1, 10))

		if movePastRetweets {
			// We were just trying to get past the retweets and replies.
			// Reenable the filters.
			//v.Set("include_rts", "false")
			//v.Set("exclude_replies", "true")
			movePastRetweets = false
		} else {
			out = append(out, tweets...)
		}

		t, err := lastTweet.CreatedAtTime()
		if err != nil {
			return nil, err
		}

		// Stop if the last tweet is older than 90 days.
		if t.Before(limit) {
			break
		}
	}

	// Remove tweets older than 90 days.
	return filterOlderTweets(out, limit)
	/*
		var i = len(out)
		for ; i > 0; i-- {
			t, err := out[i-1].CreatedAtTime()
			if err != nil {
				return nil, err
			}
			if t.After(nintyDaysAgo) {
				break
			}
		}
		return out[:i], nil
	*/
}

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

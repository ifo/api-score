package main

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

func main() {
	conf, err := setup()
	if err != nil {
		log.Fatal(err)
	}

	anaconda.SetConsumerKey(conf.Key)
	anaconda.SetConsumerSecret(conf.Secret)
	api := anaconda.NewTwitterApi(conf.Token, conf.TokenSecret)

	fmt.Println("starting")

	tweets, err := getTweets(api, conf.User, time.Now().AddDate(0, -6, 0))
	if err != nil {
		log.Fatal(err)
	}
	// TODO get last 30 days of tweets
	// TODO count tweets
	// compare numbers
	// determine output

	fmt.Println("tweet length", len(tweets))
	twts, err := filterOlderTweets(tweets, time.Now().AddDate(0, -1, 0))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("30 days tweet length", len(twts))

	orig := []anaconda.Tweet{}
	rtw := []anaconda.Tweet{}
	quote := []anaconda.Tweet{}
	replyCount := 0

	user, _ := api.GetUsersShow(conf.User, nil)
	fmt.Println("user", user.IdStr)

	for _, t := range tweets {
		switch {
		case t.RetweetedStatus != nil:
			rtw = append(rtw, t)
		case t.QuotedStatus != nil:
			quote = append(quote, t)
		default:
			orig = append(orig, t)
		}
		if t.User.Id == user.Id && t.InReplyToUserID != 0 {
			replyCount++
		}

		if t.RetweetedStatus != nil && t.QuotedStatus != nil {
			fmt.Println("retweeted and quoted")
		}
	}

	fmt.Println("original", len(orig))
	fmt.Println("quoted", len(quote))
	fmt.Println("retweet", len(rtw))
	fmt.Println("replies", replyCount)

	/*
		twJson, err := json.MarshalIndent(tweets, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("tweets.json", twJson, 0644)
		if err != nil {
			log.Fatal(err)
		}
	*/
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

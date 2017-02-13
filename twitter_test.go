package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

func Test_tweetRateAdjustment(t *testing.T) {
	limit := time.Now().AddDate(0, 0, -6)
	newTweet := makeTweet(0, limit.AddDate(0, 0, 5))
	oldTweet := makeTweet(0, limit.AddDate(0, 0, 1))

	cases := map[int]struct {
		In         []anaconda.Tweet
		DayLimit   int
		Adjustment string
		Err        error
	}{
		1: {In: []anaconda.Tweet{newTweet, newTweet, oldTweet, oldTweet, oldTweet, oldTweet}, DayLimit: 6, Adjustment: "0.15", Err: nil},
	}

	for num, c := range cases {
		adj, err := tweetRateAdjustment(c.In, c.DayLimit)
		if fmt.Sprintf("%0.2f", adj) != c.Adjustment {
			t.Errorf("Actual: %0.2f, Expected: %0.2f, Case: %d", adj, c.Adjustment, num)
		}
		if err != c.Err {
			t.Errorf("Actual: %+v, Expected: %+v, Case: %d", err, c.Err, num)
		}
	}
}

func Test_tweetReplyRatioAdjustment(t *testing.T) {
	replyTweet := makeTweet(2, time.Now())
	normalTweet := makeTweet(0, time.Now())

	cases := map[int]struct {
		In         []anaconda.Tweet
		Adjustment float64
	}{
		1: {In: []anaconda.Tweet{replyTweet, replyTweet, normalTweet}, Adjustment: 0.0},
		2: {In: []anaconda.Tweet{replyTweet}, Adjustment: 0.3},
		3: {In: []anaconda.Tweet{normalTweet}, Adjustment: 0.5},
	}

	for num, c := range cases {
		adj := tweetReplyRatioAdjustment(c.In, 1)
		if adj != c.Adjustment {
			t.Errorf("Actual: %0.1f, Expected: %0.1f, Case: %d", adj, c.Adjustment, num)
		}
	}
}

func Test_getTweets(t *testing.T) {
	t.Skip("TODO: Modify anaconda.TwitterApi.HttpClient to be a test client")
}

func Test_filterOlderTweets(t *testing.T) {
	limit := time.Now().AddDate(0, 0, -5)
	newTweet := makeTweet(0, limit.AddDate(0, 0, 2))
	oldTweet := makeTweet(0, limit.AddDate(0, 0, -3))

	cases := map[int]struct {
		In        []anaconda.Tweet
		Limit     time.Time
		OutLength int
		Err       error
	}{
		1: {In: []anaconda.Tweet{newTweet, oldTweet}, Limit: limit, OutLength: 1, Err: nil},
		2: {In: []anaconda.Tweet{newTweet}, Limit: limit, OutLength: 1, Err: nil},
		3: {In: []anaconda.Tweet{oldTweet}, Limit: limit, OutLength: 0, Err: nil},
		4: {In: []anaconda.Tweet{}, Limit: limit, OutLength: 0, Err: nil},
	}

	for num, c := range cases {
		out, err := filterOlderTweets(c.In, c.Limit)
		if len(out) != c.OutLength {
			t.Errorf("Actual: %d, Expected: %d, Case: %d", len(out), c.OutLength, num)
		}
		if err != c.Err {
			t.Errorf("Actual: %+v, Expected: %+v, Case: %d", err, c.Err, num)
		}
	}
}

func makeTweet(replyId int64, date time.Time) anaconda.Tweet {
	return anaconda.Tweet{
		InReplyToUserID: replyId,
		User:            anaconda.User{Id: 1},
		CreatedAt:       date.Format(time.RubyDate),
	}
}

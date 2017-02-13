package main

import (
	"testing"
	"time"

	"github.com/ChimeraCoder/anaconda"
)

func Test_tweetRateAdjustment(t *testing.T) {
	t.Skip()
}

func Test_tweetReplyRatioAdjustment(t *testing.T) {
	t.Skip()
}

func Test_getTweets(t *testing.T) {
	t.Skip("TODO: Modify anaconda.TwitterApi to send to getTweets")
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
		CreatedAt:       date.Format(time.RubyDate),
	}
}

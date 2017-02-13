package main

import (
	"fmt"
	"log"
)

func main() {
	conf, err := setup()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("starting")

	twScore, err := twitterScore(conf)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(twScore * conf.TwitterScoreWeight)

	/*
		origCount := 0
		replyCount := 0
		for _, t := range tweets {
			if t.User.Id == user.Id {
				if t.InReplyToUserID != 0 {
					replyCount++
				} else {
					origCount++
				}
			}
		}
		fmt.Println(replyCount)
		fmt.Println(origCount)
	*/
}

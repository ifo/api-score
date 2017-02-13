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

	fmt.Println("Calculating...\n")

	twScore, err := twitterScore(conf)
	if err != nil {
		log.Fatal(err)
	}

	totalScore := twScore * conf.TwitterScoreWeight

	fmt.Printf("Business:      %s\n", conf.Business)
	if conf.Verbose {
		fmt.Printf("Owner:         %s\n", conf.Owner)
		fmt.Printf("Twitter Score: %0.2f\n", twScore)
	}
	fmt.Printf("Total Score:   %0.2f\n", totalScore)
}

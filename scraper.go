package main

import (
	"context"
	"log"
	"time"

	"github.com/ClemSK/blog_aggregator/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequests time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	// initialize the ticker immediately, then wait for the next interval
	// "for range  <- ticker.C" would wait for the 1st min then execute
	for ; ; <-ticker.C {
		db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
	}
}

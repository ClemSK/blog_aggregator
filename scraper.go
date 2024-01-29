package main

import (
	"context"
	"log"
	"sync"
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
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue // if we return it would stop the function instead of continuing the scraping
		}

		waitGroup := &sync.WaitGroup{}
		for _, feed := range feeds {
			waitGroup.Add(1)

			go scrapeFeeds(db, waitGroup, feed)
		}
		waitGroup.Wait()
	}
}

func scrapeFeeds(db *database.Queries, waitGroup *sync.WaitGroup, feed database.Feed) {
	defer waitGroup.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("Found post: ", item.Title, "on feed ", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

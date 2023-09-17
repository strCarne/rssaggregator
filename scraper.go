package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/strCarne/rssaggregator/internal/database"
)

func startScraping(
	db *database.Queries,
	concurency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %v duration\n", concurency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurency),
		)
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, db, feed)
		}
		wg.Wait()

	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error while marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	newPosts := len(rssFeed.Channel.Item)

	for _, item := range rssFeed.Channel.Item {

		pubAt, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("couldn't parse date %v: %v\n", item.PubDate, err)
			newPosts--
			continue
		}

		_, err = db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID: uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Title: item.Title,
				Description: makeNullString(item.Description),
				PublishedAt: pubAt,
				Url: item.Link,
				FeedID: feed.ID,
			},
		)

		if err != nil {
			newPosts--
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("couldn't create a post:", err)
		}
	}
	log.Printf("Have been collected a feed %v, %v posts found", feed.Name, newPosts)
}

func makeNullString(s string) sql.NullString {
	res := sql.NullString{String: s}
	if s == "" {
		res.Valid = false
	} else {
		res.Valid = true
	}
	return res
}
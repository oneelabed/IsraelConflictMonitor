package scraper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oneelabed/RSSAggregator/internal/database"
)

func StartScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v goroutines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error fetching feeds:", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go ScrapeFeed(db, wg, feed)
		}

		wg.Wait()
	}
}

func ScrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := UrlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
		return
	}

	notRelevant := 0

	for _, item := range rssFeed.Channel.Item {
		if !IsRelevant(item.Title, item.Description) {
			notRelevant++
			continue
		}

		desc := sql.NullString{}
		if len(item.Description) != 0 {
			desc.String = item.Description
			desc.Valid = true
		}

		pubDate, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Println("Couldn't parse date:", err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: desc,
			PublishedAt: pubDate,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Error creating post:", err)
			return
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item)-notRelevant)
}

func IsRelevant(title, description string) bool {
	var keywords = []string{
		"israel",
		"idf",
		"gaza",
		"hamas",
		"hezbollah",
		"lebanon",
		"west bank",
		"jerusalem",
		"tel aviv",
		"iran",
	}

	text := strings.ToLower(title + " " + description)

	for _, k := range keywords {
		if strings.Contains(text, k) {
			return true
		}
	}

	return false
}

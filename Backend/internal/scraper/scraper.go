package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/oneelabed/IsraelConflictMonitor/internal/database"
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
		log.Printf("Error marking feed %v as fetched: %v", feed.Name, err)
		return
	}

	rssFeed, err := UrlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed %v: %v\n", feed.Name, err)
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

		pubDate, err := flexibleDate(item.PubDate)

		if err != nil {
			log.Printf("Couldn't parse date for feed %v: %v", feed.Name, err)
			continue
		}

		switch feed.Name {
		case "Jerusalem Post":
			desc.String = GetBetween(item.Description, `alt='`, `' title`)
		case "Middle East Eye":
			desc.String = GetBetween(item.Description, `<p>`, `</p>`)
		case "The Times Of Israel":
			desc.String = GetBetween(item.Description, `<p>`, `</p>`)
		case "Al Monitor":
			desc.String = GetBetween(item.Description, `<p>`, `</p>`)
		}

		desc.String = stripTags(desc.String)

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
	fullText := strings.ToLower(title + " " + description)

	// 1. Geography/Participants (The "Who/Where")
	regionKeys := []string{
		"israel", "gaza", "lebanon", "hezbollah", "hamas", "idf", "iran", "jerusalem", "tel aviv", "yemen",
		"palestine", "west bank", "ירושלים", "תל אביב", "לבנון", "איראן", "פלסטין", `איו"ש`, "גדה המערבית",
		"חמאס", "חזבאללה", `צה"ל`, "צהל", "ישראל", "עזה", "איראן",
	}

	// 2. High Intensity (The "What")
	intensityKeys := []string{
		"rocket", "missile", "airstrike", "interception", "siren", "red alert", "iron dome", "uav", "drone",
		"fire", "shelling", "atgm", "rpg", "explosion", "bombardment", "artillery", "casualty", "injured",
		"wounded", "fatality", "critical", "hospitalized", "evacuation", "emergency", "paramedics", "mda",
		"operation", "maneuver", "division", "brigade", "counter-terrorism", "special forces",
		"reserves", "miluim", "security forces", "border", "tunnel", "ambush", "attack", "stabbing", "shooting",
		"terrorist", "car-ramming", "hostages", "kidnapped", "negotiation", "ceasefire",
		"רקטה", "אזעקה", "יירוט", "תקיפה", "צבע אדום", "טיל", "ירי", `פצמ"ר`, "כיפת ברזל", `כטב"ם`,
		`נ"ט`, "פיצוץ", "נפגעים", "פצועים", "אנוש", "פינוי", "מבצע", "תמרון", "פיצוץ", "הפגזה", "ארטילריה",
		"אוגדה", "חטיבה", "סיכול", "כוחות מיוחדים", "פיגוע", "הרוגים", "חירום", "מדא", `מד"א`, "בית חולים",
		"דקירה", "מחבל", "דריסה", "מארב", "גבול", "מנהרה", "חילופי אש", "מילואים", "כוחות הביטחון", "הפסקת אש",
		"חטיפה", "חטופים",
	}

	// 3. Non Related subjects
	var negativeKeywords = []string{
		"sports", "football", "soccer", "basketball", "eurovision", "stock market",
		"weather", "entertainment", "lifestyle", "tourism", "restaurant", "recipe",
		"ספורט", "כדורגל", "כדורסל", "אירוויזיון", "בורסה", "מזג אוויר", "תיירות",
	}

	hasNonRelated := false
	for _, k := range negativeKeywords {
		if strings.Contains(fullText, k) {
			hasNonRelated = true
			return false
		}
	}

	hasRegion := false
	for _, k := range regionKeys {
		if strings.Contains(fullText, k) {
			hasRegion = true
			break
		}
	}

	hasIntensity := false
	for _, k := range intensityKeys {
		if strings.Contains(fullText, k) {
			hasIntensity = true
			break
		}
	}

	// Only return true if ALL are present
	return hasRegion && hasIntensity && !hasNonRelated
}

func flexibleDate(dateStr string) (time.Time, error) {
	layouts := []string{time.RFC1123, time.RFC1123Z, time.RFC3339,
		"2006-01-02T15:04:05-0700", "Mon, 02 Jan 2006 15:04:05", "Mon, 02 Jan 06 15:04:05 -0700"}

	dateStr = strings.TrimSpace(dateStr) // remove trailing/leading spaces

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown date format: %s", dateStr)
}

// GetBetween returns the string found between 'start' and 'finish'.
// If either boundary is not found, it returns an empty string.
func GetBetween(input, start, finish string) string {
	// Find the position of the start string
	startIndex := strings.Index(input, start)
	if startIndex == -1 {
		return ""
	}

	// Move the index to the end of the start string
	startIndex += len(start)

	// Find the position of the finish string, starting from the new startIndex
	endOffset := strings.Index(input[startIndex:], finish)
	if endOffset == -1 {
		return ""
	}

	// Calculate the actual end index
	endIndex := startIndex + endOffset

	return input[startIndex:endIndex]
}

func stripTags(input string) string {
	// This regex finds anything between < and > and deletes it
	re := regexp.MustCompile(`<[^>]*>`)
	return re.ReplaceAllString(input, "")
}

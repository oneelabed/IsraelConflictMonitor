package scraper

import (
	"encoding/xml"
	//"html"
	"io"
	"net/http"
	"strings"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func UrlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	// 1. Create a custom request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RSSFeed{}, err
	}

	// 2. Add a User-Agent to "pretend" we are a browser
	// This stops Cloudflare from blocking your Render server
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml")
	req.Header.Set("Accept-Language", "he-IL,he;q=0.9,en-US;q=0.8,en;q=0.7") // <--- Tell them you speak Hebrew
	req.Header.Set("Referer", "https://www.timesofisrael.com/")              // <--- Pretend you came from their homepage

	// 3. Execute the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	// 4. Read the data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	// Handle special characters
	xmlString := string(data)
	xmlString = strings.ReplaceAll(xmlString, "& ", "&amp; ")

	reader := strings.NewReader(xmlString)
	decoder := xml.NewDecoder(reader)
	decoder.Strict = false
	decoder.Entity = xml.HTMLEntity

	rssFeed := RSSFeed{}
	err = decoder.Decode(&rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}

	return rssFeed, nil
}

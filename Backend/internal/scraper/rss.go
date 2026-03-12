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

	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	/*sanitizedData := html.UnescapeString(string(data))
	sanitizedData = strings.ReplaceAll(sanitizedData, "&bull;", "•")
	sanitizedData = strings.ReplaceAll(string(data), "&", "&amp;")*/
	xmlString := string(data)
	xmlString = strings.ReplaceAll(xmlString, "& ", "&amp; ")

	reader := strings.NewReader(xmlString)
	decoder := xml.NewDecoder(reader)

	decoder.Strict = false
	// This tells the decoder to automatically handle HTML entities
	decoder.Entity = xml.HTMLEntity

	rssFeed := RSSFeed{}

	/*err = xml.Unmarshal([]byte(sanitizedData), &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}*/
	err = decoder.Decode(&rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}

	return rssFeed, nil
}

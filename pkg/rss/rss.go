package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"github.com/microcosm-cc/bluemonday"
)


type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}


type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}


func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("get: %w", err)
	}

	req.Header.Set("User-Agent", "gatorcli")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("http client do: %w", err)
	}
	defer res.Body.Close()

	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("read all: %w", err)
	}

	var htmx RSSFeed
	if err := xml.Unmarshal(byteData, &htmx); err != nil {
		return &RSSFeed{}, fmt.Errorf("unmarshal htmx: %w", err)
	}

	cleanText(&htmx)

	return &htmx, nil
}

func cleanText (feed *RSSFeed) error {
	p := bluemonday.StrictPolicy()

	feed.Channel.Title = html.UnescapeString(p.Sanitize(feed.Channel.Title))
	feed.Channel.Description = html.UnescapeString(p.Sanitize(feed.Channel.Description))
	
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(p.Sanitize(feed.Channel.Item[i].Title))
		feed.Channel.Item[i].Description = html.UnescapeString(p.Sanitize(feed.Channel.Item[i].Description))
	}

	return nil
}
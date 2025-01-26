package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/manzil-infinity180/golang-webrss/internal/database"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}
		log.Printf("Found %v feeds to fetch!", len(feeds))
		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
			wg.Add(1)
			// go routine scrape rss
			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}
	// fetch feed
	feedData, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		item.Title = strings.TrimSpace(item.Title)
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}

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

func fetchFeed(feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}

	return &rssFeed, nil
}

type RSSRemoteOk struct {
	Channel struct {
		Title       string            `xml:"title"`
		Description string            `xml:"description"`
		Item        []RSSRemoteOkItem `xml:"item"`
	} `xml:"channel"`
}
type RSSRemoteOkItem struct {
	Title       string `xml:"title"`
	Company     string `xml:"company"`
	Link        string `xml:"link"`
	Image       string `xml:"image"`
	Tag         string `xml:"tag"`
	Location    string `xml:"location"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func startScrapingRemoteOk(db *database.Queries, concurrency int, timeBetweenRequest time.Duration, url string) {
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := scrapeRemoteFeed(context.Background(), db, url); err != nil {
				log.Printf("Scrape feed error: %v", err)
			}
		}()
		wg.Wait()
	}
}

func scrapeRemoteFeed(ctx context.Context, db *database.Queries, url string) error {
	feedData, err := fetchRemoteOkFeed(url)
	if err != nil {
		return fmt.Errorf("fetch feed failed: %w", err)
	}
	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err = db.CreateRemoteJob(ctx, database.CreateRemoteJobParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Company: item.Company,
			Url:     item.Link,
			Image:   item.Image,
			Tag: sql.NullString{
				String: item.Tag,
				Valid:  true,
			},
			Location:    item.Location,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed jobs collected, %v posts found", len(feedData.Channel.Item))
	return nil
}

func fetchRemoteOkFeed(feedURL string) (*RSSRemoteOk, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	fmt.Println(feedURL)
	resp, err := httpClient.Get(feedURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var rssFeed RSSRemoteOk
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		log.Printf("Failed to parse RSS feed: %v", err)
		log.Printf("Raw RSS feed data: %s", string(dat))
		return nil, fmt.Errorf("XML unmarshal failed: %v", err)
	}

	return &rssFeed, nil
}

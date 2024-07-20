package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/bhashimoto/blog-aggregator-bootdev/internal/database"
	"github.com/google/uuid"
)

type RssFeedXml struct {
	XMLName          xml.Name `xml:"rss"`
	Version          string   `xml:"version,attr"`
	ContentNamespace string   `xml:"xmlns:content,attr"`
	Channel          *RssFeed
}

type RssFeed struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`       // required
	Link           string   `xml:"link"`        // required
	Description    string   `xml:"description"` // required
	Language       string   `xml:"language,omitempty"`
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"` // updated used
	Generator      string   `xml:"generator,omitempty"`
	Items          []*RssItem `xml:"item"`
}
type RssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`       // required
	Link        string   `xml:"link"`        // required
	Description string   `xml:"description"` // required
	Author      string `xml:"author,omitempty"`
	Category    string `xml:"category,omitempty"`
	PubDate     string   `xml:"pubDate,omitempty"` // created or updated
	Source      string   `xml:"source,omitempty"`
}

func (cfg *apiConfig) FetchFeedsRoutine(nFeeds int32, loopDuration time.Duration) {
	for {
		go cfg.FetchFeeds(nFeeds)
		time.Sleep(loopDuration)
	}
}


func (cfg *apiConfig) FetchFeeds(n int32) error {
	log.Println("Getting feeds to fetch")
	dbFeeds, err := cfg.db.GetNextFeedsToFetch(context.Background(), n)
	if err != nil {
		return err
	}

	feeds := []Feed{}
	for _, dbFeed := range dbFeeds {
		feeds = append(feeds, databaseFeedToFeed(dbFeed))
	}

	var wg sync.WaitGroup

	for _, feed := range feeds {
		log.Println("Starting to process feed", feed.Name)
		wg.Add(1)
		go func(Feed Feed){
			defer wg.Done()
			rssFeed, err := cfg.fetchFeed(feed)
			if err != nil {
				log.Println(err.Error())
				return
			}
			cfg.processFeed(rssFeed)
			log.Println("Finished processing feed", feed.Name)
		}(feed)
	}
	return nil
}

func (cfg *apiConfig) fetchFeed(feed Feed) (RssFeed, error){
	resp, err := http.Get(feed.Url)
	if err != nil {
		log.Println("error in http.Get")
		return RssFeed{}, err
	}

	bodyString, err := io.ReadAll(resp.Body)
	if err != nil {
		return RssFeed{}, err
	}

	data := RssFeedXml{}

	err = xml.Unmarshal(bodyString, &data)
	if err != nil {
		return RssFeed{}, err 
	}

	err = cfg.db.MarkFeedFetched(context.Background(), feed.ID)
	return *data.Channel, nil

}

func (cfg *apiConfig) processFeed(feed RssFeed) {
	dbFeed, err := cfg.db.GetFeedFromURL(context.Background(), feed.Link)
	if err != nil {
		log.Println("feed not found:", feed.Link)
		return
	}
	for _, item := range feed.Items {
		log.Printf("[%s] Fetching post: %s", feed.Title, item.Title)
		pubDate := sql.NullTime{
			Time: time.Now(),
			Valid: false,
		}
		if item.PubDate !=  "" {
			parsedDate, err := time.Parse("Wed, 03 Apr 2024 00:00:00 +0000", item.PubDate)
			if err != nil {
				log.Println("error parsing date:", item.PubDate)
			} else {
				pubDate = sql.NullTime{
					Time: parsedDate,
					Valid: true,
				}	
			}
		} 
		dbPost, err := cfg.db.CreatePost(context.Background(), database.CreatePostParams{
			ID: uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title: item.Title,
			Url: item.Link,
			FeedID: dbFeed.ID,
			Description: item.Description,
			PublishedAt: pubDate,
		})
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Inserted post:",dbPost.Title)

	}
}

func getURLsFromFeed(feed RssFeed) ([]string, error) {
	urls := []string{}
	for _, item := range feed.Items {
		urls = append(urls, item.Link)
	}
	return urls, nil
}

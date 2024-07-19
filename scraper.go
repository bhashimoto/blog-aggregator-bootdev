package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
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
			rssFeed, err := fetchFeed(feed)
			if err != nil {
				log.Println(err.Error())
				return
			}
			processFeed(rssFeed)
			log.Println("Finished processing feed", feed.Name)
		}(feed)
	}
	return nil
}

func fetchFeed(feed Feed) (RssFeed, error){
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
	return *data.Channel, nil

}

func processFeed(feed RssFeed) {
	for _, item := range feed.Items {
		fmt.Println(item.Title)
	}
}

func getURLsFromFeed(feed RssFeed) ([]string, error) {
	urls := []string{}
	for _, item := range feed.Items {
		urls = append(urls, item.Link)
	}
	return urls, nil
}

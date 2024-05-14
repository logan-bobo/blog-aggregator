package scraper

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/logan-bobo/blog-aggregator/internal/database"
)

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func Scrape(inputURL string, waitGroup *sync.WaitGroup, ch chan<- RSSFeed) {
	defer waitGroup.Done()

	returnFeed := RSSFeed{}

	_, err := url.Parse(inputURL)

	if err != nil {
		log.Println(err)
		ch <- returnFeed
		return
	}

	resp, err := http.Get(inputURL)

	if err != nil {
		log.Println(err)
		ch <- returnFeed
		return
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&returnFeed)

	if err != nil {
		log.Println(err)
		ch <- returnFeed
		return
	}

	ch <- returnFeed
	return
}

func Worker(numberOfFeedsToPull int32, db *database.Queries) error {
	ticker := time.NewTicker((1 * 10) * time.Second)

	for {
		_ = <-ticker.C

		log.Println("Starting scraper")

		// TODO: Learn about context
		feeds, err := db.GetFeedsToFetch(context.TODO(), numberOfFeedsToPull)
		log.Println(feeds)

		if err != nil {
			log.Println(err)
			return err
		}

		var wg sync.WaitGroup
		messages := make(chan RSSFeed)

		for _, feed := range feeds {
			wg.Add(1)

			log.Printf("Pulling Feed: %s", feed.Name)

			go Scrape(feed.Url, &wg, messages)

			msg := <-messages
			fmt.Println(msg.Channel.Title)
		}

		wg.Wait()

		log.Println("Scrape finished")
	}
}

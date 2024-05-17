package scraper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
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

		log.Println("starting scraper")

		// TODO: Learn about context
		feeds, err := db.GetFeedsToFetch(context.TODO(), numberOfFeedsToPull)

		if err != nil {
			log.Println(err)
			return err
		}

		var wg sync.WaitGroup
		messages := make(chan RSSFeed)

		for _, feed := range feeds {
			wg.Add(1)

			log.Printf("pulling feed: %s", feed.Name)

			go Scrape(feed.Url, &wg, messages)

			msg := <-messages
			log.Printf("processing feed: %s", msg.Channel.Title)

			for _, post := range msg.Channel.Item {
				now := time.Now()

				timeFormat := "Mon, 02 Jan 2006 15:04:05 -0700"
				parsedTime, err := time.Parse(timeFormat, post.PubDate)

				if err != nil {
					log.Printf("can not parse time %s with %s", post.PubDate, err)
				}

				outputLayout := "2006-01-02 15:04:05.000000"
				formattedTime := parsedTime.Format(outputLayout)
				formattedparsedTime, err := time.Parse(outputLayout, formattedTime)

				if err != nil {
					log.Printf("can not parse time to correct format %s, with %s", formattedTime, err)
				}

				params := database.CreatePostParams{
					ID:          uuid.New(),
					CreatedAt:   now,
					UpdatedAt:   now,
					Title:       sql.NullString{String: post.Title, Valid: true},
					Url:         sql.NullString{String: post.Link, Valid: true},
					Description: sql.NullString{String: post.Description, Valid: true},
					PublishedAt: sql.NullTime{Time: formattedparsedTime, Valid: true},
					FeedID:      uuid.NullUUID{UUID: feed.ID, Valid: true},
				}

				_, err = db.CreatePost(context.TODO(), params)

				// we dont really care if the entry already exists so skip logging that it cant be created
				if err != nil && err.Error() != "pq: duplicate key value violates unique constraint \"unique_url\"" {
					log.Printf("can not create post %s in database with error %s", post.Title, err)
				}
			}

			wg.Wait()

			log.Println("scrape finished")
		}
	}
}

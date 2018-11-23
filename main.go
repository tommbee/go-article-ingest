package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/tommbee/go-article-ingest/ingestor"
	"github.com/tommbee/go-article-ingest/normaliser"
)

var (
	mediumIngestor ingestor.Medium
)

func crawl(url string, ch chan string, chFinished chan bool) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".streamItem--postPreview").Each(func(i int, s *goquery.Selection) {
		a := s.Find("a[data-action-source]").First()
		d := s.Find("time").First()
		t := s.Find("h3.graf--leading").First()
		link, ok := a.Attr("href")
		if ok {
			date := d.Text()
			title := t.Text()
			a, err := normaliser.Normalise(title, link, date)
			if err == nil {
				fmt.Print(a)
				// Save to DB
			}
		}
	})
}

func main() {
	seedUrls := os.Args[1:]

	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)

	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}
}

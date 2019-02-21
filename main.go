package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tommbee/go-article-ingest/model"
	"github.com/tommbee/go-article-ingest/normaliser"
	"github.com/tommbee/go-article-ingest/poller"
	"github.com/tommbee/go-article-ingest/repository"
)

var repo *repository.ArticleRepository

const INTERVAL_PERIOD time.Duration = 2 * time.Hour

type jobTicker struct {
	t *time.Timer
}

func newRepo() *repository.MongoArticleRepository {
	ro := &repository.MongoArticleRepository{
		Server:       os.Getenv("SERVER"),
		DatabaseName: os.Getenv("DB"),
		AuthDatabase: os.Getenv("AUTH_DB"),
		DBSSL:        os.Getenv("DB_SSL"),
		Collection:   os.Getenv("ARTICLE_COLLECTION"),
		Username:     os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
	}
	return ro
}

func getNextTickDuration() time.Duration {
	now := time.Now()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
	if nextTick.Before(now) {
		nextTick = nextTick.Add(INTERVAL_PERIOD)
	}
	return nextTick.Sub(time.Now())
}

func NewJobTicker() jobTicker {
	fmt.Println("new tick here")
	return jobTicker{time.NewTimer(getNextTickDuration())}
}

func poll(url string, ch chan string, chFinished chan bool) {
	success := true

	// Factory generate poller
	p, err := poller.Generate(url)
	if err != nil {
		success = false
		log.Printf("Error finding poller %s: %s", url, err)
	}

	// Poll URL
	articleData, err := p.Poll(url)
	if err != nil || len(articleData) == 0 {
		success = false
		log.Printf("No articlces found on %s: %s", url, err)
	}

	/** @todo: Add rabbit queue to handle normalise separately */
	// Factory generate normaliser
	n, err := normaliser.Generate(url)
	if err != nil {
		success = false
		log.Printf("Error finding normaliser %s: %s", url, err)
	}

	// Normalise articles
	articles := []model.Article{}
	for _, ad := range articleData {
		a, err := n.Normalise(ad.Title, ad.Link, ad.Date)
		if err == nil {
			articles = append(articles, a)
		} else {
			log.Printf("Error normalising article %s: %s", a, err)
		}
	}

	/** @todo: Add rabbit queue to handle ingest separately */
	// Ingest
	repo := newRepo()
	for _, ar := range articles {
		_, err := repo.Insert(ar)
		if err != nil {
			log.Printf("Error on %s: %s", url, err)
		}
	}

	defer func() {
		chFinished <- success
	}()

	if success {
		ch <- url
	}
}

func (jt jobTicker) updateJobTicker() {
	fmt.Println("next tick here")
	jt.t.Reset(getNextTickDuration())
}

func initPolling() {
	fmt.Println(time.Now(), " - init polling")
	sourceEnv := os.Getenv("SOURCES")
	sourceUrls := strings.Split(sourceEnv, ",")

	log.Printf("Parsing URLs: %s", sourceUrls)

	chUrls := make(chan string)
	chFinished := make(chan bool)

	for _, url := range sourceUrls {
		go poll(url, chUrls, chFinished)
	}

	successUrls := make([]string, 0)
	for i := 0; i < len(sourceUrls); {
		select {
		case url := <-chUrls:
			successUrls = append(successUrls, url)
		case <-chFinished:
			i++
		}
	}
	log.Printf("Succesfully parsed: %s", successUrls)
}

func main() {
	// jt := NewJobTicker()
	// for {
	// 	<-jt.t.C
	// 	initPolling()
	// 	jt.updateJobTicker()
	// }
	initPolling()
}

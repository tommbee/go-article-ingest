package main

import (
	"log"
	"os"
	"strings"

	"github.com/tommbee/go-article-ingest/model"
	"github.com/tommbee/go-article-ingest/normaliser"
	"github.com/tommbee/go-article-ingest/poller"
	"github.com/tommbee/go-article-ingest/repository"
)

var repo *repository.ArticleRepository

func newRepo() *repository.MongoArticleRepository {
	ro := &repository.MongoArticleRepository{
		Server:       os.Getenv("SERVER"),
		DatabaseName: os.Getenv("DB"),
		Collection:   os.Getenv("ARTICLE_COLLECTION"),
	}
	return ro
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
		err = repo.Insert(ar)
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

func main() {
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

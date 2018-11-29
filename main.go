package main

import (
	"log"
	"os"
	"strings"

	"github.com/tommbee/go-article-ingest/parser"
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

func parse(url string, ch chan string, chFinished chan bool) {
	// Factory generate parser
	success := true
	p, err := parser.Generate(url)
	if err != nil {
		success = false
		log.Printf("Error on %s: %s", url, err)
	}

	/** @todo: Add rabbit queue to handle normalise */
	// Parse + Normalise
	articles, err := p.Parse(url)
	if err != nil || len(articles) == 0 {
		success = false
		log.Printf("No articlces found on %s: %s", url, err)
	}

	/** @todo: Add rabbit queue to handle ingest */
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
		go parse(url, chUrls, chFinished)
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

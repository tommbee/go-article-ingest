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
	p, err := parser.Generate(url)
	if err != nil {
		log.Fatal(err)
	}

	// Parse
	articles, err := p.Parse(url)
	if err != nil {
		log.Fatal(err)
	}

	// Ingest
	repo := newRepo()
	for i, ar := range articles {
		err = repo.Insert(ar)
		if err != nil {
			log.Fatalf("insert fail #%d %v\n", i, err)
		}
	}
}

func main() {
	sourceEnv := os.Getenv("SOURCES")
	sourceUrls := strings.Split(sourceEnv, ",")

	chUrls := make(chan string)
	chFinished := make(chan bool)

	for _, url := range sourceUrls {
		go parse(url, chUrls, chFinished)
	}
}

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
		log.Fatal(err)
	}

	// Parse
	articles, err := p.Parse(url)
	if err != nil || len(articles) == 0 {
		//log.Fatal(err)
		success = false
	}

	// *Should send to rabbit queue for ingest*

	// Ingest
	repo := newRepo()
	for _, ar := range articles {
		err = repo.Insert(ar)
		if err != nil {
			//log.Fatalf("insert fail #%d %v\n", i, err)
			success = false
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

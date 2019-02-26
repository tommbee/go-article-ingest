package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tommbee/go-article-ingest/model"
	"github.com/tommbee/go-article-ingest/normaliser"
	"github.com/tommbee/go-article-ingest/poller"
	"github.com/tommbee/go-article-ingest/repository"
)

var repo *repository.ArticleRepository
var opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
	Name: "article_app_processed_ops_total",
	Help: "The total number of parsed articles",
})
var normaliserErrors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "article_app_normaliser_errors",
	Help: "The total number of errors from normaliser",
})
var noArticleErrors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "article_app_no_article_errors",
	Help: "The total number of no article errors",
})
var dbErrors = promauto.NewCounter(prometheus.CounterOpts{
	Name: "article_app_db_errors",
	Help: "The total number of DB errors",
})

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
		noArticleErrors.Inc()
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
			normaliserErrors.Inc()
			log.Printf("Error normalising article %s: %s", a, err)
		}
	}

	/** @todo: Add rabbit queue to handle ingest separately */
	// Ingest
	repo := newRepo()
	for _, ar := range articles {
		_, err := repo.Insert(ar)
		if err != nil {
			dbErrors.Inc()
			log.Printf("Error on %s: %s", url, err)
		}
	}

	defer func() {
		chFinished <- success
	}()

	if success {
		log.Println("Success")
		opsProcessed.Inc()
		ch <- url
	}
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
	gocron.Every(1).Hour().Do(initPolling)
	_, time := gocron.NextRun()
	fmt.Println(time)
	// function Start start all the pending jobs
	<-gocron.Start()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}

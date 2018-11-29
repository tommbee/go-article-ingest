package parser

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/tommbee/go-article-ingest/model"
	"github.com/tommbee/go-article-ingest/normaliser"
)

// Parser object parses content from a URL
type Parser struct {
	config     model.Config
	normaliser *normaliser.Normaliser
}

// Parse the webpage
func (p *Parser) Parse(URL string) ([]model.Article, error) {
	log.Printf("Attempting request %s", URL)

	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	articles := []model.Article{}

	p.normaliser, err = normaliser.Generate(URL)
	if err != nil {
		return articles, err
	}

	doc.Find(p.config.BaseElement).Each(func(i int, s *goquery.Selection) {
		// Add checks here to determine if these els have been found
		a := s.Find(p.config.LinkElement).First()
		d := s.Find(p.config.DateElement).First()
		t := s.Find(p.config.TitleElement).First()
		link, ok := a.Attr("href")
		if ok {
			date := d.Text()
			title := t.Text()
			article, err := p.normaliser.Normalise(title, link, date)
			if err == nil {
				articles = append(articles, article)
			} else {
				log.Printf("Error when parsing: %s", err)
			}
		}
	})

	return articles, nil
}

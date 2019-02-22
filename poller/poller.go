package poller

import (
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tommbee/go-article-ingest/model"
)

// Poller object parses content from a URL
type Poller struct {
	config model.Config
}

// ArticleData Basic struct to hold article data
type ArticleData struct {
	Title string
	Link  string
	Date  string
}

// Poll the URL
func (p *Poller) Poll(URL string) ([]ArticleData, error) {
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

	articles := []ArticleData{}

	doc.Find(p.config.BaseElement).Each(func(i int, s *goquery.Selection) {
		// Add checks here to determine if these els have been found
		a := s.Find(p.config.LinkElement).First()
		d := s.Find(p.config.DateElement).First()
		t := s.Find(p.config.TitleElement).First()
		link, ok := a.Attr("href")
		date, dok := d.Attr("datetime")
		if ok && dok {
			if !strings.HasPrefix(link, p.config.BaseURL) {
				link = p.config.BaseURL + link
			}
			title := t.Text()
			article := ArticleData{
				Title: title,
				Date:  date,
				Link:  link,
			}
			articles = append(articles, article)
		}
	})

	return articles, nil
}

package normaliser

import (
	"log"
	"time"

	"github.com/tommbee/go-article-ingest/model"
)

// Normaliser takes data and makes it readable
type Normaliser struct {
	DateFormat string
}

// Normalise the article
func (n *Normaliser) Normalise(title string, url string, date string) (model.Article, error) {
	log.Printf("Normalising data %s", url)

	dt, err := time.Parse(n.DateFormat, date)
	if err != nil {
		return model.Article{}, err
	}
	a := model.Article{}.NewArticle(title, url, dt)

	log.Printf("Normalised data %s", a)
	return a, nil
}

package normaliser

import (
	"time"

	"github.com/tommbee/go-article-ingest/model"
)

// Normaliser takes data and makes it readable
type Normaliser struct {
	DateFormat string
}

// Normalise the article
func (n *Normaliser) Normalise(title string, url string, date string) (model.Article, error) {
	dt, err := time.Parse(n.DateFormat, date)
	if err != nil {
		return model.Article{}, err
	}
	a := model.Article{}.NewArticle(title, url, dt)
	return a, nil
}

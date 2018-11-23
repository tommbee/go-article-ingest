package normaliser

import (
	"time"

	"github.com/tommbee/go-article-ingest/model"
)

type reddit struct{}

func (r reddit) normalise(t string, u string, d string) (model.Article, error) {
	// 2018-10-26T03:35:46.274Z
	dt, err := time.Parse("", d)
	if err != nil {
		return model.Article{}, err
	}
	a := model.Article{}.NewArticle(t, u, dt)
	return a, nil
}

package normaliser

import (
	"time"

	"../model"
)

type medium struct{}

func (m medium) normalise(t string, u string, d string) (model.Article, error) {
	// 2018-10-26T03:35:46.274Z
	dt, err := time.Parse("", d)
	if err != nil {
		return model.Article{}, err
	}
	a := model.Article{}.NewArticle(t, u, dt)
	return a, nil
}

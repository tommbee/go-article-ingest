package repository

import (
	"github.com/tommbee/go-article-feed/model"
)

// ArticleRepository handles the interface to persistant storage
type ArticleRepository interface {
	Insert(*model.Article) error
}

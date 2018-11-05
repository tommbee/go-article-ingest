package model

import (
	"time"
)

// Article - an article object
type Article struct {
	title         string
	url           string
	publishedDate time.Time
}

// NewArticle construct a new Article object
func (a Article) NewArticle(title string, url string, publishedDate time.Time) Article {
	return Article{
		title:         title,
		url:           url,
		publishedDate: publishedDate,
	}
}

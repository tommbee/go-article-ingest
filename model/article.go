package model

import (
	"time"
)

// Article - an article object
type Article struct {
	Title     string    `json:"title" bson:"title"`
	URL       string    `json:"url" bson:"url"`
	Published time.Time `json:"published" bson:"published"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// NewArticle construct a new Article object
func (a Article) NewArticle(title string, url string, publishedDate time.Time) Article {
	return Article{
		Title:     title,
		URL:       url,
		Published: publishedDate,
	}
}

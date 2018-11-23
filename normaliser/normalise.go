package normaliser

import (
	"errors"
	"strings"

	"github.com/tommbee/go-article-ingest/model"
)

var (
	mediumNormaliser medium
	redditNormaliser reddit
)

// Normalise the article
func Normalise(t string, u string, d string) (model.Article, error) {
	switch {
	case strings.Contains(u, "medium"):
		return mediumNormaliser.normalise(t, u, d)
	case strings.Contains(u, "reddit"):
		return redditNormaliser.normalise(t, u, d)
	}
	return model.Article{}, errors.New("No normaliser found")
}

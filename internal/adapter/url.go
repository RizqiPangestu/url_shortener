package adapter

import "github.com/RizqiPangestu/url_shortener/internal/core"

type urlMongoAdapter struct {
}

func NewURLMongoAdapter() core.URLPort {
	return &urlMongoAdapter{}
}

func (a *urlMongoAdapter) Shorten(url string) (string, error) {
	return "https://short.url/" + url, nil
}

func (a *urlMongoAdapter) Expand(shortURL string) (string, error) {
	return "https://" + shortURL, nil
}

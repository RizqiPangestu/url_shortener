package app

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type ExpandRequest struct {
	ShortURL string `query:"short_url" validate:"required"`
}

type ExpandResponse struct {
	URL string `json:"url"`
}

package app

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type RedirectRequest struct {
	ShortPath string `param:"short_path" validate:"required"`
}

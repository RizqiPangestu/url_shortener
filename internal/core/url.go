package core

type URLService interface {
	Shorten(url string) (string, error)
	Expand(shortURL string) (string, error)
}

type URLPort interface {
	Shorten(url string) (string, error)
	Expand(shortURL string) (string, error)
}

type urlService struct {
	port URLPort
}

func NewURLService(port URLPort) URLService {
	return &urlService{
		port: port,
	}
}

func (s *urlService) Shorten(url string) (string, error) {
	return s.port.Shorten(url)
}

func (s *urlService) Expand(shortURL string) (string, error) {
	return s.port.Expand(shortURL)
}

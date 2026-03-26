package core

import "github.com/google/uuid"

type URLService interface {
	Shorten(url string) (string, error)
	Expand(shortURL string) (string, error)
}

type URLPort interface {
	SavePath(shortPath string, originURL string) error
	GetOriginURL(shortPath string) (string, error)
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
	shortPath := uuid.New().String()[:8]
	if err := s.port.SavePath(shortPath, url); err != nil {
		return "", err
	}

	return "https://" + shortPath, nil
}

func (s *urlService) Expand(shortPath string) (string, error) {
	originURL, err := s.port.GetOriginURL(shortPath)
	if err != nil {
		return "", err
	}

	return originURL, nil
}

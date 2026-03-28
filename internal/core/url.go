package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type URLService interface {
	Shorten(originalURL string) (string, error)
	Expand(shortURL string) (string, error)
}

type URLPort interface {
	SavePath(shortPath string, originURL string, ttl time.Duration) error
	GetByOriginalURL(originalURL string) (URL, error)
	GetByShortPath(shortPath string) (URL, error)
	UpdateLastAccessedAt(shortPath string) error
	DeleteByShortPath(shortPath string) error
}

type urlService struct {
	port URLPort
}

func NewURLService(port URLPort) URLService {
	return &urlService{
		port: port,
	}
}

func (s *urlService) Shorten(originalURL string) (string, error) {
	// reuse existing short path if exists
	url, err := s.port.GetByOriginalURL(originalURL)
	if err != nil {
		if errors.Is(err, ErrURLNotFound) {
			return s.shortenWithRetry(originalURL, 0)
		}
		return "", err
	}

	return url.ShortPath, nil
}

func (s *urlService) shortenWithRetry(originalURL string, attempt int) (string, error) {
	shortPath := uuid.NewString()[:8]
	if err := s.port.SavePath(shortPath, originalURL, DefaultTTL); err != nil {
		if errors.Is(err, ErrURLAlreadyExists) && attempt < 3 { // retry 3 times
			return s.shortenWithRetry(originalURL, attempt+1)
		}
		return "", err
	}

	return shortPath, nil
}

func (s *urlService) Expand(shortPath string) (string, error) {
	url, err := s.port.GetByShortPath(shortPath)
	if err != nil {
		return "", err
	}

	if time.Since(url.LastAccessedAt) > url.TTL {
		if err := s.port.DeleteByShortPath(shortPath); err != nil {
			return "", err
		}
		return "", ErrURLExpired
	}

	if err := s.port.UpdateLastAccessedAt(shortPath); err != nil {
		return "", err
	}

	return url.OriginalURL, nil
}

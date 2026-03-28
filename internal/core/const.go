package core

import (
	"errors"
	"time"
)

const (
	DefaultTTL = time.Hour * 24 * 7 // 7 days
)

var (
	ErrURLExpired       = errors.New("URL has expired")
	ErrURLAlreadyExists = errors.New("URL already exists")
	ErrURLNotFound      = errors.New("URL not found")
	ErrSystemError      = errors.New("SYSTEM_ERROR")
)

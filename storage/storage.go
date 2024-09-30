package storage

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
)

var ErrNoSavedPages = errors.New("no saved pages")

type Storage interface {
	Save(p *Page) error
	PickRandom(UserName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

// func Hash hashing unique links
func (p Page) Hash() (string, error) {
	hash := sha256.New()

	// hash URL + UserName to prevent collisions on identical links from diff users
	if _, err := io.WriteString(hash, p.URL); err != nil {
		return "", fmt.Errorf("failed to hash %w", err)
	}

	if _, err := io.WriteString(hash, p.UserName); err != nil {
		return "", fmt.Errorf("failed to hash %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

package storage

import (
	"context"
	"errors"
)

var ErrNoSavedPages = errors.New("no saved pages")

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, UserName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

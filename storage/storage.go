package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"

	"github.com/RSheremeta/read-adviser-bot/lib/e"
)

var ErrNoSavedPages = errors.New("no saved pages")
var ErrNoUserSavingsHistory = errors.New("no user savings history")

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, username string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExist(ctx context.Context, p *Page) (bool, error)
}

type Page struct {
	URL      string
	Username string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("cannot calculate hash", err)
	}
	if _, err := io.WriteString(h, p.Username); err != nil {
		return "", e.Wrap("cannot calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

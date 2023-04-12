package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/RSheremeta/read-adviser-bot/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("cannot open the database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to the database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`

	if _, err := s.db.ExecContext(ctx, q); err != nil {
		return fmt.Errorf("cannot create the table: %w", err)
	}

	return nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.Username); err != nil {
		return fmt.Errorf("cannot save the page %w", err)
	}

	return nil
}
func (s *Storage) PickRandom(ctx context.Context, username string) (*storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, q, username).Scan(&url)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("no saved pages found for user %s\n", username)
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("cannot pick a random page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		Username: username,
	}, nil
}

func (s *Storage) Remove(ctx context.Context, p *storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name =?`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.Username); err != nil {
		return fmt.Errorf("cannot remove the page: %w", err)
	}

	return nil
}

func (s *Storage) IsExist(ctx context.Context, p *storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, q, p.URL, p.Username).Scan(&count); err != nil {
		return false, fmt.Errorf("cannot check whether the page exists: %w", err)
	}

	return count > 0, nil
}

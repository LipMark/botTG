package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"TGBot/storage"
)

type Storage struct {
	db *sql.DB
}

// NewStorage connects to DB.
func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}
	return &Storage{db: db}, nil
}

// Save adds page to storage.
func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	query := `INSERT INTO pages (url,user_name) VALUES (?,?)`

	if _, err := s.db.ExecContext(ctx, query, p.URL, p.UserName); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

// PickRandom returns a random page from the vault.
func (s *Storage) PickRandom(ctx context.Context, userName string) (*storage.Page, error) {
	query := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(ctx, query, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, fmt.Errorf("can't pick random page: %w", err)
	}

	return &storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove is used to remove an existing storage page.
func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	query := `DELETE FROM pages WHERE url = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, query, page.URL, page.UserName); err != nil {
		return fmt.Errorf("can't remove  page: %w", err)
	}
	return nil
}

// IsExists check the existence of requested page.
func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	query := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	if err := s.db.QueryRowContext(ctx, query, page.URL, page.UserName).Scan(&count); err != nil {
		return false, fmt.Errorf("can't check for page existence %w", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`
	_, err := s.db.ExecContext(ctx, query)

	if err != nil {
		return fmt.Errorf("can't create table %w", err)
	}

	return nil
}

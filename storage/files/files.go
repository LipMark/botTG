package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"TGBot/storage"
)

const ReadWritePerm = 0774

type Storage struct {
	basePath string
}

// NewPath creates a new storage with given path.
func NewPath(basePath string) Storage {
	return Storage{basePath: basePath}
}

// Save adds page to storage.
func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't save this page %w", err)
		}
	}()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, ReadWritePerm); err != nil {
		return err
	}

	fileName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("can't create file name %w", err)
	}

	filePath = filepath.Join(filePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("can't create file  %w", err)
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("failed to encode  %w", err)
	}

	return nil
}

// PickRandom returns a random page from the vault.
func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("can't get random page %w", err)
		}
	}()

	path := filepath.Join(s.basePath, userName)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

// Remove is used to remove an existing storage page.
func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("can't remove file %w", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return fmt.Errorf("%s %w", msg, err)
	}

	return nil
}

// IsExists check the existence of requested page.
func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("can't check if file exists %w", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, fmt.Errorf("%s %w", msg, err)
	}

	return true, nil
}

// decodePage decoding a page using gob Decoder.
func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't decode this page %w", err)
	}
	defer func() { _ = f.Close() }()

	var page storage.Page

	if err := gob.NewDecoder(f).Decode(&page); err != nil {
		return nil, fmt.Errorf("can't decode this page %w", err)
	}

	return &page, nil
}

// fileName receives page and returns a new file name.
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}

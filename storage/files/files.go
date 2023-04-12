package files

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/RSheremeta/read-adviser-bot/lib/e"
	"github.com/RSheremeta/read-adviser-bot/storage"
)

const defaultPerm = 0774

type Storage struct {
	basePath string
}

func New(basepath string) Storage {
	return Storage{basePath: basepath}
}

func (s Storage) Save(_ context.Context, page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("cannot save", err) }()

	fPath := filepath.Join(s.basePath, page.Username)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := filename(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(_ context.Context, username string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("cannot pick random page", err) }()

	path := filepath.Join(s.basePath, username)

	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, storage.ErrNoUserSavingsHistory
		}
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(_ context.Context, p *storage.Page) error {
	fName, err := filename(p)
	if err != nil {
		return e.Wrap("cannot remove file", err)
	}

	path := filepath.Join(s.basePath, p.Username, fName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("cannot remove file %s", path)
		return e.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExist(_ context.Context, p *storage.Page) (bool, error) {
	fName, err := filename(p)
	if err != nil {
		return false, e.Wrap("cannot check whether file exists", err)
	}

	path := filepath.Join(s.basePath, p.Username, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("cannot check whether file %s exists", path)
		return false, e.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodePage(filepath string) (*storage.Page, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, e.Wrap("cannot decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("cannot decode page", err)
	}

	return &p, nil
}

func filename(p *storage.Page) (string, error) {
	return p.Hash()
}

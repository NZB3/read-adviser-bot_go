package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"read-adviser-bot/lib/errwrap"
	"read-adviser-bot/storage"
	"time"
)

type Storage struct {
	basePath string
}

// 0774 - hex number, allows read and write to all users
// read more about permissions
const defaultPerm = 0774

var ErrNoSavedPages = errors.New("no saved pages")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = errwrap.WrapIfErr("can't save", err) }()
	fPath := filepath.Join(s.basePath, page.UserName)

	if err = os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
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

func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = errwrap.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	rand.NewSource(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodPage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fName, err := fileName(page)
	if err != nil {
		return errwrap.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove file %s", path)

		return errwrap.Wrap(msg, err)
	}

	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		return false, errwrap.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, errwrap.Wrap(msg, err)
	}

	return true, nil
}

func (s Storage) decodPage(filePath string) (page *storage.Page, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errwrap.Wrap("can't decode page", err)
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, errwrap.Wrap("can't decode page", err)
	}

	return page, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}

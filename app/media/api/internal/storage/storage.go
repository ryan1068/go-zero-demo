package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var ErrTooLarge = errors.New("file too large")

type Storage interface {
	Save(ctx context.Context, name string, r io.Reader, maxBytes int64) (int64, error)
	Open(name string) (io.ReadSeekCloser, error)
	Delete(name string) error
}

type LocalStorage struct {
	dir string
}

func NewLocal(dir string) *LocalStorage {
	return &LocalStorage{dir: dir}
}

func (s *LocalStorage) Save(ctx context.Context, name string, r io.Reader, maxBytes int64) (int64, error) {
	if err := os.MkdirAll(s.dir, 0o755); err != nil {
		return 0, err
	}

	tmp, err := os.CreateTemp(s.dir, ".upload-*")
	if err != nil {
		return 0, err
	}
	tmpName := tmp.Name()
	defer func() {
		if tmpName != "" {
			if removeErr := os.Remove(tmpName); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
				_ = removeErr
			}
		}
	}()

	written, err := io.Copy(tmp, io.LimitReader(r, maxBytes+1))
	if err != nil {
		if closeErr := tmp.Close(); closeErr != nil {
			return 0, closeErr
		}
		return 0, err
	}
	if err := tmp.Close(); err != nil {
		return 0, err
	}
	if written > maxBytes {
		return 0, ErrTooLarge
	}

	if err := os.Rename(tmpName, filepath.Join(s.dir, name)); err != nil {
		return 0, err
	}
	tmpName = ""
	return written, nil
}

func (s *LocalStorage) Open(name string) (io.ReadSeekCloser, error) {
	return os.Open(filepath.Join(s.dir, filepath.Base(name)))
}

func (s *LocalStorage) Delete(name string) error {
	return os.Remove(filepath.Join(s.dir, filepath.Base(name)))
}

package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	BaseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	return &LocalStorage{BaseDir: baseDir}
}

func (l *LocalStorage) Save(id string, data io.Reader) (string, error) {
	path := filepath.Join(l.BaseDir, id+".dat")

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	return path, err
}

func (l *LocalStorage) Get(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (l *LocalStorage) Delete(path string) error {
	return os.Remove(path)
}

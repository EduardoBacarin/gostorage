package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	rootPath string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		panic("failed to create storage base directory: " + err.Error())
	}
	return &LocalStorage{rootPath: baseDir}
}

func (l *LocalStorage) Save(id string, data io.Reader) (string, error) {
	path := filepath.Join(l.rootPath, id)

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

func (l *LocalStorage) BaseDir() string {
	return l.rootPath
}

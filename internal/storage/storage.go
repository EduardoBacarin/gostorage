package storage

import "io"

type Engine interface {
	Save(id string, data io.Reader) (string, error)
	Get(path string) (io.ReadCloser, error)
	Delete(path string) error
	BaseDir() string
}

package storage

import (
	"context"
	"errors"
	"io"
)

type Storage interface {
	Put(ctx context.Context, data io.ReadCloser) (string, error)
	Get(ctx context.Context, fd string) (io.ReadCloser, error)
	Delete(ctx context.Context, fd string) error
}

type Config struct {
	LocalPath string `yaml:"local_path"`
}

func New(cfg *Config) (Storage, error) {
	if cfg == nil {
		return nil, errors.New("empty storage config")
	}
	if cfg.LocalPath != "" {
		return NewLocal(cfg.LocalPath)
	}
	return NewFS()
}

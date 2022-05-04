package disk

import (
	"errors"

	"github.com/ocean-rw/ocean/internal/ocean-osd/disk/common"
	"github.com/ocean-rw/ocean/internal/ocean-osd/disk/fs"
	"github.com/ocean-rw/ocean/internal/ocean-osd/disk/raw"
)

var ErrInvalidConfig = errors.New("invalid config")

type Config struct {
	FS  *fs.Config  `yaml:"fs"`
	Raw *raw.Config `yaml:"raw"`
}

func Open(cfg *Config) ([]common.DiskIF, error) {
	if cfg.FS != nil {
		return fs.Open(cfg.FS)
	}
	if cfg.Raw != nil {
		return raw.Open(cfg.Raw)
	}
	return nil, ErrInvalidConfig
}

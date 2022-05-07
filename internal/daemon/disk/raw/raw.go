package raw

import (
	"github.com/ocean-rw/ocean/internal/daemon/disk/common"
)

type Config struct {
}

type RawDisk struct {
}

func Open(cfg *Config) ([]common.DiskIF, error) {
	return nil, nil
}

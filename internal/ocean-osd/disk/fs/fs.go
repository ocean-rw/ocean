package fs

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	json "github.com/json-iterator/go"

	"github.com/ocean-rw/ocean/internal/ocean-osd/disk/common"
	"github.com/ocean-rw/ocean/pkg/proto"
)

const metaFilename = ".meta"

type Config struct {
	Dir  string `yaml:"dir"`
	Host string `yaml:"host"`
}

type Disk struct {
	*proto.Disk
}

func Open(cfg *Config) ([]common.DiskIF, error) {
	dir, err := filepath.Abs(cfg.Dir)
	if err != nil {
		return nil, err
	}
	disks, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	ret := make([]common.DiskIF, 0)
	for _, disk := range disks {
		if !disk.IsDir() {
			continue
		}
		path := fmt.Sprintf("%s/%s", dir, disk.Name())
		label, err := readMeta(path)
		if err != nil {
			return nil, err
		}
		d := &proto.Disk{
			DiskLabel: label,
			Host:      cfg.Host,
			Path:      path,
			State:     proto.DiskStateNormal,
		}
		ret = append(ret, &Disk{Disk: d})
	}
	return ret, nil
}

func readMeta(path string) (*proto.DiskLabel, error) {
	label := new(proto.DiskLabel)
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", path, metaFilename))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	err = json.Unmarshal(data, label)
	if err != nil {
		return nil, err
	}
	return label, nil
}

func writeMeta(path string, label *proto.DiskLabel) error {
	data, err := json.Marshal(label)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.tmp", path, metaFilename), data, 0666)
	if err != nil {
		return err
	}
	return os.Rename(fmt.Sprintf("%s/%s.tmp", path, metaFilename), fmt.Sprintf("%s/%s", path, metaFilename))
}

func (d *Disk) Init(ctx context.Context, label *proto.DiskLabel) error {
	err := writeMeta(d.Path, label)
	if err != nil {
		return err
	}
	d.DiskLabel = label
	return nil
}

func (d *Disk) Stat(ctx context.Context) (*proto.Disk, error) {
	return d.Disk, nil
}

func (d *Disk) Put(ctx context.Context, fd string, data io.ReadCloser) error {
	filename := fmt.Sprintf("%s/%s", d.Path, fd)
	f, err := os.OpenFile(filename+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, data)
	if err != nil {
		return err
	}
	defer data.Close()

	return os.Rename(fd+".tmp", filename)
}

func (d *Disk) Get(ctx context.Context, fd string) (io.ReadCloser, error) {
	filename := fmt.Sprintf("%s/%s", d.Path, fd)
	f, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f, nil
}

func (d *Disk) Delete(ctx context.Context, fd string) error {
	filename := fmt.Sprintf("%s/%s", d.Path, fd)
	return os.Remove(filename)
}

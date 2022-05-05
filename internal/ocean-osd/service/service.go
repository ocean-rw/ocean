package service

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
	json "github.com/json-iterator/go"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/ocean-osd/disk"
	"github.com/ocean-rw/ocean/internal/ocean-osd/disk/common"
	"github.com/ocean-rw/ocean/pkg/proto"
)

type Config struct {
	Host    string       `yaml:"host"`
	MgrHost string       `yaml:"mgr_host"`
	Disk    *disk.Config `yaml:"disk"`
}

type Service struct {
	cfg    *Config
	logger *zap.SugaredLogger

	disks map[uint32]common.DiskIF
}

func New(cfg *Config, logger *zap.SugaredLogger) (*Service, error) {
	disks, err := disk.Open(cfg.Disk)
	if err != nil {
		return nil, err
	}

	err = registerNewDisks(cfg.Host, cfg.MgrHost, disks)
	if err != nil {
		return nil, err
	}

	disksMap, err := filterDisks(cfg.Host, cfg.MgrHost, disks)
	if err != nil {
		return nil, err
	}

	return &Service{cfg: cfg, logger: logger, disks: disksMap}, nil
}

func registerNewDisks(host, mgrHost string, disks []common.DiskIF) error {
	ctx := context.TODO()

	newDisks := make([]*proto.Disk, 0)
	for _, d := range disks {
		stat, err := d.Stat(ctx)
		if err != nil {
			return err
		}
		if stat.DiskLabel == nil {
			label, err := allocDiskLabel(mgrHost)
			if err != nil {
				return err
			}
			err = d.Init(ctx, label)
			if err != nil {
				return err
			}
			dd := &proto.Disk{
				DiskLabel: label,
				Host:      host,
				Path:      stat.Path,
				Capacity:  stat.Capacity,
				Available: stat.Capacity,
				State:     proto.DiskStateNormal,
			}
			newDisks = append(newDisks, dd)
		}
	}
	if len(newDisks) == 0 {
		return nil
	}
	failedDisk, err := registerDisks(mgrHost, newDisks)
	if err != nil {
		return err
	}
	if len(failedDisk) > 0 {
		zap.S().Error(failedDisk)
	}
	return nil
}

func filterDisks(host, mgrHost string, disks []common.DiskIF) (map[uint32]common.DiskIF, error) {
	ctx := context.TODO()

	registeredDisks, err := listDisks(host, mgrHost)
	if err != nil {
		return nil, err
	}
	disksMap := make(map[uint32]common.DiskIF)
	for _, d := range disks {
		stat, err := d.Stat(ctx)
		if err != nil {
			return nil, err
		}
		if findDisk(registeredDisks, stat.DiskLabel) {
			disksMap[stat.DiskID] = d
		}
	}
	return disksMap, nil
}

func findDisk(disks []*proto.Disk, label *proto.DiskLabel) bool {
	for _, d := range disks {
		if label.ClusterID == d.ClusterID && label.DiskID == d.DiskID {
			return true
		}
	}
	return false
}

func listDisks(host, mgrHost string) ([]*proto.Disk, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/disks?host=%s", mgrHost, host))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := make([]*proto.Disk, 0)
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func allocDiskLabel(host string) (*proto.DiskLabel, error) {
	resp, err := http.Post(fmt.Sprintf("http://%s/allocdisklabel", host), "", nil)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := new(proto.DiskLabel)
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func registerDisks(host string, disks []*proto.Disk) ([]*proto.Disk, error) {
	diskBuf, err := json.Marshal(disks)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/registerdisks", host), "application/json", bytes.NewBuffer(diskBuf))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	ret := make([]*proto.Disk, 0)
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (s *Service) RegisterRouters(r *chi.Mux) {
	r.Get("/get", s.Get)
	r.Post("/put", s.Put)
	r.Post("/delete", s.Delete)
}

package service

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/ocean-rw/ocean/internal/daemon/disk"
	"github.com/ocean-rw/ocean/internal/daemon/disk/common"
	"github.com/ocean-rw/ocean/pkg/api/master"
	"github.com/ocean-rw/ocean/pkg/proto"
)

type Config struct {
	MyHost string         `yaml:"my_host"`
	Disk   *disk.Config   `yaml:"disk"`
	Master *master.Config `yaml:"master"`
}

type Service struct {
	cfg    *Config
	logger *zap.SugaredLogger

	mgr   *master.Master
	disks map[uint32]common.DiskIF
}

func New(cfg *Config, logger *zap.SugaredLogger) (*Service, error) {
	mgr, err := master.New(cfg.Master)
	if err != nil {
		return nil, err
	}

	disks, err := disk.Open(cfg.Disk)
	if err != nil {
		return nil, err
	}

	svc := &Service{cfg: cfg, logger: logger, mgr: mgr}

	failedDisks, err := svc.registerNewDisks(disks)
	if err != nil {
		return nil, err
	}
	if len(failedDisks) != 0 {
		for _, d := range failedDisks {
			logger.Errorf("failed to register disks %d of cluster %s", d.DiskID, d.ClusterID)
		}
	}

	disksMap, err := svc.filterDisks(disks)
	if err != nil {
		return nil, err
	}

	return &Service{cfg: cfg, logger: logger, disks: disksMap}, nil
}

func (s *Service) registerNewDisks(disks []common.DiskIF) ([]*proto.Disk, error) {
	ctx := context.TODO()
	newDisks := make([]*proto.Disk, 0)
	for _, d := range disks {
		stat, err := d.Stat(ctx)
		if err != nil {
			return nil, err
		}
		if stat.DiskLabel == nil {
			label, err := s.mgr.AllocDiskLabel()
			if err != nil {
				return nil, err
			}
			err = d.Init(ctx, label)
			if err != nil {
				return nil, err
			}
			dd := &proto.Disk{
				DiskLabel: label,
				Host:      s.cfg.MyHost,
				Path:      stat.Path,
				Capacity:  stat.Capacity,
				Available: stat.Capacity,
				State:     proto.DiskStateNormal,
			}
			newDisks = append(newDisks, dd)
		}
	}
	if len(newDisks) == 0 {
		return nil, nil
	}
	return s.mgr.RegisterDisks(newDisks)
}

func (s *Service) filterDisks(disks []common.DiskIF) (map[uint32]common.DiskIF, error) {
	ctx := context.TODO()

	registeredDisks, err := s.mgr.ListDisks(&master.ListDisksArgs{Host: s.cfg.MyHost})
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

func (s *Service) RegisterRouters(r *chi.Mux) {
	r.Get("/get", s.Get)
	r.Post("/put", s.Put)
	r.Post("/delete", s.Delete)
}

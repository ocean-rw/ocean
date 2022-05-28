package service

import (
	"errors"
	"net/http"

	json "github.com/json-iterator/go"
	"github.com/momaek/easybind"

	"github.com/ocean-rw/ocean/internal/lib/httputil"
	"github.com/ocean-rw/ocean/pkg/api/master"
	"github.com/ocean-rw/ocean/pkg/proto"
)

func (s *Service) AllocDiskLabel(w http.ResponseWriter, r *http.Request) {
	diskID, err := s.db.IDTable.AllocDiskID(r.Context())
	if err != nil {
		s.logger.Errorf("failed to alloc disk id, err: %s", err)
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	err = httputil.ReplyJSON(w, 0, &proto.DiskLabel{ClusterID: s.clusterID, DiskID: diskID})
	if err != nil {
		s.logger.Errorf("failed to reply data, err: %s", err)
	}
}

func (s *Service) RegisterDisks(w http.ResponseWriter, r *http.Request) {
	disks := make([]*proto.Disk, 0)
	err := json.NewDecoder(r.Body).Decode(&disks)
	if err != nil {
		s.logger.Errorf("failed to read body, err: %s", err)
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	defer r.Body.Close()

	failedDisks := make([]*proto.Disk, 0)
	for _, disk := range disks {
		if disk.ClusterID != s.clusterID {
			s.logger.Errorf("unmatched cluster id, disk: %s, cluster: %s", disk.ClusterID, s.clusterID)
			failedDisks = append(failedDisks, disk)
			continue
		}
		err = s.db.DiskTable.Insert(r.Context(), disk)
		if err != nil {
			s.logger.Errorf("failed to insert disk to database, err: %s", err)
			failedDisks = append(failedDisks, disk)
			err = nil
			continue
		}
	}

	code := 0
	switch {
	case len(failedDisks) == 0:
		code = http.StatusOK
	case len(failedDisks) < len(disks):
		code = http.StatusPartialContent
	case len(failedDisks) == len(disks):
		code = http.StatusBadRequest
	}
	err = httputil.ReplyJSON(w, code, failedDisks)
	if err != nil {
		s.logger.Errorf("failed to reply data, err: %s", err)
	}
}

func (s *Service) ListDisks(w http.ResponseWriter, r *http.Request) {
	listArgs := new(master.ListDisksArgs)
	err := easybind.Bind(r, listArgs)
	if err != nil {
		httputil.ReplyErr(w, http.StatusBadRequest, nil)
		return
	}

	if listArgs.State >= proto.DiskStateMAX {
		httputil.ReplyErr(w, http.StatusBadRequest, errors.New("invalid disk state"))
		return
	}

	disks, err := s.db.DiskTable.List(r.Context(), listArgs)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	err = httputil.ReplyJSON(w, http.StatusOK, disks)
	if err != nil {
		s.logger.Errorf("failed to reply data, err: %s", err)
	}
}

func (s *Service) GetDisk(w http.ResponseWriter, r *http.Request) {
	args := &struct {
		DiskID uint32 `pos:"query=disk_id"`
	}{}
	err := easybind.Bind(r, args)
	if err != nil || args.DiskID == 0 {
		httputil.ReplyErr(w, http.StatusBadRequest, nil)
		return
	}

	disks, err := s.db.DiskTable.Get(r.Context(), args.DiskID)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	err = httputil.ReplyJSON(w, http.StatusOK, disks)
	if err != nil {
		s.logger.Errorf("failed to reply data, err: %s", err)
	}
}

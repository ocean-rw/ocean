package mgr

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/ocean-rw/ocean/internal/lib/httputil"
	"github.com/ocean-rw/ocean/pkg/proto"
)

func (m *Mgr) AllocDiskLabel(w http.ResponseWriter, r *http.Request) {
	diskID, err := m.db.IDTable.AllocDiskID(r.Context())
	if err != nil {
		m.logger.Errorf("failed to alloc disk id, err: %s", err)
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	err = httputil.ReplyData(w, 0, &proto.DiskLabel{ClusterID: m.clusterID, DiskID: diskID})
	if err != nil {
		m.logger.Errorf("failed to reply data, err: %s", err)
	}
}

func (m *Mgr) RegisterDisks(w http.ResponseWriter, r *http.Request) {
	disks := make([]*proto.Disk, 0)
	err := json.NewDecoder(r.Body).Decode(&disks)
	if err != nil {
		m.logger.Errorf("failed to read body, err: %s", err)
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	defer r.Body.Close()

	failedDisks := make([]*proto.Disk, 0)
	for _, disk := range disks {
		if disk.ClusterID != m.clusterID {
			m.logger.Errorf("unmatched cluster id, disk: %s, cluster: %s", disk.ClusterID, m.clusterID)
			failedDisks = append(failedDisks, disk)
			continue
		}
		err = m.db.DiskTable.Insert(r.Context(), disk)
		if err != nil {
			m.logger.Errorf("failed to insert disk to database, err: %s", err)
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
	err = httputil.ReplyData(w, code, failedDisks)
	if err != nil {
		m.logger.Errorf("failed to reply data, err: %s", err)
	}
}

func (m *Mgr) ListDisks(w http.ResponseWriter, r *http.Request) {
	listArgs := &proto.ListDisksArgs{}

	host := r.URL.Query().Get("host")
	if host != "" {
		listArgs.Host = host
	}
	state := r.URL.Query().Get("state")
	if state != "" {
		s, err := strconv.ParseUint(state, 10, 8)
		if err != nil {
			httputil.ReplyErr(w, http.StatusBadRequest, err)
			return
		}
		listArgs.State = (proto.DiskState)(s)
	}
	if listArgs.State >= proto.DiskStateMAX {
		httputil.ReplyErr(w, http.StatusBadRequest, errors.New("invalid disk state"))
		return
	}

	disks, err := m.db.DiskTable.List(r.Context(), listArgs)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	err = httputil.ReplyData(w, http.StatusOK, disks)
	if err != nil {
		m.logger.Errorf("failed to reply data, err: %s", err)
	}
}

func (m *Mgr) GetDisk(w http.ResponseWriter, r *http.Request) {
	diskID := chi.URLParam(r, "disk_id")
	id, err := strconv.ParseUint(diskID, 10, 32)
	if err != nil {
		httputil.ReplyErr(w, http.StatusBadRequest, err)
		return
	}
	disks, err := m.db.DiskTable.Get(r.Context(), uint32(id))
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	err = httputil.ReplyData(w, http.StatusOK, disks)
	if err != nil {
		m.logger.Errorf("failed to reply data, err: %s", err)
	}
}

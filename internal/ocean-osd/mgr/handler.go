package mgr

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/ocean-rw/ocean/internal/lib/httputil"
)

func (m *Mgr) Get(w http.ResponseWriter, r *http.Request) {
	diskID, fd := chi.URLParam(r, "disk_id"), chi.URLParam(r, "fd")

	id, err := strconv.ParseUint(diskID, 10, 32)
	if err != nil {
		httputil.ReplyErr(w, http.StatusBadRequest, errors.New("invalid disk id"))
		return
	}
	disk := m.disks[uint32(id)]

	ctx := context.TODO()
	data, err := disk.Get(ctx, fd)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	_, err = io.Copy(w, data)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
}

func (m *Mgr) Put(w http.ResponseWriter, r *http.Request) {
	diskID, fd := chi.URLParam(r, "disk_id"), chi.URLParam(r, "fd")

	id, err := strconv.ParseUint(diskID, 10, 32)
	if err != nil {
		httputil.ReplyErr(w, http.StatusBadRequest, errors.New("invalid disk id"))
		return
	}
	disk := m.disks[uint32(id)]

	ctx := context.TODO()
	err = disk.Put(ctx, fd, r.Body)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Mgr) Delete(w http.ResponseWriter, r *http.Request) {
	diskID, fd := chi.URLParam(r, "disk_id"), chi.URLParam(r, "fd")

	id, err := strconv.ParseUint(diskID, 10, 32)
	if err != nil {
		httputil.ReplyErr(w, http.StatusBadRequest, errors.New("invalid disk id"))
		return
	}
	disk := m.disks[uint32(id)]

	ctx := context.TODO()
	err = disk.Delete(ctx, fd)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

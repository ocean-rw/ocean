package service

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/momaek/easybind"

	"github.com/ocean-rw/ocean/internal/lib/httputil"
	"github.com/ocean-rw/ocean/pkg/api/osd"
)

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	args := new(osd.Args)
	err := easybind.Bind(r, args)
	if err != nil {
		httputil.ReplyErr(w, http.StatusBadRequest, nil)
		return
	}
	disk, ok := s.disks[args.DiskID]
	if !ok {
		httputil.ReplyErr(w, http.StatusBadRequest, errors.New("no such disk"))
		return
	}

	ctx := context.TODO()
	data, err := disk.Get(ctx, args.FD)
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

func (s *Service) Put(w http.ResponseWriter, r *http.Request) {
	args := new(osd.Args)
	err := easybind.Bind(r, args)
	if err != nil {
		httputil.ReplyErr(w, http.StatusBadRequest, nil)
		return
	}
	disk, ok := s.disks[args.DiskID]
	if !ok {
		httputil.ReplyErr(w, http.StatusBadRequest, errors.New("no such disk"))
		return
	}

	ctx := context.TODO()
	err = disk.Put(ctx, args.FD, r.Body)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	args := new(osd.Args)
	err := easybind.Bind(r, args)
	if err != nil {
		httputil.ReplyErr(w, http.StatusBadRequest, nil)
		return
	}
	disk, ok := s.disks[args.DiskID]
	if !ok {
		httputil.ReplyErr(w, http.StatusBadRequest, errors.New("no such disk"))
		return
	}

	ctx := context.TODO()
	err = disk.Delete(ctx, args.FD)
	if err != nil {
		httputil.ReplyErr(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

package storage

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Local struct {
	fid  uint64
	path string

	sync.Mutex
}

func NewLocal(path string) (Storage, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	fid := uint64(0)
	for _, f := range files {
		idStr := strings.Split(f.Name(), ".")[0]
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			continue
		}
		if id > fid {
			fid = id
		}
	}
	return &Local{fid: fid, path: path}, nil
}

func (l *Local) Put(_ context.Context, data io.Reader) (string, int64, error) {
	fd := strconv.FormatUint(l.FID(), 10)
	filename := fmt.Sprintf("%s/%s", l.path, fd)
	f, err := os.OpenFile(filename+".tmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()

	n, err := io.Copy(f, data)
	if err != nil {
		return "", 0, err
	}

	err = os.Rename(filename+".tmp", filename)
	if err != nil {
		return "", 0, err
	}
	return fd, n, nil
}

func (l *Local) Get(_ context.Context, fd string) (io.ReadCloser, error) {
	filename := fmt.Sprintf("%s%s", l.path, fd)
	f, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (l *Local) Delete(_ context.Context, fd string) error {
	filename := fmt.Sprintf("%s/%s", l.path, fd)
	return os.Remove(filename)
}

func (l *Local) FID() uint64 {
	l.Lock()
	l.Unlock()
	l.fid++
	return l.fid
}

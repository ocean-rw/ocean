package daemon

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeoutMS = 1000

var (
	ErrEmptyConfig  = errors.New("empty config")
	ErrGetFailed    = errors.New("upload failed")
	ErrUploadFailed = errors.New("upload failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type Config struct {
	Host      string `yaml:"host"`
	TimeoutMS int64  `yaml:"timeout_ms"`
}

type Daemon struct {
	Host   string
	Client *http.Client
}

func New(cfg *Config) (*Daemon, error) {
	if cfg == nil {
		return nil, ErrEmptyConfig
	}
	if cfg.TimeoutMS <= 0 {
		cfg.TimeoutMS = defaultTimeoutMS
	}
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Duration(cfg.TimeoutMS) * time.Millisecond,
	}
	return &Daemon{Client: client, Host: cfg.Host}, nil
}

type Args struct {
	DiskID uint32 `pos:"query:disk_id"`
	FD     string `pos:"query:fd"`
}

func (d *Daemon) Get(args *Args) (io.ReadCloser, error) {
	url := fmt.Sprintf("http://%s/get?disk_id=%d&fd=%s", d.Host, args.DiskID, args.FD)
	resp, err := d.Client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrGetFailed
	}
	return resp.Body, nil
}

func (d *Daemon) Put(args *Args, data io.ReadCloser) error {
	url := fmt.Sprintf("http://%s/put?disk_id=%d&fd=%s", d.Host, args.DiskID, args.FD)
	resp, err := d.Client.Post(url, "application/octet-stream", data)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return ErrUploadFailed
	}
	return nil
}

func (d *Daemon) Delete(args *Args) error {
	url := fmt.Sprintf("http://%s/delete?disk_id=%d&fd=%s", d.Host, args.DiskID, args.FD)
	resp, err := d.Client.Post(url, "", nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return ErrDeleteFailed
	}
	return nil
}

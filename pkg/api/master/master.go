package master

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"time"

	json "github.com/json-iterator/go"

	"github.com/ocean-rw/ocean/pkg/proto"
)

const defaultTimeoutMS = 1000

var ErrEmptyConfig = errors.New("empty config")

type Config struct {
	Host      string `yaml:"host"`
	TimeoutMS int64  `yaml:"timeout_ms"`
}

type Master struct {
	Host string
	*http.Client
}

func New(cfg *Config) (*Master, error) {
	if cfg == nil {
		return nil, ErrEmptyConfig
	}
	if cfg.TimeoutMS <= 0 {
		cfg.TimeoutMS = defaultTimeoutMS
	}
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   time.Duration(cfg.TimeoutMS) * time.Millisecond,
	}
	return &Master{Client: client, Host: cfg.Host}, nil
}

func (m *Master) AllocStripes() ([]*proto.Stripe, error) {
	url := fmt.Sprintf("http://%s/allocstripes", m.Host)
	resp, err := m.Post(url, "", nil)
	if err != nil {
		return nil, err
	}
	ret := make([]*proto.Stripe, 0)
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (m *Master) AllocDiskLabel() (*proto.DiskLabel, error) {
	url := fmt.Sprintf("http://%s/allocdisklabel", m.Host)
	resp, err := m.Post(url, "", nil)
	if err != nil {
		return nil, err
	}
	ret := new(proto.DiskLabel)
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (m *Master) RegisterDisks(disks []*proto.Disk) ([]*proto.Disk, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&disks)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s/registerdisks", m.Host)
	resp, err := m.Post(url, "application/json", &buf)
	if err != nil {
		return nil, err
	}
	ret := make([]*proto.Disk, 0)
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type ListDisksArgs struct {
	Host  string          `in:"query=host"`
	State proto.DiskState `in:"query=state"`
}

func (m *Master) ListDisks(args *ListDisksArgs) ([]*proto.Disk, error) {
	url := fmt.Sprintf("http://%s/listdisks?host=%s&state=%d", m.Host, args.Host, args.State)
	resp, err := m.Get(url)
	if err != nil {
		return nil, err
	}
	ret := make([]*proto.Disk, 0)
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (m *Master) GetDisk(diskID uint32) (*proto.Disk, error) {
	url := fmt.Sprintf("http://%s/getdisk?disk_id=%d", m.Host, diskID)
	resp, err := m.Get(url)
	if err != nil {
		return nil, err
	}
	ret := new(proto.Disk)
	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type AllocDisksArgs struct {
}

func (m *Master) AllocDisks() {

}

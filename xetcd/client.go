package xetcd

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"strings"
)

type XClientV3 struct {
	Client *clientv3.Client
	root   string
}

func NewEtcdCli(opt *EtcdConfig) (*XClientV3, error) {
	cli, err := NewEtcdClientV3(opt)
	if err != nil {
		return nil, err
	}
	return &XClientV3{
		Client: cli,
		root:   opt.Root,
	}, nil
}
func NewCli(cli *clientv3.Client, root string) *XClientV3 {
	return &XClientV3{
		Client: cli,
		root:   root,
	}
}
func (h *XClientV3) GetRoot() string {
	if h.Client == nil {
		return ""
	}
	return strings.TrimSuffix(h.root, "/")
}
func (h *XClientV3) Get(key string) (string, error) {
	if h.Client == nil {
		return "", errors.New("没有连接etcd")
	}

	res, err := h.Client.Get(context.TODO(), h.GetRoot()+"/"+key)
	if err != nil {
		return "", err
	}
	if res.Count < 1 {
		return "", errors.New("empty")
	}
	return string(res.Kvs[0].Value), nil

}
func (h *XClientV3) Set(key, val string) error {
	if h.Client == nil {
		return errors.New("没有连接etcd")
	}
	_, err := h.Client.Put(context.TODO(), h.GetRoot()+"/"+key, val)
	return err
}
func (h *XClientV3) WatchDirNode(key string, action WatchDirNodeFunc) error {
	path := fmt.Sprintf("%s/%s", h.GetRoot(), key)
	return WatchDirNode(context.Background(), h.Client, path, action)
}

func (h *XClientV3) WatchNode(key string, action WatchNodeFunc) error {
	path := fmt.Sprintf("%s/%s", h.GetRoot(), key)
	return WatchNode(context.Background(), h.Client, path, action)
}

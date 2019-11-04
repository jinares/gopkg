package xetcd

import (
	"context"
	"errors"
	"go.etcd.io/etcd/clientv3"
	"strings"
)

type XClientV3 struct {
	Client *clientv3.Client
	conf   *EtcdConfig
}

func NewEtcdCli(opt *EtcdConfig) (*XClientV3, error) {
	cli, err := NewEtcdClientV3(opt)
	if err != nil {
		return nil, err
	}
	return &XClientV3{
		Client: cli,
		conf:   opt,
	}, nil
}
func (h *XClientV3) GetRoot() string {
	if h.Client == nil {
		return ""
	}
	return strings.TrimSuffix(h.conf.Root, "/")
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
		return "", nil
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

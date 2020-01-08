package xetcd

import (
	"context"
	"errors"

	"github.com/jinares/gopkg/xlog"
	"go.etcd.io/etcd/clientv3/concurrency"
)

//NewMutex NewMutex
func NewMutex(cli *XClientV3, key string, f func() error, sopts ...concurrency.SessionOption) error {
	if cli == nil || cli.Client == nil {
		return errors.New("etcdv3 cli nil")
	}
	session, err := concurrency.NewSession(cli.Client, sopts...)
	if err != nil {
		return err
	}

	defer session.Close()

	mu := concurrency.NewMutex(session, cli.GetRoot()+"/"+key)

	if err := mu.Lock(context.TODO()); err != nil {
		return err
	}
	defer func() {
		err := mu.Unlock(context.TODO())
		if err != nil {

		}
	}()

	return f()

}

//NewElection NewElection
func NewElection(cli *XClientV3, key, val string, f func() error, sopts ...concurrency.SessionOption) error {
	if cli == nil || cli.Client == nil {
		return errors.New("etcdv3 cli nil")
	}
	session, err := concurrency.NewSession(cli.Client, sopts...)
	if err != nil {
		return err
	}
	defer session.Close()
	e1 := concurrency.NewElection(session, cli.GetRoot()+"/"+key)
	if err := e1.Campaign(context.Background(), val); err != nil {
		return err
	}
	defer func() {
		err := e1.Resign(context.TODO())
		if err != nil {
			xlog.GetLog().Error(err.Error())
		}
	}()

	return f()
}

package xetcd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jinares/gopkg/xtools"
	"go.etcd.io/etcd/clientv3"
)

type (
	Option struct {
		Key     string
		Convert ConvertFunc
	}
	tmpOption struct {
		Op   Option
		Data map[string]string
	}

	ConvertFunc func(data interface{}) error
	StrConvert  func(data string)
	MapConvert  func(data map[string]string)

	ObjectConvert    func(data map[string]interface{})
	WatchNodeFunc    func(val string) error
	WatchDirNodeFunc func(key, val string) error
)

func Watch(client *clientv3.Client, root string, op ...Option) error {
	root = strings.TrimSuffix(root, "/")
	mop := map[string]Option{}
	for _, val := range op {
		mop[val.Key] = val
	}
	res, err := client.Get(context.Background(), root, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	match(root, mop, res)
	go func() {
		for {
			cc := client.Watch(context.Background(), root+"/", clientv3.WithPrefix())
			for wresp := range cc {

				for _, ev := range wresp.Events {
					k, _ := split(root, string(ev.Kv.Key))
					ress, err := client.Get(context.Background(), root+"/"+k, clientv3.WithPrefix())
					if err != nil {
						continue
					}
					match(root, mop, ress)
				}
			}
		}
	}()
	return nil
}

func WatchNode(ctx context.Context, client *clientv3.Client, path string, action WatchNodeFunc) error {
	val, err := GetNode(ctx, client, path)
	if err == nil && val != "" {
		action(val)
	}
	go func() {
		for {
			cc := client.Watch(context.Background(), path)
			for wres := range cc {
				fmt.Println(wres)
				val, err := GetNode(ctx, client, path)
				if err == nil && val != "" {
					action(val)
				}
			}
		}
	}()
	return nil

}

func WatchDirNode(ctx context.Context, client *clientv3.Client, path string, action WatchDirNodeFunc) error {
	val, err := GetDirNode(ctx, client, path)
	if err == nil {
		for k, v := range val {
			action(k, v)
		}
	}
	go func() {
		for {
			cc := client.Watch(context.Background(), path, clientv3.WithPrefix())
			for wresp := range cc {

				for _, ev := range wresp.Events {

					if len(ev.Kv.Value) < 1 {
						continue
					}
					k, sk := split(path, string(ev.Kv.Key))
					if sk != "" {
						fmt.Println(ev.Kv)
						continue
					}
					action(k, string(ev.Kv.Value))

				}
			}
		}
	}()
	return nil

}

//GetNode
func GetNode(ctx context.Context, client *clientv3.Client, key string) (string, error) {
	res, err := client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if res.Count < 1 {
		return "", errors.New("empty")
	}
	return string(res.Kvs[0].Value), nil
}
func GetDirNode(ctx context.Context, client *clientv3.Client, path string) (map[string]string, error) {

	res, err := client.Get(context.Background(), path, clientv3.WithPrefix())
	if err != nil {
		return map[string]string{}, err
	} else {
		data := map[string]string{}
		for _, v := range res.Kvs {
			key, subkey := split(path, string(v.Key))
			if key == "" {
				continue
			}
			if subkey != "" {
				continue
			}
			if strings.HasSuffix(string(v.Key), "/") {
				continue
			}
			data[key] = string(v.Value)
		}
		return data, nil
	}
	return nil, errors.New("")
}

func match(root string, mop map[string]Option, res *clientv3.GetResponse) {
	tmp := map[string]tmpOption{}
	for _, v := range res.Kvs {
		key, subkey := split(root, string(v.Key))
		op, isok := mop[key]
		if isok == false {
			continue
		}
		if op.Convert == nil {
			continue
		}
		if subkey == "" {
			op.Convert(string(v.Value))
		} else {
			if vop, isok := tmp[key]; isok {
				vop.Data[subkey] = string(v.Value)
			} else {
				tmp[key] = tmpOption{Op: op, Data: map[string]string{
					subkey: string(v.Value),
				}}
			}
		}
	}
	for _, val := range tmp {
		val.Op.Convert(val.Data)
	}
}
func split(root, path string) (key string, subkey string) {

	data := strings.TrimPrefix(path, root)
	data = strings.TrimSuffix(strings.TrimPrefix(data, "/"), "/")
	if data == "" {
		return
	}
	arr := strings.Split(data, "/")
	if len(arr) < 1 || len(arr) > 2 {
		return
	}
	key = arr[0]
	subkey = strings.Join(arr[1:], "/")
	return
}

func StringTo(action StrConvert) ConvertFunc {
	return func(data interface{}) error {
		sdata, isok := data.(string)
		if isok == false {
			return errors.New("类型错误")
		}
		action(sdata)
		return nil
	}
}

func ObjectTo(out interface{}) ConvertFunc {
	return func(data interface{}) error {
		switch sdata := data.(type) {
		case string:
			xtools.ToJSON(sdata, out)
			return nil
		case map[string]string:
			vdata := map[string]interface{}{}
			for k, v := range sdata {
				var i interface{}
				if xtools.ToJSON(v, &i) == nil {
					vdata[k] = i
				} else {
					vdata[k] = v
				}
			}
			xtools.ToJSON(xtools.JSONToStr(vdata), out)

			return nil
		default:

			return errors.New("类型错误")

		}

	}
}

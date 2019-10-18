package xetcd

import (
	"context"
	"errors"
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

	ObjectConvert func(data map[string]interface{})
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
func MapTo(action MapConvert) ConvertFunc {
	return func(data interface{}) error {
		sdata, isok := data.(map[string]string)
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

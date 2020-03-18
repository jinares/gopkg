package xtools

import (
	"errors"
	"sync"
)

type (
	//MapFunc UpdateMapFunc
	MapFunc func() error
)

var (
	maplock sync.RWMutex
)

//GetMapVal GetMapVal
func GetMapVal(ac MapFunc) error {
	if ac == nil {
		return errors.New("empty")
	}
	maplock.RLock()
	defer maplock.RUnlock()
	return ac()
}

//UpdateMap UpdateMap
func UpdateMap(ac MapFunc) error {
	if ac == nil {
		return nil
	}
	maplock.Lock()
	defer maplock.Unlock()
	return ac()
}

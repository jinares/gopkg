package xtools

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type (
	//LazyCache
	lazyCache struct {
		Key string
		Val chan interface{}
	}
	// LazyCacheArray  lazy array
	lazyCacheArray []*lazyCache
	//LazyList list
	lazyList struct {
		data *list.List
		lock sync.Mutex
	}
)

//GetKeys GetKeys
func (h *lazyCacheArray) GetKeys() []string {
	keys := []string{}
	for _, item := range *h {
		keys = append(keys, item.Key)
	}
	return keys
}

//Len len
func (h *lazyList) Len() int {
	return h.data.Len()
}

//Push Push
func (h *lazyList) Push(item *lazyCache) error {
	if h.data == nil {
		return errors.New("init")
	}
	h.lock.Lock()
	h.data.PushBack(item)
	h.lock.Unlock()
	return nil
}

//PopItem pop item
func (h *lazyList) PopItem() *lazyCache {

	if h.data == nil {
		return nil
	}
	item := h.data.Front()
	if item == nil {
		return nil
	}
	h.data.Remove(item)
	switch vv := item.Value.(type) {
	case *lazyCache:
		return vv

	default:
		return nil
	}
}

//PopItems Pop
func (h *lazyList) PopItems(n int) lazyCacheArray {

	if n < 1 {
		return []*lazyCache{}
	}
	h.lock.Lock()
	defer h.lock.Unlock()

	data := []*lazyCache{}
	for i := 0; i < n; i++ {
		item := h.PopItem()
		if item == nil {
			return data
		}
		data = append(data, item)
	}
	return data

}

//NewLazyList new list
func newLazyList() *lazyList {
	return &lazyList{
		data: list.New(),
	}
}

// NewLazy   new lazy
func newLazy(key string) *lazyCache {
	return &lazyCache{
		Key: key, Val: make(chan interface{}),
	}
}

type (

	//LazyTopics topics
	LazyTopics struct {
		list    *lazyList
		handler LazyHandler
	}
	//LazyHandler handler
	LazyHandler func(key ...string) (map[string]interface{}, error)
)

//NewTopics new lazy
func NewTopics(batch, goroutine int, action LazyHandler) *LazyTopics {
	topic := &LazyTopics{list: newLazyList(), handler: action}
	// goroutine
	if goroutine < 1 {
		goroutine = 1
	}
	if batch < 1 {
		batch = 1
	}
	lazyAryChan := make(chan lazyCacheArray, batch*goroutine*2)
	free := func(items lazyCacheArray) {
		for _, item := range items {
			item.Val <- ""
			close(item.Val)
		}
	}
	pushfunc := func(items lazyCacheArray) {

		t := time.NewTimer(200 * time.Millisecond)
		defer t.Stop()
		select {
		case lazyAryChan <- items:
			return
		case <-t.C:
			all := topic.list.Len()
			if all > (batch * 10) {
				data := topic.list.PopItems(all - batch*10)
				free(data)
			}
			free(items)

		}
	}
	go func() {
		for {
			data := topic.list.PopItems(batch)
			if len(data) < 1 {
				time.Sleep(time.Millisecond * 20)
				continue
			}
			pushfunc(data)
		}
	}()

	for i := 0; i < goroutine; i++ {
		go func() {
			for {
				data := <-lazyAryChan
				//fmt.Println("len:", len(data))
				keys := data.GetKeys()
				ret, err := action(keys...)
				if err != nil {
					continue
				}
				for _, item := range data {
					if val, isok := ret[item.Key]; isok {
						item.Val <- val
					} else {
						item.Val <- ""
					}
					close(item.Val)
				}
			}

		}()
	}
	return topic
}

//Get get val
func (h *LazyTopics) Get(key string, d time.Duration) (interface{}, error) {
	if h.list == nil {
		return nil, errors.New("没有初始化")
	}
	item := newLazy(key)
	err := h.list.Push(item)
	if err != nil {
		return nil, err
	}
	t := time.NewTimer(d)
	defer func() {
		t.Stop()
	}()
	select {
	case data := <-item.Val:
		return data, nil
	case <-t.C:
		return nil, errors.New("timeout")

	}
}

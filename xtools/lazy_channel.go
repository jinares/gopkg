package xtools

import (
	"errors"
	"time"
)

type (
	lazyChan struct {
		data chan *lazyCache
	}
)

func newLazyChan(num int) *lazyChan {
	return &lazyChan{data: make(chan *lazyCache, num)}
}
func (h *lazyChan) Push(item *lazyCache) error {

	select {
	case h.data <- item:

		return nil
	default:
		return errors.New("timeout")

	}
}
func (h *lazyChan) Pop(num int) lazyCacheArray {
	t := time.NewTimer(20 * time.Millisecond)
	defer func() {
		t.Stop()
	}()
	data := lazyCacheArray{}
	for i := 0; i < num; i++ {
		select {
		case item := <-h.data:
			data = append(data, item)
		case <-t.C:
			return data

		}
	}
	return data
}

type (
	//LazyTopics topics
	LazyTopics2 struct {
		list    *lazyChan
		handler LazyHandler
		batch   int
	}
)

func (h *LazyTopics2) getBatch() int {
	if h.batch < 1 {
		return 1
	}
	return h.batch
}
func (h *LazyTopics2) push(item *lazyCache) error {
	num := 0
	for {
		if num > h.getBatch() {
			return errors.New("timeout")
		}
		num++
		if h.list.Push(item) == nil {
			return nil
		}
		data := h.list.Pop(h.getBatch())
		if len(data) < 1 {
			continue
		}
		h.run(data)
	}
}
func (h *LazyTopics2) run(data lazyCacheArray) error {
	if h.handler == nil {
		h.handler = func(key ...string) (m map[string]interface{}, err error) {
			return nil, errors.New("handler is nil")
		}
	}
	ret, err := h.handler(data.GetKeys()...)
	if err != nil {
		return err
	}
	for _, item := range data {
		if val, isok := ret[item.Key]; isok {
			item.Val <- val
		} else {
			item.Val <- nil
		}
		close(item.Val)
	}
	return nil
}

//NewTopics new lazy
func NewTopics2(batch, goroutine int, action LazyHandler) *LazyTopics2 {
	topic := &LazyTopics2{list: newLazyChan(batch * 100), batch: batch, handler: action}
	// goroutine
	if goroutine < 1 {
		goroutine = 1
	}
	for i := 0; i < goroutine; i++ {
		go func() {
			for {
				data := topic.list.Pop(topic.getBatch())
				if len(data) < 1 {
					//time.Sleep(time.Millisecond * 200)
					continue
				}

				//fmt.Println("len:", len(data))
				topic.run(data)
			}
		}()
	}

	return topic
}

//Get get val
func (h *LazyTopics2) Get(key string, d time.Duration) (interface{}, error) {
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

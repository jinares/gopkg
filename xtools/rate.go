package xtools

import (
	"sync"
	"time"

	"github.com/bluele/gcache"
)

type (
	RateInter interface {
		Rate() float64
		Succ()
		Fail()
	}
	AcRate struct {
		lock    *sync.RWMutex
		succ    gcache.Cache
		fail    gcache.Cache
		keyfunc RateKeyFunc
	}
	RateKeyFunc func() interface{}
	ActionFunc  func() error

	Breaker struct {
		id          string
		status      int     // 0 熔断器关闭状态 1 熔断器打开 2  半熔断
		circuitTime int     // 熔断时间 s
		CircuitRate float64 //熔断的成功率
		recoverNum  int     //半熔断 成功N 次恢复
		ctrl        RateInter
	}
)

func NewAcRate(sec int) *AcRate {
	return &AcRate{
		succ: gcache.New(sec * 2).Expiration(time.Duration(sec) * time.Second).LRU().Build(),
		fail: gcache.New(sec * 2).Expiration(time.Duration(sec) * time.Second).LRU().Build(),
		keyfunc: func() interface{} {
			return time.Now().Unix()
		},
		lock: &sync.RWMutex{},
	}
}
func (h *AcRate) SetKeyFunc(ac RateKeyFunc) {
	if ac == nil {
		return
	}
	h.keyfunc = ac
}
func (h *AcRate) Rate() float64 {
	succdata := h.succ.GetALL(true)
	faildata := h.fail.GetALL(true)

	all, fail := int64(0), int64(0)
	for _, v := range succdata {
		val, isok := v.(int64)
		if isok {
			all = all + val
		}
	}
	for _, v := range faildata {
		val, isok := v.(int64)
		if isok {
			all = all + val
			fail = fail + val
		}
	}
	if all < 1 {
		return 100
	} else if fail == all {
		return 0
	}
	return float64((all-fail)*100) / float64(all)
}
func (h *AcRate) Succ() {
	h.incr(h.succ)
}

func (h *AcRate) Fail() {
	h.incr(h.fail)
}
func (h *AcRate) incr(c gcache.Cache) {
	dt := h.keyfunc()
	h.lock.Lock()
	defer h.lock.Unlock()
	val, err := c.Get(dt)
	if err != nil {
		c.Set(dt, int64(1))
		return
	}
	v, isok := val.(int64)
	if !isok {
		v = 0
	}
	c.Set(dt, v+1)
}

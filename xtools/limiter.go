package xtools

import (
	"time"

	"github.com/bluele/gcache"
	"golang.org/x/time/rate"
)

type (
	AcLimiter struct {
		gc gcache.Cache
	}
)

func (h *AcLimiter) Allow(key string) bool {
	val, err := h.gc.Get(key)
	if err != nil {
		return true
	}
	l, isok := val.(*rate.Limiter)
	if isok == false {
		return true
	}
	return l.Allow()
}
func (h *AcLimiter) Allows(key ...string) bool {
	for _, v := range key {
		if h.Allow(v) == false {
			return false
		}
	}
	return true
}

func NewAcLimiter(maxEvent int, expire time.Duration, tokennum int) *AcLimiter {
	return &AcLimiter{
		gc: gcache.New(maxEvent).LRU().Expiration(expire).LoaderFunc(func(i2 interface{}) (i interface{}, err error) {
			return rate.NewLimiter(rate.Every(expire), tokennum), nil
		}).Build(),
	}
}

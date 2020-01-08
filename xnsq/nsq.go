package xnsq

import (
	"errors"

	"github.com/nsqio/go-nsq"
)

//NSQInit NSQInit
func NSQInit(conf *ProducerConfig) {
	if conf == nil {
		return
	}
	nsqconf = conf

}

//NSQ NSQ
func NSQ(name string) (*ProducerClient, error) {

	if nsqconf == nil {
		return &ProducerClient{}, errors.New("conf is empty")
	}
	conf := nsqconf.Get(name)
	if conf == nil {
		return &ProducerClient{}, errors.New("conf is not found :" + name)
	}
	nsqProducerLock.Lock()
	pr, isok := nsqProducer[conf.Host]
	nsqProducerLock.Unlock()
	if isok == false {
		tmp, err := nsq.NewProducer(conf.GetHost(), nsq.NewConfig())
		if err != nil {
			return &ProducerClient{}, err
		}
		nsqProducer[conf.Host] = tmp
		return &ProducerClient{producer: tmp, conf: conf}, nil
	} else {
		err := pr.Ping()
		if err != nil {
			return &ProducerClient{}, err
		}
	}

	return &ProducerClient{producer: pr, conf: conf}, nil
}

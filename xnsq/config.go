package xnsq

import (
	"encoding/json"
	"errors"
	"github.com/jinares/gopkg/xtools"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/nsqio/go-nsq"
)

type (
	//NsqNodeConfig NsqNodeConfig
	NsqNodeConfig struct {
		Host  string `yaml:"Host" json:"Host"`
		Topic string `yaml:"Topic" json:"Topic"`
		IsRun int    `yaml:"IsRun" json:"IsRun"` //0 # 1 停止推送
	}
	ProducerConfig map[string]*NsqNodeConfig
	//ProducerClient ProducerClient
	ProducerClient struct {
		producer *nsq.Producer
		conf     *NsqNodeConfig
	}
)

var (
	nsqconf         = &ProducerConfig{}
	nsqProducer     = map[string]*nsq.Producer{}
	nsqProducerLock sync.Mutex
)

func getIPPort(remoteaddr string) (string, string) {
	parts := strings.Split(remoteaddr, ":")

	partsLen := len(parts)
	if partsLen == 1 {
		return parts[0], "0"
	}

	port, _ := strconv.Atoi(parts[partsLen-1])
	ipStr := strings.Join(parts[:partsLen-1], "")

	if ip := net.ParseIP(ipStr); ip == nil && strings.Index(ipStr, "$") == 0 {
		oldIPStr := ipStr
		ipStr = os.Getenv(oldIPStr[1:])

	}

	return ipStr, xtools.ToStr(port)
}

//GetHost GetHost
func (h *NsqNodeConfig) GetHost() string {
	ip, port := getIPPort(h.Host)
	return ip + ":" + port
}

//Get get conf
func (h ProducerConfig) Get(name string) *NsqNodeConfig {
	nsqProducerLock.Lock()
	defer nsqProducerLock.Unlock()
	if item, isok := h[name]; isok {

		return item
	}

	return nil
}

//PublishJSON PublishJSON
func (p *ProducerClient) PublishJSON(v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return p.Publish(body)
}

//Publish PublishJSON
func (p *ProducerClient) Publish(v []byte) error {
	if p.producer == nil {
		return errors.New("NSQ生产者不存在")
	}
	if p.conf.IsRun == 1 {
		return errors.New("NSQ停止推送")
	}
	err := p.producer.Publish(p.conf.Topic, v)
	return err
}

//PublishAsyncJSON PublishAsync
func (p *ProducerClient) PublishAsyncJSON(v interface{}) error {
	if p.producer == nil {
		return errors.New("NSQ生产者不存在")
	}
	if p.conf.IsRun == 1 {
		return errors.New("NSQ停止推送")
	}
	body, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = p.producer.PublishAsync(p.conf.Topic, []byte(body), nil)

	return err
}

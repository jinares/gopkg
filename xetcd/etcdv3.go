package xetcd

import (
	"io/ioutil"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
	yaml "gopkg.in/yaml.v2"
)

//EtcdConfig EtcdConfig
type EtcdConfig struct {
	Peers          []string `yaml:"Peers" json:"Peers"`
	TimeoutSeconds int64    `yaml:"TimeoutSeconds" json:"TimeoutSeconds"`
	Root           string   `yaml:"Root" json:"Root"`
	UserName       string   `yaml:"Username" json:"Username"`
	Password       string   `yaml:"Password" json:"Password"`
	Cert           string   `yaml:"Cert" json:"Cert"`
	Key            string   `yaml:"Key" json:"Key"`
	CA             string   `yaml:"CA" json:"CA"`
}

//NewEtcdClientV3 NewEtcdClientV3
func NewEtcdClientV3(opt *EtcdConfig) (*clientv3.Client, error) {

	c := clientv3.Config{Username: opt.UserName, Password: opt.Password, Endpoints: opt.Peers, DialTimeout: time.Duration(opt.TimeoutSeconds) * time.Second}
	if opt.Cert != "" {
		tlsInfo := transport.TLSInfo{
			CertFile:      opt.Cert,
			KeyFile:       opt.Key,
			TrustedCAFile: opt.CA,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			panic(err)
		}
		c.TLS = tlsConfig
	}
	client, err := clientv3.New(c)
	if err != nil {
		return nil, err
	}

	return client, nil

}

//yaml "gopkg.in/yaml.v2"
func BindYaml(path string, bind interface{}) error {
	yamlFile, err := ioutil.ReadFile(path)

	err = yaml.Unmarshal(yamlFile, bind)
	return err
}

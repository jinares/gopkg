module github.com/jinares/gopkg

go 1.13

require (
	github.com/Shopify/sarama v1.24.1
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/grpc-ecosystem/grpc-gateway v1.11.2
	github.com/klauspost/cpuid v1.2.1 // indirect
	github.com/nsqio/go-nsq v1.0.7
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/uber/jaeger-client-go v2.20.1+incompatible
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	go.etcd.io/etcd v0.0.0-20191011172313-6d8052314b9e
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	google.golang.org/genproto v0.0.0-20191108220845-16a3f7862a1a
	google.golang.org/grpc v1.24.0
	gopkg.in/yaml.v2 v2.2.4
)

replace go.etcd.io/etcd => github.com/etcd-io/etcd v0.0.0-20191021022006-5dc12f27251a

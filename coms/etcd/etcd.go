package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/hulklab/yago"
	"go.etcd.io/etcd/clientv3"
	"io/ioutil"
	"log"
	"time"
)

type Etcd struct {
	*clientv3.Client
}

func Ins(id ...string) *Etcd {
	var name string

	if len(id) == 0 {
		name = "etcd"
	} else {
		name = id[0]
	}

	v := yago.Component.Ins(name, func() interface{} {
		val := initEtcdConn(name)

		return val
	})

	client := v.(*Etcd)

	return client
}

// init for etcd
func initEtcdConn(name string) *Etcd {
	endpoints := yago.Config.GetStringSlice(name + ".endpoints")
	dialTimeout := yago.Config.GetDuration(name+".dial_timeout") * time.Second
	username := yago.Config.GetString(name + ".username")
	password := yago.Config.GetString(name + ".password")
	etcdCert := yago.Config.GetString(name + ".cert_file")
	etcdCertKey := yago.Config.GetString(name + ".cert_key_file")
	etcdCa := yago.Config.GetString(name + ".cert_ca_file")
	maxCallRecvMsgSize := yago.Config.GetInt(name + ".max_call_recv_msgsize_byte")
	maxCallSendMsgSize := yago.Config.GetInt(name + ".max_call_send_msgsize_byte")

	if len(endpoints) == 0 {
		log.Fatalf("Fatal error: etcd endpoints is empty")
	}

	config := clientv3.Config{}
	config.Endpoints = endpoints
	config.DialTimeout = dialTimeout
	config.Username = username
	config.Password = password
	config.MaxCallRecvMsgSize = maxCallRecvMsgSize
	config.MaxCallSendMsgSize = maxCallSendMsgSize

	// tls
	if etcdCert != "" && etcdCertKey != "" {
		cert, err := tls.LoadX509KeyPair(etcdCert, etcdCertKey)
		if err != nil {
			log.Fatal(err)
		}

		pool := x509.NewCertPool()
		if etcdCa != "" {
			caData, err := ioutil.ReadFile(etcdCa)
			if err != nil {
				log.Fatal(err)
			}
			pool.AppendCertsFromPEM(caData)
		}

		config.TLS = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
	}

	etcd, err := clientv3.New(config)
	if err != nil {
		log.Fatal(err)
	}

	t := new(Etcd)
	t.Client = etcd

	return t
}

// +build go1.9

package etcd

import (
	"bytes"
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/namespace"
	"github.com/coreos/etcd/mvcc/mvccpb"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"strings"
	"time"
)

const (
//ETCD_WAIT_TIME = time.Second * 10
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type EtcdV3PropsConfigSource struct {
	kvs.MapProperties
	name    string
	root    string
	prefix  string
	client  *clientv3.Client
	kv      clientv3.KV
	watcher clientv3.Watcher
	config  *clientv3.Config
}

func NewEtcdV3PropsConfigSource(address, root string) *EtcdV3PropsConfigSource {
	return NewEtcdV3PropsConfigSourceByName("", address, root, ETCD_WAIT_TIME)
}

func NewEtcdV3PropsConfigSourceByName(name, urls, root string, timeout time.Duration) *EtcdV3PropsConfigSource {
	s := &EtcdV3PropsConfigSource{}
	if name == "" {
		name = strings.Join([]string{"etcd", urls, root}, ":")
	}
	s.name = name
	s.Values = make(map[string]string)
	s.root = root
	if strings.LastIndex(s.root, "") > 0 {
		s.root = s.root[:len(s.root)-1]
	}
	endpoints := strings.Split(urls, ",")
	cfg := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
		// set timeout per request to fail fast when the target endpoint is unavailable
	}
	c, err := clientv3.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	s.client = c
	s.kv = namespace.NewKV(c, root)
	s.prefix = "/"
	s.watcher = namespace.NewWatcher(c, root)
	s.init()
	return s
}

func NewEtcdV3PropsCompositeConfigSource(contexts []string, address string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "EtcdKevValue"
	for _, context := range contexts {
		c := NewEtcdV3PropsConfigSource(address, context)
		s.Add(c)
	}

	return s
}

func (s *EtcdV3PropsConfigSource) init() {
	s.findProperties(s.root, nil)
}

func (s *EtcdV3PropsConfigSource) watchContext() {

}

func (s *EtcdV3PropsConfigSource) Close() {
	s.client.Close()
}

func (s *EtcdV3PropsConfigSource) findProperties(parentPath string, children []*mvccpb.KeyValue) {
	prefix := s.prefix

	res, err := s.kv.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Error(err)
		return
	}

	for _, kv := range res.Kvs {
		//value := string(kv.Value)
		k := string(kv.Key)
		props := kvs.NewProperties()
		props.Load(bytes.NewReader(kv.Value))
		for _, key := range props.Keys() {
			val := props.GetDefault(key, "")
			pkey := strings.Join([]string{k, key}, ".")
			s.registerKeyValue(pkey, val)
		}

	}

}

func (s *EtcdV3PropsConfigSource) sanitizeKey(path string, context string) string {
	//key := strings.Replace(path, context+"/", "", -1)
	key := strings.Replace(path, "/", ".", -1)
	return key
}

func (s *EtcdV3PropsConfigSource) registerKeyValue(path, value string) {
	key := s.sanitizeKey(path, s.root)
	s.Set(key, value)

}

func (s *EtcdV3PropsConfigSource) Name() string {
	return s.name
}

package consul

import (
	"bytes"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"strings"
	"time"
)

//通过key/properties, key作为section，value为props格式内容，类似ini文件格式
//Deprecated
//看看ConsulConfigSource
type ConsulPropsConfigSource struct {
	kvs.MapProperties
	name   string
	root   string
	client *api.Client
	kv     *api.KV
	config *api.Config
}

//Deprecated
func NewConsulPropsConfigSource(address, root string) *ConsulPropsConfigSource {
	return NewConsulPropsConfigSourceByName("consul", address, root, CONSUL_WAIT_TIME)
}

//Deprecated
func NewConsulPropsConfigSourceByName(name, address, root string, timeout time.Duration) *ConsulPropsConfigSource {
	s := &ConsulPropsConfigSource{}
	if name == "" {
		name = strings.Join([]string{"consul", address, root}, ":")
	}
	s.name = name
	s.Values = make(map[string]string)
	s.root = root
	s.config = api.DefaultConfig()
	s.config.Address = address
	s.config.WaitTime = timeout
	client, err := api.NewClient(s.config)
	if err != nil {
		panic(err)
	}
	s.client = client
	s.kv = client.KV()
	s.init()
	return s
}

func (s *ConsulPropsConfigSource) init() {
	s.findProperties(s.root, nil)
}

func (s *ConsulPropsConfigSource) watchContext() {

}

func (s *ConsulPropsConfigSource) Close() {
}

func (s *ConsulPropsConfigSource) findProperties(parentPath string, children []string) {
	prefix := s.root
	q := &api.QueryOptions{}
	keys, _, err := s.kv.Keys(prefix, "", q)
	if err != nil {
		log.Error(err)
		return
	}
	for _, k := range keys {
		kv, _, err := s.kv.Get(k, q)
		if err != nil {
			log.Error(err)
			continue
		}
		//value := string(kv.Value)
		props := kvs.NewProperties()
		props.Load(bytes.NewReader(kv.Value))
		for _, key := range props.Keys() {
			val := props.GetDefault(key, "")
			pkey := strings.Join([]string{k, key}, ".")
			s.registerKeyValue(pkey, val)
		}

	}

}

func (s *ConsulPropsConfigSource) sanitizeKey(path string, context string) string {
	key := strings.Replace(path, context+"/", "", -1)
	key = strings.Replace(key, "/", ".", -1)
	return key
}

func (s *ConsulPropsConfigSource) registerKeyValue(path, value string) {
	key := s.sanitizeKey(path, s.root)
	s.Set(key, value)

}

func (s *ConsulPropsConfigSource) Name() string {
	return s.name
}

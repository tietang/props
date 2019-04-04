package consul

import (
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"github.com/tietang/props/yam"
	"strings"
	"time"
)

//通过key/properties, key所谓section，value为props格式内容，类似ini文件格式
type ConsulConfigSource struct {
	kvs.MapProperties
	name        string
	root        string
	client      *api.Client
	kv          *api.KV
	config      *api.Config
	ContentType kvs.ContentType
}

func NewConsulConfigSource(address, root string, contentType kvs.ContentType) *ConsulConfigSource {
	conf := NewConsulConfigSourceByName("consul", address, root, contentType, CONSUL_WAIT_TIME)

	return conf
}

func NewConsulConfigSourceByName(name, address, root string, contentType kvs.ContentType, timeout time.Duration) *ConsulConfigSource {
	s := &ConsulConfigSource{}
	s.Values = make(map[string]string)
	s.ContentType = contentType
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

func (s *ConsulConfigSource) init() {
	s.findProperties(s.root, nil)

}

func (s *ConsulConfigSource) watchContext() {

}

func (s *ConsulConfigSource) Close() {
}

func (s *ConsulConfigSource) findProperties(parentPath string, children []string) {
	if s.ContentType == kvs.ContentKV {
		s.findKeyValue(parentPath, children)
		return
	}
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
		content := string(kv.Value)

		if s.ContentType == kvs.ContentProps {
			s.findProps(content)
		} else if s.ContentType == kvs.ContentIni {
			s.findIni(content)
		} else if s.ContentType == kvs.ContentYaml {
			s.findYaml(content)
		} else {
			log.Warn("Unsupported format：", s.ContentType)
		}
		//value := string(kv.Value)
		//props := kvs.NewProperties()
		//props.Load(bytes.NewReader(kv.Value))
		//for _, key := range props.Keys() {
		//	val := props.GetDefault(key, "")
		//	pkey := strings.Join([]string{k, key}, ".")
		//	s.registerKeyValue(pkey, val)
		//}

	}

}
func (s *ConsulConfigSource) findYaml(content string) {
	props := yam.ByYaml(content)
	s.SetAll(props.Values)
}

func (s *ConsulConfigSource) findIni(content string) {
	props := ini.ByIni(content)
	s.SetAll(props.Values)
}

func (s *ConsulConfigSource) findProps(content string) {
	props := kvs.ByProperties(content)
	s.SetAll(props.Values)
}

func (s *ConsulConfigSource) findKeyValue(parentPath string, children []string) {
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
		value := string(kv.Value)
		s.registerKeyValue(k, value)
	}

}

func (s *ConsulConfigSource) sanitizeKey(path string, context string) string {
	key := strings.Replace(path, context+"/", "", -1)
	key = strings.Replace(key, "/", ".", -1)
	return key
}

func (s *ConsulConfigSource) registerKeyValue(path, value string) {
	key := s.sanitizeKey(path, s.root)
	s.Set(key, value)

}

func (s *ConsulConfigSource) Name() string {
	return s.name
}

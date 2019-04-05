package consul

import (
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"github.com/tietang/props/yam"
	"path"
	"strings"
	"time"
)

//通过key/ini_props, key所谓section，value为props格式内容，类似ini文件格式

//配置为ContentIniProps模式
// key作为section，value为props格式内容，类似ini文件格式
// key为实际key的prefix，会添加到前面
// 比如 root=configs/dev/app
// consul
// 		full key=configs/dev/app/mysql
// 		value=(x1=0 \n x2=1)
// 实际key/value为： mysql.x1=0 mysql.x2=1
//ContentProps,ContentYamlContentIni 模式时，
// 其 key无实际配置意义，只作为配置分组标识，
// 值为对应的内容格式类型，读取时会将对应的内容转换这种类型
// 可以通过key后缀来标识格式类型，默认按照properties来读取

type ConsulConfigSource struct {
	kvs.MapProperties
	name        string
	root        string
	client      *api.Client
	kv          *api.KV
	config      *api.Config
	ContentType kvs.ContentType
}

func NewConsulConfigSource(address, root string) *ConsulConfigSource {
	conf := NewConsulConfigSourceByName("consul", address, root, kvs.ContentAuto, CONSUL_WAIT_TIME)
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
	q := &api.QueryOptions{}
	s.kv.Get("", q)
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
		var ctype kvs.ContentType
		if s.ContentType == kvs.ContentAuto {
			key := path.Base(k)
			idx := strings.LastIndex(key, ".")
			if idx == -1 || idx == len(key)-1 {
				ctype = kvs.ContentProps
			} else {
				ctype = kvs.ContentType(key[idx+1:])
			}
		} else {
			ctype = s.ContentType
		}

		if ctype == kvs.ContentProps || ctype == kvs.ContentProperties {
			s.findProps(content)
		} else if ctype == kvs.ContentIniProps {
			s.findIniProps(k, content)
		} else if ctype == kvs.ContentIni {
			s.findIni(content)
		} else if ctype == kvs.ContentYaml || ctype == kvs.ContentYam || ctype == kvs.ContentYml {
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
	if props != nil {
		s.SetAll(props.Values)
	}
}

func (s *ConsulConfigSource) findIni(content string) {
	props := ini.ByIni(content)
	if props != nil {
		s.SetAll(props.Values)
	}
}

func (s *ConsulConfigSource) findProps(content string) {
	props := kvs.ByProperties(content)
	if props != nil {
		s.SetAll(props.Values)
	}
}
func (s *ConsulConfigSource) findIniProps(key, content string) {
	props := kvs.ByProperties(content)
	if props != nil {

		prefix := path.Base(key)
		for key, value := range props.Values {
			k := prefix + "." + key
			s.Set(k, value)
		}
	}
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

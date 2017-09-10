package props

import (
    "strings"
    "github.com/hashicorp/consul/api"
    log "github.com/sirupsen/logrus"
    "bytes"
)

//通过key/properties, key所谓section，value为props格式内容，类似ini文件格式
type ConsulIniConfigSource struct {
    MapProperties
    name   string
    root   string
    client *api.Client
    kv     *api.KV
    config *api.Config
}

func NewConsulIniConfigSource(address, root string) *ConsulIniConfigSource {
    return NewConsulIniConfigSourceByName("", address, root)
}

func NewConsulIniConfigSourceByName(name, address, root string) *ConsulIniConfigSource {
    s := &ConsulIniConfigSource{}
    if name == "" {
        name = strings.Join([]string{"consul", address, root}, ":")
    }
    s.name = name
    s.values = make(map[string]string)
    s.root = root
    s.config = api.DefaultConfig()
    s.config.Address = address
    client, err := api.NewClient(s.config)
    if err != nil {
        panic(err)
    }
    s.client = client
    s.kv = client.KV()
    s.init()
    return s
}

func (s *ConsulIniConfigSource) init() {
    s.findProperties(s.root, nil)
}

func (s *ConsulIniConfigSource) watchContext() {

}

func (s *ConsulIniConfigSource) Close() {
}

func (s *ConsulIniConfigSource) findProperties(parentPath string, children []string) {
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
        props := NewProperties()
        props.Load(bytes.NewReader(kv.Value))
        for _, key := range props.Keys() {
            val := props.GetDefault(key, "")
            pkey := strings.Join([]string{k, key}, ".")
            s.registerKeyValue(pkey, val)
        }

    }

}

func (s *ConsulIniConfigSource) sanitizeKey(path string, context string) string {
    key := strings.Replace(path, context+"/", "", -1)
    key = strings.Replace(key, "/", ".", -1)
    return key
}

func (s *ConsulIniConfigSource) registerKeyValue(path, value string) {
    key := s.sanitizeKey(path, s.root)
    s.Set(key, value)

}

func (s *ConsulIniConfigSource) Name() string {
    return s.name
}

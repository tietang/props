package props

import (
    "strings"
    "github.com/hashicorp/consul/api"
)

//
type ConsulKeyValueConfigSource struct {
    MapProperties
    name   string
    root   string
    client *api.Client
    kv     *api.KV
    config *api.Config
}

func NewConsulKeyValueConfigSource(name, address, root string) *ConsulKeyValueConfigSource {
    s := &ConsulKeyValueConfigSource{}
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

func (s *ConsulKeyValueConfigSource) init() {
    s.findProperties(s.root, nil)
}

func (s *ConsulKeyValueConfigSource) watchContext() {

}

func (s *ConsulKeyValueConfigSource) Close() {
}

func (s *ConsulKeyValueConfigSource) findProperties(parentPath string, children []string) {
    prefix := s.root
    q := &api.QueryOptions{}
    keys, _, err := s.kv.Keys(prefix, "", q)
    if err != nil {
        return
    }
    for _, k := range keys {
        kv, _, err := s.kv.Get(k, q)
        if err != nil {
            continue
        }
        value := string(kv.Value)
        s.registerKeyValue(k, value)
    }

}

func (s *ConsulKeyValueConfigSource) sanitizeKey(path string, context string) string {
    key := strings.Replace(path, context+"/", "", -1)
    key = strings.Replace(key, "/", ".", -1)
    return key
}

func (s *ConsulKeyValueConfigSource) registerKeyValue(path, value string) {
    key := s.sanitizeKey(path, s.root)
    s.Set(key, value)

}

func (s *ConsulKeyValueConfigSource) Name() string {
    return s.name
}

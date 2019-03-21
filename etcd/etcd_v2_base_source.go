package etcd

import (
    "strings"
    "time"
    "github.com/coreos/etcd/client"
    "context"
    "github.com/tietang/props/kvs"
    "github.com/Unknwon/log"
)

const (
    ETCD_WAIT_TIME = time.Second * 10
    DEFAULT_WATCH_KEY = "__notice"
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type EtcdV2ConfigSource struct {
    kvs.MapProperties
    name    string
    root    string
    client  client.Client
    kapi    client.KeysAPI
    watcher client.Watcher
    config  *client.Config
    Watched bool
}

func (s *EtcdV2ConfigSource) init() {
}

func (s *EtcdV2ConfigSource) WatchKey(key string, f func(node *client.Node)) {
    opts := &client.WatcherOptions{}

    for {
        w := s.kapi.Watcher(key, opts)
        res, err := w.Next(context.Background())
        if err != nil {
            log.Error("watch", err)
        }
        f(res.Node)
    }

}

func (s *EtcdV2ConfigSource) Close() {

}

func (s *EtcdV2ConfigSource) GetChildrenNodes(path string) client.Nodes {
    q := &client.GetOptions{}
    res, err := s.kapi.Get(context.Background(), path, q)
    if err != nil {
        return make(client.Nodes, 0)
    }
    node := res.Node
    return node.Nodes
}

func (s *EtcdV2ConfigSource) SanitizeKey(path string, context string) string {
    key := strings.Replace(path, context+"/", "", -1)
    key = strings.Replace(key, "/", ".", -1)
    return key
}

func (s *EtcdV2ConfigSource) RegisterKeyValue(path, value string) {
    key := s.SanitizeKey(path, s.root)
    s.Set(key, value)

}

func (s *EtcdV2ConfigSource) Name() string {
    return s.name
}

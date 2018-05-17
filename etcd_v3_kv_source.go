// +build go1.9

package props

import (
    "strings"
    log "github.com/sirupsen/logrus"
    "time"
    "context"
    "github.com/coreos/etcd/clientv3"
    "github.com/coreos/etcd/clientv3/namespace"
    "github.com/coreos/etcd/mvcc/mvccpb"
)

const (
//ETCD_WAIT_TIME = time.Second * 10
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type EtcdV3KeyValueConfigSource struct {
    MapProperties
    name    string
    root    string
    prefix  string
    client  *clientv3.Client
    kv      clientv3.KV
    watcher clientv3.Watcher
    config  *clientv3.Config
}

func NewEtcdV3KeyValueConfigSource(address, root string) *EtcdV3KeyValueConfigSource {
    return NewEtcdV3KeyValueConfigSourceByName("", address, root, CONSUL_WAIT_TIME)
}

func NewEtcdV3KeyValueConfigSourceByName(name, urls, root string, timeout time.Duration) *EtcdV3KeyValueConfigSource {
    s := &EtcdV3KeyValueConfigSource{}
    if name == "" {
        name = strings.Join([]string{"etcd", urls, root}, ":")
    }
    s.name = name
    s.values = make(map[string]string)
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

func NewEtcdV3KeyValueCompositeConfigSource(contexts []string, address string) *CompositeConfigSource {
    s := NewEmptyNoSystemEnvCompositeConfigSource()
    s.name = "EtcdKevValue"
    for _, context := range contexts {
        c := NewEtcdV3KeyValueConfigSource(address, context)
        s.Add(c)
    }

    return s
}

func (s *EtcdV3KeyValueConfigSource) init() {
    s.findProperties(s.root, nil)
}

func (s *EtcdV3KeyValueConfigSource) watchContext() {

}

func (s *EtcdV3KeyValueConfigSource) Close() {
    s.client.Close()
}

func (s *EtcdV3KeyValueConfigSource) findProperties(parentPath string, children []*mvccpb.KeyValue) {
    prefix := s.prefix

    res, err := s.kv.Get(context.Background(), prefix, clientv3.WithPrefix())
    if err != nil {
        log.Error(err)
        return
    }

    for _, kv := range res.Kvs {
        value := string(kv.Value)
        key := string(kv.Key)
        s.registerKeyValue(key, value)
    }

}

func (s *EtcdV3KeyValueConfigSource) sanitizeKey(path string, context string) string {
    //key := strings.Replace(path, context+"/", "", -1)
    key := strings.Replace(path, "/", ".", -1)
    return key
}

func (s *EtcdV3KeyValueConfigSource) registerKeyValue(path, value string) {
    key := s.sanitizeKey(path, s.root)
    s.Set(key, value)

}

func (s *EtcdV3KeyValueConfigSource) Name() string {
    return s.name
}

package etcd

import (
    "strings"
    log "github.com/sirupsen/logrus"
    "time"
    "github.com/coreos/etcd/client"
    "context"
    "github.com/tietang/props/kvs"
)

const (
    ETCD_WAIT_TIME = time.Second * 10
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type EtcdV2KeyValueConfigSource struct {
    kvs.MapProperties
    name    string
    root    string
    client  client.Client
    kapi    client.KeysAPI
    watcher client.Watcher
    config  *client.Config
}

func NewEtcdKeyValueConfigSource(address, root string) *EtcdV2KeyValueConfigSource {
    return NewEtcdKeyValueConfigSourceByName("", address, root, ETCD_WAIT_TIME)
}

func NewEtcdKeyValueConfigSourceByName(name, urls, root string, timeout time.Duration) *EtcdV2KeyValueConfigSource {
    s := &EtcdV2KeyValueConfigSource{}
    if name == "" {
        name = strings.Join([]string{"etcd", urls, root}, ":")
    }
    s.name = name
    s.Values = make(map[string]string)
    s.root = root
    endpoints := strings.Split(urls, ",")
    cfg := client.Config{
        Endpoints: endpoints,
        Transport: client.DefaultTransport,
        // set timeout per request to fail fast when the target endpoint is unavailable
        HeaderTimeoutPerRequest: timeout,
    }
    c, err := client.New(cfg)
    if err != nil {
        log.Fatal(err)
    }
    s.client = c

    s.kapi = client.NewKeysAPI(c)
    wo := &client.WatcherOptions{
        AfterIndex: 10,
        Recursive:  true,
    }
    s.watcher = s.kapi.Watcher(root, wo)
    s.init()
    return s
}

func NewEtcdKeyValueCompositeConfigSource(contexts []string, address string) *kvs.CompositeConfigSource {
    s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
    s.ConfName = "EtcdKevValue"
    for _, context := range contexts {
        c := NewEtcdKeyValueConfigSource(address, context)
        s.Add(c)
    }

    return s
}

func (s *EtcdV2KeyValueConfigSource) init() {
    s.findProperties(s.root, nil)
}

func (s *EtcdV2KeyValueConfigSource) watchContext() {

}

func (s *EtcdV2KeyValueConfigSource) Close() {

}

func (s *EtcdV2KeyValueConfigSource) findProperties(parentPath string, children client.Nodes) {
    if len(children) == 0 {
        children = s.getChildrenNodes(parentPath)
    }
    if len(children) == 0 {
        return
    }
    for _, node := range children {

        //fp := path.Join(parentPath, node.Key)
        fp := node.Key
        //fmt.Println(fp)
        chnodes := s.getChildrenNodes(fp)
        value := node.Value
        if !node.Dir {
            s.registerKeyValue(fp, value)
        } else {
            s.findProperties(fp, chnodes)
        }
        //
    }

}

func (s *EtcdV2KeyValueConfigSource) getChildrenNodes(path string) client.Nodes {
    q := &client.GetOptions{}
    res, err := s.kapi.Get(context.Background(), path, q)
    if err != nil {
        return make(client.Nodes, 0)
    }
    node := res.Node
    return node.Nodes
}

func (s *EtcdV2KeyValueConfigSource) sanitizeKey(path string, context string) string {
    key := strings.Replace(path, context+"/", "", -1)
    key = strings.Replace(key, "/", ".", -1)
    return key
}

func (s *EtcdV2KeyValueConfigSource) registerKeyValue(path, value string) {
    key := s.sanitizeKey(path, s.root)
    s.Set(key, value)

}

func (s *EtcdV2KeyValueConfigSource) Name() string {
    return s.name
}

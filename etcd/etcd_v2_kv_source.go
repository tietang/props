package etcd

import (
	"github.com/coreos/etcd/client"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"strings"
	"time"
)

var _ kvs.ConfigSource = new(EtcdV2KeyValueConfigSource)

// 通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type EtcdV2KeyValueConfigSource struct {
	EtcdV2ConfigSource
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

func (s *EtcdV2KeyValueConfigSource) findProperties(parentPath string, children client.Nodes) {
	if len(children) == 0 {
		children = s.GetChildrenNodes(parentPath)
	}
	if len(children) == 0 {
		return
	}

	for _, node := range children {
		fp := node.Key
		if s.Watched && strings.HasSuffix(fp, DEFAULT_WATCH_KEY) {
			log.Debug("WatchNodeDataChange: ", fp)
			s.WatchKey(fp, func(node *client.Node) {
				s.findProperties(fp, node.Nodes)
			})
		}
		//fp := filepath.Join(parentPath, node.Key)

		//fmt.Println(fp)
		chnodes := s.GetChildrenNodes(fp)
		value := node.Value
		if !node.Dir {
			s.RegisterKeyValue(fp, value)
		} else {
			s.findProperties(fp, chnodes)
		}
		//
	}

}

package etcd

import (
	"context"
	"github.com/coreos/etcd/client"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"strings"
	"time"
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type EtcdV2PropsConfigSource struct {
	kvs.MapProperties
	name   string
	root   string
	client client.Client
	kapi   client.KeysAPI
	config *client.Config
}

func NewEtcdPropsConfigSource(address, root string) *EtcdV2PropsConfigSource {
	return NewEtcdPropsConfigSourceByName("", address, root, ETCD_WAIT_TIME)
}

func NewEtcdPropsConfigSourceByName(name, urls, root string, timeout time.Duration) *EtcdV2PropsConfigSource {
	s := &EtcdV2PropsConfigSource{}
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
	s.init()
	return s
}

func NewEtcdPropsCompositeConfigSource(contexts []string, address string) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "EtcdKevValue"
	for _, context := range contexts {
		c := NewEtcdPropsConfigSource(address, context)
		s.Add(c)
	}

	return s
}

func (s *EtcdV2PropsConfigSource) init() {
	s.findProperties()
}

func (s *EtcdV2PropsConfigSource) watchContext() {

}

func (s *EtcdV2PropsConfigSource) Close() {

}

func (s *EtcdV2PropsConfigSource) findProperties() {
	children := s.getChildrenNodes(s.root)
	if len(children) == 0 {
		return
	}
	for _, p := range children {

		value := p.Value

		props := kvs.NewProperties()
		props.Load(strings.NewReader(value))
		for _, key := range props.Keys() {
			val := props.GetDefault(key, "")
			pkey := strings.Join([]string{p.Key, key}, ".")
			s.registerKeyValue(pkey, val)
		}
	}

}

func (s *EtcdV2PropsConfigSource) getChildrenNodes(path string) client.Nodes {
	q := &client.GetOptions{}
	res, err := s.kapi.Get(context.Background(), path, q)

	if err != nil {
		return make(client.Nodes, 0)
	}
	node := res.Node
	return node.Nodes
}

func (s *EtcdV2PropsConfigSource) sanitizeKey(path string, context string) string {
	key := strings.Replace(path, context+"/", "", -1)
	key = strings.Replace(key, "/", ".", -1)
	return key
}

func (s *EtcdV2PropsConfigSource) registerKeyValue(path, value string) {
	key := s.sanitizeKey(path, s.root)
	s.Set(key, value)

}

func (s *EtcdV2PropsConfigSource) Name() string {
	return s.name
}

package zk

import (
	"github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"strings"
)

/*
通过key/properties, key就是section，value为props格式内容，类似ini文件格式
例如：
context: /config/demo
zk nodes:

/config/demo
    - /jdbc:
       ```
        url=tcp(127.0.0.1:3306)/Test?charset=utf8
        username=root
        password=root
        timeout=6s

        ```
    - /redis:
       ```
        host=192.168.1.123
        port=6379
        database=d
        timeout=6s
        password=password
        ```


*/
//Deprecated
//看看ZookeeperConfigSource
type ZookeeperPropsConfigSource struct {
	ZookeeperSource
}

func NewZookeeperPropsConfigSource(name string, watched bool, context string, conn *zk.Conn) *ZookeeperPropsConfigSource {
	s := &ZookeeperPropsConfigSource{}
	s.Watched = watched
	s.name = name
	s.Values = make(map[string]string)
	s.conn = conn
	s.context = context
	s.initProperties()
	return s
}

func (s *ZookeeperPropsConfigSource) initProperties() {
	s.findProperties(s.context)
}

func (s *ZookeeperPropsConfigSource) findProperties(root string) {
	children := s.getChildren(root)
	if len(children) == 0 {
		return
	}
	for _, p := range children {

		fp := filepath.Join(root, p)
		value, err := s.getPropertiesValue(fp)
		if s.Watched && strings.HasSuffix(fp, DEFAULT_WATCH_KEY) {
			log.Debug("WatchNodeDataChange: ", fp)
			s.watchGet(fp)
		}

		if err == nil {
			props := kvs.NewProperties()
			props.Load(strings.NewReader(value))
			for _, key := range props.Keys() {
				val := props.GetDefault(key, "")
				pkey := strings.Join([]string{p, key}, ".")
				s.registerKeyValue(pkey, val)
			}
		}

	}

}

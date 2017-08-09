package props

import (
    "github.com/samuel/go-zookeeper/zk"
    "strings"
    "errors"
    "path"
)

const (
    ENCODING = "UTF-8"
)

type ZookeeperConfigSource struct {
    MapProperties
    name    string
    conn    *zk.Conn
    context string
}

func NewZookeeperConfigSource(name string, context string, conn *zk.Conn) *ZookeeperConfigSource {
    s := &ZookeeperConfigSource{}
    s.name = name
    s.values = make(map[string]string)
    s.conn = conn
    s.context = context
    s.init()
    return s
}

func (s *ZookeeperConfigSource) init() {
    s.findProperties(s.context, nil)
}

func (s *ZookeeperConfigSource) watchContext() {

}

func (s *ZookeeperConfigSource) Close() {
    s.conn.Close()
}

func (s *ZookeeperConfigSource) findProperties(parentPath string, children []string) {

    if len(children) == 0 {
        children = s.getChildren(parentPath)
    }
    if len(children) == 0 {
        return
    }
    for _, p := range children {

        fp := path.Join(parentPath, p)
        //fmt.Println(fp)
        chpath := s.getChildren(fp)
        value, err := s.getPropertiesValue(fp)
        if err == nil {
            s.registerKeyValue(fp, value)
        }
        //
        s.findProperties(fp, chpath)

    }

}

func (s *ZookeeperConfigSource) getPropertiesValue(path string) (string, error) {
    d, _, err := s.conn.Get(path)
    if err != nil || len(d) == 0 {
        return "", errors.New("not value")
    }
    return string(d), nil
}

func (s *ZookeeperConfigSource) getChildren(childPath string) []string {
    children, _, err := s.conn.Children(childPath)
    if err != nil {
        return make([]string, 0)
    }
    return children
}

func (s *ZookeeperConfigSource) sanitizeKey(path string, context string) string {
    key := strings.Replace(path, context+"/", "", -1)
    key = strings.Replace(key, "/", ".", -1)
    return key
}

func (s *ZookeeperConfigSource) registerKeyValue(path, value string) {
    key := s.sanitizeKey(path, s.context)
    s.Set(key, value)

}

func (s *ZookeeperConfigSource) Name() string {
    return s.name
}

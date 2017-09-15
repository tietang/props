package props

import (
    "github.com/samuel/go-zookeeper/zk"
    "strings"
    "errors"
    "path"
    "fmt"
    log "github.com/sirupsen/logrus"
    "bytes"
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
        database=2
        timeout=6s
        password=password
        ```


*/
type ZookeeperPropsConfigSource struct {
    MapProperties
    name    string
    conn    *zk.Conn
    context string
}

func NewZookeeperPropsConfigSource(name string, context string, conn *zk.Conn) *ZookeeperPropsConfigSource {
    s := &ZookeeperPropsConfigSource{}
    s.name = name
    s.values = make(map[string]string)
    s.conn = conn
    s.context = context
    s.init()
    return s
}

func (s *ZookeeperPropsConfigSource) init() {
    s.findProperties(s.context)
}

func (s *ZookeeperPropsConfigSource) Close() {
    s.conn.Close()
}

func (s *ZookeeperPropsConfigSource) findProperties(root string) {
    children := s.getChildren(root)
    if len(children) == 0 {
        return
    }
    for _, p := range children {

        fp := path.Join(root, p)
        value, err := s.getPropertiesValue(fp)

        if err == nil {
            props := NewProperties()
            props.Load(bytes.NewReader(value))
            for _, key := range props.Keys() {
                val := props.GetDefault(key, "")
                pkey := strings.Join([]string{p, key}, ".")
                s.registerKeyValue(pkey, val)
            }
        }

    }

}

func (s *ZookeeperPropsConfigSource) getPropertiesValue(path string) ([]byte, error) {
    d, _, err := s.conn.Get(path)
    if err != nil || len(d) == 0 {
        return nil, errors.New("not value")
    }
    return d, nil
}

func (s *ZookeeperPropsConfigSource) getChildren(childPath string) []string {
    children, _, err := s.conn.Children(childPath)
    if err != nil {
        return make([]string, 0)
    }
    return children
}

func (s *ZookeeperPropsConfigSource) sanitizeKey(path string, context string) string {
    key := strings.Replace(path, context+"/", "", -1)
    key = strings.Replace(key, "/", ".", -1)
    return key
}

func (s *ZookeeperPropsConfigSource) registerKeyValue(path, value string) {
    key := s.sanitizeKey(path, s.context)
    s.Set(key, value)

}

func (s *ZookeeperPropsConfigSource) Name() string {
    return s.name
}

func (s *ZookeeperPropsConfigSource) Watch(key string, handlers ... func([]string, zk.Event)) {
    go s.watchGet(path.Join(s.context, key, KEY_NOTIFY_NODE), handlers...)
}

func (s *ZookeeperPropsConfigSource) WatchChildren(key string, handlers ... func([]string, zk.Event)) {
    pathStr := path.Join(s.context, key)
    s.watchChildren(pathStr, handlers...)
}

func (s *ZookeeperPropsConfigSource) watchChildren(pathStr string, handlers ... func([]string, zk.Event)) {
    children, stat, ch, err := s.conn.ChildrenW(pathStr)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v %+v\n", children, stat)
    e := <-ch

    s.findProperties(e.Path)
    for _, handler := range handlers {
        handler(children, e)
    }
    fmt.Printf("%+v\n", e)
    s.watchChildren(pathStr, handlers...)
}

func (g *ZookeeperPropsConfigSource) watchGet(pathStr string, handlers ... func([]string, zk.Event)) {
    log.Info(pathStr)
    exists, _, _ := g.conn.Exists(pathStr)
    if !exists {
        g.conn.Create(pathStr, []byte("1"), 1, nil)
    }
    _, stat, ch, err := g.conn.GetW(pathStr)
    children, _, err := g.conn.Children(pathStr)
    if err != nil {
        panic(err)
    }
    log.Infof("watch: %+v %+v\n", children, stat)
    e := <-ch

    //pPath:=path.Dir(e.Path)
    g.findProperties(e.Path)
    for _, handler := range handlers {
        handler(children, e)
    }
    log.Infof("notify event: %+v\n ", e)
    g.watchGet(pathStr, handlers...)
}

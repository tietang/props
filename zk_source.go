package props

import (
    "github.com/samuel/go-zookeeper/zk"
    "strings"
    "errors"
    "path"
    "fmt"
    log "github.com/sirupsen/logrus"
    "time"
    "sync"
)

const (
    ENCODING = "UTF-8"

    KEY_NOTIFY_NODE = "notify"
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type ZookeeperConfigSource struct {
    MapProperties
    name     string
    conn     *zk.Conn
    context  string
    watchers sync.Map
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

func NewZookeeperCompositeConfigSource(contexts []string, connStr []string, timeout time.Duration) *CompositeConfigSource {

    conn, ch, err := zk.Connect(connStr, timeout)
    if err != nil {
        log.Error(err)
        panic(err)
    }
    for {
        event := <-ch
        fmt.Println(event)
        if event.State == zk.StateConnected {
            break
        }
    }
    return NewZookeeperCompositeConfigSourceByConn(contexts, conn)
}

func NewZookeeperCompositeConfigSourceByConn(contexts []string, conn *zk.Conn) *CompositeConfigSource {
    s := NewEmptyCompositeConfigSource()
    s.name = "Zookeeper"
    for _, context := range contexts {
        zkms := NewZookeeperConfigSource("zk:"+context, context, conn)
        s.Add(zkms)
    }
    return s
}

func (s *ZookeeperConfigSource) init() {
    s.findProperties(s.context, nil)
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

func (s *ZookeeperConfigSource) Watch(key string, handlers ... func(children []string, event zk.Event)) {
    go s.watchGet(path.Join(s.context, key, KEY_NOTIFY_NODE), handlers...)
}

func (s *ZookeeperConfigSource) WatchChildren(key string, handlers ... func(children []string, event zk.Event)) {
    pathStr := path.Join(s.context, key)
    s.watchChildren(pathStr, handlers...)
}

func (s *ZookeeperConfigSource) watchChildren(pathStr string, handlers ... func(children []string, event zk.Event)) {
    children, stat, ch, err := s.conn.ChildrenW(pathStr)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v %+v\n", children, stat)
    e := <-ch

    s.findProperties(path.Dir(e.Path), nil)
    for _, handler := range handlers {
        handler(children, e)
    }
    fmt.Printf("%+v\n", e)
    s.watchChildren(pathStr, handlers...)
}

func (g *ZookeeperConfigSource) watchGet(pathStr string, handlers ... func(children []string, event zk.Event)) {
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
    g.findProperties(path.Dir(e.Path), nil)
    for _, handler := range handlers {
        handler(children, e)
    }
    log.Infof("notify event: %+v\n ", e)
    g.watchGet(pathStr, handlers...)
}

func (g *ZookeeperConfigSource) WatchAndRefresh(key string, t interface{}) {
    //g.watchers.Store(key, t)
    g.Watch(key, func(children []string, event zk.Event) {
        g.Unmarshal(t)
    })

}

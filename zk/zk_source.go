package zk

import (
    "github.com/samuel/go-zookeeper/zk"
)

//通过key/value来组织，过滤root prefix后，替换/为.作为properties key
type ZookeeperConfigSource struct {
    ZookeeperSource
}

func NewZookeeperConfigSource(name string, watched bool, context string, conn *zk.Conn) *ZookeeperConfigSource {
    s := &ZookeeperConfigSource{}
    s.name = name
    s.Watched = watched
    s.Values = make(map[string]string)
    s.conn = conn
    s.context = context
    s.init()
    return s
}

func (s *ZookeeperConfigSource) init() {
    s.findChildProperties(s.context, nil)
}

func (s *ZookeeperConfigSource) Close() {
    s.conn.Close()
}

func (s *ZookeeperConfigSource) findChildProperties(parentPath string, children []string) {
    s.findProperties(s.Watched, parentPath, children)

    //if len(children) == 0 {
    //    children = s.getChildren(parentPath)
    //}
    //if len(children) == 0 {
    //    return
    //}
    //for _, p := range children {
    //
    //    fp := path.Join(parentPath, p)
    //    //fmt.Println(fp)
    //    chpath := s.getChildren(fp)
    //    value, err := s.getPropertiesValue(fp)
    //    if err == nil {
    //        s.registerKeyValue(fp, value)
    //    }
    //    //
    //    s.findChildProperties(fp, chpath)
    //
    //}

}

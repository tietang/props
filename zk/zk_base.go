package zk

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"path"
	"path/filepath"
	"strings"
)

const (
	ENCODING = "UTF-8"

	KEY_NOTICE_NODE = "notice"
)

var _ kvs.ConfigSource = new(ZookeeperSource)

type ZookeeperSource struct {
	kvs.MapProperties
	name    string
	conn    *zk.Conn
	context string
	Watched bool
}

func (s *ZookeeperSource) Close() {
	s.conn.Close()
}

func (s *ZookeeperSource) getPropertiesValue(path string) (string, error) {
	d, _, err := s.conn.Get(path)
	if err != nil || len(d) == 0 {
		return "", errors.New("not value")
	}
	return string(d), nil
}

func (s *ZookeeperSource) getChildren(childPath string) []string {
	children, _, err := s.conn.Children(childPath)
	if err != nil {
		return make([]string, 0)
	}
	return children
}

func (s *ZookeeperSource) sanitizeKey(keyPath string, context string) string {
	context = filepath.Join(context) + "/"
	key := strings.TrimPrefix(keyPath, context)
	key = strings.Replace(key, "/", ".", -1)

	//key := strings.Replace(keyPath, context+"/", "", -1)
	//key = strings.Replace(key, "/", ".", -1)
	return key
}

func (s *ZookeeperSource) registerKeyValue(keyPath, value string) {
	key := s.sanitizeKey(keyPath, s.context)
	s.Set(key, value)

}

func (s *ZookeeperSource) Name() string {
	return s.name
}

func (s *ZookeeperSource) Watch(key string, handlers ...func(children []string, event zk.Event)) {
	go s.watchGet(path.Join(s.context, key, KEY_NOTICE_NODE), handlers...)
}

func (s *ZookeeperSource) WatchChildren(key string, handlers ...func(children []string, event zk.Event)) {
	pathStr := filepath.Join(s.context, key)
	s.watchChildren(pathStr, handlers...)
}

func (s *ZookeeperSource) watchChildren(pathStr string, handlers ...func(children []string, event zk.Event)) {
	children, stat, ch, err := s.conn.ChildrenW(pathStr)
	if err != nil {
		//panic(err)
		log.Error(err)
		return
	}
	fmt.Printf("%+v %+v\n", children, stat)
	e := <-ch

	s.findProperties(false, filepath.Dir(e.Path), nil)
	for _, handler := range handlers {
		handler(children, e)
	}
	fmt.Printf("%+v\n", e)
	s.watchChildren(pathStr, handlers...)
}

func (g *ZookeeperSource) watchGet(pathStr string, handlers ...func(children []string, event zk.Event)) {
	log.Info(pathStr)
	exists, _, _ := g.conn.Exists(pathStr)
	if !exists {
		g.conn.Create(pathStr, []byte("1"), 1, nil)
	}
	_, stat, ch, err := g.conn.GetW(pathStr)
	children, _, err := g.conn.Children(pathStr)
	if err != nil {
		//panic(err)
		log.Error(err)
	}
	log.Infof("watch: %+v %+v\n", children, stat)
	e := <-ch

	//pPath:=path.Dir(e.Path)
	g.findProperties(false, filepath.Dir(e.Path), nil)
	for _, handler := range handlers {
		handler(children, e)
	}
	log.Infof("notify event: %+v\n ", e)
	g.watchGet(pathStr, handlers...)
}

func (g *ZookeeperSource) WatchAndRefresh(key string, t interface{}) {
	g.Watch(key, func(children []string, event zk.Event) {
		g.Unmarshal(t)
	})

}

func (s *ZookeeperSource) findProperties(isWatched bool, parentPath string, children []string) {

	if isWatched && strings.HasSuffix(parentPath, DEFAULT_WATCH_KEY) {
		log.Debug("WatchNodeDataChange: ", parentPath)
		s.watchGet(parentPath)
	}

	if len(children) == 0 {
		children = s.getChildren(parentPath)
	}
	if len(children) == 0 {
		return
	}
	for _, p := range children {

		fp := filepath.Join(parentPath, p)
		//fmt.Println(fp)
		chpath := s.getChildren(fp)
		value, err := s.getPropertiesValue(fp)
		if err == nil {
			s.registerKeyValue(fp, value)
		}
		//
		s.findProperties(isWatched, fp, chpath)

	}

}

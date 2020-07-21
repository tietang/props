package zk

import (
	"github.com/samuel/go-zookeeper/zk"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/v3/kvs"
	"time"
)

func NewZookeeperCompositeConfigSource(contexts []string, connStr []string, timeout time.Duration) *kvs.CompositeConfigSource {
	conn, ch, err := zk.Connect(connStr, timeout)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	for {
		event := <-ch
		log.Info(event)
		if event.State == zk.StateConnected {
			log.Info("zookeeper connected. ", connStr)
			break
		}
	}
	return NewZookeeperCompositeConfigSourceByConn(contexts, conn)
}

func NewZookeeperCompositeConfigSourceByConn(contexts []string, conn *zk.Conn) *kvs.CompositeConfigSource {
	s := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
	s.ConfName = "Zookeeper"
	for _, context := range contexts {
		zkms := NewZookeeperConfigSource(false, context, conn)
		s.Add(zkms)
	}
	return s
}

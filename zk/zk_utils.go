package zk

import (
    "github.com/samuel/go-zookeeper/zk"
    "path"
    log "github.com/sirupsen/logrus"
)

const (
    DEFAULT_WATCH_KEY = "__notice"
)

var flags = int32(0)
var acl = zk.WorldACL(zk.PermAll)

func ZkCreateString(conn *zk.Conn, path string, value string) (string, error) {
    return ZkCreate(conn, path, []byte(value))
}

func ZkCreate(conn *zk.Conn, nodePath string, value []byte) (string, error) {

    d, _ := path.Split(nodePath)
    ppath := path.Clean(d)
    if !ZkExits(conn, ppath) {
        ZkCreate(conn, ppath, []byte(""))
    }
    return conn.Create(nodePath, []byte(value), flags, acl)
}

func ZkExits(conn *zk.Conn, path string) bool {
    b, _, e := conn.Exists(path)
    if e != nil {
        return false
    }
    return b
}

func ZkWatchNodeCreated(conn *zk.Conn, path string) {
    log.Println("watchNodeCreated")
    for {
        _, _, ch, _ := conn.ExistsW(path)
        e := <-ch
        log.Println("ExistsW:", e.Type, "Event:", e)
        if e.Type == zk.EventNodeCreated {
            log.Println("NodeCreated ")
            return
        }
    }
}
func ZkWatchNodeDeleted(conn *zk.Conn, path string) {
    log.Println("watchNodeDeleted")
    for {
        _, _, ch, _ := conn.ExistsW(path)
        e := <-ch
        log.Println("ExistsW:", e.Type, "Event:", e)
        if e.Type == zk.EventNodeDeleted {
            log.Println("NodeDeleted ")
            return
        }
    }
}

func ZkWatchNodeDataChange(conn *zk.Conn, path string) {
    for {
        _, _, ch, _ := conn.GetW(path)
        e := <-ch
        log.Println("GetW('"+path+"'):", e.Type, "Event:", e)
    }
}

func ZkWatchChildrenChanged(conn *zk.Conn, path string) {
    for {
        _, _, ch, _ := conn.ChildrenW(path)
        e := <-ch
        log.Println("ChildrenW:", e.Type, "Event:", e)
    }
}
